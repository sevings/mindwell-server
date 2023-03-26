package helper

import (
	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/utils"
	"log"
	"time"
)

func SendReminders(tx *utils.AutoTx, mail *utils.Postman) {
	q := sqlf.Select("show_name, gender.type, email, last_seen_at").
		From("users").
		Join("gender", "users.gender = gender.id").
		Where("verified").
		Where("age(last_seen_at) > interval '90 days'").
		Where("karma >= 0").
		OrderBy("last_seen_at")

	tx.QueryStmt(q)

	var users []user
	for {
		var usr user
		if !tx.Scan(&usr.name, &usr.gender, &usr.email, &usr.at) {
			break
		}

		users = append(users, usr)
	}

	log.Printf("Sending %d emails...\n", len(users))

	for i, usr := range users {
		date := usr.at.Format(time.RFC3339)
		log.Printf("%d. Sending to %s (%s)...\n", i+1, usr.name, date)
		mail.SendReminder(usr.email, usr.name, usr.gender)
		time.Sleep(time.Second)
	}
}

func SendSurvey(tx *utils.AutoTx, mail *utils.Postman, surveyUrl string) {
	if surveyUrl == "" {
		log.Println("The survey url is empty")
		return
	}

	sub := sqlf.Select("author_id, max(created_at) as at").
		From("entries").
		GroupBy("author_id")

	q := sqlf.Select("show_name, email, last_entries.at").
		From("").SubQuery("(", ") as last_entries", sub).
		Join("users", "users.id = last_entries.author_id").
		Where("verified").
		Where("invited_by is not null").
		Where("entries_count > 5").
		Where("age(last_entries.at) < interval '18 months'").
		Where("age(last_entries.at) > interval '2 months'").
		Where("karma >= 0").
		OrderBy("at")

	tx.QueryStmt(q)

	var users []user
	for {
		var usr user
		if !tx.Scan(&usr.name, &usr.email, &usr.at) {
			break
		}

		users = append(users, usr)
	}

	log.Printf("Sending %d emails...\n", len(users))

	for i, usr := range users {
		date := usr.at.Format(time.RFC3339)
		log.Printf("%d. Sending to %s (%s)...\n", i+1, usr.name, date)
		mail.SendSurvey(usr.email, usr.name, surveyUrl)
		time.Sleep(time.Second)
	}
}
