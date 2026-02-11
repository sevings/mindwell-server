package server

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/leporo/sqlf"
	"go.uber.org/zap"
	"github.com/sevings/mindwell-server/lib/auth"
	"github.com/sevings/mindwell-server/lib/database"
	"github.com/sevings/mindwell-server/lib/helpers"
	"github.com/sevings/mindwell-server/lib/notifications"
	"github.com/sevings/mindwell-server/lib/pubsub"
	"github.com/sevings/mindwell-server/lib/validation"

	"github.com/BurntSushi/toml"
	"github.com/carlescere/scheduler"
	"github.com/go-openapi/runtime/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	cache "github.com/patrickmn/go-cache"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations"
	goconf "github.com/zpatrick/go-config"
	"golang.org/x/text/language"
)

const (
	logApi      = "api"
	logTelegram = "telegram"
	logEmail    = "email"
	logSystem   = "system"
	logRequest  = "request"
)

type MindwellServer struct {
	DB    *sql.DB
	PS    *pubsub.PubSub
	API   *operations.MindwellAPI
	Ntf   *notifications.CompositeNotifier
	Eac   validation.EmailAllowedChecker
	Imgs  *cache.Cache
	log   *zap.Logger
	cfg   *goconf.Config
	local *i18n.Localizer
	errs  map[string]*i18n.Message
}

func NewMindwellServer(api *operations.MindwellAPI, configPath string) *MindwellServer {
	logger, err := zap.NewProduction(zap.WithCaller(false))
	if err != nil {
		log.Println(err)
	}

	config := helpers.LoadConfig(configPath)
	db := database.OpenDatabase(config)

	trFile, err := config.String("server.tr_file")
	if err != nil {
		logger.Error(err.Error(), zap.String("type", "system"))
	}

	bundle := i18n.NewBundle(language.Russian)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile(trFile)

	srv := &MindwellServer{
		DB:    db,
		API:   api,
		Imgs:  cache.New(2*24*time.Hour, 24*time.Hour),
		log:   logger,
		cfg:   config,
		local: i18n.NewLocalizer(bundle),
		errs: map[string]*i18n.Message{
			"no_entry":       {ID: "no_entry", Other: "Entry not found or you have no access rights."},
			"no_comment":     {ID: "no_comment", Other: "Comment not found or you have no access rights."},
			"no_tlog":        {ID: "no_tlog", Other: "Tlog not found or you have no access rights."},
			"no_theme":       {ID: "no_theme", Other: "Theme not found or you have no access rights."},
			"no_chat":        {ID: "no_chat", Other: "Chat not found or you have no access rights."},
			"no_message":     {ID: "no_message", Other: "Message not found or you have no access rights."},
			"no_wish":        {ID: "no_wish", Other: "Wish not found or it's not for you."},
			"no_request":     {ID: "no_friend_request", Other: "You have no friend request from this user."},
			"invalid_invite": {ID: "invalid_invite", Other: "Invite is invalid."},
		},
	}

	srv.PS = pubsub.NewPubSub(database.ConnectionString(config), srv.LogSystem())
	srv.PS.Start()

	srv.Ntf = notifications.NewCompositeNotifier(srv, srv.PS)
	srv.Eac = validation.NewEmailChecker(srv)

	if _, err := scheduler.Every().Day().At("03:10").Run(srv.recalcKarma); err != nil {
		srv.LogSystem().Error(err.Error())
	}
	if _, err := scheduler.Every().Day().At("03:15").Run(srv.giveInvites); err != nil {
		srv.LogSystem().Error(err.Error())
	}
	if _, err := scheduler.Every(8).Hours().Run(srv.searchAlts); err != nil {
		srv.LogSystem().Error(err.Error())
	}

	return srv
}

func (srv *MindwellServer) GetDB() *sql.DB {
	return srv.DB
}

func (srv *MindwellServer) ConfigStrings(field string) []string {
	value, err := srv.cfg.String(field)
	if err != nil {
		srv.LogSystem().Warn(err.Error())
	}

	if value == "" {
		return nil
	}

	return strings.Split(value, ";")
}

func (srv *MindwellServer) ConfigString(field string) string {
	value, err := srv.cfg.String(field)
	if err != nil {
		srv.LogSystem().Warn(err.Error())
	}

	return value
}

func (srv *MindwellServer) ConfigOptString(field string) string {
	value, err := srv.cfg.StringOr(field, "")
	if err != nil {
		srv.LogSystem().Warn(err.Error())
	}

	return value
}

func (srv *MindwellServer) ConfigBytes(field string) []byte {
	return []byte(srv.ConfigString(field))
}

func (srv *MindwellServer) ConfigInt(field string) int {
	value, err := srv.cfg.Int(field)
	if err != nil {
		srv.LogSystem().Warn(err.Error())
	}

	return value
}

func (srv *MindwellServer) ConfigInts(field string) []int {
	var values []int

	for _, str := range srv.ConfigStrings(field) {
		value, err := strconv.ParseInt(str, 10, 32)
		if err == nil {
			values = append(values, int(value))
		} else {
			srv.LogSystem().Warn(err.Error())
		}
	}

	return values
}

func (srv *MindwellServer) ConfigInt64(field string) int64 {
	strValue := srv.ConfigString(field)
	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		srv.LogSystem().Warn(err.Error())
	}

	return value
}

func (srv *MindwellServer) ConfigInt64s(field string) []int64 {
	var values []int64

	for _, str := range srv.ConfigStrings(field) {
		value, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			values = append(values, value)
		} else {
			srv.LogSystem().Warn(err.Error())
		}
	}

	return values
}

func (srv *MindwellServer) ConfigBool(field string) bool {
	value, err := srv.cfg.Bool(field)
	if err != nil {
		srv.LogSystem().Warn(err.Error())
	}

	return value
}

func (srv *MindwellServer) NewAvatar(avatar string) *models.Avatar {
	base := srv.ConfigString("images.base_url") + "avatars/"

	return &models.Avatar{
		X42:  base + "42/" + avatar,
		X92:  base + "92/" + avatar,
		X124: base + "124/" + avatar,
	}
}

func (srv *MindwellServer) NewCover(id int64, cover string) *models.Cover {
	if len(cover) == 0 {
		cnt := int64(srv.ConfigInt("images.covers"))
		n := int(id % cnt)
		cover = "default/" + strconv.Itoa(n) + ".jpg"
	}

	base := srv.ConfigString("images.base_url")

	return &models.Cover{
		ID:    id,
		X1920: base + "covers/1920/" + cover,
		X318:  base + "covers/318/" + cover,
	}
}

func (srv *MindwellServer) ImagesFolder() string {
	return srv.ConfigString("images.folder")
}

func (srv *MindwellServer) Transact(txFunc func(*database.AutoTx) middleware.Responder) middleware.Responder {
	return database.Transact(srv.DB, txFunc)
}

func (srv *MindwellServer) TokenHash() auth.TokenHash {
	return auth.NewTokenHash(srv)
}

// NewError returns error object with some message
func (srv *MindwellServer) NewError(msg *i18n.Message) *models.Error {
	if msg == nil {
		msg = &i18n.Message{ID: "internal_error", Other: "Internal server error."}
	}

	message, err := srv.local.Localize(&i18n.LocalizeConfig{DefaultMessage: msg})
	if err != nil {
		srv.LogApi().Error(err.Error())
	}

	return &models.Error{Message: message}
}

func (srv *MindwellServer) StandardError(name string) *models.Error {
	msg := srv.errs[name]
	if msg == nil {
		srv.LogApi().Sugar().Error("Standard error not found:", name)
	}

	return srv.NewError(msg)
}

func (srv *MindwellServer) Log(tpe string) *zap.Logger {
	return srv.log.With(zap.String("type", tpe))
}

func (srv *MindwellServer) LogApi() *zap.Logger {
	return srv.Log(logApi)
}

func (srv *MindwellServer) LogTelegram() *zap.Logger {
	return srv.Log(logTelegram)
}

func (srv *MindwellServer) LogEmail() *zap.Logger {
	return srv.Log(logEmail)
}

func (srv *MindwellServer) LogSystem() *zap.Logger {
	return srv.Log(logSystem)
}

func (srv *MindwellServer) LogRequest() *zap.Logger {
	return srv.Log(logRequest)
}

func (srv *MindwellServer) recalcKarma() {
	start := time.Now()

	_, err := srv.DB.Exec("SELECT recalc_karma()")
	if err != nil {
		srv.LogSystem().Error(err.Error())
	}

	srv.LogSystem().Info("Karma recalculated",
		zap.Int64("duration", time.Since(start).Microseconds()),
	)
}

func (srv *MindwellServer) giveInvites() {
	start := time.Now()

	tx := database.NewAutoTx(srv.DB)
	defer tx.Finish()

	tx.Query("SELECT user_id FROM give_invites()")

	var ids []int
	for {
		var id int
		if !tx.Scan(&id) {
			break
		}

		ids = append(ids, id)
	}

	for _, id := range ids {
		srv.Ntf.SendNewInvite(tx, id)
	}

	srv.LogSystem().Info("Given invites",
		zap.Int64("duration", time.Since(start).Microseconds()),
		zap.Int("invites", len(ids)),
	)
}

func (srv *MindwellServer) searchAlts() {
	start := time.Now()

	tx := database.NewAutoTx(srv.DB)
	defer tx.Finish()

	namesQ := sqlf.Select("name").
		From("users").
		Where("now() - last_seen_at < interval '8 hours'").
		Where("alt_of IS NULL").
		Where("NOT confirmed_alt")

	var found []string
	names := tx.QueryStmt(namesQ).ScanStrings()
	for _, name := range names {
		altQ := sqlf.Select("ul.name, COUNT(*) AS cnt").
			From("user_log AS ul").
			Join("user_log AS ol", "ul.uid2 = ol.uid2").
			Where("ul.name <> lower(?)", name).
			Where("ol.name = lower(?)", name).
			Where("ol.uid2 <> 0").
			GroupBy("ul.name").
			OrderBy("cnt DESC").
			Limit(1)

		var alt string
		var cnt int64
		tx.QueryStmt(altQ).Scan(&alt, &cnt)
		if cnt == 0 {
			oldQ := sqlf.Select("TRUE").
				From("users").
				Where("lower(name) = lower(?)", name).
				Where("age(created_at) > interval '1 year'").
				Where("karma > 1")

			isOld := tx.QueryStmt(oldQ).ScanBool()
			if isOld {
				updQ := sqlf.Update("users").
					Set("confirmed_alt", true).
					Where("lower(name) = lower(?)", name)

				tx.ExecStmt(updQ)
				if tx.Error() != nil {
					srv.LogSystem().Error(tx.Error().Error())
				}
			}
		} else {
			const idQ = "SELECT id FROM users WHERE lower(name) = lower($1)"
			id1 := tx.QueryInt64(idQ, alt)
			id2 := tx.QueryInt64(idQ, name)
			if id1 > id2 {
				continue
			}

			found = append(found, name)

			updQ := sqlf.Update("users").
				Set("alt_of", alt).
				Where("lower(name) = lower(?)", name)

			tx.ExecStmt(updQ)
			if tx.Error() != nil {
				srv.LogSystem().Error(tx.Error().Error())
			}
		}
	}

	for _, user := range found {
		srv.Ntf.Tg.SendPossibleAlts(tx, user)
	}

	srv.LogSystem().Info("Searched alts",
		zap.Int64("duration", time.Since(start).Microseconds()),
		zap.Int("alts", len(found)),
	)
}
