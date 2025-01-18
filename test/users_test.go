package test

import (
	"strings"
	"testing"

	"github.com/sevings/mindwell-server/restapi/operations/me"
	"github.com/sevings/mindwell-server/restapi/operations/users"

	"github.com/stretchr/testify/require"

	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/utils"
)

func TestGetMe(t *testing.T) {
	req := require.New(t)

	get := func(i int) *models.AuthProfile {
		load := api.MeGetMeHandler.Handle
		resp := load(me.GetMeParams{}, userIDs[i])
		body, ok := resp.(*me.GetMeOK)
		if !ok {
			t.Fatal("error get me")
		}

		return body.Payload
	}

	for i, user := range profiles {
		myProfile := get(i)

		req.Equal(user.ID, myProfile.ID)
		req.Equal(user.Name, myProfile.Name)

		req.Zero(myProfile.Ban.Invite)
		req.Zero(myProfile.Ban.Vote)
		req.Zero(myProfile.Ban.Comment)
		req.Zero(myProfile.Ban.Live)
	}

	banInvite(db, userIDs[0])
	myProfile := get(0)
	req.NotZero(myProfile.Ban.Invite)
	req.Zero(myProfile.Ban.Vote)
	removeUserRestrictions(db, userIDs)

	banVote(db, userIDs[0])
	myProfile = get(0)
	req.Zero(myProfile.Ban.Invite)
	req.NotZero(myProfile.Ban.Vote)
	removeUserRestrictions(db, userIDs)

	banComment(db, userIDs[0])
	myProfile = get(0)
	req.Zero(myProfile.Ban.Invite)
	req.Zero(myProfile.Ban.Vote)
	req.NotZero(myProfile.Ban.Comment)
	removeUserRestrictions(db, userIDs)

	banLive(db, userIDs[0])
	myProfile = get(0)
	req.Zero(myProfile.Ban.Invite)
	req.Zero(myProfile.Ban.Vote)
	req.Zero(myProfile.Ban.Comment)
	req.NotZero(myProfile.Ban.Live)
	removeUserRestrictions(db, userIDs)

	banShadow(db, userIDs[0])
	myProfile = get(0)
	req.Zero(myProfile.Ban.Invite)
	req.Zero(myProfile.Ban.Vote)
	req.Zero(myProfile.Ban.Comment)
	req.Zero(myProfile.Ban.Live)
	removeUserRestrictions(db, userIDs)
}

func compareUsers(t *testing.T, user *models.AuthProfile, profile *models.Profile) {
	req := require.New(t)

	req.Equal(user.ID, profile.ID)
	req.Equal(user.Name, profile.Name)
	req.Equal(user.ShowName, profile.ShowName)
	req.Equal(user.IsOnline, profile.IsOnline)
	req.Equal(user.Avatar, profile.Avatar)

	req.Equal(user.Gender, profile.Gender)
	req.Equal(user.IsDaylog, profile.IsDaylog)
	req.Equal(user.Privacy, profile.Privacy)
	req.Equal(user.ChatPrivacy, profile.ChatPrivacy)
	req.Equal(user.Title, profile.Title)
	req.Equal(user.Rank, profile.Rank)
	req.Equal(user.CreatedAt, profile.CreatedAt)
	req.Equal(user.LastSeenAt, profile.LastSeenAt)
	req.Equal(user.InvitedBy, profile.InvitedBy)
	req.Equal(user.AgeLowerBound, profile.AgeLowerBound)
	req.Equal(user.AgeUpperBound, profile.AgeUpperBound)
	req.Equal(user.Country, profile.Country)
	req.Equal(user.City, profile.City)
	req.Equal(user.Cover, profile.Cover)
	req.NotEmpty(user.Cover)
}

func TestGetUser(t *testing.T) {
	get := api.UsersGetUsersNameHandler.Handle
	params := users.GetUsersNameParams{}

	for i, user := range profiles {
		params.Name = strings.ToUpper(user.Name)
		resp := get(params, userIDs[i])
		body, ok := resp.(*users.GetUsersNameOK)
		if !ok {
			badBody, ok := resp.(*users.GetUsersNameNotFound)
			if ok {
				t.Fatal(badBody.Payload.Message)
			}

			t.Fatalf("error get user by name %s", user.Name)
		}

		compareUsers(t, user, body.Payload)
	}

	params.Name = "trolol not found"
	resp := get(params, userIDs[0])
	_, ok := resp.(*users.GetUsersNameNotFound)
	require.True(t, ok)

	noAuthUser := utils.NoAuthUser()
	params.Name = userIDs[0].Name

	setUserPrivacy(t, userIDs[0], "registered")
	resp = get(params, noAuthUser)
	_, ok = resp.(*users.GetUsersNameNotFound)
	require.True(t, ok)

	setUserPrivacy(t, userIDs[0], "all")
	resp = get(params, noAuthUser)
	body, ok := resp.(*users.GetUsersNameOK)
	require.True(t, ok)
	require.Equal(t, userIDs[0].ID, body.Payload.ID)
}

func checkEditProfile(t *testing.T, user *models.AuthProfile, params me.PutMeParams) {
	edit := api.MePutMeHandler.Handle
	id := models.UserID{
		ID:   user.ID,
		Name: user.Name,
		Ban:  &models.UserIDBan{},
	}
	resp := edit(params, &id)
	body, ok := resp.(*me.PutMeOK)
	require.True(t, ok)

	profile := body.Payload
	compareUsers(t, user, profile)
}

func setUserPrivacy(t *testing.T, userID *models.UserID, privacy string) {
	params := me.PutMeParams{
		Privacy:     privacy,
		ChatPrivacy: "invited",
		ShowName:    userID.Name,
	}
	edit := api.MePutMeHandler.Handle
	resp := edit(params, userID)
	_, ok := resp.(*me.PutMeOK)
	require.True(t, ok)
}

func setUserChatPrivacy(t *testing.T, userID *models.UserID, chatPrivacy, privacy string) {
	params := me.PutMeParams{
		ChatPrivacy: chatPrivacy,
		Privacy:     privacy,
		ShowName:    userID.Name,
	}
	edit := api.MePutMeHandler.Handle
	resp := edit(params, userID)
	_, ok := resp.(*me.PutMeOK)
	require.True(t, ok)
}

func TestEditProfile(t *testing.T) {
	user := *profiles[0]
	user.AgeLowerBound = 35
	user.AgeUpperBound = 40
	user.Birthday = "1990-01-01T20:01:31.844+03:00"
	user.City = "city edit"
	user.Country = "country edit"
	user.Gender = "female"
	user.IsDaylog = true
	user.Privacy = "followers"
	user.ChatPrivacy = "me"
	user.Title = "title edit"
	user.ShowInTops = false
	user.ShowName = "showname edit"

	params := me.PutMeParams{
		Birthday:    &user.Birthday,
		City:        &user.City,
		Country:     &user.Country,
		Gender:      &user.Gender,
		IsDaylog:    &user.IsDaylog,
		Privacy:     user.Privacy,
		ChatPrivacy: user.ChatPrivacy,
		Title:       &user.Title,
		ShowInTops:  &user.ShowInTops,
		ShowName:    user.ShowName,
	}

	checkEditProfile(t, &user, params)

	user.Privacy = "all"
	params.Privacy = user.Privacy
	user.ChatPrivacy = "invited"
	params.ChatPrivacy = user.ChatPrivacy
	checkEditProfile(t, &user, params)
}

func TestIsOpenForMe(t *testing.T) {
	req := require.New(t)

	check := func(userID *models.UserID, name string, res bool) {
		tx := utils.NewAutoTx(db)
		defer tx.Finish()
		req.Equal(res, utils.CanViewTlogName(tx, userID, name))
	}

	noAuthUser := utils.NoAuthUser()

	check(userIDs[0], userIDs[0].Name, true)
	check(userIDs[1], userIDs[0].Name, true)
	check(userIDs[2], userIDs[0].Name, true)
	check(userIDs[3], userIDs[0].Name, true)
	check(noAuthUser, userIDs[0].Name, true)

	checkFollow(t, userIDs[1], userIDs[0], profiles[0], models.RelationshipRelationFollowed, true)
	setUserPrivacy(t, userIDs[0], "followers")
	checkFollow(t, userIDs[2], userIDs[0], profiles[0], models.RelationshipRelationRequested, true)

	check(userIDs[0], userIDs[0].Name, true)
	check(userIDs[1], userIDs[0].Name, true)
	check(userIDs[2], userIDs[0].Name, false)
	check(userIDs[3], userIDs[0].Name, false)
	check(noAuthUser, userIDs[0].Name, false)

	setUserPrivacy(t, userIDs[0], "invited")

	check(userIDs[0], userIDs[0].Name, true)
	check(userIDs[1], userIDs[0].Name, true)
	check(userIDs[2], userIDs[0].Name, true)
	check(userIDs[3], userIDs[0].Name, false)
	check(noAuthUser, userIDs[0].Name, false)

	checkFollow(t, userIDs[0], userIDs[1], profiles[1], models.RelationshipRelationIgnored, true)

	check(userIDs[0], userIDs[0].Name, true)
	check(userIDs[1], userIDs[0].Name, false)
	check(userIDs[2], userIDs[0].Name, true)
	check(userIDs[3], userIDs[0].Name, false)
	check(noAuthUser, userIDs[0].Name, false)

	setUserPrivacy(t, userIDs[0], "registered")

	check(userIDs[0], userIDs[0].Name, true)
	check(userIDs[1], userIDs[0].Name, false)
	check(userIDs[2], userIDs[0].Name, true)
	check(userIDs[3], userIDs[0].Name, true)
	check(noAuthUser, userIDs[0].Name, false)

	setUserPrivacy(t, userIDs[0], "all")

	check(userIDs[0], userIDs[0].Name, true)
	check(userIDs[1], userIDs[0].Name, false)
	check(userIDs[2], userIDs[0].Name, true)
	check(userIDs[3], userIDs[0].Name, true)
	check(noAuthUser, userIDs[0].Name, true)

	checkUnfollow(t, userIDs[0], userIDs[1])
	checkUnfollow(t, userIDs[1], userIDs[0])
	checkUnfollow(t, userIDs[2], userIDs[0])

	banShadow(db, userIDs[0])
	banShadow(db, userIDs[1])
	setUserPrivacy(t, userIDs[1], "invited")
	setUserPrivacy(t, userIDs[2], "invited")

	check(userIDs[0], userIDs[0].Name, true)
	check(userIDs[1], userIDs[0].Name, true)
	check(userIDs[2], userIDs[0].Name, true)
	check(userIDs[3], userIDs[0].Name, true)
	check(noAuthUser, userIDs[0].Name, true)

	check(userIDs[0], userIDs[0].Name, true)
	check(userIDs[0], userIDs[1].Name, true)
	check(userIDs[0], userIDs[2].Name, false)
	check(userIDs[0], userIDs[3].Name, true)

	setUserPrivacy(t, userIDs[2], "registered")
	check(userIDs[0], userIDs[2].Name, true)

	setUserPrivacy(t, userIDs[2], "followers")
	check(userIDs[0], userIDs[2].Name, false)

	checkFollow(t, userIDs[0], userIDs[3], profiles[3], models.RelationshipRelationFollowed, true)
	setUserPrivacy(t, userIDs[3], "followers")
	check(userIDs[0], userIDs[3].Name, true)

	setUserPrivacy(t, userIDs[2], "all")
	setUserPrivacy(t, userIDs[3], "all")
	checkUnfollow(t, userIDs[0], userIDs[3])
	removeUserRestrictions(db, userIDs)
}

func TestIsChatAllowed(t *testing.T) {
	req := require.New(t)

	check := func(userID, to *models.UserID, res bool) {
		load := api.UsersGetUsersNameHandler.Handle
		resp := load(users.GetUsersNameParams{Name: to.Name}, userID)
		body, ok := resp.(*users.GetUsersNameOK)
		if !ok {
			req.False(res)
			return
		}

		req.Equal(res, body.Payload.Rights.Chat)
	}

	noAuthUser := utils.NoAuthUser()

	check(userIDs[0], userIDs[0], true)
	check(userIDs[1], userIDs[0], true)
	check(userIDs[2], userIDs[0], true)
	check(userIDs[3], userIDs[0], false)
	check(noAuthUser, userIDs[0], false)

	checkFollow(t, userIDs[1], userIDs[0], profiles[0], models.RelationshipRelationFollowed, true)
	setUserChatPrivacy(t, userIDs[0], "followers", "followers")
	checkFollow(t, userIDs[2], userIDs[0], profiles[0], models.RelationshipRelationRequested, true)

	check(userIDs[0], userIDs[0], true)
	check(userIDs[1], userIDs[0], true)
	check(userIDs[2], userIDs[0], false)
	check(userIDs[3], userIDs[0], false)
	check(noAuthUser, userIDs[0], false)

	setUserChatPrivacy(t, userIDs[0], "friends", "all")
	checkFollow(t, userIDs[0], userIDs[2], profiles[2], models.RelationshipRelationFollowed, true)
	checkFollow(t, userIDs[2], userIDs[0], profiles[0], models.RelationshipRelationFollowed, true)
	checkFollow(t, userIDs[0], userIDs[3], profiles[3], models.RelationshipRelationFollowed, true)
	checkFollow(t, userIDs[3], userIDs[0], profiles[0], models.RelationshipRelationFollowed, true)

	check(userIDs[0], userIDs[0], true)
	check(userIDs[1], userIDs[0], false)
	check(userIDs[2], userIDs[0], true)
	check(userIDs[3], userIDs[0], true)
	check(noAuthUser, userIDs[0], false)

	checkFollow(t, userIDs[0], userIDs[2], profiles[2], models.RelationshipRelationIgnored, true)

	check(userIDs[0], userIDs[0], true)
	check(userIDs[1], userIDs[0], false)
	check(userIDs[2], userIDs[0], false)
	check(userIDs[3], userIDs[0], true)
	check(noAuthUser, userIDs[0], false)

	setUserChatPrivacy(t, userIDs[0], "me", "all")

	check(userIDs[0], userIDs[0], true)
	check(userIDs[1], userIDs[0], false)
	check(userIDs[2], userIDs[0], false)
	check(userIDs[3], userIDs[0], false)
	check(noAuthUser, userIDs[0], false)

	setUserChatPrivacy(t, userIDs[0], "invited", "all")

	check(userIDs[0], userIDs[0], true)
	check(userIDs[1], userIDs[0], true)
	check(userIDs[2], userIDs[0], false)
	check(userIDs[3], userIDs[0], false)
	check(noAuthUser, userIDs[0], false)

	checkUnfollow(t, userIDs[1], userIDs[0])
	checkUnfollow(t, userIDs[2], userIDs[0])
	checkUnfollow(t, userIDs[3], userIDs[0])
	checkUnfollow(t, userIDs[0], userIDs[2])
	checkUnfollow(t, userIDs[0], userIDs[3])

	banShadow(db, userIDs[0])
	banShadow(db, userIDs[1])
	check(userIDs[0], userIDs[0], true)
	check(userIDs[1], userIDs[0], true)
	check(userIDs[2], userIDs[0], true)
	check(userIDs[0], userIDs[2], false)

	setUserChatPrivacy(t, userIDs[1], "followers", "all")
	check(userIDs[0], userIDs[1], false)
	checkFollow(t, userIDs[0], userIDs[1], profiles[1], models.RelationshipRelationFollowed, true)
	check(userIDs[0], userIDs[1], true)

	setUserChatPrivacy(t, userIDs[2], "followers", "all")
	check(userIDs[0], userIDs[2], false)
	checkFollow(t, userIDs[0], userIDs[2], profiles[2], models.RelationshipRelationFollowed, true)
	check(userIDs[0], userIDs[2], false)

	setUserChatPrivacy(t, userIDs[2], "friends", "all")
	check(userIDs[0], userIDs[2], false)
	checkFollow(t, userIDs[2], userIDs[0], profiles[0], models.RelationshipRelationFollowed, true)
	check(userIDs[0], userIDs[2], true)

	setUserChatPrivacy(t, userIDs[1], "invited", "all")
	setUserChatPrivacy(t, userIDs[2], "invited", "all")

	checkUnfollow(t, userIDs[0], userIDs[1])
	checkUnfollow(t, userIDs[0], userIDs[2])
	checkUnfollow(t, userIDs[2], userIDs[0])
}
