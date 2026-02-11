package userutil

import (
	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/lib/database"
)

// LoadUser loads a user by ID from the database.
func LoadUser(tx *database.AutoTx, id int64) *models.User {
	const query = `
SELECT id, name, show_name,
is_online(last_seen_at) AND creator_id IS NULL, creator_id IS NOT NULL
FROM users
WHERE id = $1`

	return queryUser(tx, query, id)
}

// LoadUserByName loads a user by name from the database.
func LoadUserByName(tx *database.AutoTx, name string) *models.User {
	const query = `
SELECT id, name, show_name,
is_online(last_seen_at) AND creator_id IS NULL, creator_id IS NOT NULL
FROM users
WHERE lower(name) = lower($1)`

	return queryUser(tx, query, name)
}

// LoadOnlineUsers loads all currently online users.
func LoadOnlineUsers(tx *database.AutoTx) []*models.User {
	const query = `
SELECT id, name, show_name,
is_online(last_seen_at) AND creator_id IS NULL, creator_id IS NOT NULL
FROM users
WHERE is_online(last_seen_at)
`

	tx.Query(query)

	var users []*models.User
	for {
		user := &models.User{}
		if !tx.Scan(&user.ID, &user.Name, &user.ShowName,
			&user.IsOnline, &user.IsTheme) {
			break
		}

		users = append(users, user)
	}

	return users
}

// queryUser is an internal helper for loading a single user.
func queryUser(tx *database.AutoTx, query string, arg any) *models.User {
	var user models.User

	tx.Query(query, arg).Scan(&user.ID, &user.Name, &user.ShowName,
		&user.IsOnline, &user.IsTheme)

	return &user
}

// LoadRelation loads the relationship type between two users.
func LoadRelation(tx *database.AutoTx, from, to int64) string {
	if from == 0 || to == 0 {
		return models.RelationshipRelationNone
	}

	relationQuery := sqlf.Select("relation.type").
		From("relations").
		Join("relation", "relation.id = relations.type").
		Where("from_id = ?", from).
		Where("to_id = ?", to)

	var relation string
	tx.QueryStmt(relationQuery).Scan(&relation)

	if relation == "" {
		return models.RelationshipRelationNone
	}

	return relation
}

// CanViewTlogName checks if the user can view the specified tlog.
func CanViewTlogName(tx *database.AutoTx, userID *models.UserID, tlog string) bool {
	if tlog == "" {
		return false
	}

	q := sqlf.Select("TRUE").
		From("users AS authors").
		LeftJoin("user_privacy", "authors.privacy = user_privacy.id").
		Where("lower(name) = lower(?)", tlog)

	AddCanViewAuthorQuery(q, userID)
	return tx.QueryStmt(q).ScanBool()
}

// CanViewEntry returns true if the user is allowed to read the entry.
func CanViewEntry(tx *database.AutoTx, userID *models.UserID, entryID int64) bool {
	if entryID == 0 {
		return false
	}

	q := sqlf.Select("TRUE").
		From("entries").
		Join("entry_privacy", "entries.visible_for = entry_privacy.id").
		Join("users AS authors", "entries.author_id = authors.id").
		Join("users AS entry_users", "entries.user_id = entry_users.id").
		Join("user_privacy", "authors.privacy = user_privacy.id").
		Where("entries.id = ?", entryID)

	AddCanViewEntryQuery(q, userID)

	return tx.QueryStmt(q).ScanBool()
}
