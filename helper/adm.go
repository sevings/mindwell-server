package helper

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/sevings/mindwell-server/utils"
)

type user struct {
	email  string
	name   string
	gender string
	at     time.Time // used in mailing
}

func genderNames(tx *utils.AutoTx) [][]string {
	names := make([][]string, 3)

	tx.Query("SELECT users.name, users.gender FROM users, adm " +
		"WHERE lower(adm.name) = lower(users.name)")

	for {
		var name string
		var gender int64
		ok := tx.Scan(&name, &gender)
		if !ok {
			break
		}

		names[gender] = append(names[gender], name)
	}

	return names
}

func loadIgnored(genderNames [][]string, tx *utils.AutoTx) map[string][]string {
	const q = `
SELECT t.name
FROM relations
JOIN relation ON relations.type = relation.id
JOIN users AS f ON relations.from_id = f.id
JOIN users AS t ON relations.to_id = t.id
WHERE lower(f.name) = lower($1)
	AND relation.type = 'ignored'
`

	ignored := map[string][]string{}

	for _, names := range genderNames {
		for _, from := range names {
			tx.Query(q, from)
			for {
				var to string
				ok := tx.Scan(&to)
				if !ok {
					break
				}

				ignored[from] = append(ignored[from], to)
				ignored[to] = append(ignored[to], from)
			}
		}
	}

	return ignored
}

func mixNames(names [][]string, ignored map[string][]string) []string {
	cnt := len(names[0]) + len(names[1]) + len(names[2])
	log.Printf("Found %d adms...\n", cnt)

	adm := make([]string, cnt+1)

remix:
	for {
		i := 0
		for j := 0; j <= 2; j++ {
			curNames := make([]string, len(names[j]))
			copy(curNames, names[j])

			for ; len(curNames) > 0; i += 2 {
				if i >= cnt {
					i = 1
				}

				adm[i], curNames = takeRandom(curNames)
			}
		}

		adm[cnt] = adm[0]

		for j := 0; j < len(adm)-1; j++ {
			gs := adm[j]
			gf := adm[j+1]
			ign := ignored[gs]
			for _, name := range ign {
				if name == gf {
					log.Printf("Found ignoring pair: %s, %s. Remixing.", gs, gf)
					continue remix
				}
			}
		}

		break
	}

	return adm
}

func setAdm(adm []string, tx *utils.AutoTx) {
	for i := 0; i < len(adm)-1; i++ {
		gs := adm[i]
		gf := adm[i+1]
		tx.Exec("UPDATE adm SET grandfather = $2 WHERE lower(name) = lower($1)", gs, gf)

		rows := tx.RowsAffected()
		if rows != 1 {
			log.Printf("Couldn't set grandfather for %s\n", gs)
		}
	}
}

func loadUsers(adm []string, tx *utils.AutoTx) []user {
	var users []user

	for _, name := range adm {
		tx.Query("SELECT show_name, gender.type, email, verified FROM users, gender "+
			"WHERE lower(name) = lower($1) AND users.gender = gender.id", name)
		for {
			var verified bool
			var usr user
			ok := tx.Scan(&usr.name, &usr.gender, &usr.email, &verified)
			if !ok {
				break
			}
			if verified {
				users = append(users, usr)
			}
		}
	}

	return users
}

func takeRandom(s []string) (string, []string) {
	var result string
	i := rand.Intn(len(s))
	result, s[i] = s[i], s[len(s)-1]
	return result, s[:len(s)-1]
}

func UpdateAdm(tx *utils.AutoTx, mail *utils.Postman) {
	rand.Seed(time.Now().UTC().UnixNano())

	names := genderNames(tx)
	if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
		return
	}

	ignored := loadIgnored(names, tx)
	if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
		return
	}

	adm := mixNames(names, ignored)
	setAdm(adm, tx)
	if tx.Error() != nil {
		return
	}

	users := loadUsers(adm[:len(adm)-1], tx)
	if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
		return
	}

	for _, usr := range users {
		mail.SendAdm(usr.email, usr.name, usr.gender)
	}

	log.Println("Completed. Sending emails...")
}
