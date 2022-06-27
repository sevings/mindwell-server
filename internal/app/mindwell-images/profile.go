package images

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi_images/operations/me"
	"github.com/sevings/mindwell-server/restapi_images/operations/themes"
	"github.com/sevings/mindwell-server/utils"
	"io"
)

func updateProfilePhoto(mi *MindwellImages, profileID int64, file io.ReadCloser, action string) bool {
	store := newImageStore(mi)
	store.ReadImage(file)
	store.SetID(profileID)

	if store.Error() != nil {
		mi.LogApi().Error(store.Error().Error())
		return false
	}

	mi.QueueAction(store, profileID, action)
	return true
}

func NewAvatarUpdater(mi *MindwellImages) func(me.PutMeAvatarParams, *models.UserID) middleware.Responder {
	return func(params me.PutMeAvatarParams, userID *models.UserID) middleware.Responder {
		ok := updateProfilePhoto(mi, userID.ID, params.File, ActionAvatar)
		if !ok {
			return me.NewPutMeAvatarBadRequest()
		}

		return me.NewPutMeAvatarNoContent()
	}
}

func NewCoverUpdater(mi *MindwellImages) func(me.PutMeCoverParams, *models.UserID) middleware.Responder {
	return func(params me.PutMeCoverParams, userID *models.UserID) middleware.Responder {
		ok := updateProfilePhoto(mi, userID.ID, params.File, ActionCover)
		if !ok {
			return me.NewPutMeCoverBadRequest()
		}

		return me.NewPutMeCoverNoContent()
	}
}

func checkThemeAdmin(tx *utils.AutoTx, userID *models.UserID, name string) (int64, bool) {
	const q = `SELECT id, creator_id FROM users WHERE lower(name) = lower($1)`

	var themeID, creatorID int64
	tx.Query(q, name).Scan(&themeID, &creatorID)

	return themeID, creatorID == userID.ID
}

func NewThemeAvatarUpdater(mi *MindwellImages) func(themes.PutThemesNameAvatarParams, *models.UserID) middleware.Responder {
	return func(params themes.PutThemesNameAvatarParams, userID *models.UserID) middleware.Responder {
		tx := utils.NewAutoTx(mi.DB())
		defer tx.Finish()

		themeID, allowed := checkThemeAdmin(tx, userID, params.Name)
		if !allowed {
			return themes.NewPutThemesNameAvatarForbidden()
		}

		ok := updateProfilePhoto(mi, themeID, params.File, ActionAvatar)
		if !ok {
			return themes.NewPutThemesNameAvatarBadRequest()
		}

		return themes.NewPutThemesNameAvatarNoContent()
	}
}

func NewThemeCoverUpdater(mi *MindwellImages) func(themes.PutThemesNameCoverParams, *models.UserID) middleware.Responder {
	return func(params themes.PutThemesNameCoverParams, userID *models.UserID) middleware.Responder {
		tx := utils.NewAutoTx(mi.DB())
		defer tx.Finish()

		themeID, allowed := checkThemeAdmin(tx, userID, params.Name)
		if !allowed {
			return themes.NewPutThemesNameCoverForbidden()
		}

		ok := updateProfilePhoto(mi, themeID, params.File, ActionCover)
		if !ok {
			return themes.NewPutThemesNameCoverBadRequest()
		}

		return themes.NewPutThemesNameCoverNoContent()
	}
}
