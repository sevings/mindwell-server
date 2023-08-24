package test

import (
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/me"
	"github.com/sevings/mindwell-server/restapi/operations/themes"
	"github.com/sevings/mindwell-server/restapi/operations/users"
	"github.com/sevings/mindwell-server/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func checkLoadMyImages(t *testing.T, user *models.UserID) *models.ImageList {
	str := ""
	var limit int64 = 10
	params := me.GetMeImagesParams{
		After:  &str,
		Before: &str,
		Limit:  &limit,
	}
	load := api.MeGetMeImagesHandler.Handle
	resp := load(params, user)
	body, ok := resp.(*me.GetMeImagesOK)
	if !ok {
		t.Fatal("error load images")
	}

	return body.Payload
}

func TestLoadMyImages(t *testing.T) {
	feed := checkLoadMyImages(t, userIDs[0])
	req := require.New(t)
	req.Empty(feed.Data)
	req.False(feed.HasAfter)
	req.False(feed.HasBefore)
}

func checkLoadTlogImages(t *testing.T, tlog, user *models.UserID, success bool) *models.ImageList {
	str := ""
	var limit int64 = 10
	params := users.GetUsersNameImagesParams{
		After:  &str,
		Before: &str,
		Limit:  &limit,
		Name:   tlog.Name,
	}
	load := api.UsersGetUsersNameImagesHandler.Handle
	resp := load(params, user)
	body, ok := resp.(*users.GetUsersNameImagesOK)
	require.Equal(t, success, ok)
	if !ok {
		return nil
	}

	return body.Payload
}

func TestLoadTlogImages(t *testing.T) {
	req := require.New(t)
	feed := checkLoadTlogImages(t, userIDs[0], userIDs[0], true)
	req.Empty(feed.Data)
	req.False(feed.HasAfter)
	req.False(feed.HasBefore)

	feed = checkLoadTlogImages(t, userIDs[1], userIDs[0], true)
	req.Empty(feed.Data)
	req.False(feed.HasAfter)
	req.False(feed.HasBefore)

	setUserPrivacy(t, userIDs[0], "followers")
	checkLoadTlogImages(t, userIDs[0], userIDs[1], false)
	checkFollow(t, userIDs[1], userIDs[0], profiles[0], models.RelationshipRelationRequested, true)
	checkPermitFollow(t, userIDs[0], userIDs[1], true)
	checkLoadTlogImages(t, userIDs[0], userIDs[1], true)
	checkUnfollow(t, userIDs[1], userIDs[0])

	setUserPrivacy(t, userIDs[0], "invited")
	checkLoadTlogImages(t, userIDs[0], userIDs[1], true)
	checkLoadTlogImages(t, userIDs[0], userIDs[3], false)

	setUserPrivacy(t, userIDs[0], "registered")
	checkLoadTlogImages(t, userIDs[0], userIDs[1], true)
	checkLoadTlogImages(t, userIDs[0], userIDs[3], true)
	checkLoadTlogImages(t, userIDs[0], utils.NoAuthUser(), false)

	setUserPrivacy(t, userIDs[0], "all")
}

func checkLoadThemeImages(t *testing.T, user *models.UserID, name string, success bool) *models.ImageList {
	str := ""
	var limit int64 = 10
	params := themes.GetThemesNameImagesParams{
		After:  &str,
		Before: &str,
		Limit:  &limit,
		Name:   name,
	}
	load := api.ThemesGetThemesNameImagesHandler.Handle
	resp := load(params, user)
	body, ok := resp.(*themes.GetThemesNameImagesOK)
	require.Equal(t, success, ok)
	if !ok {
		return nil
	}

	return body.Payload
}

func TestLoadThemeImages(t *testing.T) {
	req := require.New(t)
	theme := createTestTheme(t, userIDs[0])
	feed := checkLoadThemeImages(t, userIDs[0], theme.Name, true)
	req.Empty(feed.Data)
	req.False(feed.HasAfter)
	req.False(feed.HasBefore)

	feed = checkLoadThemeImages(t, userIDs[1], theme.Name, true)
	req.Empty(feed.Data)
	req.False(feed.HasAfter)
	req.False(feed.HasBefore)

	setThemePrivacy(t, userIDs[0], theme, "invited")
	checkLoadThemeImages(t, userIDs[1], theme.Name, true)
	checkLoadThemeImages(t, userIDs[3], theme.Name, false)

	setThemePrivacy(t, userIDs[0], theme, "registered")
	checkLoadThemeImages(t, userIDs[1], theme.Name, true)
	checkLoadThemeImages(t, userIDs[3], theme.Name, true)
	checkLoadThemeImages(t, utils.NoAuthUser(), theme.Name, false)

	deleteTheme(t, theme)
}
