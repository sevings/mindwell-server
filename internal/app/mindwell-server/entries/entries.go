package entries

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/golang-commonmark/markdown"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-server/restapi/operations/entries"

	"github.com/sevings/mindwell-server/internal/app/mindwell-server/comments"
	"github.com/sevings/mindwell-server/internal/app/mindwell-server/users"
	"github.com/sevings/mindwell-server/internal/app/mindwell-server/watchings"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations"
	"github.com/sevings/mindwell-server/utils"
)

// ConfigureAPI creates operations handlers
func ConfigureAPI(db *sql.DB, api *operations.MindwellAPI) {
	api.EntriesPostEntriesUsersMeHandler = entries.PostEntriesUsersMeHandlerFunc(newMyTlogPoster(db))

	api.EntriesGetEntriesIDHandler = entries.GetEntriesIDHandlerFunc(newEntryLoader(db))
	api.EntriesPutEntriesIDHandler = entries.PutEntriesIDHandlerFunc(newEntryEditor(db))
	api.EntriesDeleteEntriesIDHandler = entries.DeleteEntriesIDHandlerFunc(newEntryDeleter(db))

	api.EntriesGetEntriesLiveHandler = entries.GetEntriesLiveHandlerFunc(newLiveLoader(db))
	api.EntriesGetEntriesAnonymousHandler = entries.GetEntriesAnonymousHandlerFunc(newAnonymousLoader(db))
	api.EntriesGetEntriesBestHandler = entries.GetEntriesBestHandlerFunc(newBestLoader(db))
	api.EntriesGetEntriesUsersIDHandler = entries.GetEntriesUsersIDHandlerFunc(newTlogLoader(db))
	api.EntriesGetEntriesUsersMeHandler = entries.GetEntriesUsersMeHandlerFunc(newMyTlogLoader(db))
	api.EntriesGetEntriesFriendsHandler = entries.GetEntriesFriendsHandlerFunc(newFriendsFeedLoader(db))
}

var wordRe *regexp.Regexp
var md *markdown.Markdown

func init() {
	wordRe = regexp.MustCompile("[a-zA-Zа-яА-ЯёЁ0-9]+")
	md = markdown.New(markdown.Typographer(false), markdown.Breaks(true), markdown.Tables(false))
}

func wordCount(content, title string) int64 {
	var wc int64
	contentWords := wordRe.FindAllStringIndex(content, -1)
	wc += int64(len(contentWords))

	titleWords := wordRe.FindAllStringIndex(title, -1)
	wc += int64(len(titleWords))

	return wc
}

func cutText(text, format string, size int) (string, bool) {
	if len(text) <= size {
		return text, false
	}

	text = fmt.Sprintf(format, text)

	var isSpace bool
	trim := func(r rune) bool {
		if isSpace {
			return unicode.IsSpace(r)
		}

		isSpace = unicode.IsSpace(r)
		return true
	}
	text = strings.TrimRightFunc(text, trim)

	text += "&hellip;"

	return text, true
}

func cutEntry(title, content string) (cutTitle string, cutContent string, hasCut bool) {
	const titleLength = 60
	const titleFormat = "%.60s"
	cutTitle, isTitleCut := cutText(title, titleFormat, titleLength)

	const contentLength = 500
	const contentFormat = "%.500s"
	cutContent, isContentCut := cutText(content, contentFormat, contentLength)

	hasCut = isTitleCut || isContentCut
	if !hasCut {
		cutTitle = ""
		cutContent = ""
	}

	return
}

func entryCategory(entry *models.Entry) string {
	if entry.WordCount > 100 {
		return "longread"
	}

	// media
	return "tweet"
}

func createEntry(tx *utils.AutoTx, userID int64, title, content, privacy string, isVotable bool) *models.Entry {
	if privacy == "followers" {
		privacy = models.EntryPrivacySome //! \todo add users to list
	}

	cutTitle, cutContent, hasCut := cutEntry(title, content)

	entry := models.Entry{
		Title:       title,
		CutTitle:    cutTitle,
		Content:     md.RenderToString([]byte(content)),
		CutContent:  md.RenderToString([]byte(cutContent)),
		EditContent: content,
		HasCut:      hasCut,
		WordCount:   wordCount(content, title),
		Privacy:     privacy,
		Vote:        models.EntryVoteBan,
		IsWatching:  true,
	}

	category := entryCategory(&entry)

	const q = `
	INSERT INTO entries (author_id, title, cut_title, content, cut_content, edit_content, 
		has_cut, word_count, visible_for, is_votable, category)
	VALUES ($1, $2, $3, $4, $5, $6,$7, $8,
		(SELECT id FROM entry_privacy WHERE type = $9), 
		$10, (SELECT id from categories WHERE type = $11))
	RETURNING id, extract(epoch from created_at)`

	tx.Query(q, userID, title, cutTitle, entry.Content, entry.CutContent, entry.EditContent,
		hasCut, entry.WordCount, privacy, isVotable, category).Scan(&entry.ID, &entry.CreatedAt)

	watchings.AddWatching(tx, userID, entry.ID)
	author := users.LoadUserByID(tx, userID)
	entry.Author = author

	return &entry
}

func newMyTlogPoster(db *sql.DB) func(entries.PostEntriesUsersMeParams, *models.UserID) middleware.Responder {
	return func(params entries.PostEntriesUsersMeParams, uID *models.UserID) middleware.Responder {
		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			entry := createEntry(tx, int64(*uID),
				*params.Title, params.Content, params.Privacy, *params.IsVotable)

			if tx.Error() != nil {
				return entries.NewPostEntriesUsersMeForbidden()
			}

			return entries.NewPostEntriesUsersMeCreated().WithPayload(entry)
		})
	}
}

func editEntry(tx *utils.AutoTx, entryID, userID int64, title, content, privacy string, isVotable bool) *models.Entry {
	if privacy == "followers" {
		privacy = models.EntryPrivacySome //! \todo add users to list
	}

	cutTitle, cutContent, hasCut := cutEntry(title, content)

	entry := models.Entry{
		ID:          entryID,
		Title:       title,
		CutTitle:    cutTitle,
		Content:     md.RenderToString([]byte(content)),
		CutContent:  md.RenderToString([]byte(cutContent)),
		EditContent: content,
		HasCut:      hasCut,
		WordCount:   wordCount(content, title),
		Privacy:     privacy,
		Vote:        models.EntryVoteBan,
		IsWatching:  true,
		IsVotable:   isVotable,
	}

	category := entryCategory(&entry)

	const q = `
	UPDATE entries
	SET title = $1, cut_title = $2, content = $3, cut_content = $4, edit_content = $5, has_cut = $6, 
	word_count = $7, 
	visible_for = (SELECT id FROM entry_privacy WHERE type = $8), 
	is_votable = $9,
	category = (SELECT id from categories WHERE type = $10)
	WHERE id = $11 AND author_id = $12
	RETURNING extract(epoch from created_at)`

	tx.Query(q, title, cutTitle, entry.Content, entry.CutContent, entry.EditContent, hasCut,
		entry.WordCount, privacy, isVotable, category, entryID, userID).Scan(&entry.CreatedAt)

	watchings.AddWatching(tx, userID, entry.ID)

	author := users.LoadUserByID(tx, userID)
	entry.Author = author

	return &entry
}

func newEntryEditor(db *sql.DB) func(entries.PutEntriesIDParams, *models.UserID) middleware.Responder {
	return func(params entries.PutEntriesIDParams, uID *models.UserID) middleware.Responder {
		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			entry := editEntry(tx, params.ID, int64(*uID),
				*params.Title, params.Content, params.Privacy, *params.IsVotable)

			if tx.Error() != nil {
				return entries.NewPutEntriesIDForbidden()
			}

			return entries.NewPutEntriesIDOK().WithPayload(entry)
		})
	}
}

func entryVoteStatus(authorID, userID int64, vote sql.NullFloat64) string {
	switch {
	case authorID == userID:
		return models.EntryVoteBan
	case !vote.Valid:
		return models.EntryVoteNot
	case vote.Float64 > 0:
		return models.EntryVotePos
	default:
		return models.EntryVoteNeg
	}
}

func loadEntry(tx *utils.AutoTx, entryID, userID int64) *models.Entry {
	const q = tlogFeedQueryStart + `
		WHERE entries.id = $2
			AND (entries.author_id = $1
				OR entry_privacy.type = 'all' 
				OR (entry_privacy.type = 'some' 
					AND EXISTS(SELECT 1 from entries_privacy WHERE user_id = $1 AND entry_id = entries.id)))
		`

	var entry models.Entry
	var author models.User
	var vote sql.NullFloat64
	var avatar string
	tx.Query(q, userID, entryID).Scan(&entry.ID, &entry.CreatedAt, &entry.Rating, &entry.Votes,
		&entry.Title, &entry.CutTitle, &entry.Content, &entry.CutContent, &entry.EditContent,
		&entry.HasCut, &entry.WordCount, &entry.Privacy,
		&entry.IsVotable, &entry.CommentCount,
		&author.ID, &author.Name, &author.ShowName,
		&author.IsOnline,
		&avatar,
		&vote, &entry.IsFavorited, &entry.IsWatching)

	if author.ID != userID {
		entry.EditContent = ""
	}

	entry.Vote = entryVoteStatus(author.ID, userID, vote)

	author.Avatar = utils.NewAvatar(avatar)
	entry.Author = &author

	cmt := comments.LoadEntryComments(tx, userID, entryID, 5, "", "")
	entry.Comments = cmt

	return &entry
}

func newEntryLoader(db *sql.DB) func(entries.GetEntriesIDParams, *models.UserID) middleware.Responder {
	return func(params entries.GetEntriesIDParams, uID *models.UserID) middleware.Responder {
		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			entry := loadEntry(tx, params.ID, int64(*uID))

			if entry.ID == 0 {
				return entries.NewGetEntriesIDNotFound()
			}

			return entries.NewGetEntriesIDOK().WithPayload(entry)
		})
	}
}

func deleteEntry(tx *utils.AutoTx, entryID, userID int64) bool {
	var authorID int64
	tx.Query("SELECT author_id FROM entries WHERE id = $1", entryID).Scan(&authorID)
	if authorID != userID {
		return false
	}

	tx.Exec("DELETE from entries WHERE id = $1", entryID)
	return true
}

func newEntryDeleter(db *sql.DB) func(entries.DeleteEntriesIDParams, *models.UserID) middleware.Responder {
	return func(params entries.DeleteEntriesIDParams, uID *models.UserID) middleware.Responder {
		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			ok := deleteEntry(tx, params.ID, int64(*uID))
			if ok {
				return entries.NewDeleteEntriesIDOK()
			}

			if tx.Error() == sql.ErrNoRows {
				return entries.NewDeleteEntriesIDNotFound()
			}

			return entries.NewDeleteEntriesIDForbidden()
		})
	}
}
