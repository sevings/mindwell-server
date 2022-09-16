package test

import (
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/account"
	"github.com/sevings/mindwell-server/restapi/operations/themes"
	"github.com/sevings/mindwell-server/utils"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func loadInvites(t *testing.T, userID *models.UserID) int {
	load := api.AccountGetAccountInvitesHandler.Handle
	resp := load(account.GetAccountInvitesParams{}, userID)
	body, ok := resp.(*account.GetAccountInvitesOK)
	require.True(t, ok)
	return len(body.Payload.Invites)
}

func checkCreateTheme(t *testing.T, success bool, user *models.UserID, name, showName string) *models.Profile {
	req := require.New(t)

	inv := loadInvites(t, user)

	post := api.ThemesPostThemesHandler.Handle
	params := themes.PostThemesParams{
		Name:     name,
		ShowName: showName,
	}

	resp := post(params, user)
	body, ok := resp.(*themes.PostThemesOK)
	req.Equal(success, ok)
	if !ok {
		return nil
	}

	theme := body.Payload
	req.True(theme.IsTheme)
	req.Empty(theme.Gender)
	req.Zero(theme.LastSeenAt)
	req.False(theme.IsOnline)
	req.Equal(name, theme.Name)
	req.Equal(showName, theme.ShowName)
	req.NotNil(theme.CreatedBy)
	req.Equal(theme.CreatedBy.ID, user.ID)
	req.Equal(theme.CreatedBy.Name, user.Name)
	req.Nil(theme.InvitedBy)

	req.Equal(inv-1, loadInvites(t, user))

	return theme
}

func deleteTheme(t *testing.T, theme *models.Profile) {
	_, err := db.Exec("SELECT delete_user($1)", theme.Name)
	if err != nil {
		log.Println(err)
	}
	require.Nil(t, err)
}

func giveInvite(t *testing.T, user *models.UserID) {
	_, err := db.Exec("SELECT give_invite($1)", user.Name)
	if err != nil {
		log.Println(err)
	}
	require.Nil(t, err)
}

func TestCreateTheme(t *testing.T) {
	require.Zero(t, loadInvites(t, userIDs[0]))
	checkCreateTheme(t, false, userIDs[0], "test_theme", "test test")

	giveInvite(t, userIDs[0])
	theme := checkCreateTheme(t, true, userIDs[0], "test_theme", "test test test")
	deleteTheme(t, theme)
}

func checkLoadTheme(t *testing.T, success bool, user *models.UserID, exp *models.Profile) {
	req := require.New(t)

	load := api.ThemesGetThemesNameHandler.Handle
	params := themes.GetThemesNameParams{
		Name: exp.Name,
	}

	resp := load(params, user)
	body, ok := resp.(*themes.GetThemesNameOK)
	req.Equal(success, ok)
	if !ok {
		return
	}

	theme := body.Payload
	req.True(theme.IsTheme)
	req.Empty(theme.Gender)
	req.Zero(theme.LastSeenAt)
	req.False(theme.IsOnline)
	req.Equal(exp.Name, theme.Name)
	req.Equal(exp.ShowName, theme.ShowName)
	req.Equal(exp.Title, theme.Title)
	req.Equal(exp.Privacy, theme.Privacy)
	req.NotNil(theme.CreatedBy)
	req.Equal(*exp.CreatedBy, *theme.CreatedBy)
	req.Nil(theme.InvitedBy)
}

func createTestTheme(t *testing.T, user *models.UserID) *models.Profile {
	giveInvite(t, user)

	name := utils.GenerateString(7)
	theme := checkCreateTheme(t, true, user, name, name+" "+name)
	checkLoadTheme(t, true, user, theme)

	return theme
}

func setThemePrivacy(t *testing.T, user *models.UserID, theme *models.Profile, privacy string) {
	save := api.ThemesPutThemesNameHandler.Handle
	params := themes.PutThemesNameParams{
		Name:     theme.Name,
		Privacy:  privacy,
		ShowName: theme.ShowName,
		Title:    &theme.Title,
	}

	resp := save(params, user)
	_, ok := resp.(*themes.PutThemesNameOK)
	require.True(t, ok)
}

func TestLoadTheme(t *testing.T) {
	checkLoadTheme(t, false, userIDs[0], &models.Profile{Friend: models.Friend{User: models.User{Name: "not found"}}})

	theme := createTestTheme(t, userIDs[0])

	noAuthUser := utils.NoAuthUser()

	setThemePrivacy(t, userIDs[0], theme, "registered")
	checkLoadTheme(t, false, noAuthUser, theme)

	setThemePrivacy(t, userIDs[0], theme, "all")
	checkLoadTheme(t, true, noAuthUser, theme)

	deleteTheme(t, theme)
}

func checkEditTheme(t *testing.T, success bool, user *models.UserID, name, showName, title, privacy string) *models.Profile {
	req := require.New(t)

	save := api.ThemesPutThemesNameHandler.Handle
	params := themes.PutThemesNameParams{
		Name:     name,
		Privacy:  privacy,
		ShowName: showName,
		Title:    &title,
	}

	resp := save(params, user)
	body, ok := resp.(*themes.PutThemesNameOK)
	req.Equal(success, ok)
	if !ok {
		return nil
	}

	theme := body.Payload
	req.True(theme.IsTheme)
	req.Empty(theme.Gender)
	req.Zero(theme.LastSeenAt)
	req.False(theme.IsOnline)
	req.Equal(name, theme.Name)
	req.Equal(showName, theme.ShowName)
	req.Equal(title, theme.Title)
	req.Equal(privacy, theme.Privacy)
	req.NotNil(theme.CreatedBy)
	req.Equal(theme.CreatedBy.ID, user.ID)
	req.Equal(theme.CreatedBy.Name, user.Name)
	req.Nil(theme.InvitedBy)

	checkLoadTheme(t, true, user, theme)

	return theme
}

func TestEditTheme(t *testing.T) {
	theme := createTestTheme(t, userIDs[0])

	checkEditTheme(t, false, userIDs[1], "", "edit show name", "theme title", "invited")
	checkEditTheme(t, true, userIDs[0], theme.Name, "edit show name", "theme title", "invited")

	deleteTheme(t, theme)

	checkEditTheme(t, false, userIDs[0], theme.Name, "edit show name", "theme title", "invited")
	checkEditTheme(t, false, userIDs[0], userIDs[0].Name, "edit show name", "theme title", "invited")
}

func checkTopThemes(t *testing.T, top string, size int) []*models.Friend {
	get := api.ThemesGetThemesHandler.Handle
	params := themes.GetThemesParams{Top: &top}
	resp := get(params, userIDs[0])
	body, ok := resp.(*themes.GetThemesOK)

	require.True(t, ok)

	list := body.Payload
	require.Equal(t, size, len(list.Themes))
	require.Equal(t, top, list.Top)
	require.Empty(t, list.Query)

	return list.Themes
}

func TestTopThemes(t *testing.T) {
	t1 := createTestTheme(t, userIDs[0])
	time.Sleep(time.Millisecond)
	t2 := createTestTheme(t, userIDs[0])
	time.Sleep(time.Millisecond)
	t3 := createTestTheme(t, userIDs[0])
	time.Sleep(time.Millisecond)
	t4 := createTestTheme(t, userIDs[0])

	list := checkTopThemes(t, "new", 4)

	req := require.New(t)
	req.Equal(t4.ID, list[0].ID)
	req.Equal(t3.ID, list[1].ID)
	req.Equal(t2.ID, list[2].ID)
	req.Equal(t1.ID, list[3].ID)

	checkTopThemes(t, "rank", 4)

	deleteTheme(t, t4)
	deleteTheme(t, t3)
	deleteTheme(t, t2)
	deleteTheme(t, t1)
}

func checkSearchThemes(t *testing.T, query string, size int) []*models.Friend {
	get := api.ThemesGetThemesHandler.Handle
	params := themes.GetThemesParams{Query: &query}
	resp := get(params, userIDs[0])
	body, ok := resp.(*themes.GetThemesOK)

	require.True(t, ok)

	list := body.Payload
	require.Equal(t, size, len(list.Themes))
	require.Equal(t, query, list.Query)
	require.Empty(t, list.Top)

	return list.Themes
}

func TestSearchThemes(t *testing.T) {
	giveInvite(t, userIDs[0])
	giveInvite(t, userIDs[0])
	giveInvite(t, userIDs[0])
	giveInvite(t, userIDs[0])

	t1 := checkCreateTheme(t, true, userIDs[0], "test0_theme", "test0")
	t2 := checkCreateTheme(t, true, userIDs[0], "test1_theme", "test1")
	t3 := checkCreateTheme(t, true, userIDs[0], "test2_theme", "test2")
	t4 := checkCreateTheme(t, true, userIDs[0], "mind_theme", "mind_t")

	checkSearchThemes(t, "testo", 3)
	checkSearchThemes(t, "mind", 1)
	checkSearchThemes(t, "psychotherapist", 0)

	deleteTheme(t, t4)
	deleteTheme(t, t3)
	deleteTheme(t, t2)
	deleteTheme(t, t1)
}

func TestIsThemeOpenForMe(t *testing.T) {
	theme := createTestTheme(t, userIDs[0])

	check := func(userID *models.UserID, name string, res bool) {
		tx := utils.NewAutoTx(db)
		defer tx.Finish()
		require.Equal(t, res, utils.CanViewTlogName(tx, userID, name))
	}

	noAuthUser := utils.NoAuthUser()

	check(userIDs[0], theme.Name, true)
	check(userIDs[1], theme.Name, true)
	check(userIDs[2], theme.Name, true)
	check(userIDs[3], theme.Name, true)
	check(noAuthUser, theme.Name, true)

	setThemePrivacy(t, userIDs[0], theme, "invited")

	check(userIDs[0], theme.Name, true)
	check(userIDs[1], theme.Name, true)
	check(userIDs[2], theme.Name, true)
	check(userIDs[3], theme.Name, false)
	check(noAuthUser, theme.Name, false)

	setThemePrivacy(t, userIDs[0], theme, "registered")

	check(userIDs[0], theme.Name, true)
	check(userIDs[1], theme.Name, true)
	check(userIDs[2], theme.Name, true)
	check(userIDs[3], theme.Name, true)
	check(noAuthUser, theme.Name, false)

	setThemePrivacy(t, userIDs[0], theme, "all")

	check(userIDs[0], theme.Name, true)
	check(userIDs[1], theme.Name, true)
	check(userIDs[2], theme.Name, true)
	check(userIDs[3], theme.Name, true)
	check(noAuthUser, theme.Name, true)

	deleteTheme(t, theme)
}

func TestThemeCreatorRights(t *testing.T) {
	theme := createTestTheme(t, userIDs[0])

	toTheme := &models.AuthProfile{Profile: *theme}
	checkFollow(t, userIDs[1], nil, toTheme, models.RelationshipRelationFollowed, true)
	checkFollow(t, userIDs[2], nil, toTheme, models.RelationshipRelationFollowed, true)
	checkFollow(t, userIDs[3], nil, toTheme, models.RelationshipRelationFollowed, true)

	e := createThemeEntry(t, userIDs[1], theme.Name, models.EntryPrivacyAll, true, true, true, true)
	checkDeleteEntry(t, e.ID, userIDs[1], true)

	e = createThemeEntry(t, userIDs[2], theme.Name, models.EntryPrivacyAll, true, true, true, true)
	checkDeleteEntry(t, e.ID, userIDs[1], false)
	checkDeleteEntry(t, e.ID, userIDs[0], true)

	e = createThemeEntry(t, userIDs[1], theme.Name, models.EntryPrivacyAll, true, true, true, true)

	c := postComment(userIDs[1], e.ID)
	checkDeleteComment(t, c, userIDs[1], true)

	c = postComment(userIDs[2], e.ID)
	checkDeleteComment(t, c, userIDs[1], true)

	c = postComment(userIDs[2], e.ID)
	checkDeleteComment(t, c, userIDs[2], true)

	c = postComment(userIDs[1], e.ID)
	checkDeleteComment(t, c, userIDs[2], false)
	checkDeleteComment(t, c, userIDs[3], false)
	checkDeleteComment(t, c, userIDs[0], true)

	checkDeleteEntry(t, e.ID, userIDs[1], true)

	following := &models.UserID{Name: theme.Name}
	checkUnfollow(t, userIDs[1], following)
	checkUnfollow(t, userIDs[2], following)
	checkUnfollow(t, userIDs[3], following)
}
