package wishes

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/wishes"
	"github.com/sevings/mindwell-server/utils"
	"go.uber.org/zap"
	"strings"
)

var wishesEnabled bool

// ConfigureAPI creates operations handlers
func ConfigureAPI(srv *utils.MindwellServer) {
	srv.API.WishesGetWishesIDHandler = wishes.GetWishesIDHandlerFunc(newWishLoader(srv))
	srv.API.WishesPutWishesIDHandler = wishes.PutWishesIDHandlerFunc(newWishSender(srv))
	srv.API.WishesDeleteWishesIDHandler = wishes.DeleteWishesIDHandlerFunc(newWishDecliner(srv))
	srv.API.WishesPostWishesIDThankHandler = wishes.PostWishesIDThankHandlerFunc(newWishThanker(srv))

	wishesEnabled = srv.ConfigBool("adm.wish_enabled")
	if wishesEnabled {
		srv.PS.Subscribe("online_users", createWishFromUserData(srv))
	}
}

func createWishFromUserData(srv *utils.MindwellServer) func(userData []byte) {
	return func(userData []byte) {
		var user models.User
		err := json.Unmarshal(userData, &user)
		if err != nil {
			srv.Log("wishes").Error(err.Error())
			return
		}

		createWishFromUser(srv, user.ID, user.Name)
	}
}

func newWishLoader(srv *utils.MindwellServer) func(wishes.GetWishesIDParams, *models.UserID) middleware.Responder {
	return func(params wishes.GetWishesIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			wish, found := LoadWish(tx, userID, params.ID)
			if !found {
				err := srv.StandardError("no_wish")
				return wishes.NewGetWishesIDNotFound().WithPayload(err)
			}

			return wishes.NewGetWishesIDOK().WithPayload(wish)
		})
	}
}

func newWishSender(srv *utils.MindwellServer) func(wishes.PutWishesIDParams, *models.UserID) middleware.Responder {
	return func(params wishes.PutWishesIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			content := strings.TrimSpace(params.Content)
			content = utils.ReplaceToHtml(content)
			toID, found := saveWish(tx, userID, params.ID, content)
			if !found {
				err := srv.StandardError("no_wish")
				return wishes.NewPutWishesIDNotFound().WithPayload(err)
			}

			receiver := utils.LoadUser(tx, toID)
			srv.Ntf.SendWishReceived(tx, params.ID, receiver.Name)

			return wishes.NewPutWishesIDOK()
		})
	}
}

func newWishDecliner(srv *utils.MindwellServer) func(wishes.DeleteWishesIDParams, *models.UserID) middleware.Responder {
	return func(params wishes.DeleteWishesIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			ok := declineWish(tx, userID, params.ID)
			if !ok {
				err := srv.StandardError("no_wish")
				return wishes.NewDeleteWishesIDNotFound().WithPayload(err)
			}

			return wishes.NewDeleteWishesIDOK()
		})
	}
}

func createWishFromUser(srv *utils.MindwellServer, id int64, name string) {
	tx := utils.NewAutoTx(srv.DB)
	defer tx.Finish()

	logger := srv.Log("wishes")

	wishID, created := createWish(tx, id)
	if created {
		srv.Ntf.SendWishCreated(tx, wishID, name)
		logger.Info("created", zap.String("from", name))
	}

	if tx.HasQueryError() {
		logger.Error(tx.Error().Error())
	}
}

func newWishThanker(srv *utils.MindwellServer) func(wishes.PostWishesIDThankParams, *models.UserID) middleware.Responder {
	return func(params wishes.PostWishesIDThankParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			ok := thankWish(tx, userID, params.ID)
			if !ok {
				err := srv.StandardError("no_wish")
				return wishes.NewPostWishesIDThankNotFound().WithPayload(err)
			}

			if wishesEnabled {
				go createWishFromUser(srv, userID.ID, userID.Name)
			}

			return wishes.NewPostWishesIDThankNoContent()
		})
	}
}
