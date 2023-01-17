package users

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/users"
	"github.com/sevings/mindwell-server/utils"
)

func searchUsers(srv *utils.MindwellServer, tx *utils.AutoTx, queryWhere, searchQuery string) []*models.Friend {
	query := usersQuerySelect + `, 0 
					FROM (
						SELECT *, $1 <<-> to_search_string(name, show_name, country, city) AS trgm_dist 
						FROM users 
						WHERE ` + queryWhere + `
						ORDER BY trgm_dist
						LIMIT 50
					) AS users, gender, user_privacy
					WHERE trgm_dist < 0.6` + usersQueryJoins
	tx.Query(query, searchQuery)
	list, _, _ := loadUserList(srv, tx, false)
	return list
}

func loadTopUsers(srv *utils.MindwellServer, tx *utils.AutoTx, top string, userID *models.UserID) []*models.Friend {
	query := usersQuerySelect + ", 0 "

	if top == "rank" {
		query += "FROM users, gender, user_privacy WHERE invited_by IS NOT NULL" + usersQueryJoins + "ORDER BY rank ASC"
		query += " LIMIT 50"
		tx.Query(query)
	} else if top == "new" {
		query += "FROM users, gender, user_privacy WHERE invited_by IS NOT NULL" + usersQueryJoins + " ORDER BY created_at DESC"
		query += " LIMIT 50"
		tx.Query(query)
	} else if top == "waiting" {
		query += `
			FROM (
				SELECT *, COALESCE(last_entries.last_created_at, '-infinity') AS last_entry
				FROM users
				LEFT JOIN (
					SELECT author_id, MAX(entries.created_at) AS last_created_at
					FROM entries 
					INNER JOIN users ON entries.author_id = users.id
					INNER JOIN entry_privacy ON entries.visible_for = entry_privacy.id
					INNER JOIN user_privacy ON users.privacy = user_privacy.id
					WHERE 
						CASE user_privacy.type
						WHEN 'all' THEN TRUE
						WHEN 'registered' THEN TRUE
						WHEN 'invited' THEN $1
						ELSE FALSE
						END
						AND
						CASE entry_privacy.type
						WHEN 'all' THEN TRUE
						WHEN 'registered' THEN TRUE
						WHEN 'invited' THEN $1
						ELSE FALSE
						END
					GROUP BY author_id
				) AS last_entries ON users.id = last_entries.author_id
				WHERE invited_by IS NULL AND creator_id IS NULL
				ORDER BY last_entry DESC, created_at DESC
			) as users
			INNER JOIN gender ON gender.id = users.gender
			INNER JOIN user_privacy ON users.privacy = user_privacy.id
		`
		query += " LIMIT 50"
		tx.Query(query, userID.IsInvited)
	} else {
		fmt.Printf("Unknown users top: %s\n", top)
		return nil
	}

	list, _, _ := loadUserList(srv, tx, false)
	return list
}

func newUsersLoader(srv *utils.MindwellServer) func(users.GetUsersParams, *models.UserID) middleware.Responder {
	return func(params users.GetUsersParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			body := &users.GetUsersOKBody{}

			if params.Query != nil {
				body.Users = searchUsers(srv, tx, "creator_id IS NULL", *params.Query)
				body.Query = *params.Query
			} else {
				body.Users = loadTopUsers(srv, tx, *params.Top, userID)
				body.Top = *params.Top
			}

			return users.NewGetUsersOK().WithPayload(body)
		})
	}
}
