package entries

import (
	"database/sql"
	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/internal/app/mindwell-server/comments"
	"github.com/sevings/mindwell-server/restapi/operations/themes"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/golang-commonmark/markdown"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-server/restapi/operations/entries"
	"github.com/sevings/mindwell-server/restapi/operations/me"
	usersAPI "github.com/sevings/mindwell-server/restapi/operations/users"

	"github.com/sevings/mindwell-server/internal/app/mindwell-server/users"
	"github.com/sevings/mindwell-server/internal/app/mindwell-server/watchings"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/utils"
)

// ConfigureAPI creates operations handlers
func ConfigureAPI(srv *utils.MindwellServer) {
	srv.API.MePostMeTlogHandler = me.PostMeTlogHandlerFunc(newMyTlogPoster(srv))
	srv.API.ThemesPostThemesNameTlogHandler = themes.PostThemesNameTlogHandlerFunc(newThemePoster(srv))

	srv.API.MeGetMeTlogHandler = me.GetMeTlogHandlerFunc(newMyTlogLoader(srv))
	srv.API.UsersGetUsersNameTlogHandler = usersAPI.GetUsersNameTlogHandlerFunc(newTlogLoader(srv))
	srv.API.ThemesGetThemesNameTlogHandler = themes.GetThemesNameTlogHandlerFunc(newThemeLoader(srv))

	srv.API.MeGetMeFavoritesHandler = me.GetMeFavoritesHandlerFunc(newMyFavoritesLoader(srv))
	srv.API.UsersGetUsersNameFavoritesHandler = usersAPI.GetUsersNameFavoritesHandlerFunc(newTlogFavoritesLoader(srv))

	srv.API.EntriesGetEntriesIDHandler = entries.GetEntriesIDHandlerFunc(newEntryLoader(srv))
	srv.API.EntriesPutEntriesIDHandler = entries.PutEntriesIDHandlerFunc(newEntryEditor(srv))
	srv.API.EntriesDeleteEntriesIDHandler = entries.DeleteEntriesIDHandlerFunc(newEntryDeleter(srv))
	srv.API.EntriesGetEntriesRandomHandler = entries.GetEntriesRandomHandlerFunc(newRandomEntryLoader(srv))
	srv.API.EntriesGetEntriesIDAdjacentHandler = entries.GetEntriesIDAdjacentHandlerFunc(newAdjacentLoader(srv))

	srv.API.EntriesGetEntriesLiveHandler = entries.GetEntriesLiveHandlerFunc(newLiveLoader(srv))
	srv.API.EntriesGetEntriesBestHandler = entries.GetEntriesBestHandlerFunc(newBestLoader(srv))
	srv.API.EntriesGetEntriesFriendsHandler = entries.GetEntriesFriendsHandlerFunc(newFriendsFeedLoader(srv))
	srv.API.EntriesGetEntriesWatchingHandler = entries.GetEntriesWatchingHandlerFunc(newWatchingLoader(srv))

	srv.API.MeGetMeCalendarHandler = me.GetMeCalendarHandlerFunc(newMyCalendarLoader(srv))
	srv.API.UsersGetUsersNameCalendarHandler = usersAPI.GetUsersNameCalendarHandlerFunc(newTlogCalendarLoader(srv))
	srv.API.ThemesGetThemesNameCalendarHandler = themes.GetThemesNameCalendarHandlerFunc(newThemeCalendarLoader(srv))
}

var wordRe *regexp.Regexp
var imgRe *regexp.Regexp
var ytRe *regexp.Regexp
var md *markdown.Markdown

func init() {
	wordRe = regexp.MustCompile("[a-zA-Zа-яА-ЯёЁ0-9]+")
	imgRe = regexp.MustCompile("!\\[[^\\]]*\\]\\([^\\)]+\\)")
	ytRe = regexp.MustCompile(`(?i)(?:https?://)?(?:www\.)?(?:m\.)?(?:youtube.com/watch\?\S*v=|youtu.be/)[a-z0-9\-_]+\S*`)

	markdown.RegisterCoreRule(250, appendTargetToLinks)
	md = markdown.New(markdown.Typographer(false), markdown.Breaks(true), markdown.Tables(false))
}

func appendTargetToLinks(s *markdown.StateCore) {
	for _, token := range s.Tokens {
		inline, ok := token.(*markdown.Inline)
		if !ok {
			continue
		}

		for _, tok := range inline.Children {
			link, ok := tok.(*markdown.LinkOpen)
			if !ok {
				continue
			}

			link.Target = "_blank"
		}
	}
}

func wordCount(content, title string) int64 {
	var wc int64

	content = imgRe.ReplaceAllLiteralString(content, " ")
	contentWords := wordRe.FindAllStringIndex(content, -1)
	wc += int64(len(contentWords))

	titleWords := wordRe.FindAllStringIndex(title, -1)
	wc += int64(len(titleWords))

	return wc
}

func entryCategory(entry *models.Entry) string {
	if entry.WordCount > 100 {
		return "longread"
	}

	if entry.WordCount < 50 {
		hasMedia := len(entry.Images) > 0
		if !hasMedia {
			hasMedia = len(imgRe.FindAllStringIndex(entry.EditContent, -1)) > 0
		}
		if !hasMedia {
			hasMedia = len(ytRe.FindAllStringIndex(entry.EditContent, -1)) > 0
		}
		if hasMedia {
			return "media"
		}
	}

	return "tweet"
}

func setEntryTexts(entry *models.Entry, hasAttach bool) {
	title := strings.TrimSpace(entry.Title)
	title = bluemonday.StrictPolicy().Sanitize(title)

	cutTitle, isTitleCut := utils.CutText(title, 100)

	lineCount := 15
	if hasAttach {
		lineCount = 5
	}

	editContent := strings.TrimSpace(entry.EditContent)
	content := md.RenderToString([]byte(editContent))
	cutContent, isContentCut := utils.CutHtml(content, lineCount, 44)

	hasCut := isTitleCut || isContentCut
	if !hasCut {
		cutTitle = ""
		cutContent = ""
	}

	entry.Title = title
	entry.CutTitle = cutTitle
	entry.Content = content
	entry.CutContent = cutContent
	entry.EditContent = editContent
	entry.HasCut = hasCut
}

func initMyEntry(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID,
	entry *models.Entry, hasAttach bool) *models.Error {
	isTheme := entry.Author.IsTheme

	if entry.Privacy == models.EntryPrivacyMe {
		if isTheme {
			return srv.NewError(&i18n.Message{ID: "post_me_to_theme", Other: "Private entries in themes are not allowed."})
		}

		entry.Rating.IsVotable = false
		entry.InLive = false
	} else if entry.Privacy == models.EntryPrivacyFollowers {
		entry.InLive = false
	}

	if !isTheme {
		entry.IsAnonymous = false
	}

	entry.IsWatching = true
	entry.WordCount = wordCount(entry.EditContent, entry.Title)

	if entry.Author.ID > 0 {
		entry.Author = users.LoadUserByID(srv, tx, entry.Author.ID)
	} else {
		entry.Author = users.LoadUserByName(srv, tx, entry.Author.Name)
	}

	if entry.Author.ID == 0 {
		return srv.NewError(nil)
	}

	if isTheme && !entry.IsAnonymous {
		entry.User = users.LoadUserByID(srv, tx, userID.ID)
	}

	setEntryTexts(entry, hasAttach)
	setEntryRights(entry, userID, 0)

	return nil
}

func createEntry(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, entry *models.Entry, hasAttach bool) *models.Error {
	if err := initMyEntry(srv, tx, userID, entry, hasAttach); err != nil {
		return err
	}

	category := entryCategory(entry)

	const q = `
	INSERT INTO entries (user_id, author_id, title, edit_content, word_count,
		visible_for,
		is_commentable, is_votable, is_anonymous, in_live, is_shared,
		category)
	VALUES ($1, $2, $3, $4, $5,
		(SELECT id FROM entry_privacy WHERE type = $6), 
		$7, $8, $9, $10, $11,
		(SELECT id from categories WHERE type = $12))
	RETURNING id, extract(epoch from created_at)`

	tx.Query(q, userID.ID, entry.Author.ID, entry.Title, entry.EditContent, entry.WordCount,
		entry.Privacy,
		entry.IsCommentable, entry.Rating.IsVotable, entry.IsAnonymous, entry.InLive, entry.IsShared,
		category).
		Scan(&entry.ID, &entry.CreatedAt)

	entry.Rating.ID = entry.ID
	watchings.AddWatching(tx, userID.ID, entry.ID)

	return nil
}

func loadEntryImages(srv *utils.MindwellServer, tx *utils.AutoTx, entry *models.Entry, images []int64) {
	for _, imageID := range images {
		img := utils.LoadImage(srv, tx, imageID)
		if img != nil {
			entry.Images = append(entry.Images, img)
		}
	}
}

func loadEntryTags(tx *utils.AutoTx, entry *models.Entry) {
	const q = `SELECT tag FROM entry_tags INNER JOIN tags ON tag_id = tags.id WHERE entry_id = $1 ORDER BY tag`
	entry.Tags = tx.QueryStrings(q, entry.ID)
}

func attachImages(srv *utils.MindwellServer, tx *utils.AutoTx, entry *models.Entry, images []int64) {
	if len(images) == 0 {
		return
	}

	const q = "INSERT INTO entry_images(entry_id, image_id)	VALUES($1, $2)"
	for _, imageID := range images {
		tx.Exec(q, entry.ID, imageID)
	}

	loadEntryImages(srv, tx, entry, images)
}

func setTags(tx *utils.AutoTx, entry *models.Entry) {
	realTags := make([]string, 0, len(entry.Tags))
tagLoop:
	for _, tag := range entry.Tags {
		tag = strings.TrimSpace(tag)
		tag = strings.ToLower(tag)
		if tag == "" {
			continue
		}

		for _, realTag := range realTags {
			if tag == realTag {
				continue tagLoop
			}
		}

		realTags = append(realTags, tag)
	}

	if len(realTags) == 0 {
		entry.Tags = nil
		return
	}

	for _, tag := range realTags {
		tagID := tx.QueryInt64("SELECT id FROM tags WHERE tag = $1", tag)
		if tagID == 0 {
			tagID = tx.QueryInt64("INSERT INTO tags(tag) VALUES($1) RETURNING id", tag)
		}
		tx.Exec("INSERT INTO entry_tags(entry_id, tag_id) VALUES($1, $2)", entry.ID, tagID)
	}

	entry.Tags = realTags
}

func allowedInLive(followersCount, entryCount int64) bool {
	switch {
	case followersCount < 3:
		return entryCount < 1
	case followersCount < 10:
		return entryCount < 2
	case followersCount < 50:
		return entryCount < 3
	default:
		return true
	}
}

func allowedWithoutVoting(srv *utils.MindwellServer, userID *models.UserID) *models.Error {
	if userID.IsInvited {
		return nil
	}

	return srv.NewError(&i18n.Message{ID: "post_wo_voting", Other: "You're not allowed to post without voting."})
}

func canPostInLive(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID) *models.Error {
	if userID.Ban.Live {
		return srv.NewError(&i18n.Message{ID: "post_in_live", Other: "You're not allowed to post in live."})
	}

	if userID.NegKarma {
		return srv.NewError(&i18n.Message{ID: "post_in_live_karma", Other: "You're not allowed to post in live."})
	}

	var entryCount int64
	const countQ = `SELECT count(*) FROM entries WHERE user_id = $1 
		AND date_trunc('day', created_at) = CURRENT_DATE AND in_live
	`
	tx.Query(countQ, userID.ID).Scan(&entryCount)

	if !allowedInLive(userID.FollowersCount, entryCount) {
		return srv.NewError(&i18n.Message{ID: "post_in_live_followers", Other: "You can't post in live anymore today."})
	}

	return nil
}

func canPostInTheme(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, theme string) *models.Error {
	if userID.Ban.Live {
		return srv.NewError(&i18n.Message{ID: "post_in_themes", Other: "You're not allowed to post in themes."})
	}

	if userID.NegKarma {
		return srv.NewError(&i18n.Message{ID: "post_in_themes_karma", Other: "You're not allowed to post in themes."})
	}

	postThemeErr := srv.NewError(&i18n.Message{ID: "post_in_theme", Other: "You're not allowed to post in this theme."})

	if !userID.IsInvited {
		return postThemeErr
	}

	const themeQ = `SELECT id, creator_id FROM users WHERE lower(name) = lower($1) AND creator_id IS NOT NULL`
	var themeID, creatorID int64
	tx.Query(themeQ, theme).Scan(&themeID, &creatorID)

	if themeID == 0 {
		return srv.StandardError("no_theme")
	}

	if creatorID != userID.ID {
		relationTo := utils.LoadRelation(tx, userID.ID, themeID)
		if relationTo != models.RelationshipRelationFollowed {
			return postThemeErr
		}

		relationFrom := utils.LoadRelation(tx, themeID, userID.ID)
		if relationFrom == models.RelationshipRelationIgnored {
			return postThemeErr
		}
	}

	return nil
}

func checkImagesOwner(srv *utils.MindwellServer, tx *utils.AutoTx, userID int64, images []int64) *models.Error {
	for _, imageID := range images {
		authorID := tx.QueryInt64("SELECT user_id FROM images WHERE id = $1", imageID)

		if authorID == 0 {
			return srv.NewError(&i18n.Message{ID: "attached_image_not_found", Other: "Attached image not found."})
		}

		if authorID != userID {
			return srv.NewError(&i18n.Message{ID: "attach_not_your_image", Other: "You can attach only your own images."})
		}
	}

	return nil
}

func postNewEntry(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, entry *models.Entry, images []int64) *models.Error {
	prev, found, same := checkPrev(userID, entry)
	if found {
		if same {
			*entry = *prev
			return nil
		}

		return srv.NewError(&i18n.Message{ID: "post_same_entry", Other: "You are trying to post the same entry again."})
	}

	if entry.InLive &&
		(entry.Privacy == models.EntryPrivacyAll ||
			entry.Privacy == models.EntryPrivacyRegistered ||
			entry.Privacy == models.EntryPrivacyInvited) {
		err := canPostInLive(srv, tx, userID)
		if err != nil {
			return err
		}
	}

	if entry.Author.IsTheme && entry.Privacy != models.EntryPrivacyMe {
		err := canPostInTheme(srv, tx, userID, entry.Author.Name)
		if err != nil {
			return err
		}
	}

	if !entry.Rating.IsVotable {
		err := allowedWithoutVoting(srv, userID)
		if err != nil {
			return err
		}
	}

	if entry.IsShared {
		if !userID.Verified {
			return srv.NewError(&i18n.Message{ID: "post_shared", Other: "You are now allowed to create shared posts."})
		}
	}

	imageErr := checkImagesOwner(srv, tx, userID.ID, images)
	if imageErr != nil {
		return imageErr
	}

	if err := createEntry(srv, tx, userID, entry, len(images) > 0); err != nil {
		return err
	}

	attachImages(srv, tx, entry, images)
	setTags(tx, entry)

	if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
		return srv.NewError(nil)
	}

	setPrev(entry, userID)

	return nil
}

func newMyTlogPoster(srv *utils.MindwellServer) func(me.PostMeTlogParams, *models.UserID) middleware.Responder {
	return func(params me.PostMeTlogParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			entry := &models.Entry{
				Author: &models.User{
					ID: userID.ID,
				},
				Privacy:       params.Privacy,
				IsCommentable: *params.IsCommentable,
				InLive:        *params.InLive,
				IsShared:      *params.IsShared,
				IsAnonymous:   false,
				Rating: &models.Rating{
					IsVotable: *params.IsVotable,
				},
				Title:       *params.Title,
				EditContent: params.Content,
				Tags:        params.Tags,
			}

			err := postNewEntry(srv, tx, userID, entry, params.Images)
			if err != nil {
				return me.NewPostMeTlogForbidden().WithPayload(err)
			}

			return me.NewPostMeTlogCreated().WithPayload(entry)
		})
	}
}

func newThemePoster(srv *utils.MindwellServer) func(themes.PostThemesNameTlogParams, *models.UserID) middleware.Responder {
	return func(params themes.PostThemesNameTlogParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			entry := &models.Entry{
				Author: &models.User{
					Name:    params.Name,
					IsTheme: true,
				},
				Privacy:       params.Privacy,
				IsCommentable: *params.IsCommentable,
				InLive:        *params.InLive,
				IsShared:      *params.IsShared,
				IsAnonymous:   *params.IsAnonymous,
				Rating: &models.Rating{
					IsVotable: *params.IsVotable,
				},
				Title:       *params.Title,
				EditContent: params.Content,
				Tags:        params.Tags,
			}

			err := postNewEntry(srv, tx, userID, entry, params.Images)
			if err != nil {
				return themes.NewPostThemesNameTlogForbidden().WithPayload(err)
			}

			return themes.NewPostThemesNameTlogCreated().WithPayload(entry)
		})
	}
}

func editEntry(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, entry *models.Entry,
	hasAttach bool) *models.Error {

	if err := initMyEntry(srv, tx, userID, entry, hasAttach); err != nil {
		return err
	}

	category := entryCategory(entry)

	const q = `
	UPDATE entries
	SET title = $1, edit_content = $2,
	word_count = $3, 
	visible_for = (SELECT id FROM entry_privacy WHERE type = $4), 
	is_commentable = $5, is_votable = $6, in_live = $7, is_shared = $8,
	category = (SELECT id from categories WHERE type = $9)
	WHERE id = $10 AND user_id = $11
	RETURNING extract(epoch from created_at)`

	tx.Query(q, entry.Title, entry.EditContent,
		entry.WordCount, entry.Privacy,
		entry.IsCommentable, entry.Rating.IsVotable, entry.InLive, entry.IsShared,
		category, entry.ID, userID.ID).
		Scan(&entry.CreatedAt)

	if tx.Error() == sql.ErrNoRows {
		return srv.StandardError("no_entry")
	}

	entry.Rating.ID = entry.ID
	watchings.AddWatching(tx, userID.ID, entry.ID)

	return nil
}

func reattachImages(srv *utils.MindwellServer, tx *utils.AutoTx, entry *models.Entry, images []int64) {
	tx.Exec("DELETE FROM entry_images WHERE entry_id = $1", entry.ID)

	if len(images) == 0 {
		return
	}

	const q = "INSERT INTO entry_images(entry_id, image_id)	VALUES($1, $2)"
	for _, imageID := range images {
		tx.Exec(q, entry.ID, imageID)
	}

	loadEntryImages(srv, tx, entry, images)
}

func resetTags(tx *utils.AutoTx, entry *models.Entry) {
	tx.Exec("DELETE FROM entry_tags WHERE entry_id = $1", entry.ID)

	setTags(tx, entry)
}

func canEditInLive(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, entryID int64) *models.Error {
	var inLive bool
	const entryQ = "SELECT in_live FROM entries WHERE id = $1"
	tx.Query(entryQ, entryID).Scan(&inLive)
	if inLive {
		return nil
	}

	if userID.Ban.Live {
		return srv.NewError(&i18n.Message{ID: "edit_in_live", Other: "You are not allowed to post in live."})
	}

	if userID.NegKarma {
		return srv.NewError(&i18n.Message{ID: "edit_in_live_karma", Other: "You are not allowed to post in live."})
	}

	var entryCount int64
	const countQ = `
		SELECT count(*)
		FROM entries, 
			(
				SELECT created_at
				FROM entries
				WHERE id = $2
			) AS entry
		WHERE author_id = $1 
			AND date_trunc('day', entries.created_at) = date_trunc('day', entry.created_at)
			AND in_live
	`
	tx.Query(countQ, userID.ID, entryID).Scan(&entryCount)

	if !allowedInLive(userID.FollowersCount, entryCount) {
		return srv.NewError(&i18n.Message{ID: "edit_in_live_followers", Other: "You can't post in live anymore on this day."})
	}

	return nil
}

func newEntryEditor(srv *utils.MindwellServer) func(entries.PutEntriesIDParams, *models.UserID) middleware.Responder {
	return func(params entries.PutEntriesIDParams, uID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			if *params.InLive &&
				(params.Privacy == models.EntryPrivacyAll ||
					params.Privacy == models.EntryPrivacyRegistered ||
					params.Privacy == models.EntryPrivacyInvited) {
				err := canEditInLive(srv, tx, uID, params.ID)
				if err != nil {
					return entries.NewPutEntriesIDForbidden().WithPayload(err)
				}
			}

			if !*params.IsVotable {
				err := allowedWithoutVoting(srv, uID)
				if err != nil {
					return entries.NewPutEntriesIDForbidden().WithPayload(err)
				}
			}

			imageErr := checkImagesOwner(srv, tx, uID.ID, params.Images)
			if imageErr != nil {
				return entries.NewPutEntriesIDForbidden().WithPayload(imageErr)
			}

			authorID := tx.QueryInt64("SELECT author_id FROM entries WHERE id = $1", params.ID)

			entry := &models.Entry{
				ID: params.ID,
				Author: &models.User{
					ID:      authorID,
					IsTheme: authorID != uID.ID,
				},
				Privacy:       params.Privacy,
				IsCommentable: *params.IsCommentable,
				InLive:        *params.InLive,
				IsShared:      *params.IsShared,
				IsAnonymous:   false,
				Rating: &models.Rating{
					IsVotable: *params.IsVotable,
				},
				Title:       *params.Title,
				EditContent: params.Content,
				Tags:        params.Tags,
			}

			if err := editEntry(srv, tx, uID, entry, len(params.Images) > 0); err != nil {
				return entries.NewPutEntriesIDBadRequest().WithPayload(err)
			}

			reattachImages(srv, tx, entry, params.Images)
			resetTags(tx, entry)

			if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
				err := srv.NewError(&i18n.Message{ID: "edit_not_your_entry", Other: "You can't edit someone else's entries."})
				return entries.NewPutEntriesIDForbidden().WithPayload(err)
			}

			updatePrev(entry, uID)

			return entries.NewPutEntriesIDOK().WithPayload(entry)
		})
	}
}

func entryVoteStatus(vote sql.NullFloat64) int64 {
	switch {
	case !vote.Valid:
		return 0
	case vote.Float64 > 0:
		return 1
	default:
		return -1
	}
}

func setEntryRights(entry *models.Entry, userID *models.UserID, themeCreatorID int64) {
	authorID := entry.Author.ID
	if entry.User != nil {
		authorID = entry.User.ID
	}

	entry.Rights = &models.EntryRights{
		Edit:     authorID == userID.ID,
		Delete:   authorID == userID.ID || (entry.Author.IsTheme && themeCreatorID == userID.ID),
		Comment:  authorID == userID.ID || (!userID.Ban.Comment && entry.IsCommentable),
		Vote:     authorID != userID.ID && !userID.Ban.Vote && entry.Rating.IsVotable,
		Complain: authorID != userID.ID && !userID.Ban.Complain,
	}
}

func LoadEntry(srv *utils.MindwellServer, tx *utils.AutoTx, entryID int64, userID *models.UserID) *models.Entry {
	var shared bool
	if !utils.CanViewEntry(tx, userID, entryID) {
		q := sqlf.Select("is_shared").
			From("entries").
			Where("id = ?", entryID)
		if !tx.QueryStmt(q).ScanBool() {
			return &models.Entry{}
		}

		shared = true
	}

	query := feedQuery(userID, 1).
		Where("entries.id = ?", entryID)
	tx.QueryStmt(query)

	feed := loadFeed(srv, tx, userID, false)

	if len(feed.Entries) == 0 {
		return &models.Entry{}
	}

	entry := feed.Entries[0]

	if shared {
		entry.Rights.Comment = false
		entry.Rights.Delete = false
		entry.Rights.Edit = false
		entry.Rights.Vote = false
	} else if userID.ID > 0 {
		cmt := comments.LoadEntryComments(srv, tx, userID, entryID, 5, "", "")
		entry.Comments = cmt
	}

	return entry
}

func newEntryLoader(srv *utils.MindwellServer) func(entries.GetEntriesIDParams, *models.UserID) middleware.Responder {
	return func(params entries.GetEntriesIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			entry := LoadEntry(srv, tx, params.ID, userID)

			if entry.ID == 0 {
				err := srv.StandardError("no_entry")
				return entries.NewGetEntriesIDNotFound().WithPayload(err)
			}

			return entries.NewGetEntriesIDOK().WithPayload(entry)
		})
	}
}

func deleteEntry(srv *utils.MindwellServer, tx *utils.AutoTx, entryID, userID int64) bool {
	if userID == 0 {
		return false
	}

	var entryUserID, themeCreatorID = comments.LoadEntryAuthor(tx, entryID)
	if entryUserID != userID && themeCreatorID != userID {
		return false
	}

	const commentsQuery = "SELECT id FROM comments WHERE entry_id = $1"
	commentIds := tx.QueryInt64s(commentsQuery, entryID)

	for _, id := range commentIds {
		srv.Ntf.SendRemoveComment(tx, id)
	}

	tx.Exec("DELETE from entries WHERE id = $1", entryID)

	return true
}

func newEntryDeleter(srv *utils.MindwellServer) func(entries.DeleteEntriesIDParams, *models.UserID) middleware.Responder {
	return func(params entries.DeleteEntriesIDParams, uID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			ok := deleteEntry(srv, tx, params.ID, uID.ID)
			if ok {
				removePrev(params.ID, uID)
				return entries.NewDeleteEntriesIDOK()
			}

			if tx.Error() == sql.ErrNoRows {
				err := srv.StandardError("no_entry")
				return entries.NewDeleteEntriesIDNotFound().WithPayload(err)
			}

			err := srv.NewError(&i18n.Message{ID: "delete_not_your_entry", Other: "You can't delete someone else's entries."})
			return entries.NewDeleteEntriesIDForbidden().WithPayload(err)
		})
	}
}

var (
	maxID     uint64
	maxDate   int64
	idMap     bitarray.BitArray
	randGuard sync.Mutex
)

func loadRandomEntry(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID) *models.Entry {
	randGuard.Lock()
	defer randGuard.Unlock()

	if idMap == nil {
		idMap = bitarray.NewSparseBitArray()
	}

	const oneDay = 24 * 60 * 60
	now := time.Now().Unix()
	if now > maxDate+oneDay {
		prevID := maxID
		maxDate = now

		tx.Query("SELECT coalesce(max(id), 0) FROM entries")
		tx.Scan(&maxID)

		tx.Query("SELECT distinct(id / 100) FROM entries WHERE id > $1", prevID)
		var k uint64
		for tx.Scan(&k) {
			err := idMap.SetBit(k)
			if err != nil {
				log.Println(err)
			}
		}
	}

	if maxID == 0 {
		return &models.Entry{}
	}

	const maxAttempts = 20
	for i := 0; i < maxAttempts; {
		entryID := rand.Int63n(int64(maxID))

		k := uint64(entryID) / 100
		if k < idMap.Capacity() {
			ok, err := idMap.GetBit(k)
			if err != nil {
				log.Println(err)
			}
			if !ok {
				continue
			}
		}

		entry := LoadEntry(srv, tx, entryID, userID)
		if entry.ID > 0 {
			return entry
		}

		i++
	}

	return &models.Entry{}
}

func newRandomEntryLoader(srv *utils.MindwellServer) func(entries.GetEntriesRandomParams, *models.UserID) middleware.Responder {
	return func(params entries.GetEntriesRandomParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			entry := loadRandomEntry(srv, tx, userID)

			if entry.ID == 0 {
				err := srv.StandardError("no_entry")
				return entries.NewGetEntriesRandomNotFound().WithPayload(err)
			}

			return entries.NewGetEntriesRandomOK().WithPayload(entry)
		})
	}
}
