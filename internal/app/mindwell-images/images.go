package images

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi_images/operations/images"
	"github.com/sevings/mindwell-server/utils"
)

func setProcessingImage(mi *MindwellImages, img *models.Image) {
	img.Thumbnail = &models.ImageSize{
		Width:  100,
		Height: 100,
		URL:    mi.BaseURL() + "albums/thumbnails/processing.png",
	}

	img.Small = &models.ImageSize{
		Width:  480,
		Height: 300,
		URL:    mi.BaseURL() + "albums/small/processing.png",
	}

	img.Medium = &models.ImageSize{
		Width:  960,
		Height: 600,
		URL:    mi.BaseURL() + "albums/medium/processing.png",
	}

	img.Large = &models.ImageSize{
		Width:  1440,
		Height: 900,
		URL:    mi.BaseURL() + "albums/large/processing.png",
	}
}

func NewImageUploader(mi *MindwellImages) func(images.PostImagesParams, *models.UserID) middleware.Responder {
	return func(params images.PostImagesParams, userID *models.UserID) middleware.Responder {
		store := newImageStore(mi)
		store.ReadImage(params.File)

		if store.Error() != nil {
			mi.LogApi().Error(store.Error().Error())
			return images.NewPostImagesBadRequest()
		}

		img := &models.Image{
			Author: &models.User{
				ID:   userID.ID,
				Name: userID.Name,
			},
			IsAnimated: store.IsAnimated(),
			Processing: true,
		}

		return utils.Transact(mi.DB(), func(tx *utils.AutoTx) middleware.Responder {
			tx.Query("INSERT INTO images(user_id, path, extension, processing) VALUES($1, 'processing', $2, $3) RETURNING id",
				userID.ID, store.FileExtension(), img.Processing)
			tx.Scan(&img.ID)

			store.SetID(img.ID)
			tx.Exec("UPDATE images SET path = $2 WHERE id = $1", img.ID, store.FileName())

			if tx.Error() != nil {
				return images.NewPostImagesBadRequest()
			}

			setProcessingImage(mi, img)
			mi.QueueAction(store, img.ID, ActionAlbum)

			return images.NewPostImagesOK().WithPayload(img)
		})
	}
}

func NewImageLoader(mi *MindwellImages) func(images.GetImagesIDParams, *models.UserID) middleware.Responder {
	return func(params images.GetImagesIDParams, userID *models.UserID) middleware.Responder {
		return utils.Transact(mi.DB(), func(tx *utils.AutoTx) middleware.Responder {
			var authorID int64
			var path, extension string
			var processing bool

			tx.Query("SELECT user_id, path, extension, processing FROM images WHERE id = $1", params.ID).
				Scan(&authorID, &path, &extension, &processing)
			if authorID == 0 {
				return images.NewGetImagesIDNotFound()
			}

			if authorID != userID.ID {
				entryID := tx.QueryInt64("SELECT entry_id FROM entry_images WHERE image_id = $1", params.ID)
				if entryID == 0 {
					return images.NewGetImagesIDForbidden()
				}

				allowed := utils.CanViewEntry(tx, userID, entryID)
				if !allowed {
					return images.NewGetImagesIDForbidden()
				}
			}

			img := &models.Image{
				ID: params.ID,
				Author: &models.User{
					ID: authorID,
				},
				IsAnimated: extension == imageExtensionGif && !processing,
				Processing: processing,
			}

			if processing {
				setProcessingImage(mi, img)
				return images.NewGetImagesIDOK().WithPayload(img)
			}

			var width, height int64
			var size string
			tx.Query(`
				SELECT width, height, (SELECT type FROM size WHERE size.id = image_sizes.size)
				FROM image_sizes
				WHERE image_sizes.image_id = $1
			`, params.ID)

			filePath := path + "." + extension

			var previewPath string
			if img.IsAnimated {
				previewPath = path + ".jpg"
			}

			for tx.Scan(&width, &height, &size) {
				switch size {
				case "thumbnail":
					img.Thumbnail = &models.ImageSize{
						Height: height,
						Width:  width,
						URL:    mi.BaseURL() + "albums/thumbnails/" + filePath,
					}
					if img.IsAnimated {
						img.Thumbnail.Preview = mi.BaseURL() + "albums/thumbnails/" + previewPath
					}
				case "small":
					img.Small = &models.ImageSize{
						Height: height,
						Width:  width,
						URL:    mi.BaseURL() + "albums/small/" + filePath,
					}
					if img.IsAnimated {
						img.Small.Preview = mi.BaseURL() + "albums/small/" + previewPath
					}
				case "medium":
					img.Medium = &models.ImageSize{
						Height: height,
						Width:  width,
						URL:    mi.BaseURL() + "albums/medium/" + filePath,
					}
					if img.IsAnimated {
						img.Medium.Preview = mi.BaseURL() + "albums/medium/" + previewPath
					}
				case "large":
					img.Large = &models.ImageSize{
						Height: height,
						Width:  width,
						URL:    mi.BaseURL() + "albums/large/" + filePath,
					}
					if img.IsAnimated {
						img.Large.Preview = mi.BaseURL() + "albums/large/" + previewPath
					}
				}
			}

			return images.NewGetImagesIDOK().WithPayload(img)
		})
	}
}

func NewImageDeleter(mi *MindwellImages) func(images.DeleteImagesIDParams, *models.UserID) middleware.Responder {
	return func(params images.DeleteImagesIDParams, userID *models.UserID) middleware.Responder {
		return utils.Transact(mi.DB(), func(tx *utils.AutoTx) middleware.Responder {
			authorID := tx.QueryInt64("SELECT user_id FROM images WHERE id = $1", params.ID)
			if authorID == 0 {
				return images.NewDeleteImagesIDNotFound()
			}

			if authorID != userID.ID {
				return images.NewDeleteImagesIDForbidden()
			}

			store := newImageStore(mi)
			mi.QueueAction(store, params.ID, ActionDelete)

			return images.NewDeleteImagesIDNoContent()
		})
	}
}
