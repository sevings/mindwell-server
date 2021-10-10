package utils

import (
	"database/sql"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type userRequest struct {
	user string
	ip   string
	ua   string
	dev  string
	app  string
	uid  string
	at   time.Time
}

func (req userRequest) key() string {
	var str strings.Builder
	str.Grow(68)

	str.WriteString(req.ip)
	str.WriteString(req.dev)
	str.WriteString(req.app)
	str.WriteString(req.uid)
	str.WriteString(req.user)

	return str.String()
}

type userLog struct {
	db   *sql.DB
	log  *zap.Logger
	ch   chan *userRequest
	tick *time.Ticker
	prev map[string]*userRequest
}

func CreateUserLog(db *sql.DB, log *zap.Logger) middleware.Builder {
	ul := &userLog{
		db:   db,
		log:  log,
		ch:   make(chan *userRequest, 200),
		tick: time.NewTicker(time.Hour),
		prev: make(map[string]*userRequest),
	}

	go ul.run()

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ul.ServeHTTP(w, r)
			handler.ServeHTTP(w, r)
		})
	}
}

func (ul *userLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	typeTok := strings.SplitN(token, " ", 2)
	if len(typeTok) > 1 {
		token = typeTok[1]
	}
	user := strings.SplitN(token, ".", 2)[0]

	uid := r.Header.Get("X-Uid")
	app := r.Header.Get("X-App")
	dev := r.Header.Get("X-Dev")
	ip := r.Header.Get("X-Forwarded-For")
	ua := r.UserAgent()

	at := time.Now()

	ul.ch <- &userRequest{
		user: user,
		ip:   ip,
		ua:   ua,
		dev:  dev,
		app:  app,
		uid:  uid,
		at:   at,
	}
}

func (ul *userLog) run() {
	for {
		select {
		case <-ul.tick.C:
			ul.clearOld()
		case req := <-ul.ch:
			ul.addRequest(req)
		}
	}
}

func (ul *userLog) clearOld() {
	newPrev := make(map[string]*userRequest)
	minAt := time.Now().Add(-1 * time.Hour)

	tx := NewAutoTx(ul.db)
	defer tx.Finish()

	for key, req := range ul.prev {
		if req.at.After(minAt) {
			newPrev[key] = req
		} else {
			ul.save(tx, req, false)
		}
	}

	ul.prev = newPrev
}

func (ul *userLog) addRequest(req *userRequest) {
	ul.log.Info("user",
		zap.String("user", req.user),
		zap.String("ip", req.ip),
		zap.String("user_agent", req.ua),
		zap.String("dev", req.dev),
		zap.String("app", req.app),
		zap.String("uid", req.uid),
	)

	key := req.key()
	prevReq, found := ul.prev[key]

	if found {
		prevReq.at = req.at
		return
	}

	ul.prev[key] = req

	tx := NewAutoTx(ul.db)
	defer tx.Finish()

	ul.save(tx, req, true)
}

func (ul *userLog) save(tx *AutoTx, req *userRequest, first bool) {
	if req.user == "" {
		return
	}

	app, err := strconv.ParseUint(req.app, 16, 64)
	if err != nil {
		ul.log.Warn("app id is invalid", zap.String("app", req.app))
	}

	uid, err := strconv.ParseUint(req.uid, 16, 32)
	if err != nil {
		ul.log.Warn("uid is invalid", zap.String("uid", req.uid))
	}

	dev, err := strconv.ParseUint(req.dev, 16, 32)
	if err != nil {
		ul.log.Warn("dev id is invalid", zap.String("dev", req.dev))
	}

	ip := strings.SplitN(req.ip, ",", 2)[0]

	const query = `
    INSERT INTO user_log(name, ip, user_agent, device, app, uid, at, first) 
    VALUES(lower($1), $2, $3, $4, $5, $6, $7, $8)
`

	tx.Exec(query, req.user, ip, req.ua, int32(dev), int64(app), int32(uid), req.at, first)
}
