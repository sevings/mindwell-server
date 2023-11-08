package utils

import (
	"database/sql"
	"github.com/go-openapi/errors"
	"github.com/sevings/mindwell-server/models"
	"strings"
	"time"
)

var errUnauthorized = errors.New(401, "Invalid or expired API key")
var errAccessToken = errors.New(401, "Invalid access token")
var errAppToken = errors.New(401, "Invalid app token")
var errAccessDenied = errors.New(403, "Access denied")

const userIDQuery = `
			SELECT users.id, name, followers_count, 
				invited_by is not null, karma < -1, verified,
				invite_ban > CURRENT_DATE, vote_ban > CURRENT_DATE, 
				comment_ban > CURRENT_DATE, live_ban > CURRENT_DATE,
				user_ban > CURRENT_DATE,
				authority.type
			FROM users
			JOIN authority ON users.authority = authority.id `

func LoadUserIDByID(tx *AutoTx, id int64) (*models.UserID, error) {
	const q = userIDQuery + "WHERE users.id = $1"
	tx.Query(q, id)
	return scanUserID(tx)
}

func LoadUserIDByName(tx *AutoTx, name string) (*models.UserID, error) {
	const q = userIDQuery + "WHERE lower(name) = lower($1)"
	tx.Query(q, name)
	return scanUserID(tx)
}

func scanUserID(tx *AutoTx) (*models.UserID, error) {
	var user models.UserID
	user.Ban = &models.UserIDBan{}
	tx.Scan(&user.ID, &user.Name, &user.FollowersCount,
		&user.IsInvited, &user.NegKarma, &user.Verified,
		&user.Ban.Invite, &user.Ban.Vote,
		&user.Ban.Comment, &user.Ban.Live,
		&user.Ban.Account,
		&user.Authority)
	if tx.Error() != nil {
		return nil, errUnauthorized
	}

	user.Ban.Invite = user.Ban.Invite || !user.IsInvited || !user.Verified
	user.Ban.Vote = user.Ban.Vote || !user.IsInvited || user.NegKarma || !user.Verified
	user.Ban.Comment = user.Ban.Comment || !user.IsInvited || !user.Verified
	user.Ban.Live = user.Ban.Live || !user.Verified

	return &user, nil
}

type AuthFlow uint8

const (
	AppFlow      AuthFlow = 1
	CodeFlow     AuthFlow = 2
	PasswordFlow AuthFlow = 4
)

var allScopes = [...]string{
	"account:read",
	"account:write",
	"adm:read",
	"adm:write",
	"comments:read",
	"comments:write",
	"entries:read",
	"entries:write",
	"favorites:write",
	"images:read",
	"images:write",
	"messages:read",
	"messages:write",
	"notifications:read",
	"relations:write",
	"settings:read",
	"settings:write",
	"users:read",
	"users:write",
	"votes:write",
	"watchings:write",
	"themes:read",
	"themes:write",
}

func findScope(scope string) (uint32, error) {
	for i, s := range allScopes {
		if scope == s {
			return 1 << i, nil
		}
	}

	return 0, errors.New(401, "invalid scope: %s", scope)
}

func ScopeFromString(scopes []string) (uint32, error) {
	var scope uint32

	for _, s := range scopes {
		n, err := findScope(s)
		if err != nil {
			return 0, err
		}

		scope += n
	}

	return scope, nil
}

func ScopeToString(scope uint32) []string {
	var scopes []string

	for i, s := range allScopes {
		if scope|1<<i == scope {
			scopes = append(scopes, s)
		}
	}

	return scopes
}

const AccessTokenLifetime = 60 * 60 * 24
const RefreshTokenLifetime = 60 * 60 * 24 * 30
const AppTokenLifetime = 60 * 60 * 24
const AccessTokenLength = 32
const RefreshTokenLength = 48
const AppTokenLength = 32
const CodeLength = 32

func NewOAuth2User(h TokenHash, db *sql.DB, flowReq AuthFlow) func(string, []string) (*models.UserID, error) {
	const query = `
SELECT scope, flow, ban, user_ban > CURRENT_DATE
FROM sessions
JOIN users ON users.id = user_id
JOIN apps ON apps.id = app_id
WHERE lower(users.name) = lower($1) 
	AND access_hash = $2
	AND access_thru > $3
`

	return func(token string, scopes []string) (*models.UserID, error) {
		scopeReq, err := ScopeFromString(scopes)
		if err != nil {
			return nil, err
		}

		nameToken := strings.Split(token, ".")
		if len(nameToken) < 2 {
			return nil, errAccessToken
		}

		accessToken := nameToken[1]
		if len(accessToken) != AccessTokenLength {
			return nil, errAccessToken
		}

		name := nameToken[0]
		hash := h.AccessTokenHash(token)
		now := time.Now()

		tx := NewAutoTx(db)
		defer tx.Finish()

		var scopeEx uint32
		var flowEx AuthFlow
		var appBan, userBan bool
		tx.Query(query, name, hash[:], now).Scan(&scopeEx, &flowEx, &appBan, &userBan)
		if tx.Error() != nil || userBan {
			return nil, errAccessToken
		}

		if appBan || scopeEx&scopeReq != scopeReq || flowEx&flowReq != flowReq {
			return nil, errAccessDenied
		}

		userID, err := LoadUserIDByName(tx, name)
		if userID != nil && userID.Ban.Account {
			return nil, errUnauthorized
		}

		return userID, err
	}
}

func NoAuthUser() *models.UserID {
	return &models.UserID{
		Ban: &models.UserIDBan{
			Comment: true,
			Invite:  true,
			Live:    true,
			Vote:    true,
		},
	}
}

func NewOAuth2App(h TokenHash, db *sql.DB) func(string, []string) (*models.UserID, error) {
	const query = `
SELECT ban, flow
FROM app_tokens
JOIN apps ON apps.id = app_id
WHERE lower(apps.name) = lower($1) 
	AND token_hash = $2
	AND valid_thru > $3
`

	return func(token string, scopes []string) (*models.UserID, error) {
		nameToken := strings.Split(token, "+")
		if len(nameToken) < 2 {
			return nil, errAppToken
		}

		appToken := nameToken[1]
		if len(appToken) != AppTokenLength {
			return nil, errAppToken
		}

		name := nameToken[0]
		hash := h.AppTokenHash(token)
		now := time.Now()

		tx := NewAutoTx(db)
		defer tx.Finish()

		var ban bool
		var flowEx AuthFlow
		tx.Query(query, name, hash, now).Scan(&ban, &flowEx)
		if tx.Error() != nil {
			return nil, errAppToken
		}

		if ban || flowEx&AppFlow != AppFlow {
			return nil, errAccessDenied
		}

		return NoAuthUser(), nil
	}
}
