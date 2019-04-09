package test

import (
	"testing"
	"time"

	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/comments"
	"github.com/stretchr/testify/require"
)

func checkComment(t *testing.T, cmt *models.Comment, entryID int64, userID *models.UserID, author *models.AuthProfile, content string) {
	req := require.New(t)

	req.Equal(entryID, cmt.EntryID)
	req.Equal("<p>"+content+"</p>", cmt.Content)
	req.Equal(content, cmt.EditContent)

	req.Equal(author.ID, cmt.Author.ID)
	req.Equal(author.Name, cmt.Author.Name)
	req.Equal(author.ShowName, cmt.Author.ShowName)
	req.Equal(author.IsOnline, cmt.Author.IsOnline)
	req.Equal(author.Avatar, cmt.Author.Avatar)

	req.Equal(cmt.ID, cmt.Rating.ID)
	req.True(cmt.Rating.IsVotable)
	req.Zero(cmt.Rating.Rating)
	req.Zero(cmt.Rating.UpCount)
	req.Zero(cmt.Rating.DownCount)

	if userID.ID == author.ID {
		req.Equal(models.RatingVoteBan, cmt.Rating.Vote)
	} else {
		req.Equal(models.RatingVoteNot, cmt.Rating.Vote)
	}

	req.Equal(userID.ID == author.ID, cmt.Rights.Edit)
	req.Equal(userID.ID == author.ID, cmt.Rights.Delete)
	req.Equal(userID.ID != author.ID && !userID.Ban.Vote, cmt.Rights.Vote)
}

func checkLoadComment(t *testing.T, commentID int64, userID *models.UserID, success bool,
	author *models.AuthProfile, entryID int64, content string) {

	load := api.CommentsGetCommentsIDHandler.Handle
	resp := load(comments.GetCommentsIDParams{ID: commentID}, userID)
	body, ok := resp.(*comments.GetCommentsIDOK)
	require.Equal(t, success, ok)
	if !success {
		return
	}

	cmt := body.Payload
	checkComment(t, cmt, entryID, userID, author, content)
}

func checkPostComment(t *testing.T,
	entryID int64, content string, success bool,
	author *models.AuthProfile, id *models.UserID) int64 {

	params := comments.PostEntriesIDCommentsParams{
		ID:      entryID,
		Content: content,
	}

	post := api.CommentsPostEntriesIDCommentsHandler.Handle
	resp := post(params, id)
	body, ok := resp.(*comments.PostEntriesIDCommentsCreated)
	require.Equal(t, success, ok)
	if !success {
		return 0
	}

	cmt := body.Payload
	checkComment(t, cmt, params.ID, id, author, params.Content)

	checkLoadComment(t, cmt.ID, id, true, author, params.ID, params.Content)

	return cmt.ID
}

func checkEditComment(t *testing.T,
	commentID int64, content string, entryID int64, success bool,
	author *models.AuthProfile, id *models.UserID) {

	params := comments.PutCommentsIDParams{
		ID:      commentID,
		Content: content,
	}

	edit := api.CommentsPutCommentsIDHandler.Handle
	resp := edit(params, id)
	body, ok := resp.(*comments.PutCommentsIDOK)
	require.Equal(t, success, ok)
	if !success {
		return
	}

	cmt := body.Payload
	checkComment(t, cmt, entryID, id, author, content)

	checkLoadComment(t, commentID, id, true, author, entryID, content)
}

func checkDeleteComment(t *testing.T, commentID int64, userID *models.UserID, success bool) {
	del := api.CommentsDeleteCommentsIDHandler.Handle
	resp := del(comments.DeleteCommentsIDParams{ID: commentID}, userID)
	_, ok := resp.(*comments.DeleteCommentsIDOK)
	require.Equal(t, success, ok)
}

func TestOpenComments(t *testing.T) {
	postEntry(userIDs[0], models.EntryPrivacyAll, true)
	feed := checkLoadTlog(t, userIDs[0], userIDs[0], 10, "", "", 1)
	entry := feed.Entries[0]

	var id int64

	id = checkPostComment(t, entry.ID, "blabla", true, profiles[0], userIDs[0])
	checkEditComment(t, id, "edited comment", entry.ID, true, profiles[0], userIDs[0])
	checkEntryWatching(t, userIDs[0], entry.ID, true, true)

	id = checkPostComment(t, entry.ID, "blabla", true, profiles[1], userIDs[1])
	checkEditComment(t, id, "edited comment", entry.ID, true, profiles[1], userIDs[1])
	checkEntryWatching(t, userIDs[1], entry.ID, true, true)

	banComment(db, userIDs[0])
	checkPostComment(t, entry.ID, "blabla", true, profiles[0], userIDs[0])
	removeUserRestrictions(db, userIDs)

	banComment(db, userIDs[1])
	checkPostComment(t, entry.ID, "blabla", false, profiles[1], userIDs[1])
	removeUserRestrictions(db, userIDs)

	checkDeleteEntry(t, entry.ID, userIDs[0], true)
}

func TestPrivateComments(t *testing.T) {
	postEntry(userIDs[0], models.EntryPrivacyMe, true)
	feed := checkLoadTlog(t, userIDs[0], userIDs[0], 10, "", "", 1)
	entry := feed.Entries[0]

	var id int64

	id = checkPostComment(t, entry.ID, "blabla", true, profiles[0], userIDs[0])
	checkEditComment(t, id, "edited comment", entry.ID, true, profiles[0], userIDs[0])
	checkEntryWatching(t, userIDs[0], entry.ID, true, true)

	checkEditComment(t, id, "edited comment", entry.ID, false, profiles[1], userIDs[1])
	id = checkPostComment(t, entry.ID, "blabla", false, profiles[1], userIDs[1])
	checkEntryWatching(t, userIDs[1], entry.ID, false, false)

	checkDeleteEntry(t, entry.ID, userIDs[0], true)
}

func postComment(id *models.UserID, entryID int64) int64 {
	params := comments.PostEntriesIDCommentsParams{
		ID:      entryID,
		Content: "test comment",
	}

	post := api.CommentsPostEntriesIDCommentsHandler.Handle
	resp := post(params, id)
	body := resp.(*comments.PostEntriesIDCommentsCreated)
	cmt := body.Payload

	time.Sleep(10 * time.Millisecond)

	return cmt.ID
}

func TestCommentHTML(t *testing.T) {
	req := require.New(t)
	entry := postEntry(userIDs[0], models.EntryPrivacyAll, false)

	post := func(content string) *models.Comment {
		params := comments.PostEntriesIDCommentsParams{
			ID:      entry.ID,
			Content: content,
		}

		post := api.CommentsPostEntriesIDCommentsHandler.Handle
		resp := post(params, userIDs[0])
		body, ok := resp.(*comments.PostEntriesIDCommentsCreated)
		req.True(ok)

		return body.Payload
	}

	content := "http://ex.com/im.jpg"
	cmt := post(content)

	req.Equal(content, cmt.EditContent)
	req.Equal("<p><img src=\""+content+"\"></p>", cmt.Content)

	edit := func(content string) *models.Comment {
		params := comments.PutCommentsIDParams{
			ID:      cmt.ID,
			Content: content,
		}

		edit := api.CommentsPutCommentsIDHandler.Handle
		resp := edit(params, userIDs[0])
		body, ok := resp.(*comments.PutCommentsIDOK)
		req.True(ok)
		return body.Payload
	}

	checkImage := func(content string) {
		cmt = edit(content)
		req.Equal(content, cmt.EditContent)
		req.Equal("<p><img src=\""+content+"\"></p>", cmt.Content)
	}

	checkImage("http://ex.com/im.jpg?trolo")
	checkImage("hTTps://ex.com/im.GIf?oooo#aaa")

	checkURL := func(content string) {
		cmt = edit(content)
		req.Equal(content, cmt.EditContent)
		req.Equal("<p><a href=\""+content+"\" target=\"_blank\" rel=\"noopener nofollow\">"+content+"</a></p>", cmt.Content)
	}

	checkURL("http://ex.com/im.ajpg")
	checkURL("tg://resolve?domain=telegram")
	checkURL("https://ex.com/im?oooo#aaa")

	checkText := func(content string) {
		cmt = edit(content)
		req.Equal(content, cmt.EditContent)
		req.Equal("<p>"+content+"</p>", cmt.Content)
	}

	checkText("http://")
	checkText("://a")
	checkText("aa:// a")

	content = "<>&\n\"'\t"
	cmt = edit(content)
	req.Equal("<p>&lt;&gt;&amp;<br>&#34;&#39;</p>", cmt.Content)
}
