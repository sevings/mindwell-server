package test

import (
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/comments"
	"github.com/sevings/mindwell-server/restapi/operations/themes"
	"github.com/sevings/mindwell-server/restapi/operations/users"
	"github.com/stretchr/testify/require"
	"testing"
)

func checkCanLoadComment(t *testing.T, commentID int64, userID *models.UserID) {
	load := api.CommentsGetCommentsIDHandler.Handle
	resp := load(comments.GetCommentsIDParams{ID: commentID}, userID)
	_, ok := resp.(*comments.GetCommentsIDOK)
	require.True(t, ok)
}

func checkLoadUserComments(t *testing.T, user *models.UserID, author *models.AuthProfile, limit int64, before, after string, size int) *models.CommentList {
	req := require.New(t)
	params := users.GetUsersNameCommentsParams{
		After:  &after,
		Before: &before,
		Limit:  &limit,
		Name:   author.Name,
	}
	load := api.UsersGetUsersNameCommentsHandler.Handle
	resp := load(params, user)
	body, ok := resp.(*users.GetUsersNameCommentsOK)
	if !ok {
		req.Zero(size)
		return nil
	}

	list := body.Payload
	req.Equal(size, len(list.Data))

	for _, cmt := range list.Data {
		req.Equal(author.ID, cmt.Author.ID)

		checkCanLoadComment(t, cmt.ID, user)
	}

	return list
}

func TestUserComments(t *testing.T) {
	theme := createTestTheme(t, userIDs[0])
	e1 := postThemeEntry(userIDs[0], theme.Name, models.EntryPrivacyAll, true).ID
	postComment(userIDs[0], e1)
	postComment(userIDs[1], e1)

	e2 := postEntry(userIDs[0], models.EntryPrivacyAll, true).ID
	id1 := postComment(userIDs[0], e2)
	id2 := postComment(userIDs[0], e2)
	postComment(userIDs[1], e2)

	e3 := postEntry(userIDs[0], models.EntryPrivacyFollowers, true).ID
	id3 := postComment(userIDs[0], e3)

	follow(userIDs[1], userIDs[0].Name, models.RelationshipRelationFollowed)
	postComment(userIDs[1], e3)

	e4 := postEntry(userIDs[1], models.EntryPrivacyAll, true).ID
	id4 := postComment(userIDs[0], e4)
	follow(userIDs[1], userIDs[2].Name, models.RelationshipRelationIgnored)

	e5 := postEntry(userIDs[0], models.EntryPrivacyMe, false).ID
	id5 := postComment(userIDs[0], e5)

	req := require.New(t)
	list := checkLoadUserComments(t, userIDs[0], profiles[0], 10, "", "", 5)
	req.Equal(id5, list.Data[0].ID)
	req.Equal(id4, list.Data[1].ID)
	req.Equal(id3, list.Data[2].ID)
	req.Equal(id2, list.Data[3].ID)
	req.Equal(id1, list.Data[4].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadUserComments(t, userIDs[1], profiles[0], 10, "", "", 4)
	req.Equal(id4, list.Data[0].ID)
	req.Equal(id3, list.Data[1].ID)
	req.Equal(id2, list.Data[2].ID)
	req.Equal(id1, list.Data[3].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadUserComments(t, userIDs[2], profiles[0], 10, "", "", 2)
	req.Equal(id2, list.Data[0].ID)
	req.Equal(id1, list.Data[1].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	checkLoadUserComments(t, userIDs[0], profiles[1], 10, "", "", 3)
	checkLoadUserComments(t, userIDs[0], profiles[2], 10, "", "", 0)
	checkLoadUserComments(t, userIDs[2], profiles[1], 10, "", "", 0)

	list = checkLoadUserComments(t, userIDs[1], profiles[0], 3, "", "", 3)
	req.Equal(id4, list.Data[0].ID)
	req.Equal(id3, list.Data[1].ID)
	req.Equal(id2, list.Data[2].ID)
	req.False(list.HasAfter)
	req.True(list.HasBefore)

	list = checkLoadUserComments(t, userIDs[1], profiles[0], 10, list.NextBefore, "", 1)
	req.Equal(id1, list.Data[0].ID)
	req.True(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadUserComments(t, userIDs[1], profiles[0], 10, "", list.NextAfter, 3)
	req.Equal(id4, list.Data[0].ID)
	req.Equal(id3, list.Data[1].ID)
	req.Equal(id2, list.Data[2].ID)
	req.False(list.HasAfter)
	req.True(list.HasBefore)

	checkUnfollow(t, userIDs[1], userIDs[0])
	checkUnfollow(t, userIDs[1], userIDs[2])

	checkDeleteEntry(t, e1, userIDs[0], true)
	checkDeleteEntry(t, e2, userIDs[0], true)
	checkDeleteEntry(t, e3, userIDs[0], true)
	checkDeleteEntry(t, e4, userIDs[1], true)
	checkDeleteEntry(t, e5, userIDs[0], true)

	deleteTheme(t, theme)
}

func checkLoadThemeComments(t *testing.T, user *models.UserID, theme *models.Profile, limit int64, before, after string, size int) *models.CommentList {
	req := require.New(t)
	params := themes.GetThemesNameCommentsParams{
		After:  &after,
		Before: &before,
		Limit:  &limit,
		Name:   theme.Name,
	}
	load := api.ThemesGetThemesNameCommentsHandler.Handle
	resp := load(params, user)
	body, ok := resp.(*themes.GetThemesNameCommentsOK)
	if !ok {
		req.Zero(size)
		return nil
	}

	list := body.Payload
	req.Equal(size, len(list.Data))

	for _, cmt := range list.Data {
		entry := loadEntry(user, cmt.EntryID)
		req.Equal(theme.ID, entry.Author.ID)

		checkCanLoadComment(t, cmt.ID, user)
	}

	return list
}

func TestThemeComments(t *testing.T) {
	theme := createTestTheme(t, userIDs[0])
	e1 := postThemeEntry(userIDs[0], theme.Name, models.EntryPrivacyAll, true).ID
	id1 := postComment(userIDs[0], e1)
	id2 := postComment(userIDs[1], e1)

	e2 := postEntry(userIDs[0], models.EntryPrivacyAll, true).ID
	postComment(userIDs[0], e2)
	postComment(userIDs[1], e2)

	follow(userIDs[1], theme.Name, models.RelationshipRelationFollowed)
	e3 := postThemeEntry(userIDs[1], theme.Name, models.EntryPrivacyFollowers, false).ID
	id3 := postComment(userIDs[1], e3)

	req := require.New(t)
	list := checkLoadThemeComments(t, userIDs[0], theme, 10, "", "", 3)
	req.Equal(id3, list.Data[0].ID)
	req.Equal(id2, list.Data[1].ID)
	req.Equal(id1, list.Data[2].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadThemeComments(t, userIDs[1], theme, 10, "", "", 3)
	req.Equal(id3, list.Data[0].ID)
	req.Equal(id2, list.Data[1].ID)
	req.Equal(id1, list.Data[2].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadThemeComments(t, userIDs[2], theme, 10, "", "", 2)
	req.Equal(id2, list.Data[0].ID)
	req.Equal(id1, list.Data[1].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadThemeComments(t, userIDs[1], theme, 2, "", "", 2)
	req.Equal(id3, list.Data[0].ID)
	req.Equal(id2, list.Data[1].ID)
	req.False(list.HasAfter)
	req.True(list.HasBefore)

	list = checkLoadThemeComments(t, userIDs[1], theme, 10, list.NextBefore, "", 1)
	req.Equal(id1, list.Data[0].ID)
	req.True(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadThemeComments(t, userIDs[1], theme, 10, "", list.NextAfter, 2)
	req.Equal(id3, list.Data[0].ID)
	req.Equal(id2, list.Data[1].ID)
	req.False(list.HasAfter)
	req.True(list.HasBefore)

	unfollow(userIDs[1], theme.Name)

	checkDeleteEntry(t, e1, userIDs[0], true)
	checkDeleteEntry(t, e2, userIDs[0], true)
	checkDeleteEntry(t, e3, userIDs[0], true)

	deleteTheme(t, theme)
}

func checkLoadAllComments(t *testing.T, user *models.UserID, limit int64, before, after string, size int) *models.CommentList {
	req := require.New(t)
	params := comments.GetCommentsParams{
		After:  &after,
		Before: &before,
		Limit:  &limit,
	}
	load := api.CommentsGetCommentsHandler.Handle
	resp := load(params, user)
	body, ok := resp.(*comments.GetCommentsOK)
	if !ok {
		req.Zero(size)
		return nil
	}

	list := body.Payload
	req.Equal(size, len(list.Data))

	for _, cmt := range list.Data {
		checkCanLoadComment(t, cmt.ID, user)
	}

	return list
}

func TestAllComments(t *testing.T) {
	theme := createTestTheme(t, userIDs[0])
	e1 := postThemeEntry(userIDs[0], theme.Name, models.EntryPrivacyAll, true).ID
	id1 := postComment(userIDs[0], e1)
	id2 := postComment(userIDs[1], e1)

	e2 := postEntry(userIDs[0], models.EntryPrivacyAll, true).ID
	id3 := postComment(userIDs[0], e2)
	id4 := postComment(userIDs[1], e2)

	e3 := postEntry(userIDs[0], models.EntryPrivacyFollowers, true).ID
	id5 := postComment(userIDs[0], e3)

	follow(userIDs[1], userIDs[0].Name, models.RelationshipRelationFollowed)
	id6 := postComment(userIDs[1], e3)

	e4 := postEntry(userIDs[1], models.EntryPrivacyAll, true).ID
	id7 := postComment(userIDs[0], e4)
	follow(userIDs[1], userIDs[2].Name, models.RelationshipRelationIgnored)

	e5 := postEntry(userIDs[0], models.EntryPrivacyMe, false).ID
	id8 := postComment(userIDs[0], e5)

	req := require.New(t)
	list := checkLoadAllComments(t, userIDs[0], 10, "", "", 8)
	req.Equal(id8, list.Data[0].ID)
	req.Equal(id7, list.Data[1].ID)
	req.Equal(id6, list.Data[2].ID)
	req.Equal(id5, list.Data[3].ID)
	req.Equal(id4, list.Data[4].ID)
	req.Equal(id3, list.Data[5].ID)
	req.Equal(id2, list.Data[6].ID)
	req.Equal(id1, list.Data[7].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadAllComments(t, userIDs[1], 10, "", "", 7)
	req.Equal(id7, list.Data[0].ID)
	req.Equal(id6, list.Data[1].ID)
	req.Equal(id5, list.Data[2].ID)
	req.Equal(id4, list.Data[3].ID)
	req.Equal(id3, list.Data[4].ID)
	req.Equal(id2, list.Data[5].ID)
	req.Equal(id1, list.Data[6].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadAllComments(t, userIDs[2], 10, "", "", 4)
	req.Equal(id4, list.Data[0].ID)
	req.Equal(id3, list.Data[1].ID)
	req.Equal(id2, list.Data[2].ID)
	req.Equal(id1, list.Data[3].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadAllComments(t, userIDs[1], 3, "", "", 3)
	req.Equal(id7, list.Data[0].ID)
	req.Equal(id6, list.Data[1].ID)
	req.Equal(id5, list.Data[2].ID)
	req.False(list.HasAfter)
	req.True(list.HasBefore)

	list = checkLoadAllComments(t, userIDs[1], 10, list.NextBefore, "", 4)
	req.Equal(id4, list.Data[0].ID)
	req.Equal(id3, list.Data[1].ID)
	req.Equal(id2, list.Data[2].ID)
	req.Equal(id1, list.Data[3].ID)
	req.True(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadAllComments(t, userIDs[1], 10, "", list.NextAfter, 3)
	req.Equal(id7, list.Data[0].ID)
	req.Equal(id6, list.Data[1].ID)
	req.Equal(id5, list.Data[2].ID)
	req.False(list.HasAfter)
	req.True(list.HasBefore)

	checkUnfollow(t, userIDs[1], userIDs[0])
	checkUnfollow(t, userIDs[1], userIDs[2])

	checkDeleteEntry(t, e1, userIDs[0], true)
	checkDeleteEntry(t, e2, userIDs[0], true)
	checkDeleteEntry(t, e3, userIDs[0], true)
	checkDeleteEntry(t, e4, userIDs[1], true)
	checkDeleteEntry(t, e5, userIDs[0], true)

	deleteTheme(t, theme)
}

func checkLoadEntryComments(t *testing.T, user *models.UserID, entryID, limit int64, before, after string, size int) *models.CommentList {
	req := require.New(t)
	params := comments.GetEntriesIDCommentsParams{
		After:  &after,
		Before: &before,
		Limit:  &limit,
		ID:     entryID,
	}
	load := api.CommentsGetEntriesIDCommentsHandler.Handle
	resp := load(params, user)
	body, ok := resp.(*comments.GetEntriesIDCommentsOK)
	if !ok {
		req.Zero(size)
		return nil
	}

	list := body.Payload
	req.Equal(size, len(list.Data))

	for _, cmt := range list.Data {
		req.Equal(entryID, cmt.EntryID)

		checkCanLoadComment(t, cmt.ID, user)
	}

	return list
}

func TestEntryComments(t *testing.T) {
	theme := createTestTheme(t, userIDs[0])
	e1 := postThemeEntry(userIDs[0], theme.Name, models.EntryPrivacyAll, true).ID
	id1 := postComment(userIDs[0], e1)
	id2 := postComment(userIDs[1], e1)

	follow(userIDs[1], userIDs[0].Name, models.RelationshipRelationFollowed)
	follow(userIDs[0], userIDs[2].Name, models.RelationshipRelationIgnored)

	req := require.New(t)
	list := checkLoadEntryComments(t, userIDs[0], e1, 10, "", "", 2)
	req.Equal(id1, list.Data[0].ID)
	req.Equal(id2, list.Data[1].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	checkLoadEntryComments(t, userIDs[1], e1, 10, "", "", 2)
	checkLoadEntryComments(t, userIDs[2], e1, 10, "", "", 2)

	e2 := postEntry(userIDs[0], models.EntryPrivacyAll, true).ID
	id3 := postComment(userIDs[0], e2)
	id4 := postComment(userIDs[0], e2)
	id5 := postComment(userIDs[1], e2)

	list = checkLoadEntryComments(t, userIDs[1], e2, 10, "", "", 3)
	req.Equal(id3, list.Data[0].ID)
	req.Equal(id4, list.Data[1].ID)
	req.Equal(id5, list.Data[2].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	checkLoadEntryComments(t, userIDs[2], e2, 10, "", "", 0)
	checkLoadEntryComments(t, userIDs[3], e2, 10, "", "", 3)

	e3 := postEntry(userIDs[0], models.EntryPrivacyFollowers, true).ID
	id6 := postComment(userIDs[0], e3)

	list = checkLoadEntryComments(t, userIDs[1], e3, 10, "", "", 1)
	req.Equal(id6, list.Data[0].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	checkLoadEntryComments(t, userIDs[2], e3, 10, "", "", 0)
	checkLoadEntryComments(t, userIDs[3], e3, 10, "", "", 0)

	e4 := postEntry(userIDs[1], models.EntryPrivacyAll, true).ID
	id7 := postComment(userIDs[0], e4)
	id8 := postComment(userIDs[2], e4)

	list = checkLoadEntryComments(t, userIDs[1], e4, 10, "", "", 2)
	req.Equal(id7, list.Data[0].ID)
	req.Equal(id8, list.Data[1].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadEntryComments(t, userIDs[0], e4, 10, "", "", 1)
	req.Equal(id7, list.Data[0].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadEntryComments(t, userIDs[2], e4, 10, "", "", 1)
	req.Equal(id8, list.Data[0].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	e5 := postEntry(userIDs[2], models.EntryPrivacyAll, true).ID
	id9 := postComment(userIDs[0], e5)
	idA := postComment(userIDs[2], e5)

	list = checkLoadEntryComments(t, userIDs[0], e5, 10, "", "", 2)
	req.Equal(id9, list.Data[0].ID)
	req.Equal(idA, list.Data[1].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	list = checkLoadEntryComments(t, userIDs[2], e5, 10, "", "", 2)
	req.Equal(id9, list.Data[0].ID)
	req.Equal(idA, list.Data[1].ID)
	req.False(list.HasAfter)
	req.False(list.HasBefore)

	checkUnfollow(t, userIDs[1], userIDs[0])
	checkUnfollow(t, userIDs[0], userIDs[2])

	banShadow(db, userIDs[0])
	banShadow(db, userIDs[1])

	list = checkLoadEntryComments(t, userIDs[0], e1, 10, "", "", 2)
	list = checkLoadEntryComments(t, userIDs[1], e1, 10, "", "", 2)
	list = checkLoadEntryComments(t, userIDs[2], e1, 10, "", "", 1)
	req.Equal(id1, list.Data[0].ID)

	e6 := postEntry(userIDs[2], models.EntryPrivacyAll, true).ID
	postComment(userIDs[0], e6)
	checkLoadEntryComments(t, userIDs[0], e6, 10, "", "", 1)
	checkLoadEntryComments(t, userIDs[1], e6, 10, "", "", 1)
	checkLoadEntryComments(t, userIDs[2], e6, 10, "", "", 0)

	e7 := postEntry(userIDs[0], models.EntryPrivacyAll, true).ID
	postComment(userIDs[0], e7)
	checkLoadEntryComments(t, userIDs[0], e7, 10, "", "", 1)
	checkLoadEntryComments(t, userIDs[1], e7, 10, "", "", 1)
	checkLoadEntryComments(t, userIDs[2], e7, 10, "", "", 1)

	removeUserRestrictions(db, userIDs)

	checkDeleteEntry(t, e1, userIDs[0], true)
	checkDeleteEntry(t, e2, userIDs[0], true)
	checkDeleteEntry(t, e3, userIDs[0], true)
	checkDeleteEntry(t, e4, userIDs[1], true)
	checkDeleteEntry(t, e5, userIDs[2], true)
	checkDeleteEntry(t, e6, userIDs[2], true)
	checkDeleteEntry(t, e7, userIDs[0], true)

	deleteTheme(t, theme)
}
