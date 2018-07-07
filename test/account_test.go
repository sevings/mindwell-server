package test

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	accountImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/account"
	commentsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/comments"
	designImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/design"
	entriesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/entries"
	favoritesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/favorites"
	relationsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/relations"
	usersImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/users"
	votesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/votes"
	watchingsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/watchings"
	"github.com/sevings/mindwell-server/utils"

	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations"
	"github.com/sevings/mindwell-server/restapi/operations/account"
	"github.com/stretchr/testify/require"
)

var srv *utils.MindwellServer
var api *operations.MindwellAPI
var db *sql.DB
var userIDs []*models.UserID
var profiles []*models.AuthProfile
var esm EmailSenderMock

func TestMain(m *testing.M) {
	api = &operations.MindwellAPI{}
	srv = utils.NewMindwellServer(api, "../configs/server")
	db = srv.DB

	srv.Mail = &esm

	utils.ClearDatabase(db)

	accountImpl.ConfigureAPI(srv)
	usersImpl.ConfigureAPI(srv)
	entriesImpl.ConfigureAPI(srv)
	votesImpl.ConfigureAPI(srv)
	favoritesImpl.ConfigureAPI(srv)
	watchingsImpl.ConfigureAPI(srv)
	commentsImpl.ConfigureAPI(srv)
	designImpl.ConfigureAPI(srv)
	relationsImpl.ConfigureAPI(srv)

	userIDs, profiles = registerTestUsers(db)

	if len(esm.Emails) != 3 {
		log.Fatal("Email count")
	}

	for i := 0; i < 3; i++ {
		email := "test" + strconv.Itoa(i)
		if esm.Emails[i] != email {
			log.Fatal("Greeting has not sent to ", email)
		}
	}

	esm.Clear()

	os.Exit(m.Run())
}

func checkEmail(t *testing.T, email string, free bool) {
	check := api.AccountGetAccountEmailEmailHandler.Handle
	resp := check(account.GetAccountEmailEmailParams{Email: email})
	body, ok := resp.(*account.GetAccountEmailEmailOK)

	require.True(t, ok, email)
	require.Equal(t, email, *body.Payload.Email)
	require.Equal(t, free, *body.Payload.IsFree, email)
}

func TestCheckEmail(t *testing.T) {
	checkEmail(t, "123", true)
}

func checkName(t *testing.T, name string, free bool) {
	check := api.AccountGetAccountNameNameHandler.Handle
	resp := check(account.GetAccountNameNameParams{Name: name})
	body, ok := resp.(*account.GetAccountNameNameOK)

	require.True(t, ok, name)
	require.Equal(t, name, *body.Payload.Name)
	require.Equal(t, free, *body.Payload.IsFree, name)
}

func TestCheckName(t *testing.T) {
	checkName(t, "HaveANICEDay", false)
	checkName(t, "nAMe", true)
}

func checkInvites(t *testing.T, userID *models.UserID, size int) {
	load := api.AccountGetAccountInvitesHandler.Handle
	resp := load(account.GetAccountInvitesParams{}, userID)
	body, ok := resp.(*account.GetAccountInvitesOK)

	require.True(t, ok, "user %d", userID.ID)
	require.Equal(t, size, len(body.Payload.Invites), "user %d", userID.ID)
}

func checkLogin(t *testing.T, user *models.AuthProfile, name, password string) {
	params := account.PostAccountLoginParams{
		Name:     name,
		Password: password,
	}

	login := api.AccountPostAccountLoginHandler.Handle
	resp := login(params)
	body, ok := resp.(*account.PostAccountLoginOK)
	if !ok {
		badBody, ok := resp.(*account.PostAccountLoginBadRequest)
		if ok {
			t.Fatal(badBody.Payload.Message)
		}

		t.Fatal("login error")
	}

	require.Equal(t, user, body.Payload)
}

func changePassword(t *testing.T, userID *models.UserID, old, upd string, ok bool) {
	params := account.PostAccountPasswordParams{
		OldPassword: old,
		NewPassword: upd,
	}

	update := api.AccountPostAccountPasswordHandler.Handle
	resp := update(params, userID)
	switch resp.(type) {
	case *account.PostAccountPasswordOK:
		require.True(t, ok)
		return
	case *account.PostAccountPasswordForbidden:
		body := resp.(*account.PostAccountPasswordForbidden)
		require.False(t, ok, body.Payload.Message)
	default:
		t.Fatalf("set password user %d", userID.ID)
	}
}

func checkVerify(t *testing.T, userID *models.UserID, email string) {
	esm.CheckEmail(t, email)

	request := api.AccountPostAccountVerificationHandler.Handle
	resp := request(account.PostAccountVerificationParams{}, userID)
	_, ok := resp.(*account.PostAccountVerificationOK)

	req := require.New(t)
	req.True(ok, "user %d", userID.ID)
	req.Equal(1, len(esm.Emails))
	req.Equal(email, esm.Emails[0])
	req.Equal(1, len(esm.Codes))

	verify := api.AccountGetAccountVerificationEmailHandler.Handle
	params := account.GetAccountVerificationEmailParams{
		Code:  esm.Codes[0],
		Email: esm.Emails[0],
	}
	resp = verify(params)
	_, ok = resp.(*account.GetAccountVerificationEmailOK)

	req.True(ok)

	esm.Clear()
}

func TestRegister(t *testing.T) {
	{
		const q = "INSERT INTO invites(referrer_id, word1, word2, word3) VALUES(1, 1, 1, 1)"
		for i := 0; i < 3; i++ {
			db.Exec(q)
		}
	}

	inviter := &models.UserID{
		ID:   1,
		Name: "haveniceday",
	}

	checkInvites(t, inviter, 3)
	checkName(t, "testtEst", true)
	checkEmail(t, "testeMAil", true)

	params := account.PostAccountRegisterParams{
		Name:     "testtest",
		Email:    "testemail",
		Password: "test123",
		Invite:   "acknown acknown acknown",
		Referrer: "HaveANiceDay",
	}

	register := api.AccountPostAccountRegisterHandler.Handle
	resp := register(params)
	body, ok := resp.(*account.PostAccountRegisterCreated)
	if !ok {
		badBody, ok := resp.(*account.PostAccountRegisterBadRequest)
		if ok {
			t.Fatal(badBody.Payload.Message)
		}

		t.Fatal("reg error")
	}

	user := body.Payload
	userID := &models.UserID{
		ID:   user.ID,
		Name: user.Name,
	}

	checkInvites(t, inviter, 2)
	checkName(t, "testtEst", false)
	checkEmail(t, "testeMAil", false)
	checkLogin(t, user, params.Name, params.Password)
	checkLogin(t, user, strings.ToUpper(params.Name), params.Password)

	checkVerify(t, userID, user.Account.Email)
	user.Account.Verified = true
	checkLogin(t, user, params.Name, params.Password)

	changePassword(t, userID, "test123", "new123", true)
	changePassword(t, userID, "test123", "new123", false)
	checkLogin(t, user, params.Name, "new123")

	req := require.New(t)
	req.Equal(params.Name, user.Name)
	req.Equal(params.Email, user.Account.Email)
	req.Equal(params.Referrer, user.InvitedBy.Name)

	req.Equal(user.Name, user.ShowName)
	req.True(user.IsOnline)
	// req.Empty(user.Avatar)

	req.Equal("not set", user.Gender)
	req.False(user.IsDaylog)
	req.Equal("all", user.Privacy)
	req.Empty(user.Title)
	req.Zero(user.Karma)
	req.NotEmpty(user.CreatedAt)
	req.Equal(user.CreatedAt, user.LastSeenAt)
	req.Zero(user.AgeLowerBound)
	req.Zero(user.AgeLowerBound)
	req.Empty(user.Country)
	req.Empty(user.City)

	cnt := user.Counts
	req.Zero(cnt.Entries)
	req.Zero(cnt.Followings)
	req.Zero(cnt.Followers)
	req.Zero(cnt.Ignored)
	req.Zero(cnt.Invited)
	req.Zero(cnt.Comments)
	req.Zero(cnt.Favorites)
	req.Zero(cnt.Tags)

	req.Empty(user.Birthday)
	req.False(user.ShowInTops)

	acc := user.Account
	req.Equal(32, len(acc.APIKey))
	req.NotEmpty(acc.ValidThru)
	req.True(acc.Verified)

	resp = register(params)
	_, ok = resp.(*account.PostAccountRegisterBadRequest)
	req.True(ok)
	checkInvites(t, inviter, 2)

	gender := "female"
	city := "Moscow"
	country := "Russia"
	bday := "01.06.1992"

	params = account.PostAccountRegisterParams{
		Name:     "testtest2",
		Email:    "testemail2",
		Password: "test123",
		Gender:   &gender,
		City:     &city,
		Country:  &country,
		Birthday: &bday,
		Invite:   "acknown acknown acknown",
		Referrer: "HaveANiceDay",
	}

	resp = register(params)
	body, ok = resp.(*account.PostAccountRegisterCreated)
	req.True(ok)

	user = body.Payload
	userID = &models.UserID{
		ID:   user.ID,
		Name: user.Name,
	}

	checkInvites(t, inviter, 1)
	checkName(t, "testtEst2", false)
	checkEmail(t, "testeMAil2", false)
	checkLogin(t, user, params.Name, params.Password)

	checkVerify(t, userID, user.Account.Email)
	user.Account.Verified = true
	checkLogin(t, user, params.Name, params.Password)

	changePassword(t, userID, "test123", "new123", true)
	changePassword(t, userID, "test123", "new123", false)
	checkLogin(t, user, params.Name, "new123")

	req.Equal(gender, user.Gender)
	req.Equal(city, user.City)
	req.Equal(country, user.Country)

	req.Equal(int64(25), user.AgeLowerBound)
	req.Equal(int64(29), user.AgeUpperBound)
	req.Equal("1992-01-06T00:00:00Z", user.Birthday)

	req.NotEqual(acc.APIKey, user.Account.APIKey)
}
