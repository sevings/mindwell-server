package oauth2

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/patrickmn/go-cache"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/oauth2"
	"github.com/sevings/mindwell-server/utils"
	"strconv"
	"strings"
	"time"
)

type flow uint8

const (
	appFlow      flow = 1
	codeFlow     flow = 2
	passwordFlow flow = 4
)

// ConfigureAPI creates operations handlers
func ConfigureAPI(srv *utils.MindwellServer) {
	apiSecret := []byte(srv.ConfigString("server.api_secret"))

	srv.API.APIKeyHeaderAuth = utils.NewKeyAuth(srv.DB, apiSecret)
	srv.API.NoAPIKeyAuth = utils.NoApiKeyAuth
	srv.API.OAuth2PasswordAuth = newOAuth2User(srv.DB, passwordFlow)
	srv.API.OAuth2CodeAuth = newOAuth2User(srv.DB, codeFlow)

	srv.API.Oauth2GetOauth2AuthHandler = oauth2.GetOauth2AuthHandlerFunc(newOAuth2Auth(srv))
	srv.API.Oauth2PostOauth2TokenHandler = oauth2.PostOauth2TokenHandlerFunc(newOAuth2Token(srv))

	srv.API.Oauth2GetOauth2AppsIDHandler = oauth2.GetOauth2AppsIDHandlerFunc(newAppLoader(srv))
}

const accessTokenLifetime = 60 * 60 * 24
const refreshTokenLifetime = 60 * 60 * 24 * 30
const appTokenLifetime = 60 * 60 * 24
const accessTokenLength = 32
const refreshTokenLength = 48
const appTokenLength = 32
const codeLength = 32

type authData struct {
	appID       int64
	appSecret   string
	redirectUri string
	userID      int64
	userName    string
	scope       uint32
	challenge   string
	method      string
}

var authCache = cache.New(15*time.Minute, time.Hour)

var allScopes = [...]string{
	"account:read",
	"account:write",
	"adm:read",
	"adm:write",
	"comments:read",
	"comments:write",
	"entries:read",
	"entries:write",
	"favotites:write",
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
}

func findScope(scope string) (uint32, error) {
	for i, s := range allScopes {
		if scope == s {
			return 1 << i, nil
		}
	}

	return 0, fmt.Errorf("scope is invalid: %s", scope)
}

func scopeFromString(scopes []string) (uint32, error) {
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

func scopeToString(scope uint32) []string {
	var scopes []string

	for i, s := range allScopes {
		if scope|1<<i == scope {
			scopes = append(scopes, s)
		}
	}

	return scopes
}

func newOAuth2User(db *sql.DB, flowReq flow) func(string, []string) (*models.UserID, error) {
	const query = `
SELECT scope, flow
FROM access_tokens
JOIN users ON users.id = user_id
JOIN apps ON apps.id = app_id
WHERE lower(users.name) = lower($1) 
	AND token_hash = $2
	AND access_tokens.valid_thru > $3
`

	return func(token string, scopes []string) (*models.UserID, error) {
		scopeReq, err := scopeFromString(scopes)
		if err != nil {
			return nil, err
		}

		nameToken := strings.Split(token, ".")
		if len(nameToken) < 2 {
			return nil, fmt.Errorf("access token is invalid: %s", token)
		}

		accessToken := nameToken[1]
		if len(accessToken) != accessTokenLength {
			return nil, fmt.Errorf("access token is invalid: %s", token)
		}

		name := nameToken[0]
		hash := sha256.Sum256([]byte(accessToken))
		now := time.Now()

		tx := utils.NewAutoTx(db)
		defer tx.Finish()

		var scopeEx uint32
		var flowEx flow
		tx.Query(query, name, hash[:], now).Scan(&scopeEx, &flowEx)
		if tx.Error() != nil {
			return nil, fmt.Errorf("access token is invalid: %s", token)
		}

		if scopeEx&scopeReq != scopeReq || flowEx&flowReq != flowReq {
			return nil, errors.New("access denied")
		}

		return utils.LoadUserIDByName(tx, name)
	}
}

func newOAuth2App(db *sql.DB) func(string) error {
	const query = `
SELECT ban, flow
FROM app_tokens
JOIN apps ON apps.id = app_id
WHERE lower(apps.name) = lower($1) 
	AND token_hash = $2
	AND valid_thru > $3
`

	return func(token string) error {
		nameToken := strings.Split(token, "+")
		if len(nameToken) < 2 {
			return fmt.Errorf("app token is invalid: %s", token)
		}

		appToken := nameToken[1]
		if len(appToken) != appTokenLength {
			return fmt.Errorf("access token is invalid: %s", token)
		}

		name := nameToken[0]
		hash := sha256.Sum256([]byte(appToken))
		now := time.Now()

		tx := utils.NewAutoTx(db)
		defer tx.Finish()

		var ban bool
		var flowEx flow
		tx.Query(query, name, hash, now).Scan(&ban, &flowEx)
		if tx.Error() != nil {
			return fmt.Errorf("app token is invalid: %s", token)
		}

		if ban || flowEx&appFlow != appFlow {
			return errors.New("access denied")
		}

		return nil
	}
}

func createAccessToken(tx *utils.AutoTx, appID, userID int64, scope uint32, name string) (string, error) {
	token := utils.GenerateString(accessTokenLength)
	hash := sha256.Sum256([]byte(token))
	thru := time.Now().Add(accessTokenLifetime * time.Second)

	const query = `
INSERT INTO access_tokens(app_id, user_id, token_hash, scope, valid_thru)
VALUES($1, $2, $3, $4, $5)
`

	tx.Exec(query, appID, userID, hash[:], scope, thru)

	return name + "." + token, tx.Error()
}

func createRefreshToken(tx *utils.AutoTx, appID, userID int64, scope uint32) (string, error) {
	token := utils.GenerateString(refreshTokenLength)
	hash := sha256.Sum256([]byte(token))
	thru := time.Now().Add(refreshTokenLifetime * time.Second)

	const query = `
INSERT INTO refresh_tokens(app_id, user_id, token_hash, scope, valid_thru)
VALUES($1, $2, $3, $4, $5)
`

	tx.Exec(query, appID, userID, hash[:], scope, thru)

	id := strconv.FormatInt(userID, 32)
	return id + "." + token, tx.Error()
}

func createTokens(tx *utils.AutoTx, appID, userID int64, scope uint32, userName string) *oauth2.PostOauth2TokenOKBody {
	refreshToken, err := createRefreshToken(tx, appID, userID, scope)
	if err != nil {
		return nil
	}

	accessToken, err := createAccessToken(tx, appID, userID, scope, userName)
	if err != nil {
		return nil
	}

	resp := &oauth2.PostOauth2TokenOKBody{
		AccessToken:  accessToken,
		ExpiresIn:    accessTokenLifetime,
		RefreshToken: refreshToken,
		TokenType:    oauth2.PostOauth2TokenOKBodyTokenTypeBearer,
		Scope:        scopeToString(scope),
	}

	return resp
}

func createAppToken(tx *utils.AutoTx, appID int64, appName string) (string, error) {
	const query = `
INSERT INTO app_tokens(app_id, token_hash, valid_thru)
VALUES($1, $2, $3)
`

	token := utils.GenerateString(appTokenLength)
	hash := sha256.Sum256([]byte(token))
	thru := time.Now().Add(appTokenLifetime * time.Second)

	tx.Exec(query, appID, hash[:], thru)

	return appName + "+" + token, tx.Error()
}

func checkCodeGrant(tx *utils.AutoTx, appID int64) (authData, bool, error) {
	const query = `
SELECT secret, redirect_uri, flow, ban
FROM apps
WHERE id = $1
`

	auth := authData{appID: appID}
	var ban bool
	var f flow
	tx.Query(query, appID).Scan(&auth.appSecret, &auth.redirectUri, &f, &ban)

	granted := !ban && f&codeFlow == codeFlow
	return auth, granted, tx.Error()
}

func checkPasswordGrant(tx *utils.AutoTx, appID int64, appSecret string) (bool, error) {
	const grantQuery = `
SELECT flow, ban
FROM apps
WHERE id = $1 AND secret = $2
`
	var ban bool
	var f flow
	tx.Query(grantQuery, appID, appSecret).Scan(&f, &ban)

	granted := !ban && f&passwordFlow == passwordFlow
	return granted, tx.Error()
}

func checkAppGrant(tx *utils.AutoTx, appID int64, appSecret string) (string, bool, error) {
	const query = `
SELECT name, flow, ban
FROM apps
WHERE id = $1 AND secret = $2
`
	var name string
	var ban bool
	var f flow
	tx.Query(query, appID, appSecret).Scan(&name, &f, &ban)

	granted := !ban && f&appFlow == appFlow
	return name, granted, tx.Error()
}

func checkRefreshGrant(tx *utils.AutoTx, appID int64, appSecret string) (bool, error) {
	const query = `
SELECT ban
FROM apps
WHERE id = $1 AND secret = $2
`

	ban := tx.QueryBool(query, appID, appSecret)
	return !ban, tx.Error()
}

func getAuthBadRequest(err string) middleware.Responder {
	body := models.OAuth2Error{Error: err}
	return oauth2.NewGetOauth2AuthBadRequest().WithPayload(&body)
}

func newOAuth2Auth(srv *utils.MindwellServer) func(oauth2.GetOauth2AuthParams, *models.UserID) middleware.Responder {
	return func(params oauth2.GetOauth2AuthParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			scope, err := scopeFromString(params.Scope)
			if err != nil {
				return getAuthBadRequest(models.OAuth2ErrorErrorInvalidScope)
			}

			auth, granted, err := checkCodeGrant(tx, params.ClientID)
			if err != nil {
				return getAuthBadRequest(models.OAuth2ErrorErrorUnrecognizedClient)
			}
			if auth.redirectUri != params.RedirectURI {
				return getAuthBadRequest(models.OAuth2ErrorErrorInvalidRedirect)
			}
			if !granted {
				return getAuthBadRequest(models.OAuth2ErrorErrorInvalidGrant)
			}

			if auth.appSecret == "" && params.CodeChallenge == nil {
				return getAuthBadRequest(models.OAuth2ErrorErrorInvalidRequest)
			}

			resp := &oauth2.GetOauth2AuthOKBody{
				Code: utils.GenerateString(codeLength),
			}

			if params.State != nil {
				resp.State = *params.State
			}

			auth.userID = userID.ID
			auth.userName = userID.Name
			auth.scope = scope

			if params.CodeChallenge != nil {
				auth.challenge = *params.CodeChallenge
			}

			if params.CodeChallengeMethod != nil {
				auth.method = *params.CodeChallengeMethod
			}

			authCache.SetDefault(resp.Code, auth)

			return oauth2.NewGetOauth2AuthOK().WithPayload(resp)
		})
	}
}

func loadUserByPassword(srv *utils.MindwellServer, tx *utils.AutoTx, name, password string) (int64, string) {
	name = strings.TrimSpace(name)
	password = strings.TrimSpace(password)
	hash := srv.PasswordHash(password)

	const userIdQuery = `
SELECT id, name
FROM users
WHERE password_hash = $2
	AND (lower(name) = lower($1) OR lower(email) = lower($1))
`

	var userID int64
	var userName string
	tx.Query(userIdQuery, name, hash).Scan(&userID, &userName)

	return userID, userName
}

func postTokenBadRequest(err string) middleware.Responder {
	body := models.OAuth2Error{Error: err}
	return oauth2.NewPostOauth2TokenBadRequest().WithPayload(&body)
}

func requestPasswordToken(srv *utils.MindwellServer, tx *utils.AutoTx, appID int64, appSecret, name, password string) middleware.Responder {
	userID, userName := loadUserByPassword(srv, tx, name, password)
	if userID == 0 {
		return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
	}

	granted, err := checkPasswordGrant(tx, appID, appSecret)
	if err != nil {
		return postTokenBadRequest(models.OAuth2ErrorErrorUnrecognizedClient)
	}
	if !granted {
		return postTokenBadRequest(models.OAuth2ErrorErrorInvalidGrant)
	}

	var scope uint32 = 1<<31 - 1
	resp := createTokens(tx, appID, userID, scope, userName)
	if resp == nil {
		return postTokenBadRequest(models.OAuth2ErrorErrorServerError)
	}

	return oauth2.NewPostOauth2TokenOK().WithPayload(resp)
}

func requestAppToken(tx *utils.AutoTx, appID int64, appSecret string) middleware.Responder {
	appName, granted, err := checkAppGrant(tx, appID, appSecret)
	if err != nil {
		return postTokenBadRequest(models.OAuth2ErrorErrorUnrecognizedClient)
	}
	if !granted {
		return postTokenBadRequest(models.OAuth2ErrorErrorInvalidGrant)
	}

	appToken, err := createAppToken(tx, appID, appName)
	if err != nil {
		return postTokenBadRequest(models.OAuth2ErrorErrorServerError)
	}

	resp := &oauth2.PostOauth2TokenOKBody{
		AccessToken: appToken,
		ExpiresIn:   accessTokenLifetime,
		TokenType:   oauth2.PostOauth2TokenOKBodyTokenTypeBearer,
	}

	return oauth2.NewPostOauth2TokenOK().WithPayload(resp)
}

func requestAccessToken(tx *utils.AutoTx, appID int64, code, redirectUri string, appSecret, verifier *string) middleware.Responder {
	authValue, found := authCache.Get(code)
	if !found {
		return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
	}

	auth := authValue.(authData)
	if auth.appID != appID {
		return postTokenBadRequest(models.OAuth2ErrorErrorUnrecognizedClient)
	}
	if auth.redirectUri != redirectUri {
		return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRedirect)
	}

	if auth.appSecret != "" {
		if appSecret == nil {
			return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
		}

		if auth.appSecret != *appSecret {
			return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
		}
	}

	if auth.challenge != "" {
		if verifier == nil {
			return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
		}

		switch auth.method {
		case "S256":
			sum := sha256.Sum256([]byte(*verifier))
			ch := base64.URLEncoding.EncodeToString(sum[:])
			if auth.challenge != ch {
				return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
			}
		default:
			return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
		}
	}

	resp := createTokens(tx, auth.appID, auth.userID, auth.scope, auth.userName)
	if resp == nil {
		return postTokenBadRequest(models.OAuth2ErrorErrorServerError)
	}

	authCache.Delete(code)

	return oauth2.NewPostOauth2TokenOK().WithPayload(resp)
}

func requestRefreshToken(tx *utils.AutoTx, appID int64, appSecret, token string) middleware.Responder {
	const query = `
DELETE FROM refresh_tokens
WHERE app_id = $1 AND user_id = $2 AND token_hash = $3
RETURNING scope, valid_thru
`

	granted, err := checkRefreshGrant(tx, appID, appSecret)
	if err != nil {
		return postTokenBadRequest(models.OAuth2ErrorErrorUnrecognizedClient)
	}
	if !granted {
		return postTokenBadRequest(models.OAuth2ErrorErrorInvalidGrant)
	}

	idToken := strings.Split(token, ".")
	if len(idToken) != 2 {
		return postTokenBadRequest(models.OAuth2ErrorErrorInvalidToken)
	}

	userID, err := strconv.ParseInt(idToken[0], 32, 32)
	if err != nil {
		return postTokenBadRequest(models.OAuth2ErrorErrorInvalidToken)
	}

	hash := sha256.Sum256([]byte(idToken[1]))

	var scope uint32
	var thru time.Time
	tx.Query(query, appID, userID, hash[:]).Scan(&scope, &thru)

	if scope == 0 || time.Now().After(thru) {
		return postTokenBadRequest(models.OAuth2ErrorErrorInvalidToken)
	}

	userName := tx.QueryString("SELECT name FROM users WHERE id = $1", userID)
	resp := createTokens(tx, appID, userID, scope, userName)
	if resp == nil {
		return postTokenBadRequest(models.OAuth2ErrorErrorServerError)
	}

	return oauth2.NewPostOauth2TokenOK().WithPayload(resp)
}

func newOAuth2Token(srv *utils.MindwellServer) func(oauth2.PostOauth2TokenParams) middleware.Responder {
	return func(params oauth2.PostOauth2TokenParams) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			if params.GrantType == "password" {
				if params.Username == nil || params.Password == nil {
					return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
				}

				return requestPasswordToken(srv, tx, params.ClientID, *params.ClientSecret, *params.Username, *params.Password)
			}

			if params.GrantType == "authorization_code" {
				if params.Code == nil || params.RedirectURI == nil {
					return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
				}

				return requestAccessToken(tx, params.ClientID, *params.Code, *params.RedirectURI, params.ClientSecret, params.CodeVerifier)
			}

			if params.GrantType == "client_credentials" {
				if *params.ClientSecret == "" {
					return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
				}

				return requestAppToken(tx, params.ClientID, *params.ClientSecret)
			}

			if params.GrantType == "refresh_token" {
				if params.RefreshToken == nil {
					return postTokenBadRequest(models.OAuth2ErrorErrorInvalidRequest)
				}

				return requestRefreshToken(tx, params.ClientID, *params.ClientSecret, *params.RefreshToken)
			}

			return postTokenBadRequest(models.OAuth2ErrorErrorUnsupportedGrantType)
		})
	}
}

func loadApp(tx *utils.AutoTx, appID int64) (*models.App, bool) {
	const query = `
SELECT name, show_name, platform, info
FROM apps
WHERE id = $1
`

	app := &models.App{ID: appID}
	tx.Query(query, appID).Scan(&app.Name, &app.ShowName, &app.Platform, &app.Info)

	return app, tx.Error() == nil
}

func newAppLoader(srv *utils.MindwellServer) func(oauth2.GetOauth2AppsIDParams, *models.UserID) middleware.Responder {
	return func(params oauth2.GetOauth2AppsIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			app, ok := loadApp(tx, params.ID)
			if !ok {
				err := &i18n.Message{ID: "no_app", Other: "App not found."}
				return oauth2.NewGetOauth2AppsIDNotFound().WithPayload(srv.NewError(err))
			}

			return oauth2.NewGetOauth2AppsIDOK().WithPayload(app)
		})
	}
}
