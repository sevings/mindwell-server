package comments

import (
	"database/sql"
	"fmt"
	"github.com/sevings/mindwell-server/restapi/operations/me"
	"github.com/sevings/mindwell-server/restapi/operations/themes"
	"github.com/sevings/mindwell-server/restapi/operations/users"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	usersIml "github.com/sevings/mindwell-server/internal/app/mindwell-server/users"
	"github.com/sevings/mindwell-server/restapi/operations/comments"

	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/utils"
)

// ConfigureAPI creates operations handlers
func ConfigureAPI(srv *utils.MindwellServer) {
	srv.API.CommentsGetCommentsIDHandler = comments.GetCommentsIDHandlerFunc(newCommentLoader(srv))
	srv.API.CommentsPutCommentsIDHandler = comments.PutCommentsIDHandlerFunc(newCommentEditor(srv))
	srv.API.CommentsDeleteCommentsIDHandler = comments.DeleteCommentsIDHandlerFunc(newCommentDeleter(srv))

	srv.API.CommentsGetEntriesIDCommentsHandler = comments.GetEntriesIDCommentsHandlerFunc(newEntryCommentsLoader(srv))
	srv.API.CommentsPostEntriesIDCommentsHandler = comments.PostEntriesIDCommentsHandlerFunc(newCommentPoster(srv))

	srv.API.MeGetMeCommentsHandler = me.GetMeCommentsHandlerFunc(newMyCommentsLoader(srv))
	srv.API.UsersGetUsersNameCommentsHandler = users.GetUsersNameCommentsHandlerFunc(newUserCommentsLoader(srv))
	srv.API.ThemesGetThemesNameCommentsHandler = themes.GetThemesNameCommentsHandlerFunc(newThemeCommentsLoader(srv))
	srv.API.CommentsGetCommentsHandler = comments.GetCommentsHandlerFunc(newAllCommentsLoader(srv))

	srv.API.CommentsGetEntriesIDCommentatorHandler = comments.GetEntriesIDCommentatorHandlerFunc(newCommentatorLoader(srv))
}

var imgRe *regexp.Regexp
var iMdRe *regexp.Regexp
var urlRe *regexp.Regexp

func init() {
	imgRe = regexp.MustCompile(`(?i)^https?.+\.(?:png|jpg|jpeg|gif)(?:\?\S*)?$`)
	iMdRe = regexp.MustCompile(`!\[[^]]*]\(([^)]+)\)`)
	urlRe = regexp.MustCompile(`([a-zA-Z][a-zA-Z\d+\-.]*://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=%]+)`)
}

func HtmlContent(content string) string {
	replaceURL := func(href string) string {
		if imgRe.MatchString(href) {
			return fmt.Sprintf("<img src=\"%s\">", href)
		}

		text, err := url.QueryUnescape(href)
		if err != nil {
			text = href
		} else {
			href = text
		}
		if len(text) > 40 {
			text = text[:40] + "..."
		}
		return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener nofollow">%s</a>`, href, text)
	}

	content = strings.TrimSpace(content)
	content = utils.ReplaceToHtml(content)
	content = iMdRe.ReplaceAllString(content, "$1")
	content = urlRe.ReplaceAllStringFunc(content, replaceURL)

	return "<p>" + content + "</p>"
}

func commentVote(vote sql.NullFloat64) int64 {
	switch {
	case !vote.Valid:
		return 0
	case vote.Float64 > 0:
		return 1
	default:
		return -1
	}
}

func setCommentRights(comment *models.Comment, userID *models.UserID, cmtUserID, entryUserID, themeCreatorID int64) {
	comment.Rights = &models.CommentRights{
		Edit:     cmtUserID == userID.ID,
		Delete:   cmtUserID == userID.ID || entryUserID == userID.ID || themeCreatorID == userID.ID,
		Vote:     cmtUserID != userID.ID && !userID.Ban.Vote,
		Complain: cmtUserID != userID.ID,
	}
}

func setCommentText(comment *models.Comment) {
	comment.Content = HtmlContent(comment.EditContent)
}

func LoadEntryAuthor(tx *utils.AutoTx, entryID int64) (entryUserID int64, themeCreatorID int64) {
	const entryAuthorQuery = `
SELECT user_id, COALESCE(creator_id, 0)
FROM entries
JOIN users ON author_id = users.id
WHERE entries.id = $1`
	tx.Query(entryAuthorQuery, entryID).Scan(&entryUserID, &themeCreatorID)
	return entryUserID, themeCreatorID
}

func LoadComment(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, commentID int64) *models.Comment {
	q := baseFeedQuery(userID, 1).
		Where("comments.id = ?", commentID)

	tx.QueryStmt(q)

	cmts := loadFeed(srv, tx, userID, false)
	if len(cmts.Data) > 0 {
		return cmts.Data[0]
	}

	return nil
}

func newCommentLoader(srv *utils.MindwellServer) func(comments.GetCommentsIDParams, *models.UserID) middleware.Responder {
	return func(params comments.GetCommentsIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			comment := LoadComment(srv, tx, userID, params.ID)
			if tx.Error() != nil {
				err := srv.StandardError("no_comment")
				return comments.NewGetCommentsIDNotFound().WithPayload(err)
			}

			canView := utils.CanViewEntry(tx, userID, comment.EntryID)
			if !canView {
				err := srv.StandardError("no_entry")
				return comments.NewGetCommentsIDNotFound().WithPayload(err)
			}

			return comments.NewGetCommentsIDOK().WithPayload(comment)
		})
	}
}

func editComment(tx *utils.AutoTx, comment *models.Comment) {
	const q = `
		UPDATE comments
		SET edit_content = $2
		WHERE id = $1`

	tx.Exec(q, comment.ID, comment.EditContent)
}

func newCommentEditor(srv *utils.MindwellServer) func(comments.PutCommentsIDParams, *models.UserID) middleware.Responder {
	return func(params comments.PutCommentsIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			comment := LoadComment(srv, tx, userID, params.ID)
			if tx.Error() != nil {
				err := srv.StandardError("no_comment")
				return comments.NewGetCommentsIDNotFound().WithPayload(err)
			}

			if !comment.Rights.Edit {
				err := srv.NewError(&i18n.Message{ID: "edit_not_your_comment", Other: "You can't edit someone else's comments."})
				return comments.NewGetCommentsIDForbidden().WithPayload(err)
			}

			comment.Content = HtmlContent(params.Content)
			comment.EditContent = params.Content
			editComment(tx, comment)
			if tx.Error() != nil {
				err := srv.NewError(nil)
				return comments.NewGetCommentsIDNotFound().WithPayload(err)
			}

			srv.Ntf.SendUpdateComment(tx, comment)

			updatePrev(comment, userID)

			return comments.NewPutCommentsIDOK().WithPayload(comment)
		})
	}
}

func canDeleteComment(tx *utils.AutoTx, userID *models.UserID, commentID int64) bool {
	if userID.ID == 0 {
		return false
	}

	const q = `
		SELECT comments.user_id, entries.user_id, COALESCE(users.creator_id, 0)
		FROM comments
		JOIN entries on entries.id = comments.entry_id
		JOIN users on entries.author_id = users.id
		WHERE comments.id = $1`

	var commentUserID, entryUserID, themeCreatorID int64
	tx.Query(q, commentID).Scan(&commentUserID, &entryUserID, &themeCreatorID)

	return commentUserID == userID.ID || entryUserID == userID.ID || themeCreatorID == userID.ID
}

func deleteComment(tx *utils.AutoTx, commentID int64) {
	const q = `
		DELETE FROM comments
		WHERE id = $1`

	tx.Exec(q, commentID)
}

func newCommentDeleter(srv *utils.MindwellServer) func(comments.DeleteCommentsIDParams, *models.UserID) middleware.Responder {
	return func(params comments.DeleteCommentsIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			allowed := canDeleteComment(tx, userID, params.ID)
			if tx.Error() != nil {
				err := srv.NewError(nil)
				return comments.NewDeleteCommentsIDNotFound().WithPayload(err)
			}
			if !allowed {
				err := srv.NewError(&i18n.Message{ID: "delete_not_your_comment", Other: "You can't delete someone else's comments."})
				return comments.NewDeleteCommentsIDForbidden().WithPayload(err)
			}

			deleteComment(tx, params.ID)
			if tx.Error() != nil {
				err := srv.NewError(nil)
				return comments.NewDeleteCommentsIDNotFound().WithPayload(err)
			}

			srv.Ntf.SendRemoveComment(tx, params.ID)

			removePrev(params.ID, userID)

			return comments.NewDeleteCommentsIDOK()
		})
	}
}

func newEntryCommentsLoader(srv *utils.MindwellServer) func(comments.GetEntriesIDCommentsParams, *models.UserID) middleware.Responder {
	return func(params comments.GetEntriesIDCommentsParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewEntry(tx, userID, params.ID)
			if !canView {
				err := srv.StandardError("no_entry")
				return comments.NewGetEntriesIDCommentsNotFound().WithPayload(err)
			}

			data := LoadEntryComments(srv, tx, userID, params.ID, *params.Limit, *params.After, *params.Before)
			if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
				err := srv.NewError(nil)
				return comments.NewGetEntriesIDCommentsNotFound().WithPayload(err)
			}

			return comments.NewGetEntriesIDCommentsOK().WithPayload(data)
		})
	}
}

func newMyCommentsLoader(srv *utils.MindwellServer) func(me.GetMeCommentsParams, *models.UserID) middleware.Responder {
	return func(params me.GetMeCommentsParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			data := loadUserComments(srv, tx, userID, *params.Limit, userID.Name, *params.After, *params.Before)
			return me.NewGetMeCommentsOK().WithPayload(data)
		})
	}
}

func newUserCommentsLoader(srv *utils.MindwellServer) func(users.GetUsersNameCommentsParams, *models.UserID) middleware.Responder {
	return func(params users.GetUsersNameCommentsParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewTlogName(tx, userID, params.Name)
			if !canView {
				err := srv.StandardError("no_tlog")
				return users.NewGetUsersNameCommentsNotFound().WithPayload(err)
			}

			data := loadUserComments(srv, tx, userID, *params.Limit, params.Name, *params.After, *params.Before)
			if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
				err := srv.NewError(nil)
				return users.NewGetUsersNameCommentsNotFound().WithPayload(err)
			}

			return users.NewGetUsersNameCommentsOK().WithPayload(data)
		})
	}
}

func newThemeCommentsLoader(srv *utils.MindwellServer) func(themes.GetThemesNameCommentsParams, *models.UserID) middleware.Responder {
	return func(params themes.GetThemesNameCommentsParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewTlogName(tx, userID, params.Name)
			if !canView {
				err := srv.StandardError("no_theme")
				return themes.NewGetThemesNameCommentsNotFound().WithPayload(err)
			}

			data := loadThemeComments(srv, tx, userID, *params.Limit, params.Name, *params.After, *params.Before)
			if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
				err := srv.NewError(nil)
				return themes.NewGetThemesNameCommentsNotFound().WithPayload(err)
			}

			return themes.NewGetThemesNameCommentsOK().WithPayload(data)
		})
	}
}

func newAllCommentsLoader(srv *utils.MindwellServer) func(comments.GetCommentsParams, *models.UserID) middleware.Responder {
	return func(params comments.GetCommentsParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			data := loadAllComments(srv, tx, userID, *params.Limit, *params.After, *params.Before)
			return comments.NewGetCommentsOK().WithPayload(data)
		})
	}
}

func commentatorID(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, entryID int64) (int64, *models.Error) {
	const q = `
		SELECT author_id, user_id, is_commentable, is_anonymous
		FROM entries
		WHERE id = $1
	`

	var entryAuthorID, entryUserID int64
	var commentable, anonymous bool
	tx.Query(q, entryID).Scan(&entryAuthorID, &entryUserID, &commentable, &anonymous)
	if entryUserID == userID.ID {
		if anonymous {
			return entryAuthorID, nil
		}

		return entryUserID, nil
	}

	if !commentable || userID.Ban.Comment || !utils.CanViewEntry(tx, userID, entryID) {
		err := srv.NewError(&i18n.Message{ID: "cant_comment", Other: "You can't comment this entry."})
		return 0, err
	}

	return userID.ID, nil
}

func postComment(tx *utils.AutoTx, userID *models.UserID, comment *models.Comment) {
	const q = `
		INSERT INTO comments (user_id, author_id, entry_id, edit_content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, extract(epoch from created_at)`

	comment.Content = HtmlContent(comment.EditContent)
	comment.Rating = &models.Rating{
		IsVotable: true,
	}
	comment.Rights = &models.CommentRights{
		Edit:   true,
		Delete: true,
		Vote:   false,
	}

	tx.Query(q, userID.ID, comment.Author.ID, comment.EntryID, comment.EditContent)
	tx.Scan(&comment.ID, &comment.CreatedAt)

	comment.Rating.ID = comment.ID
}

func newCommentPoster(srv *utils.MindwellServer) func(comments.PostEntriesIDCommentsParams, *models.UserID) middleware.Responder {
	return func(params comments.PostEntriesIDCommentsParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			prev, found := checkPrev(params, userID)
			if found {
				return comments.NewPostEntriesIDCommentsCreated().WithPayload(prev)
			}

			authorID, err := commentatorID(srv, tx, userID, params.ID)
			if err != nil {
				return comments.NewPostEntriesIDCommentsNotFound().WithPayload(err)
			}

			comment := &models.Comment{
				Author:      usersIml.LoadUserByID(srv, tx, authorID),
				EditContent: params.Content,
				EntryID:     params.ID,
			}

			postComment(tx, userID, comment)
			if tx.Error() != nil {
				err := srv.NewError(nil)
				return comments.NewPostEntriesIDCommentsNotFound().WithPayload(err)
			}

			srv.Ntf.SendNewComment(tx, comment)

			setPrev(comment, userID)

			return comments.NewPostEntriesIDCommentsCreated().WithPayload(comment)
		})
	}
}

func newCommentatorLoader(srv *utils.MindwellServer) func(comments.GetEntriesIDCommentatorParams, *models.UserID) middleware.Responder {
	return func(params comments.GetEntriesIDCommentatorParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			authorID, err := commentatorID(srv, tx, userID, params.ID)
			if err != nil {
				return comments.NewGetEntriesIDCommentatorNotFound().WithPayload(err)
			}

			author := usersIml.LoadUserByID(srv, tx, authorID)
			return comments.NewGetEntriesIDCommentatorOK().WithPayload(author)
		})
	}
}
