package utils

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	cache "github.com/patrickmn/go-cache"
	"github.com/sevings/mindwell-server/models"
	"golang.org/x/net/proxy"
)

const errorText = "Что-то пошло не так…"

type TelegramBot struct {
	srv    *MindwellServer
	api    *tgbotapi.BotAPI
	url    string
	logins *cache.Cache
	send   chan *tgbotapi.MessageConfig
	stop   chan interface{}
}

func connectToProxy(srv *MindwellServer) *http.Client {
	auth := proxy.Auth{
		User:     srv.ConfigString("proxy.user"),
		Password: srv.ConfigString("proxy.password"),
	}

	if len(auth.User) == 0 {
		return nil
	}

	url := srv.ConfigString("proxy.url")
	dialer, err := proxy.SOCKS5("tcp", url, &auth, proxy.Direct)
	if err != nil {
		log.Println(err)
		return nil
	}

	tr := &http.Transport{Dial: dialer.Dial}
	return &http.Client{
		Transport: tr,
	}
}

func NewTelegramBot(srv *MindwellServer) *TelegramBot {
	bot := &TelegramBot{
		srv:    srv,
		url:    srv.ConfigString("server.base_url"),
		logins: cache.New(10*time.Minute, 10*time.Minute),
		send:   make(chan *tgbotapi.MessageConfig, 200),
		stop:   make(chan interface{}),
	}

	proxy := connectToProxy(srv)
	if proxy == nil {
		return bot
	}

	token := srv.ConfigString("telegram.token")
	if len(token) == 0 {
		return bot
	}

	api, err := tgbotapi.NewBotAPIWithClient(token, proxy)
	if err != nil {
		log.Print(err)
		return bot
	}

	bot.api = api
	// api.Debug = true

	log.Printf("Running Telegram bot %s", api.Self.UserName)

	go bot.run()

	return bot
}

func (bot *TelegramBot) sendMessage(chat int64, text string) {
	msg := tgbotapi.NewMessage(chat, text)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "Markdown"
	bot.send <- &msg
}

func (bot *TelegramBot) run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.api.GetUpdatesChan(u)
	if err != nil {
		log.Print(err)
	}

	for {
		select {
		case msg := <-bot.send:
			_, err := bot.api.Send(msg)
			if err != nil {
				log.Println("Telegram: ", err)
			}
		case <-bot.stop:
			return
		case upd := <-updates:
			if upd.Message == nil || !upd.Message.IsCommand() {
				continue
			}

			cmd := upd.Message.Command()
			log.Printf("Telegram: [%s] %s", upd.Message.From.UserName, cmd)

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
			default:
				reply = "Неизвестная команда. Попробуй /help."
			}

			bot.sendMessage(upd.Message.Chat.ID, reply)
		}
	}
}

func (bot *TelegramBot) Stop() {
	if bot.api == nil {
		return
	}

	close(bot.stop)
	bot.api.StopReceivingUpdates()
}

func (bot *TelegramBot) BuildToken(userID int64) string {
	token := GenerateString(8)
	bot.logins.Set(token, userID, cache.DefaultExpiration)
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
		return bot.help(upd)
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
		return "Скопируй верный ключ со [своей страницы настроек](https://mindwell.win/account/notifications)."
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
		log.Print(err)
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
		log.Print(err)
		return errorText
	}
}

func (bot *TelegramBot) help(upd *tgbotapi.Update) string {
	return "Привет! Я могу отправлять тебе уведомления из Mindwell.\n" +
		"Чтобы начать, скопируй ключ со [страницы настроек](https://mindwell.win/account/notifications). " +
		"Отправь его мне, используя команду `/login <ключ>`. Так ты подтвердишь свой аккаунт.\n" +
		"Чтобы я забыл твой номер в телеграме, введи /logout."
}

func (bot *TelegramBot) SendNewComment(chat int64, cmt *models.Comment) {
	if bot.api == nil {
		return
	}

	text := cmt.Author.ShowName + ": \n" +
		"«" + cmt.EditContent + "».\n" +
		bot.url + "entries/" + strconv.FormatInt(cmt.EntryID, 10) + "#comments"

	bot.sendMessage(chat, text)
}
