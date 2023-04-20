package utils

import (
	"database/sql"
	"fmt"
	"github.com/leporo/sqlf"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	"github.com/sevings/mindwell-server/models"
)

const errorText = "Что-то пошло не так…"
const unrecognisedText = "Неизвестная команда. Попробуй /help."

var tgHtmlEsc = strings.NewReplacer(
	"<", "&lt;",
	">", "&gt;",
	"&", "&amp;",
	"\"", "&quot;",
	"'", "&quot;",
	"\r", "",
)

var idRe = regexp.MustCompile(`(\d+)(?:#comments)?$`)
var loginRe = regexp.MustCompile(`/users/([0-9\-_]*[a-zA-Z][a-zA-Z0-9\-_]*)(?:/entries)?$`)

func extractEntryID(url string) (int64, string) {
	match := idRe.FindStringSubmatch(url)
	if len(match) < 2 {
		return 0, "Укажи ID записи."
	}

	id, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return 0, err.Error()
	}

	return id, ""
}

func extractLogin(url string) string {
	if !strings.Contains(url, "/") {
		return url
	}

	match := loginRe.FindStringSubmatch(url)
	if len(match) < 2 {
		return ""
	}

	return match[1]
}

type TelegramBot struct {
	srv    *MindwellServer
	api    *tgbotapi.BotAPI
	url    string
	ipAPI  string
	log    *zap.Logger
	admins []int
	group  int64
	logins *cache.Cache
	cmts   *cache.Cache
	msgs   *cache.Cache
	send   chan func()
	stop   chan interface{}
}

type messageID struct {
	chat int64
	msg  int
}

type messageIDs []messageID

func NewTelegramBot(srv *MindwellServer) *TelegramBot {
	bot := &TelegramBot{
		srv:    srv,
		url:    srv.ConfigString("server.base_url"),
		ipAPI:  srv.ConfigOptString("server.ip_api"),
		log:    srv.LogTelegram(),
		admins: srv.ConfigInts("telegram.admins"),
		group:  srv.ConfigInt64("telegram.admin_group"),
		logins: cache.New(10*time.Minute, 10*time.Minute),
		cmts:   cache.New(12*time.Hour, 1*time.Hour),
		msgs:   cache.New(12*time.Hour, 1*time.Hour),
		send:   make(chan func(), 200),
		stop:   make(chan interface{}),
	}

	if bot.ipAPI == "" {
		bot.ipAPI = "http://%s/"
	}

	bot.ipAPI = `<a href="` + bot.ipAPI + `">%s</a>`

	go bot.run()

	return bot
}

func (bot *TelegramBot) sendMessageNow(chat int64, text string) tgbotapi.Message {
	msg := tgbotapi.NewMessage(chat, text)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"
	message, err := bot.api.Send(msg)
	if err != nil {
		bot.log.Error(err.Error())
	}

	return message
}

func (bot *TelegramBot) sendMessage(chat int64, text string) {
	bot.send <- func() { bot.sendMessageNow(chat, text) }
}

func (bot *TelegramBot) isAdmin(upd *tgbotapi.Update) bool {
	chat := upd.Message.From.ID

	for _, admin := range bot.admins {
		if admin == chat {
			return true
		}
	}

	return false
}

func (bot *TelegramBot) run() {
	token := bot.srv.ConfigString("telegram.token")
	if len(token) == 0 {
		return
	}

	client := http.DefaultClient

	api, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		bot.log.Error(err.Error())
		return
	}

	bot.api = api
	// api.Debug = true

	bot.log.Sugar().Infof("Running Telegram bot %s", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.api.GetUpdatesChan(u)
	if err != nil {
		bot.log.Error(err.Error())
	}

	for {
		select {
		case send := <-bot.send:
			send()
		case <-bot.stop:
			return
		case upd := <-updates:
			if upd.Message == nil || !upd.Message.IsCommand() {
				continue
			}

			bot.command(upd)
		}
	}
}

func (bot *TelegramBot) command(upd tgbotapi.Update) {
	cmd := upd.Message.Command()
	bot.log.Info("update",
		zap.String("cmd", cmd),
		zap.String("from", upd.Message.From.UserName),
	)

	var reply string
	switch cmd {
	case "start":
		reply = bot.start(&upd)
	case "login":
		reply = bot.login(&upd)
	case "logout":
		reply = bot.logout(&upd)
	case "help":
		reply = bot.help(&upd)
	case "hide":
		reply = bot.hide(&upd)
	case "ban":
		reply = bot.ban(&upd)
	case "unban":
		reply = bot.unban(&upd)
	case "create":
		reply = bot.create(&upd)
	case "info":
		reply = bot.info(&upd)
	case "alts":
		reply = bot.alts(&upd)
	case "votes":
		reply = bot.votes(&upd)
	case "stat":
		reply = bot.stat(&upd)
	default:
		reply = unrecognisedText
	}

	bot.sendMessageNow(upd.Message.Chat.ID, reply)
}

func (bot *TelegramBot) Stop() {
	if bot.api == nil {
		return
	}

	bot.api.StopReceivingUpdates()
	close(bot.stop)
}

func (bot *TelegramBot) BuildToken(userID int64) string {
	token := GenerateString(8)
	bot.logins.SetDefault(token, userID)
	return token
}

func (bot *TelegramBot) VerifyToken(token string) int64 {
	userID, found := bot.logins.Get(token)
	if !found {
		return 0
	}

	bot.logins.Delete(token)
	return userID.(int64)
}

func (bot *TelegramBot) start(upd *tgbotapi.Update) string {
	token := upd.Message.CommandArguments()
	if len(token) == 0 {
		return `Привет! Я могу отправлять тебе уведомления из Mindwell.
Чтобы начать, скопируй ключ со <a href="` + bot.url + `account/notifications">страницы настроек</a>.
Отправь его мне, используя команду <code>/login &lt;ключ&gt;</code>. Так ты подтвердишь свой аккаунт.
Чтобы я забыл твой номер в Телеграме, введи /logout.`
	}

	return bot.login(upd)
}

func (bot *TelegramBot) login(upd *tgbotapi.Update) string {
	if upd.Message.From == nil {
		return errorText
	}

	token := upd.Message.CommandArguments()
	userID := bot.VerifyToken(token)

	if userID == 0 {
		return `Скопируй верный ключ со <a href="` + bot.url + `account/notifications">своей страницы настроек</a>.`
	}

	const q = `
		UPDATE users
		SET telegram = $2
		WHERE id = $1
		RETURNING show_name
	`

	var name string
	err := bot.srv.DB.QueryRow(q, userID, upd.Message.Chat.ID).Scan(&name)
	if err != nil {
		bot.log.Error(err.Error())
		return errorText
	}

	return "Привет, " + name + "! Теперь я буду отправлять тебе уведомления из Mindwell. " +
		"Используй команду /logout, если захочешь прекратить."
}

func (bot *TelegramBot) logout(upd *tgbotapi.Update) string {
	if upd.Message.From == nil {
		return errorText
	}

	from := upd.Message.From.ID

	const q = `
		UPDATE users
		SET telegram = NULL
		WHERE telegram = $1
		RETURNING show_name
	`

	var name string
	err := bot.srv.DB.QueryRow(q, from).Scan(&name)
	if err == nil {
		return "Я больше не буду беспокоить тебя, " + name + "."
	} else if err == sql.ErrNoRows {
		return "К этому номеру не привязан аккаунт в Mindwell."
	} else {
		bot.log.Error(err.Error())
		return errorText
	}
}

func (bot *TelegramBot) help(upd *tgbotapi.Update) string {
	text := `Я бот для получения уведомлений из Mindwell. Доступные команды:
<code>/login &lt;ключ&gt;</code> — авторизоваться с использованием автоматически сгенерированного пароля. Его можно получить на <a href="` + bot.url + `account/notifications">странице настроек</a>.
/logout — не получать больше уведомления на этот аккаунт.
/help — вывести данную краткую справку.`

	if bot.isAdmin(upd) {
		text += `

Команды администрирования:
<code>/hide {id или ссылка}</code> — скрыть запись.
<code>/ban {live | vote | comment | invite | adm} {N} {login или ссылка}</code> — запретить пользователю выбранные действия на N дней, в случае adm — навсегда.
<code>/ban user {N} {login или ссылка}</code> — заблокировать пользователя на N дней.
<code>/unban {login или ссылка}</code> — разблокировать пользователя.
<code>/info {email, login или ссылка}</code> — информация о пользователе.
<code>/alts {email, login или ссылка}</code> — искать альтернативные аккаунты пользователя.
<code>/delete {email}</code> — удалить пользователя.
<code>/create app {dev_name} {public | private} {code | password} {redirect_uri} {name} {show_name} {platform} {info}</code> - создать приложение.
<code>/create theme {name} {creator}</code> — создать тему.
/stat — статистика сервера.
`
	}

	return text
}

func (bot *TelegramBot) hide(upd *tgbotapi.Update) string {
	if !bot.isAdmin(upd) {
		return unrecognisedText
	}

	if upd.Message.From == nil {
		return errorText
	}

	url := upd.Message.CommandArguments()
	id, errStr := extractEntryID(url)
	if errStr != "" {
		return errStr
	}

	atx := NewAutoTx(bot.srv.DB)
	defer atx.Finish()

	const q = "UPDATE entries SET visible_for = (SELECT id FROM entry_privacy WHERE type = 'me') WHERE id = $1"
	atx.Exec(q, id)
	if atx.RowsAffected() < 1 {
		return "Запись не найдена."
	}

	return "Запись скрыта."
}

func (bot *TelegramBot) ban(upd *tgbotapi.Update) string {
	if !bot.isAdmin(upd) {
		return unrecognisedText
	}

	if upd.Message.From == nil {
		return errorText
	}

	args := strings.Split(upd.Message.CommandArguments(), " ")
	if len(args) < 2 {
		return "Укажи необходимые аргументы."
	}

	url := args[len(args)-1]
	login := extractLogin(url)
	if login == "" {
		return "Укажи логин пользователя."
	}

	dayCount := args[len(args)-2]
	banCount := len(args) - 2
	if _, err := strconv.Atoi(dayCount); err != nil {
		dayCount = "7300"
		banCount++
	}
	banUntil := "CURRENT_DATE + interval '" + dayCount + " days'"
	banTypes := args[:banCount]
	if len(banTypes) == 0 {
		return "Укажи необходимые ограничения."
	}

	q := sqlf.Update("users").
		Where("lower(name) = lower(?)", login)
	for _, ban := range banTypes {
		switch ban {
		case "live":
			q.SetExpr("live_ban", banUntil)
		case "vote":
			q.SetExpr("vote_ban", banUntil)
		case "comment":
			q.SetExpr("comment_ban", banUntil)
		case "invite":
			q.SetExpr("invite_ban", banUntil)
		case "adm":
			q.Set("adm_ban", true)
		case "user":
			q.SetExpr("user_ban", banUntil)
		default:
			q.Close()
			return "Неизвестный аргумент: " + ban + "."
		}
	}

	atx := NewAutoTx(bot.srv.DB)
	defer atx.Finish()

	atx.ExecStmt(q)
	if atx.RowsAffected() < 1 {
		return "Пользователь " + login + " не найден."
	}

	return "Пользователь " + login + " ограничен в правах на " + dayCount + " дней."
}

func (bot *TelegramBot) unban(upd *tgbotapi.Update) string {
	if !bot.isAdmin(upd) {
		return unrecognisedText
	}

	if upd.Message.From == nil {
		return errorText
	}

	url := upd.Message.CommandArguments()
	login := extractLogin(url)
	if login == "" {
		return "Укажи логин пользователя."
	}

	atx := NewAutoTx(bot.srv.DB)
	defer atx.Finish()

	const q = "UPDATE users SET user_ban = CURRENT_DATE WHERE lower(name) = lower($1) RETURNING id"
	id := atx.QueryInt64(q, login)
	if id == 0 {
		return "Пользователь " + login + " не найден."
	}

	return "Пользователь " + login + " разблокирован."
}

func (bot *TelegramBot) create(upd *tgbotapi.Update) string {
	if !bot.isAdmin(upd) {
		return unrecognisedText
	}

	if upd.Message.From == nil {
		return errorText
	}

	args := strings.Split(upd.Message.CommandArguments(), "\n")
	if len(args) < 2 {
		return "Укажи необходимые аргументы."
	}

	switch args[0] {
	case "app":
		return bot.createApp(args[1:])
	case "theme":
		return bot.createTheme(args[1:])
	}

	return unrecognisedText
}

func (bot *TelegramBot) createApp(args []string) string {
	if len(args) < 8 {
		return "Укажи необходимые аргументы."
	}

	devName := strings.TrimSpace(args[0])
	appType := strings.TrimSpace(args[1])
	appFlow := strings.TrimSpace(args[2])
	redirectUri := strings.TrimSpace(args[3])
	appName := strings.TrimSpace(args[4])
	showName := strings.TrimSpace(args[5])
	platform := strings.TrimSpace(args[6])
	info := strings.TrimSpace(args[7])

	var secretHash []byte
	var secret string

	switch appType {
	case "private":
		secret = GenerateString(64)
		secretHash = bot.srv.TokenHash().AppSecretHash(secret)
	case "public":
		secret = "(не сгенерирован)"
	default:
		return "Тип приложения может быть только public или private."
	}

	flow := 1
	switch appFlow {
	case "code":
		flow += 2
	case "password":
		flow += 4
	default:
		return "Тип авторизации может быть только code или password."
	}

	uriRe := regexp.MustCompile(`^\w+://[^#]+$`)
	if !uriRe.MatchString(redirectUri) {
		return "Redirect uri должен содержать схему и не содержать якорь."
	}

	appName = strings.TrimSpace(appName)
	nameRe := regexp.MustCompile(`^\w+$`)
	if !nameRe.MatchString(appName) {
		return "Название приложения может содержать только латинские буквы, цифры и знак подчеркивания."
	}

	tx := NewAutoTx(bot.srv.DB)
	defer tx.Finish()

	devName = strings.TrimSpace(devName)
	dev, err := LoadUserIDByName(tx, devName)
	if err != nil {
		return "Пользователь " + devName + " не найден."
	}

	appID := int64(rand.Int31())

	const query = `
INSERT INTO apps(id, secret_hash, redirect_uri, developer_id, flow, name, show_name, platform, info)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

	tx.Exec(query, appID, secretHash, redirectUri, dev.ID, flow,
		appName, showName, platform, info)
	if tx.Error() != nil {
		bot.log.Error(tx.Error().Error())
		return "Не удалось создать приложение."
	}

	text := "Приложение создано."
	text += "\n<b>client id</b>: " + strconv.FormatInt(appID, 10)
	text += "\n<b>client secret</b>: " + secret
	text += "\n<b>app name</b>: " + appName
	text += "\n<b>app show name</b>: " + showName
	text += "\n<b>developer name</b>: " + dev.Name
	text += "\n<b>redirect uri</b>: " + redirectUri
	text += "\n<b>flow</b>: " + appFlow
	text += "\n<b>platform</b>: " + platform
	text += "\n<b>info</b>: " + info

	return text
}

func (bot *TelegramBot) createTheme(args []string) string {
	if len(args) < 2 {
		return "Укажи необходимые аргументы."
	}

	name := strings.TrimSpace(args[0])
	creator := strings.TrimSpace(args[1])

	tx := NewAutoTx(bot.srv.DB)
	defer tx.Finish()

	const rankQ = "SELECT COUNT(*) + 1 FROM users WHERE creator_id IS NOT NULL AND karma >= 0"
	rank := tx.QueryInt64(rankQ)

	const q = `
		INSERT INTO users 
		(name, show_name, email, password_hash, creator_id, rank)
		values($1, $1, $1, '', (SELECT id FROM users WHERE lower(name) = lower($2)), $3)`

	tx.Exec(q, name, creator, rank)
	if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
		bot.log.Error(tx.Error().Error())
		return "Не удалось создать тему."
	}

	return fmt.Sprintf(`Тема создана: <a href="%sthemes/%s">%s</a>`, bot.url, name, name)
}

func (bot *TelegramBot) info(upd *tgbotapi.Update) string {
	if !bot.isAdmin(upd) {
		return unrecognisedText
	}

	if upd.Message.From == nil {
		return errorText
	}

	arg := upd.Message.CommandArguments()
	arg = extractLogin(arg)
	if arg == "" {
		return "Укажи логин или адрес почты."
	}

	atx := NewAutoTx(bot.srv.DB)
	defer atx.Finish()

	const q = `
SELECT users.id, users.name, users.show_name, created_at, 
	email, verified, rank, karma,
	invited.name, invited.show_name,
	entries_count, followers_count, followings_count, comments_count, invited_count,
	invite_ban, vote_ban, comment_ban, live_ban, adm_ban
FROM users
LEFT JOIN (SELECT id, name, show_name FROM users) AS invited ON users.invited_by = invited.id
WHERE lower(users.email) = lower($1) OR lower(users.name) = lower($1)`

	atx.Query(q, arg)

	var id int64
	var name, showName, email string
	var invitedByName, invitedByShowName sql.NullString
	var verified bool
	var createdAt time.Time
	var rank int64
	var karma float64
	var entries, followers, followings, comments, invited int64
	var inviteBan, voteBan, commentBan, liveBan time.Time
	var admBan bool
	atx.Scan(&id, &name, &showName, &createdAt,
		&email, &verified, &rank, &karma,
		&invitedByName, &invitedByShowName,
		&entries, &followers, &followings, &comments, &invited,
		&inviteBan, &voteBan, &commentBan, &liveBan, &admBan)

	if id == 0 {
		return "Пользователь с логином или адресом почты " + arg + " не найден."
	}

	today := time.Now()

	var invitedByLink string
	if invitedByName.Valid && invitedByShowName.Valid {
		invitedByLink = bot.userNameLink(invitedByName.String, invitedByShowName.String)
	} else {
		invitedByLink = "(not invited)"
	}

	var text string
	text += "\n<b>id</b>: " + strconv.FormatInt(id, 10)
	text += "\n<b>url</b>: " + bot.userNameLink(name, showName)
	text += "\n<b>email</b>: " + email
	text += "\n<b>verified</b>: " + strconv.FormatBool(verified)
	text += "\n<b>created at</b>: " + createdAt.Format("15:04:05 02 Jan 2006 MST")
	text += "\n<b>rank</b>: " + strconv.FormatInt(rank, 10)
	text += "\n<b>karma</b>: " + strconv.FormatFloat(karma, 'f', 2, 64)
	text += "\n<b>invited by</b>: " + invitedByLink
	text += "\n<b>entries</b>: " + strconv.FormatInt(entries, 10)
	text += "\n<b>followers</b>: " + strconv.FormatInt(followers, 10)
	text += "\n<b>followings</b>: " + strconv.FormatInt(followings, 10)
	text += "\n<b>comments</b>: " + strconv.FormatInt(comments, 10)
	text += "\n<b>invited</b>: " + strconv.FormatInt(invited, 10)
	text += "\n<b>invite ban</b>: " + strconv.FormatBool(inviteBan.After(today)) + ", " + inviteBan.Format("02 Jan 2006")
	text += "\n<b>vote ban</b>: " + strconv.FormatBool(voteBan.After(today)) + ", " + voteBan.Format("02 Jan 2006")
	text += "\n<b>comment ban</b>: " + strconv.FormatBool(commentBan.After(today)) + ", " + commentBan.Format("02 Jan 2006")
	text += "\n<b>live ban</b>: " + strconv.FormatBool(liveBan.After(today)) + ", " + liveBan.Format("02 Jan 2006")
	text += "\n<b>adm ban</b>: " + strconv.FormatBool(admBan)

	return text
}

func (bot *TelegramBot) alts(upd *tgbotapi.Update) string {
	if !bot.isAdmin(upd) {
		return unrecognisedText
	}

	if upd.Message.From == nil {
		return errorText
	}

	arg := upd.Message.CommandArguments()
	arg = extractLogin(arg)
	if arg == "" {
		return "Укажи логин."
	}

	users := strings.Split(arg, " ")

	limit := 10
	if l, err := strconv.ParseInt(users[len(users)-1], 10, 8); err == nil {
		limit = int(l)
		users = users[:len(users)-1]
	}

	if len(users) > 2 {
		return "Укажи не более двух аккаунтов."
	}

	emails := make([]string, 2, 2)

	atx := NewAutoTx(bot.srv.DB)
	defer atx.Finish()

	q := `SELECT email FROM users WHERE lower(name) = lower($1)`
	emails[0] = atx.QueryString(q, users[0])
	if emails[0] == "" {
		return "Пользователь " + users[0] + " не найден."
	}

	if len(users) == 1 {
		return bot.possibleAlts(atx, users[0], emails[0], limit)
	}

	emails[1] = atx.QueryString(q, users[1])
	if emails[1] == "" {
		return "Пользователь " + users[1] + " не найден."
	}

	return bot.compareUsers(atx, users[0], users[1], emails[0], emails[1], limit)
}

func (bot *TelegramBot) possibleAlts(atx *AutoTx, user, email string, limit int) string {
	getAlts := func(q *sqlf.Stmt) string {
		atx.QueryStmt(q)

		var alts []string
		for {
			var cnt int64
			var alt string
			ok := atx.Scan(&alt, &cnt)
			if !ok {
				break
			}

			alt = fmt.Sprintf(`%s (%d)`, bot.userNameLink(alt, alt), cnt)
			alts = append(alts, alt)
		}

		return strings.Join(alts, ", ")
	}

	q := sqlf.Select("ul.name, COUNT(*) AS cnt").
		From("user_log AS ul").
		Where("ul.name <> lower(?)", user).
		Where("ol.name = lower(?)", user).
		Where("ul.first <> ol.first").
		GroupBy("ul.name").
		OrderBy("cnt DESC").
		Limit(limit)

	ipQuery := q.Clone().
		Join("user_log AS ol", "ul.ip = ol.ip").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'")
	ipAlts := getAlts(ipQuery)

	devQuery := q.Clone().
		Join("user_log AS ol", "ul.device = ol.device").
		Where("ol.device <> 0").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'")
	devAlts := getAlts(devQuery)

	appQuery := q.Clone().
		Join("user_log AS ol", "ul.app = ol.app AND ul.device = ol.device").
		Where("ol.app <> 0 AND ol.device <> 0").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'")
	appAlts := getAlts(appQuery)

	uidQuery := q.Clone().
		Join("user_log AS ol", "ul.uid = ol.uid").
		Where("ol.uid <> 0")
	uidAlts := getAlts(uidQuery)

	uid2Query := q.Clone().
		Join("user_log AS ol", "ul.uid2 = ol.uid2").
		Where("ol.uid2 <> 0")
	uid2Alts := getAlts(uid2Query)

	emailSubQuery := sqlf.Select("id, name").
		Select("to_search_string(?) <<-> to_search_string(email) AS dist", email).
		From("users").OrderBy("dist ASC, id DESC").Limit(limit + 1)
	emailQuery := sqlf.Select("name").
		Select("100 - round(dist * 100)").
		From("").SubQuery("(", ") AS u", emailSubQuery).
		Where("lower(name) <> lower(?)", user).
		Where("dist < 0.5")
	emailAlts := getAlts(emailQuery)

	text := "Possible accounts of" + bot.userNameLink(user, user)
	text += "\n<b>By IP</b>: " + ipAlts
	text += "\n<b>By device</b>: " + devAlts
	text += "\n<b>By app</b>: " + appAlts
	text += "\n<b>By UID</b>: " + uidAlts
	text += "\n<b>By UID2</b>: " + uid2Alts
	text += "\n<b>By email</b>: " + emailAlts

	return text
}

func (bot *TelegramBot) compareUsers(atx *AutoTx, userA, userB, emailA, emailB string, limit int) string {
	getCountsFmt := func(q *sqlf.Stmt, format func(string) string) string {
		atx.QueryStmt(q)

		var result []string
		for {
			var cnt int64
			var from, to time.Time
			var data string
			ok := atx.Scan(&cnt, &from, &to, &data)
			if !ok {
				break
			}

			data = fmt.Sprintf(format(data)+" (%d, %s — %s)", cnt,
				from.Format("02.01.2006"), to.Format("02.01.2006"))
			result = append(result, data)
		}

		return strings.Join(result, "\n")
	}

	getCounts := func(q *sqlf.Stmt) string {
		return getCountsFmt(q, func(data string) string { return data })
	}

	fmtIP := func(ip string) string {
		return fmt.Sprintf(bot.ipAPI, ip, ip)
	}

	commonQuery := func() *sqlf.Stmt {
		return sqlf.Select("COUNT(*) AS cnt, MIN(ul.at), MAX(ul.at)").
			From("user_log AS ul").
			Where("ul.name = lower(?)", userA).
			Where("ol.name = lower(?)", userB).
			Where("ul.first <> ol.first").
			OrderBy("cnt DESC").
			Limit(limit)
	}

	commonIpsQ := commonQuery().
		Join("user_log AS ol", "ul.ip = ol.ip").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'").
		Select("ul.ip").
		GroupBy("ul.ip")
	commonIps := getCountsFmt(commonIpsQ, fmtIP)

	commonAppsQ := commonQuery().
		Join("user_log AS ol", "ul.app = ol.app AND ul.device = ol.device").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'").
		Select("ul.user_agent").
		GroupBy("ul.user_agent")
	commonApps := getCounts(commonAppsQ)

	diffQuery := func(a, b string) *sqlf.Stmt {
		sub := sqlf.Select("*").From("user_log").Where("name = ?", b)
		return sqlf.Select("COUNT(*) AS cnt, MIN(ol.at), MAX(ol.at)").
			From("").SubQuery("(", ") AS ul", sub).
			Where("ol.name = lower(?)", a).
			OrderBy("cnt DESC").
			Limit(limit)
	}

	diffIpsQ := func(a, b string) *sqlf.Stmt {
		return diffQuery(a, b).
			RightJoin("user_log AS ol", "ul.ip = ol.ip").
			Where("ul.ip IS NULL").
			Select("ol.ip").
			GroupBy("ol.ip")
	}

	diffIpsA := getCountsFmt(diffIpsQ(userA, userB), fmtIP)
	diffIpsB := getCountsFmt(diffIpsQ(userB, userA), fmtIP)

	diffAppsQ := func(a, b string) *sqlf.Stmt {
		return diffQuery(a, b).
			RightJoin("user_log AS ol", "ul.app = ol.app AND ul.device = ol.device").
			Where("ul.app IS NULL").
			Select("ol.user_agent").
			GroupBy("ol.user_agent")
	}

	diffAppsA := getCounts(diffAppsQ(userA, userB))
	diffAppsB := getCounts(diffAppsQ(userB, userA))

	text := fmt.Sprintf(`Comparison of %s and %s`,
		bot.userNameLink(userA, userA), bot.userNameLink(userB, userB))
	text += "\n<b>Common IPs</b>:\n" + commonIps
	text += "\n<b>Common apps</b>:\n" + commonApps
	text += "\n<b>IPs used only by " + userA + "</b>:\n" + diffIpsA
	text += "\n<b>IPs used only by " + userB + "</b>:\n" + diffIpsB
	text += "\n<b>Apps used only by " + userA + "</b>:\n" + diffAppsA
	text += "\n<b>Apps used only by " + userB + "</b>:\n" + diffAppsB
	text += "\n<b>" + userA + "'s email</b>: " + emailA
	text += "\n<b>" + userB + "'s email</b>: " + emailB

	return text
}

func (bot *TelegramBot) votes(upd *tgbotapi.Update) string {
	if !bot.isAdmin(upd) {
		return unrecognisedText
	}

	if upd.Message.From == nil {
		return errorText
	}

	url := upd.Message.CommandArguments()
	id, errStr := extractEntryID(url)
	if errStr != "" {
		return errStr
	}

	atx := NewAutoTx(bot.srv.DB)
	defer atx.Finish()

	const entryQuery = `SELECT name
FROM entries
JOIN users on author_id = users.id
WHERE entries.id = $1`

	author := atx.QueryString(entryQuery, id)
	if author == "" {
		return fmt.Sprintf("Запись %d не найдена.", id)
	}

	text := fmt.Sprintf("Голоса за запись %d автора %s\n", id, bot.userNameLink(author, author))

	const votesQuery = `SELECT name, vote > 0, entry_votes.created_at
FROM entry_votes
JOIN users ON user_id = users.id
WHERE entry_id = $1
ORDER BY created_at DESC`

	atx.Query(votesQuery, id)
	for {
		var name string
		var positive bool
		var createdAt time.Time
		if !atx.Scan(&name, &positive, &createdAt) {
			break
		}

		vote := "за"
		if !positive {
			vote = "против"
		}
		at := createdAt.Format("02.01.2006 15:04")
		text += fmt.Sprintf(`%s %s (%s)`, bot.userNameLink(name, name), vote, at)
		text += "\n"
	}

	return text
}

func (bot *TelegramBot) stat(upd *tgbotapi.Update) string {
	if !bot.isAdmin(upd) {
		return unrecognisedText
	}

	if upd.Message.From == nil {
		return errorText
	}

	atx := NewAutoTx(bot.srv.DB)
	defer atx.Finish()

	var text string

	addInt64 := func(key string, value int64) {
		text += "\n<b>" + key + "</b>: " + strconv.FormatInt(value, 10)
	}

	addFloat64 := func(key string, value float64) {
		text += "\n<b>" + key + "</b>: " + strconv.FormatFloat(value, 'f', 2, 64)
	}

	const usersQuery = `SELECT count(*) FROM users WHERE last_seen_at > created_at`
	users := atx.QueryInt64(usersQuery)
	addInt64("users", users)

	const invitedUsersQuery = `SELECT count(*) FROM users WHERE invited_by IS NOT NULL`
	invitedUsers := atx.QueryInt64(invitedUsersQuery)
	addInt64("invited users", invitedUsers)

	const negKarmaUsersQuery = `SELECT count(*) FROM users WHERE karma < -1`
	negKarmaUsers := atx.QueryInt64(negKarmaUsersQuery)
	addInt64("users with karma &lt; -1", negKarmaUsers)

	const genderUsersQuery = `
SELECT gender.type AS sex, count(*)
FROM users
JOIN gender ON users.gender = gender.id
WHERE users.last_seen_at > users.created_at
GROUP BY sex
ORDER BY sex`
	atx.Query(genderUsersQuery)
	for {
		var gender string
		var count int64
		if !atx.Scan(&gender, &count) {
			break
		}

		addInt64(gender+" gender users", count)
	}

	const newUsersMonthQuery = `SELECT count(*) FROM users WHERE now() - created_at < interval '1 month' AND last_seen_at > created_at`
	newUsersMonth := atx.QueryInt64(newUsersMonthQuery)
	addInt64("last month new users", newUsersMonth)

	const onlineUsersNowQuery = `SELECT count(*) FROM users WHERE is_online(last_seen_at)`
	onlineUsersNow := atx.QueryInt64(onlineUsersNowQuery)
	addInt64("online users", onlineUsersNow)

	const onlineUsersWeekQuery = `SELECT count(*) FROM users WHERE now() - last_seen_at < interval '7 days' AND last_seen_at > created_at`
	onlineUsersWeek := atx.QueryInt64(onlineUsersWeekQuery)
	addInt64("last week online users", onlineUsersWeek)

	const onlineUsersMonthQuery = `SELECT count(*) FROM users WHERE now() - last_seen_at < interval '1 month' AND last_seen_at > created_at`
	onlineUsersMonth := atx.QueryInt64(onlineUsersMonthQuery)
	addInt64("last month online users", onlineUsersMonth)

	const postingUsersMonthQuery = `
SELECT count(distinct author_id)
FROM entries
WHERE now() - created_at < interval '1 month'`
	postingUsersMonth := atx.QueryInt64(postingUsersMonthQuery)
	addInt64("last month posting users", postingUsersMonth)

	const chatsQuery = `
SELECT count(*)
FROM chats
JOIN messages ON last_message = messages.id
WHERE last_message > 0 AND messages.author_id <> 1`
	chats := atx.QueryInt64(chatsQuery)
	addInt64("user chats", chats)

	const avgEntriesQuery = `
SELECT count(*) / 7.0
FROM entries
WHERE created_at::date < current_date
	AND created_at::date >= current_date - interval '7 days'`
	avgEntries := atx.QueryFloat64(avgEntriesQuery)
	addFloat64("avg entries", avgEntries)

	const avgCommentsQuery = `
SELECT count(*) / 7.0
FROM comments
WHERE created_at::date < current_date
	AND created_at::date >= current_date - interval '7 days'`
	avgComments := atx.QueryFloat64(avgCommentsQuery)
	addFloat64("avg comments", avgComments)

	const avgMessagesQuery = `
SELECT count(*) / 7.0
FROM messages
WHERE created_at::date < current_date
	AND created_at::date >= current_date - interval '7 days'
	AND author_id <> 1`
	avgMessages := atx.QueryFloat64(avgMessagesQuery)
	addFloat64("avg user messages", avgMessages)

	if atx.Error() != nil {
		return errorText
	}

	return text
}

func idToString(id int64) string {
	return strconv.FormatInt(id, 32)
}

func (bot *TelegramBot) comment(entryTitle string, cmt *models.Comment) (cmtID, text string) {
	cmtID = idToString(cmt.ID)

	link := bot.url + "entries/" + strconv.FormatInt(cmt.EntryID, 10) + "#comments"

	text = tgHtmlEsc.Replace(cmt.Author.ShowName) + " пишет: \n" +
		"«" + tgHtmlEsc.Replace(cmt.EditContent) + "»\n"

	if len(entryTitle) == 0 {
		text += `К <a href="` + link + `">записи</a>`
	} else {
		text += `<a href="` + link + `">` + entryTitle + `</a>`
	}

	return
}

func (bot *TelegramBot) SendNewComment(chat int64, entryTitle string, cmt *models.Comment) {
	if bot.api == nil {
		return
	}

	cmtID, text := bot.comment(entryTitle, cmt)

	bot.send <- func() {
		msg := bot.sendMessageNow(chat, text)
		if msg.MessageID <= 0 {
			return
		}

		var msgIDs messageIDs
		msgIDsVar, found := bot.cmts.Get(cmtID)
		if found {
			msgIDs = msgIDsVar.(messageIDs)
		}
		msgIDs = append(msgIDs, messageID{chat, msg.MessageID})

		bot.cmts.SetDefault(cmtID, msgIDs)
	}
}

func (bot *TelegramBot) SendUpdateComment(entryTitle string, cmt *models.Comment) {
	if bot.api == nil {
		return
	}

	cmtID, text := bot.comment(entryTitle, cmt)
	msgIDsVar, found := bot.cmts.Get(cmtID)
	if !found {
		return
	}

	msgIDs := msgIDsVar.(messageIDs)

	bot.send <- func() {
		for _, msgID := range msgIDs {
			msg := tgbotapi.NewEditMessageText(msgID.chat, msgID.msg, text)
			msg.DisableWebPagePreview = true
			msg.ParseMode = "HTML"
			_, err := bot.api.Send(msg)
			if err != nil {
				bot.log.Error(err.Error())
			}
		}
	}
}

func (bot *TelegramBot) SendRemoveComment(commentID int64) {
	if bot.api == nil {
		return
	}

	cmtID := idToString(commentID)
	msgIDsVar, found := bot.cmts.Get(cmtID)
	if !found {
		return
	}

	msgIDs := msgIDsVar.(messageIDs)

	bot.send <- func() {
		for _, msgID := range msgIDs {
			msg := tgbotapi.NewDeleteMessage(msgID.chat, msgID.msg)
			_, err := bot.api.DeleteMessage(msg)
			if err != nil {
				bot.log.Error(err.Error())
			}
		}
	}
}

func (bot *TelegramBot) SendPasswordChanged(chat int64) {
	if bot.api == nil {
		return
	}

	const text = "Пароль к твоему тлогу был изменен."
	bot.sendMessage(chat, text)
}

func (bot *TelegramBot) SendEmailChanged(chat int64) {
	if bot.api == nil {
		return
	}

	const text = "Твой адрес почты был изменен."
	bot.sendMessage(chat, text)
}

func (bot *TelegramBot) SendNewFollower(chat int64, fromName, fromShowName, fromGender string, toPrivate bool) {
	if bot.api == nil {
		return
	}

	var ending string
	if fromGender == "female" {
		ending = "ась"
	} else {
		ending = "ся"
	}

	link := bot.userNameLink(fromName, fromShowName)

	var text string
	if toPrivate {
		text = link + " запрашивает доступ к чтению твоего тлога."
	} else {
		text = link + " подписал" + ending + " на твой тлог."
	}

	bot.sendMessage(chat, text)
}

func (bot *TelegramBot) SendNewAccept(chat int64, fromName, fromShowName, fromGender string) {
	if bot.api == nil {
		return
	}

	var ending string
	if fromGender == "female" {
		ending = "а"
	} else {
		ending = ""
	}

	link := bot.userNameLink(fromName, fromShowName)
	text := link + " разрешил" + ending + " тебе читать свой тлог."

	bot.sendMessage(chat, text)
}

func (bot *TelegramBot) SendNewInvite(chat int64) {
	if bot.api == nil {
		return
	}

	text := `У тебя появилось новое приглашение! <a href="` + bot.url + `users?top=waiting">Используй</a> его с умом.`
	bot.sendMessage(chat, text)
}

func (bot *TelegramBot) SendInvited(chat int64, fromName, fromShowName, fromGender string) {
	if bot.api == nil {
		return
	}

	var ending string
	if fromGender == "female" {
		ending = "а"
	} else {
		ending = ""
	}

	link := bot.userNameLink(fromName, fromShowName)
	text := link + " отправил" + ending + " тебе приглашение на Mindwell. " +
		"Теперь тебе доступны все функции сайта (при отсутствии других ограничений)."

	bot.sendMessage(chat, text)
}

func (bot *TelegramBot) SendAdmSent(chat int64) {
	if bot.api == nil {
		return
	}

	text := `Твой Дед Мороз отправил тебе подарок! Когда получишь, не забудь <a href="` + bot.url +
		`adm">поставить нужный флажок</a>. И не открывай до Нового года.`
	bot.sendMessage(chat, text)
}

func (bot *TelegramBot) SendAdmReceived(chat int64) {
	if bot.api == nil {
		return
	}

	text := `Внук получил твой новогодний подарок.`
	bot.sendMessage(chat, text)
}

func (bot *TelegramBot) message(msg *models.Message) (msgID, text string) {
	msgID = idToString(msg.ID)

	link := bot.url + "chats/" + msg.Author.Name

	text = tgHtmlEsc.Replace(msg.Author.ShowName) + " пишет: \n" +
		"«" + tgHtmlEsc.Replace(msg.EditContent) + "»\n" +
		`В <a href="` + link + `">беседе</a>`

	return
}

func (bot *TelegramBot) SendNewMessage(chat int64, msg *models.Message) {
	if bot.api == nil {
		return
	}

	msgID, text := bot.message(msg)

	bot.send <- func() {
		msg := bot.sendMessageNow(chat, text)
		if msg.MessageID <= 0 {
			return
		}

		var msgIDs messageIDs
		msgIDsVar, found := bot.msgs.Get(msgID)
		if found {
			msgIDs = msgIDsVar.(messageIDs)
		}
		msgIDs = append(msgIDs, messageID{chat, msg.MessageID})

		bot.msgs.SetDefault(msgID, msgIDs)
	}
}

func (bot *TelegramBot) SendUpdateMessage(msg *models.Message) {
	if bot.api == nil {
		return
	}

	msgID, text := bot.message(msg)
	msgIDsVar, found := bot.msgs.Get(msgID)
	if !found {
		return
	}

	msgIDs := msgIDsVar.(messageIDs)

	bot.send <- func() {
		for _, msgID := range msgIDs {
			msg := tgbotapi.NewEditMessageText(msgID.chat, msgID.msg, text)
			msg.DisableWebPagePreview = true
			msg.ParseMode = "HTML"
			_, err := bot.api.Send(msg)
			if err != nil {
				bot.log.Error(err.Error())
			}
		}
	}
}

func (bot *TelegramBot) SendRemoveMessage(messageID int64) {
	if bot.api == nil {
		return
	}

	cmtID := idToString(messageID)
	msgIDsVar, found := bot.msgs.Get(cmtID)
	if !found {
		return
	}

	msgIDs := msgIDsVar.(messageIDs)

	bot.send <- func() {
		for _, msgID := range msgIDs {
			msg := tgbotapi.NewDeleteMessage(msgID.chat, msgID.msg)
			_, err := bot.api.DeleteMessage(msg)
			if err != nil {
				bot.log.Error(err.Error())
			}
		}
	}
}

func (bot *TelegramBot) userNameLink(name, showName string) string {
	return fmt.Sprintf(`<a href="%susers/%s">%s</a>`, bot.url, name, tgHtmlEsc.Replace(showName))

}

func (bot *TelegramBot) userLink(user *models.User) string {
	return bot.userNameLink(user.Name, user.ShowName)
}

func (bot *TelegramBot) SendCommentComplain(from, against *models.User, content, comment string, commentID, entryID int64) {
	if bot.api == nil || bot.group == 0 {
		return
	}

	text := "Пользователь " + bot.userLink(from) + " пожаловался на комментарий " +
		strconv.FormatInt(commentID, 10) + " от " + bot.userLink(against) + ". " +
		"Текст комментария:\n\n«" + comment + "»\n\n"

	if content != "" {
		text += "Пояснение:\n«" + content + "»\n\n"
	}

	link := bot.url + "entries/" + strconv.FormatInt(entryID, 10) + "#comments"
	text += `К <a href="` + link + `">записи</a>`

	bot.sendMessage(bot.group, text)
}

func (bot *TelegramBot) SendEntryComplain(from, against *models.User, content, entry string, entryID int64) {
	entry, _ = CutText(entry, 2048)

	text := "Пользователь " + bot.userLink(from) + " пожаловался на запись " +
		strconv.FormatInt(entryID, 10) + " от " + bot.userLink(against) + ". " +
		"Текст записи:\n\n«" + entry + "»\n\n"

	if content != "" {
		text += "Пояснение:\n«" + content + "»\n\n"
	}

	link := bot.url + "entries/" + strconv.FormatInt(entryID, 10)
	text += `<a href="` + link + `">Запись</a>`

	bot.sendMessage(bot.group, text)
}
