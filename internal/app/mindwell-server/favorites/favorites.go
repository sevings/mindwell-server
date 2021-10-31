package favorites

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/favorites"
	"github.com/sevings/mindwell-server/utils"
)

// ConfigureAPI creates operations handlers
func ConfigureAPI(srv *utils.MindwellServer) {
	srv.API.FavoritesGetEntriesIDFavoriteHandler = favorites.GetEntriesIDFavoriteHandlerFunc(newStatusLoader(srv))
	srv.API.FavoritesPutEntriesIDFavoriteHandler = favorites.PutEntriesIDFavoriteHandlerFunc(newFavoriteAdder(srv))
	srv.API.FavoritesDeleteEntriesIDFavoriteHandler = favorites.DeleteEntriesIDFavoriteHandlerFunc(newFavoriteDeleter(srv))
}

func favoriteCount(tx *utils.AutoTx, entryID int64) int64 {
	const q = `
		SELECT favorites_count
		FROM entries
		WHERE id = $1`
	return tx.QueryInt64(q, entryID)
}

func favoriteStatus(tx *utils.AutoTx, userID, entryID int64) *models.FavoriteStatus {
	status := models.FavoriteStatus{ID: entryID}
	status.Count = favoriteCount(tx, entryID)

	const q = `
		SELECT TRUE 
		FROM favorites
		WHERE user_id = $1 AND entry_id = $2`
	status.IsFavorited = tx.QueryBool(q, userID, entryID)

	return &status
}

func newStatusLoader(srv *utils.MindwellServer) func(favorites.GetEntriesIDFavoriteParams, *models.UserID) middleware.Responder {
	return func(params favorites.GetEntriesIDFavoriteParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewEntry(tx, userID, params.ID)
			if !canView {
				err := srv.StandardError("no_entry")
				return favorites.NewGetEntriesIDFavoriteNotFound().WithPayload(err)
			}

			status := favoriteStatus(tx, userID.ID, params.ID)
			return favorites.NewGetEntriesIDFavoriteOK().WithPayload(status)
		})
	}
}

func addToFavorites(tx *utils.AutoTx, userID, entryID int64) *models.FavoriteStatus {
	const q = `
		INSERT INTO favorites(user_id, entry_id)
		VALUES($1, $2)
		ON CONFLICT ON CONSTRAINT unique_user_favorite
		DO NOTHING`

	tx.Exec(q, userID, entryID)

	status := models.FavoriteStatus{
		ID:          entryID,
		IsFavorited: true,
		Count:       favoriteCount(tx, entryID),
	}

	return &status
}

func newFavoriteAdder(srv *utils.MindwellServer) func(favorites.PutEntriesIDFavoriteParams, *models.UserID) middleware.Responder {
	return func(params favorites.PutEntriesIDFavoriteParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewEntry(tx, userID, params.ID)
			if !canView {
				err := srv.StandardError("no_entry")
				return favorites.NewPutEntriesIDFavoriteNotFound().WithPayload(err)
			}

			status := addToFavorites(tx, userID.ID, params.ID)
			return favorites.NewPutEntriesIDFavoriteOK().WithPayload(status)
		})
	}
}

func removeFromFavorites(tx *utils.AutoTx, userID, entryID int64) *models.FavoriteStatus {
	const q = `
		DELETE FROM favorites
		WHERE user_id = $1 AND entry_id = $2`

	tx.Exec(q, userID, entryID)

	status := models.FavoriteStatus{
		ID:          entryID,
		IsFavorited: false,
		Count:       favoriteCount(tx, entryID),
	}

	return &status
}

func newFavoriteDeleter(srv *utils.MindwellServer) func(favorites.DeleteEntriesIDFavoriteParams, *models.UserID) middleware.Responder {
	return func(params favorites.DeleteEntriesIDFavoriteParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewEntry(tx, userID, params.ID)
			if !canView {
				err := srv.StandardError("no_entry")
				return favorites.NewDeleteEntriesIDFavoriteNotFound().WithPayload(err)
			}

			status := removeFromFavorites(tx, userID.ID, params.ID)
			return favorites.NewDeleteEntriesIDFavoriteOK().WithPayload(status)
		})
	}
}
