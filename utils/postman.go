package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/matcornic/hermes"
	"github.com/sevings/mindwell-server/models"
	"gopkg.in/mailgun/mailgun-go.v1"
)

type Postman struct {
	mg mailgun.Mailgun
	h  hermes.Hermes
	ch chan *mailgun.Message
}

func NewPostman(domain, apiKey, pubKey string) *Postman {
	pm := &Postman{
		mg: mailgun.NewMailgun(domain, apiKey, pubKey),
		h: hermes.Hermes{
			Theme: &hermes.Flat{},
			Product: hermes.Product{
				Name:      "команда Mindwell",
				Link:      "http://mindwell.win",
				Logo:      "http://www.mindwell.win/assets/olympus/img/logo-mindwell.png",
				Copyright: "© 2018 Mindwell.",
				TroubleText: "Если кнопка '{ACTION}' по какой-то причине не работает, " +
					"скопируй и вставь в адреснуню строку браузера следующую ссылку: ",
			},
		},
		ch: make(chan *mailgun.Message, 10),
	}

	go func() {
		for msg := range pm.ch {
			resp, id, err := pm.mg.Send(msg)
			if err == nil {
				fmt.Printf("ID: %s Resp: %s\n", id, resp)
			} else {
				log.Println(err)
			}
		}
	}()

	return pm
}

func (pm *Postman) send(email hermes.Email, address, subj, name string) {
	email.Body.Title = "Привет, " + name
	email.Body.Signature = "С наилучшими пожеланиями"
	email.Body.Outros = []string{
		"Появились вопросы или какая-то проблема? " +
			"Не стесняйся и просто ответь на это письмо. Мы будем рады помочь. ",
	}

	text, err := pm.h.GeneratePlainText(email)
	if err != nil {
		log.Println(err)
	}

	from := "Команда Mindwell <support@mindwell.win>"
	recp := name + " <" + address + ">"
	msg := pm.mg.NewMessage(from, subj, text, recp)

	// html, err := pm.h.GenerateHTML(email)
	// if err != nil {
	// 	log.Println(err)
	// }
	// msg.SetHtml(html)

	// err = ioutil.WriteFile("preview.html", []byte(html), 0644)
	// err = ioutil.WriteFile("preview.txt", []byte(text), 0644)

	pm.ch <- msg
}

func (pm *Postman) SendGreeting(address, name, code string) {
	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				"добро пожаловать на борт нашего корабля!",
				"Располагайся, чувствуй себя как дома. Тебе у нас понравится. ",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Открой эту ссылку, чтобы подтвердить свой email:",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Начать пользоваться Mindwell",
						Link:  "http://mindwell.win/account/verification/" + address + "?code=" + code,
					},
				},
			},
		},
	}

	subj := "Приветствуем в Mindwell, " + name + "!"
	pm.send(email, address, subj, name)
}

func (pm *Postman) SendNewComment(address, fromGender, toShowName, entryTitle string, cmt *models.Comment) {
	var ending string
	if fromGender == models.ProfileAllOf1GenderFemale {
		ending = "а"
	}

	var entry string
	if len(entryTitle) > 0 {
		entry = " «" + entryTitle + "»"
	} else {
		entry = ", за которой ты следишь"
	}

	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				cmt.Author.ShowName + " оставил" + ending + " новый комментарий к записи" + entry + ".",
				"Вот, что он" + ending + " пишет:",
				"«" + cmt.Content + "».",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Узнать подробности и ответить можно по ссылке:",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Открыть запись",
						Link:  "http://mindwell.win/entries/" + strconv.FormatInt(cmt.EntryID, 10) + "#comments",
					},
				},
			},
		},
	}

	subj := "Новый комментарий к записи" + entry
	pm.send(email, address, subj, toShowName)
}

func (pm *Postman) SendNewFollower(address, fromName, fromShowName, fromGender string, toPrivate bool, toShowName string) {
	var ending, pronoun string
	if fromGender == models.ProfileAllOf1GenderFemale {
		ending = "ась"
		pronoun = "её"
	} else {
		ending = "ся"
		pronoun = "его"
	}

	var intro, text string
	if toPrivate {
		intro = fromShowName + " запрашивает доступ к чтению твоего тлога."
		text = "Принять или отклонить запрос можно на странице " + pronoun + " профиля: "
	} else {
		intro = fromShowName + " подписал" + ending + " на твой тлог."
		text = "Ссылка на " + pronoun + " профиль: "
	}

	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				intro,
			},
			Actions: []hermes.Action{
				{
					Instructions: text,
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  fromShowName,
						Link:  "http://mindwell.win/users/" + fromName,
					},
				},
			},
		},
	}

	const subj = "Новый подписчик"
	pm.send(email, address, subj, toShowName)
}
