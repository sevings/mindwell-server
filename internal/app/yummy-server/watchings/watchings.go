package watchings

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/yummy-server/internal/app/yummy-server/utils"
	"github.com/sevings/yummy-server/models"
	"github.com/sevings/yummy-server/restapi/operations"
	"github.com/sevings/yummy-server/restapi/operations/watchings"
)

// ConfigureAPI creates operations handlers
func ConfigureAPI(db *sql.DB, api *operations.YummyAPI) {
	api.WatchingsGetEntriesIDWatchingHandler = watchings.GetEntriesIDWatchingHandlerFunc(newWatchingStatusLoader(db))
	api.WatchingsPutEntriesIDWatchingHandler = watchings.PutEntriesIDWatchingHandlerFunc(newWatchingAdder(db))
	api.WatchingsDeleteEntriesIDWatchingHandler = watchings.DeleteEntriesIDWatchingHandlerFunc(newWatchingDeleter(db))
}

func watchingStatus(tx *utils.AutoTx, userID, entryID int64) *models.WatchingStatus {
	const q = `
		SELECT TRUE 
		FROM watching
		WHERE user_id = $1 AND entry_id = $2`

	status := models.WatchingStatus{ID: entryID}

	tx.Query(q, userID, entryID).Scan(&status.IsWatching)

	return &status
}

func newWatchingStatusLoader(db *sql.DB) func(watchings.GetEntriesIDWatchingParams, *models.UserID) middleware.Responder {
	return func(params watchings.GetEntriesIDWatchingParams, uID *models.UserID) middleware.Responder {
		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			userID := int64(*uID)
			canView := utils.CanViewEntry(tx, userID, params.ID)
			if !canView {
				return watchings.NewGetEntriesIDWatchingNotFound()
			}

			status := watchingStatus(tx, userID, params.ID)
			return watchings.NewGetEntriesIDWatchingOK().WithPayload(status)
		})
	}
}

func AddWatching(tx *utils.AutoTx, userID, entryID int64) *models.WatchingStatus {
	const q = `
		INSERT INTO watching(user_id, entry_id)
		VALUES($1, $2)
		ON CONFLICT ON CONSTRAINT unique_user_watching
		DO NOTHING`

	tx.Exec(q, userID, entryID)

	status := models.WatchingStatus{
		ID:         entryID,
		IsWatching: true,
	}

	return &status
}

func newWatchingAdder(db *sql.DB) func(watchings.PutEntriesIDWatchingParams, *models.UserID) middleware.Responder {
	return func(params watchings.PutEntriesIDWatchingParams, uID *models.UserID) middleware.Responder {
		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			userID := int64(*uID)
			canView := utils.CanViewEntry(tx, userID, params.ID)
			if !canView {
				return watchings.NewPutEntriesIDWatchingNotFound()
			}

			status := AddWatching(tx, userID, params.ID)
			return watchings.NewPutEntriesIDWatchingOK().WithPayload(status)
		})
	}
}

func RemoveWatching(tx *utils.AutoTx, userID, entryID int64) *models.WatchingStatus {
	const q = `
		DELETE FROM watching
		WHERE user_id = $1 AND entry_id = $2`

	tx.Exec(q, userID, entryID)

	status := models.WatchingStatus{
		ID:         entryID,
		IsWatching: false,
	}

	return &status
}

func newWatchingDeleter(db *sql.DB) func(watchings.DeleteEntriesIDWatchingParams, *models.UserID) middleware.Responder {
	return func(params watchings.DeleteEntriesIDWatchingParams, uID *models.UserID) middleware.Responder {
		return utils.Transact(db, func(tx *utils.AutoTx) middleware.Responder {
			userID := int64(*uID)
			canView := utils.CanViewEntry(tx, userID, params.ID)
			if !canView {
				return watchings.NewDeleteEntriesIDWatchingNotFound()
			}

			status := RemoveWatching(tx, userID, params.ID)
			return watchings.NewDeleteEntriesIDWatchingOK().WithPayload(status)
		})
	}
}
