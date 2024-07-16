package images

import (
	"database/sql"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/internal/app/mindwell-server/entries"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/me"
	"github.com/sevings/mindwell-server/restapi/operations/themes"
	"github.com/sevings/mindwell-server/restapi/operations/users"
	"github.com/sevings/mindwell-server/utils"
)

// ConfigureAPI creates operations handlers
func ConfigureAPI(srv *utils.MindwellServer) {
	srv.API.MeGetMeImagesHandler = me.GetMeImagesHandlerFunc(newMyImagesLoader(srv))
	srv.API.UsersGetUsersNameImagesHandler = users.GetUsersNameImagesHandlerFunc(newTlogImagesLoader(srv))
	srv.API.ThemesGetThemesNameImagesHandler = themes.GetThemesNameImagesHandlerFunc(newThemeImagesLoader(srv))
}

func baseFeedQuery(userID *models.UserID, limit int64) *sqlf.Stmt {
	return sqlf.Select("images.id").
		From("images").
		Join("entry_images", "images.id = entry_images.image_id").
		Join("entries", "entry_images.entry_id = entries.id").
		Limit(limit)
}

func myFeedQuery(userID *models.UserID, limit int64) *sqlf.Stmt {
	return baseFeedQuery(userID, limit).
		Where("entries.author_id = ?", userID.ID)
}

func tlogFeedQuery(userID *models.UserID, limit int64, tlog string) *sqlf.Stmt {
	q := baseFeedQuery(userID, limit).
		Join("users AS authors", "entries.author_id = authors.id").
		Join("entry_privacy", "entries.visible_for = entry_privacy.id").
		Where("lower(authors.name) = lower(?)", tlog)

	entries.AddRelationToTlogQuery(q, userID, tlog)
	return utils.AddEntryOpenQuery(q, userID, false)
}

func loadFeed(srv *utils.MindwellServer, tx *utils.AutoTx,
	query, scrollQ *sqlf.Stmt, afterS, beforeS string) *models.ImageList {
	defer scrollQ.Close()
	after := utils.ParseInt64(afterS)
	before := utils.ParseInt64(beforeS)
	reverse := false
	if after > 0 {
		reverse = true
		query.Where("images.id > ?", after).
			OrderBy("images.id ASC")
	} else if before > 0 {
		reverse = false
		query.Where("images.id < ?", before).
			OrderBy("images.id DESC")
	} else {
		query.OrderBy("images.id DESC")
	}

	imageIDs := tx.QueryStmt(query).ScanInt64s()
	if tx.Error() != nil && !errors.Is(tx.Error(), sql.ErrNoRows) {
		return nil
	}

	if reverse {
		for i, j := 0, len(imageIDs)-1; i < j; i, j = i+1, j-1 {
			imageIDs[i], imageIDs[j] = imageIDs[j], imageIDs[i]
		}
	}

	feed := &models.ImageList{}
	for _, imageID := range imageIDs {
		img := utils.LoadImage(srv, tx, imageID)
		if img != nil {
			feed.Data = append(feed.Data, img)
		}
	}

	if len(feed.Data) == 0 {
		if before > 0 {
			afterQuery := scrollQ.Clone().Where("images.id >= ?", before).
				OrderBy("images.id ASC")

			tx.QueryStmt(afterQuery)
			var nextAfter int64
			tx.Scan(&nextAfter)
			feed.NextAfter = utils.FormatInt64(nextAfter)
		}

		if after > 0 {
			beforeQuery := scrollQ.Clone().Where("images.id <= ?", after).
				OrderBy("images.id DESC")

			tx.QueryStmt(beforeQuery)
			var nextBefore int64
			tx.Scan(&nextBefore)
			feed.NextBefore = utils.FormatInt64(nextBefore)
		}
	} else {
		oldest := feed.Data[len(feed.Data)-1].ID
		newest := feed.Data[0].ID
		var id int64

		feed.NextBefore = utils.FormatInt64(oldest)
		beforeQuery := scrollQ.Clone().Where("images.id < ?", oldest)
		tx.QueryStmt(beforeQuery)
		tx.Scan(&id)
		feed.HasBefore = id > 0

		feed.NextAfter = utils.FormatInt64(newest)
		afterQuery := scrollQ.Clone().Where("images.id > ?", newest)
		tx.QueryStmt(afterQuery)
		tx.Scan(&id)
		feed.HasAfter = id > 0
	}

	return feed
}

func newMyImagesLoader(srv *utils.MindwellServer) func(me.GetMeImagesParams, *models.UserID) middleware.Responder {
	return func(params me.GetMeImagesParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			query := myFeedQuery(userID, *params.Limit)
			scrollQ := myFeedQuery(userID, 1)
			feed := loadFeed(srv, tx, query, scrollQ, *params.After, *params.Before)
			if feed == nil {
				return me.NewGetMeImagesNotFound()
			}

			return me.NewGetMeImagesOK().WithPayload(feed)
		})
	}
}

func newTlogImagesLoader(srv *utils.MindwellServer) func(users.GetUsersNameImagesParams, *models.UserID) middleware.Responder {
	return func(params users.GetUsersNameImagesParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewTlogName(tx, userID, params.Name)
			if !canView {
				err := srv.StandardError("no_tlog")
				return users.NewGetUsersNameImagesNotFound().WithPayload(err)
			}

			query := tlogFeedQuery(userID, *params.Limit, params.Name)
			scrollQ := tlogFeedQuery(userID, 1, params.Name)
			feed := loadFeed(srv, tx, query, scrollQ, *params.After, *params.Before)
			if feed == nil {
				return users.NewGetUsersNameImagesNotFound()
			}

			return users.NewGetUsersNameImagesOK().WithPayload(feed)
		})
	}
}

func newThemeImagesLoader(srv *utils.MindwellServer) func(themes.GetThemesNameImagesParams, *models.UserID) middleware.Responder {
	return func(params themes.GetThemesNameImagesParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewTlogName(tx, userID, params.Name)
			if !canView {
				err := srv.StandardError("no_theme")
				return themes.NewGetThemesNameImagesNotFound().WithPayload(err)
			}

			query := tlogFeedQuery(userID, *params.Limit, params.Name)
			scrollQ := tlogFeedQuery(userID, 1, params.Name)
			feed := loadFeed(srv, tx, query, scrollQ, *params.After, *params.Before)
			if feed == nil {
				return themes.NewGetThemesNameImagesNotFound()
			}

			return themes.NewGetThemesNameImagesOK().WithPayload(feed)
		})
	}
}
