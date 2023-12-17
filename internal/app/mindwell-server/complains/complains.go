package complains

import (
	"database/sql"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/leporo/sqlf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/chats"
	"github.com/sevings/mindwell-server/restapi/operations/comments"
	"github.com/sevings/mindwell-server/restapi/operations/entries"
	"github.com/sevings/mindwell-server/restapi/operations/themes"
	"github.com/sevings/mindwell-server/restapi/operations/users"
	"github.com/sevings/mindwell-server/restapi/operations/wishes"
	"github.com/sevings/mindwell-server/utils"
)

var errCAY *models.Error
var errCC *models.Error

// ConfigureAPI creates operations handlers
func ConfigureAPI(srv *utils.MindwellServer) {
	errCAY = srv.NewError(&i18n.Message{ID: "complain_against_yourself", Other: "You can't complain against yourself."})
	errCC = srv.NewError(&i18n.Message{ID: "cant_complain", Other: "You are not allowed to complain."})

	srv.API.EntriesPostEntriesIDComplainHandler = entries.PostEntriesIDComplainHandlerFunc(newEntryComplainer(srv))
	srv.API.CommentsPostCommentsIDComplainHandler = comments.PostCommentsIDComplainHandlerFunc(newCommentComplainer(srv))
	srv.API.ChatsPostMessagesIDComplainHandler = chats.PostMessagesIDComplainHandlerFunc(newMessageComplainer(srv))
	srv.API.UsersPostUsersNameComplainHandler = users.PostUsersNameComplainHandlerFunc(newUserComplainer(srv))
	srv.API.ThemesPostThemesNameComplainHandler = themes.PostThemesNameComplainHandlerFunc(newThemeComplainer(srv))
	srv.API.WishesPostWishesIDComplainHandler = wishes.PostWishesIDComplainHandlerFunc(newWishComplainer(srv))
}

const selectPrevQuery = `
    SELECT id 
    FROM complains 
    WHERE user_id = $1 AND subject_id = $2 
        AND type = (SELECT id FROM complain_type WHERE type = $3)
`

const updateQuery = `
    UPDATE complains
    SET content = $2
    WHERE id = $1
`

const createQuery = `
    INSERT INTO complains(user_id, type, subject_id, content)
    VALUES($1, (SELECT id FROM complain_type WHERE type = $2), $3, $4)    
`

func newEntryComplainer(srv *utils.MindwellServer) func(entries.PostEntriesIDComplainParams, *models.UserID) middleware.Responder {
	return func(params entries.PostEntriesIDComplainParams, userID *models.UserID) middleware.Responder {
		if userID.Ban.Complain {
			return entries.NewGetEntriesIDForbidden().WithPayload(errCC)
		}

		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			allowed := utils.CanViewEntry(tx, userID, params.ID)
			if !allowed {
				err := srv.StandardError("no_entry")
				return entries.NewPostEntriesIDComplainNotFound().WithPayload(err)
			}

			authorID := tx.QueryInt64("SELECT author_id FROM entries WHERE id = $1", params.ID)
			if authorID == userID.ID {
				return entries.NewGetEntriesIDForbidden().WithPayload(errCAY)
			}

			id := tx.QueryInt64(selectPrevQuery, userID.ID, params.ID, "entry")
			if tx.Error() != sql.ErrNoRows {
				if tx.Error() != nil {
					err := srv.StandardError("no_entry")
					return entries.NewPostEntriesIDComplainNotFound().WithPayload(err)
				}

				tx.Exec(updateQuery, id, *params.Content)
				return entries.NewPostEntriesIDComplainNoContent()
			}

			tx.Exec(createQuery, userID.ID, "entry", params.ID, *params.Content)
			srv.Ntf.SendNewEntryComplain(tx, params.ID, userID.ID, *params.Content)
			return entries.NewPostEntriesIDComplainNoContent()
		})
	}
}

func newCommentComplainer(srv *utils.MindwellServer) func(comments.PostCommentsIDComplainParams, *models.UserID) middleware.Responder {
	return func(params comments.PostCommentsIDComplainParams, userID *models.UserID) middleware.Responder {
		if userID.Ban.Complain {
			comments.NewPostCommentsIDComplainForbidden().WithPayload(errCC)
		}

		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			var entryID, authorID int64
			tx.Query("SELECT entry_id, author_id FROM comments WHERE id = $1", params.ID)
			tx.Scan(&entryID, &authorID)

			allowed := utils.CanViewEntry(tx, userID, entryID)
			if !allowed {
				err := srv.StandardError("no_comment")
				return comments.NewPostCommentsIDComplainNotFound().WithPayload(err)
			}

			if authorID == userID.ID {
				return comments.NewPostCommentsIDComplainForbidden().WithPayload(errCAY)
			}

			id := tx.QueryInt64(selectPrevQuery, userID.ID, params.ID, "comment")
			if tx.Error() != sql.ErrNoRows {
				if tx.Error() != nil {
					err := srv.StandardError("no_comment")
					return comments.NewPostCommentsIDComplainNotFound().WithPayload(err)
				}

				tx.Exec(updateQuery, id, *params.Content)
				return comments.NewPostCommentsIDComplainNoContent()
			}

			tx.Exec(createQuery, userID.ID, "comment", params.ID, *params.Content)
			srv.Ntf.SendNewCommentComplain(tx, params.ID, userID.ID, *params.Content)
			return comments.NewPostCommentsIDComplainNoContent()
		})
	}
}

func newMessageComplainer(srv *utils.MindwellServer) func(chats.PostMessagesIDComplainParams, *models.UserID) middleware.Responder {
	return func(params chats.PostMessagesIDComplainParams, userID *models.UserID) middleware.Responder {
		if userID.Ban.Complain {
			return chats.NewPostMessagesIDComplainForbidden().WithPayload(errCC)
		}

		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			chatQuery := sqlf.Select("creator_id, partner_id, author_id").
				From("chats").
				Join("messages", "chats.id = messages.chat_id").
				Where("messages.id = ?", params.ID)

			var creatorID, partnerID, authorID int64
			tx.QueryStmt(chatQuery).Scan(&creatorID, &partnerID, &authorID)
			if creatorID != userID.ID && partnerID != userID.ID {
				err := srv.StandardError("no_message")
				return chats.NewPostMessagesIDComplainNotFound().WithPayload(err)
			}

			if authorID == userID.ID {
				return chats.NewPostMessagesIDComplainForbidden().WithPayload(errCAY)
			}

			id := tx.QueryInt64(selectPrevQuery, userID.ID, params.ID, "message")
			if !errors.Is(tx.Error(), sql.ErrNoRows) {
				if tx.Error() != nil {
					err := srv.StandardError("no_message")
					return chats.NewPostMessagesIDComplainNotFound().WithPayload(err)
				}

				tx.Exec(updateQuery, id, *params.Content)
				return chats.NewPostMessagesIDComplainNoContent()
			}

			tx.Exec(createQuery, userID.ID, "message", params.ID, *params.Content)
			srv.Ntf.SendNewMessageComplain(tx, params.ID, userID.ID, *params.Content)
			return chats.NewPostMessagesIDComplainNoContent()
		})
	}
}

func newUserComplainer(srv *utils.MindwellServer) func(users.PostUsersNameComplainParams, *models.UserID) middleware.Responder {
	return func(params users.PostUsersNameComplainParams, userID *models.UserID) middleware.Responder {
		if userID.Ban.Complain {
			return users.NewPostUsersNameComplainForbidden().WithPayload(errCC)
		}

		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			against := utils.LoadUserByName(tx, params.Name)
			if against.ID == 0 || against.IsTheme {
				err := srv.StandardError("no_tlog")
				return users.NewPostUsersNameComplainNotFound().WithPayload(err)
			}
			if against.ID == userID.ID {
				return users.NewPostUsersNameComplainForbidden().WithPayload(errCAY)
			}

			id := tx.QueryInt64(selectPrevQuery, userID.ID, against.ID, "user")
			if !errors.Is(tx.Error(), sql.ErrNoRows) {
				if tx.Error() != nil {
					err := srv.StandardError("no_tlog")
					return users.NewPostUsersNameComplainNotFound().WithPayload(err)
				}

				tx.Exec(updateQuery, id, *params.Content)
				return users.NewPostUsersNameComplainNoContent()
			}

			tx.Exec(createQuery, userID.ID, "user", against.ID, *params.Content)
			srv.Ntf.SendNewUserComplain(tx, against, userID.ID, *params.Content)
			return users.NewPostUsersNameComplainNoContent()
		})
	}
}

func newThemeComplainer(srv *utils.MindwellServer) func(themes.PostThemesNameComplainParams, *models.UserID) middleware.Responder {
	return func(params themes.PostThemesNameComplainParams, userID *models.UserID) middleware.Responder {
		if userID.Ban.Complain {
			return users.NewPostUsersNameComplainForbidden().WithPayload(errCC)
		}

		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			against := utils.LoadUserByName(tx, params.Name)
			if against.ID == 0 || !against.IsTheme {
				err := srv.StandardError("no_theme")
				return themes.NewPostThemesNameComplainNotFound().WithPayload(err)
			}

			id := tx.QueryInt64(selectPrevQuery, userID.ID, against.ID, "theme")
			if !errors.Is(tx.Error(), sql.ErrNoRows) {
				if tx.Error() != nil {
					err := srv.StandardError("no_theme")
					return themes.NewPostThemesNameComplainNotFound().WithPayload(err)
				}

				tx.Exec(updateQuery, id, *params.Content)
				return themes.NewPostThemesNameComplainNoContent()
			}

			tx.Exec(createQuery, userID.ID, "user", against.ID, *params.Content)
			srv.Ntf.SendNewThemeComplain(tx, against, userID.ID, *params.Content)
			return themes.NewPostThemesNameComplainNoContent()
		})
	}
}

func newWishComplainer(srv *utils.MindwellServer) func(wishes.PostWishesIDComplainParams, *models.UserID) middleware.Responder {
	return func(params wishes.PostWishesIDComplainParams, userID *models.UserID) middleware.Responder {
		if userID.Ban.Complain {
			return wishes.NewPostWishesIDComplainForbidden().WithPayload(errCC)
		}

		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			const wishQuery = `
UPDATE wishes
SET state = (SELECT id FROM wish_states WHERE state = 'complained')
WHERE wishes.id = $1 AND to_id = $2
	AND state IN (SELECT id FROM wish_states WHERE state = 'sent' OR state = 'complained')
`
			tx.Exec(wishQuery, params.ID, userID.ID)
			if tx.RowsAffected() < 1 {
				err := srv.StandardError("no_wish")
				return wishes.NewPostWishesIDThankNotFound().WithPayload(err)
			}

			tx.QueryInt64(selectPrevQuery, userID.ID, params.ID, "wish")
			if !errors.Is(tx.Error(), sql.ErrNoRows) {
				if tx.Error() != nil {
					err := srv.StandardError("no_wish")
					return wishes.NewPostWishesIDComplainNotFound().WithPayload(err)
				}

				return wishes.NewPostWishesIDComplainNoContent()
			}

			tx.Exec(createQuery, userID.ID, "wish", params.ID, "")
			srv.Ntf.SendNewWishComplain(tx, params.ID, userID.ID)
			return wishes.NewPostWishesIDComplainNoContent()
		})
	}
}
