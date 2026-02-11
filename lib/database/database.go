package database

import (
	"database/sql"
	"log"

	"github.com/leporo/sqlf"
	_ "github.com/lib/pq" // PostgreSQL driver
	goconf "github.com/zpatrick/go-config"
)

func init() {
	sqlf.SetDialect(sqlf.PostgreSQL)
}

func dropTable(tx *sql.Tx, table string) {
	_, err := tx.Exec("delete from " + table)
	if err != nil {
		tx.Rollback()
		log.Fatal("cannot clear table " + table + ": " + err.Error())
	}
}

// ClearDatabase drops user data tables and then creates default user
func ClearDatabase(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("cannot begin tx: ", err)
	}

	_, err = tx.Exec("UPDATE users SET pinned_entry = NULL")
	if err != nil {
		tx.Rollback()
		log.Fatal("cannot clear pinned entries: " + err.Error())
	}

	dropTable(tx, "vote_weights")
	dropTable(tx, "entries_privacy")
	dropTable(tx, "entry_tags")
	dropTable(tx, "entry_votes")
	dropTable(tx, "entry_images")
	dropTable(tx, "comment_votes")
	dropTable(tx, "comments")
	dropTable(tx, "favorites")
	dropTable(tx, "invites")
	dropTable(tx, "relations")
	dropTable(tx, "tags")
	dropTable(tx, "watching")
	dropTable(tx, "entries")
	dropTable(tx, "adm")
	dropTable(tx, "info")
	dropTable(tx, "notifications")
	dropTable(tx, "images")
	dropTable(tx, "complains")
	dropTable(tx, "talkers")
	dropTable(tx, "chats")
	//dropTable(tx, "apps")
	dropTable(tx, "sessions")
	//dropTable(tx, "app_tokens")
	dropTable(tx, "user_log")
	dropTable(tx, "wishes")
	dropTable(tx, "user_badges")
	// dropTable(tx, "badges")

	_, err = tx.Exec("delete from users where id != 1")
	if err != nil {
		tx.Rollback()
		log.Fatal("cannot clear table users: " + err.Error())
	}

	_, err = tx.Exec("delete from apps where flow < 4")
	if err != nil {
		tx.Rollback()
		log.Fatal("cannot clear table apps: " + err.Error())
	}

	tx.Exec("INSERT INTO invites (referrer_id, word1, word2, word3) VALUES(1, 1, 1, 1);")
	tx.Exec("INSERT INTO invites (referrer_id, word1, word2, word3) VALUES(1, 2, 2, 2);")
	tx.Exec("INSERT INTO invites (referrer_id, word1, word2, word3) VALUES(1, 3, 3, 3);")

	tx.Commit()
}

func ConnectionString(config *goconf.Config) string {
	host, err := config.String("database.host")
	if err != nil {
		log.Print(err)
	}

	port, err := config.String("database.port")
	if err != nil {
		log.Print(err)
	}

	user, err := config.String("database.user")
	if err != nil {
		log.Print(err)
	}

	pass, err := config.String("database.password")
	if err != nil {
		log.Print(err)
	}

	name, err := config.String("database.name")
	if err != nil {
		log.Print(err)
	}

	connStr := "user=" + user + " password=" + pass + " dbname=" + name + " host=" + host + " port=" + port
	connStr += " sslmode=disable"
	return connStr
}

// OpenDatabase returns db opened from config.
func OpenDatabase(config *goconf.Config) *sql.DB {
	driver, err := config.StringOr("database.driver", "postgres")
	if err != nil {
		log.Print(err)
	}

	connStr := ConnectionString(config)
	db, err := sql.Open(driver, connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
