package test

import (
	"testing"

	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/me"
	"github.com/sevings/mindwell-server/restapi/operations/users"
	"github.com/stretchr/testify/require"
)

func checkLoadMyBadges(t *testing.T, userID *models.UserID, cnt int) *models.BadgeList {
	var limit int64 = 100
	load := api.MeGetMeBadgesHandler.Handle
	params := me.GetMeBadgesParams{Limit: &limit}
	resp := load(params, userID)
	body, ok := resp.(*me.GetMeBadgesOK)
	require.True(t, ok)
	if !ok {
		return nil
	}

	list := body.Payload
	require.NotNil(t, list)
	require.Equal(t, cnt, len(list.Data))

	return list
}

func checkLoadTlogBadges(t *testing.T, tlog string, userID *models.UserID, success bool, cnt int) *models.BadgeList {
	var limit int64 = 100
	load := api.UsersGetUsersNameBadgesHandler.Handle
	params := users.GetUsersNameBadgesParams{Name: tlog, Limit: &limit}
	resp := load(params, userID)
	body, ok := resp.(*users.GetUsersNameBadgesOK)
	require.Equal(t, success, ok)
	if !ok {
		return nil
	}

	list := body.Payload
	require.NotNil(t, list)
	require.Equal(t, cnt, len(list.Data))

	return list
}

func TestLoadBadges(t *testing.T) {
	checkLoadMyBadges(t, userIDs[0], 0)
	checkLoadMyBadges(t, userIDs[3], 0)
	checkLoadTlogBadges(t, userIDs[0].Name, userIDs[0], true, 0)
	checkLoadTlogBadges(t, userIDs[1].Name, userIDs[0], true, 0)
	checkLoadTlogBadges(t, userIDs[0].Name, userIDs[3], true, 0)

	setUserPrivacy(t, userIDs[0], models.EntryPrivacyInvited)
	checkLoadTlogBadges(t, userIDs[0].Name, userIDs[1], true, 0)
	checkLoadTlogBadges(t, userIDs[0].Name, userIDs[3], false, 0)
	setUserPrivacy(t, userIDs[0], models.EntryPrivacyAll)

	follow(userIDs[0], userIDs[1].Name, models.RelationshipRelationIgnored)
	checkLoadTlogBadges(t, userIDs[0].Name, userIDs[1], false, 0)
	checkLoadTlogBadges(t, userIDs[0].Name, userIDs[3], true, 0)
	unfollow(userIDs[0], userIDs[1].Name)

	addBadgeQuery := `
	INSERT INTO user_badges(user_id, badge_id)
	VALUES($1, (SELECT id FROM badges WHERE code = $2))
	`
	db.Exec(addBadgeQuery, userIDs[0].ID, "test1")
	db.Exec(addBadgeQuery, userIDs[0].ID, "test2")
	db.Exec(addBadgeQuery, userIDs[1].ID, "test1")

	list := checkLoadMyBadges(t, userIDs[0], 2)
	require.True(t, list.Data[0].Code == "test2")
	require.True(t, list.Data[1].Code == "test1")

	list = checkLoadMyBadges(t, userIDs[1], 1)
	require.True(t, list.Data[0].Code == "test1")

	checkLoadTlogBadges(t, userIDs[0].Name, userIDs[0], true, 2)
	checkLoadTlogBadges(t, userIDs[1].Name, userIDs[0], true, 1)

	rmAllBadgesQuery := "DELETE FROM user_badges"
	db.Exec(rmAllBadgesQuery)
}
