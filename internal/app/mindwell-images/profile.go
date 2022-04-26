package images

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi_images/operations/me"
)

func NewAvatarUpdater(mi *MindwellImages) func(me.PutMeAvatarParams, *models.UserID) middleware.Responder {
	return func(params me.PutMeAvatarParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(mi)
		store.ReadImage(params.File)
		store.SetID(userID.ID)

		if store.Error() != nil {
			mi.LogApi().Error(store.Error().Error())
			return me.NewPutMeAvatarBadRequest()
		}

		mi.QueueAction(store, userID.ID, ActionAvatar)

		return me.NewPutMeAvatarNoContent()
	}
}

func NewCoverUpdater(mi *MindwellImages) func(me.PutMeCoverParams, *models.UserID) middleware.Responder {
	return func(params me.PutMeCoverParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(mi)
		store.ReadImage(params.File)
		store.SetID(userID.ID)

		if store.Error() != nil {
			mi.LogApi().Error(store.Error().Error())
			return me.NewPutMeCoverBadRequest()
		}

		mi.QueueAction(store, userID.ID, ActionCover)

		return me.NewPutMeCoverNoContent()
	}
}
