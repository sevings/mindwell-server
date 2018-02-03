package entries

import (
	"database/sql"
	"log"
	"regexp"

	"github.com/golang-commonmark/markdown"
	"github.com/microcosm-cc/bluemonday"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/yummy-server/restapi/operations/entries"

	"github.com/sevings/yummy-server/internal/app/yummy-server/users"
	"github.com/sevings/yummy-server/internal/app/yummy-server/utils"
	"github.com/sevings/yummy-server/internal/app/yummy-server/watchings"
	"github.com/sevings/yummy-server/models"
	"github.com/sevings/yummy-server/restapi/operations"
)

// ConfigureAPI creates operations handlers
func ConfigureAPI(db *sql.DB, api *operations.YummyAPI) {
	api.EntriesPostEntriesUsersMeHandler = entries.PostEntriesUsersMeHandlerFunc(newMyTlogPoster(db))
	api.EntriesGetEntriesLiveHandler = entries.GetEntriesLiveHandlerFunc(newLiveLoader(db))
	api.EntriesGetEntriesAnonymousHandler = entries.GetEntriesAnonymousHandlerFunc(newAnonymousLoader(db))
	api.EntriesGetEntriesBestHandler = entries.GetEntriesBestHandlerFunc(newBestLoader(db))
	api.EntriesGetEntriesUsersIDHandler = entries.GetEntriesUsersIDHandlerFunc(newTlogLoader(db))
	api.EntriesGetEntriesUsersMeHandler = entries.GetEntriesUsersMeHandlerFunc(newMyTlogLoader(db))
}

var wordRe *regexp.Regexp
var htmlPolicy *bluemonday.Policy
var md *markdown.Markdown

func init() {
	wordRe = regexp.MustCompile("[a-zA-Zа-яА-ЯёЁ0-9]+")
	htmlPolicy = bluemonday.UGCPolicy()
	md = markdown.New(markdown.Typographer(false), markdown.Breaks(true), markdown.Nofollow(true))
}

func createEntry(tx utils.AutoTx, userID int64, title, content, privacy string, isVotable bool) (*models.Entry, bool) {
	author, _ := users.LoadUserByID(tx, userID)

	var wordCount int64
	contentWords := wordRe.FindAllStringIndex(content, -1)
	wordCount += int64(len(contentWords))

	titleWords := wordRe.FindAllStringIndex(title, -1)
	wordCount += int64(len(titleWords))

	if privacy == "followers" {
		privacy = models.EntryPrivacySome //! \todo add users to list
	}

	content = md.RenderToString([]byte(content))
	content = htmlPolicy.Sanitize(content)

	entry := models.Entry{
		Title:      title,
		Content:    content,
		WordCount:  wordCount,
		Privacy:    privacy,
		Author:     author,
		Vote:       models.EntryVoteMy,
		IsWatching: true,
	}

	const q = `
	INSERT INTO entries (author_id, title, content, word_count, visible_for, is_votable)
	VALUES ($1, $2, $3, $4, (SELECT id FROM entry_privacy WHERE type = $5), $6)
	RETURNING id, created_at`

	err := tx.QueryRow(q, author.ID, title, content, wordCount,
		privacy, isVotable).Scan(&entry.ID, &entry.CreatedAt)
	if err != nil {
		log.Print(err)
		return nil, false
	}

	err = watchings.AddWatching(tx, userID, entry.ID)
	if err != nil {
		return nil, false
	}

	return &entry, true
}

func newMyTlogPoster(db *sql.DB) func(entries.PostEntriesUsersMeParams, *models.UserID) middleware.Responder {
	return func(params entries.PostEntriesUsersMeParams, uID *models.UserID) middleware.Responder {
		return utils.Transact(db, func(tx utils.AutoTx) (middleware.Responder, bool) {
			entry, created := createEntry(tx, int64(*uID),
				*params.Title, params.Content, *params.Privacy, *params.IsVotable)

			if !created {
				return entries.NewPostEntriesUsersMeForbidden(), false
			}

			return entries.NewPostEntriesUsersMeOK().WithPayload(entry), true
		})
	}
}
