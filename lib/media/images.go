package media

import (
	"strconv"

	cache "github.com/patrickmn/go-cache"
	"github.com/sevings/mindwell-server/lib/database"
	"github.com/sevings/mindwell-server/lib/userutil"
	"github.com/sevings/mindwell-server/models"
)

// ConfigProvider provides configuration access for image URLs
type ConfigProvider interface {
	ConfigString(key string) string
}

func setProcessingImage(baseURL string, img *models.Image) {
	img.Thumbnail = &models.ImageSize{
		Width:  100,
		Height: 100,
		URL:    baseURL + "albums/thumbnails/processing.png",
	}

	img.Small = &models.ImageSize{
		Width:  480,
		Height: 300,
		URL:    baseURL + "albums/small/processing.png",
	}

	img.Medium = &models.ImageSize{
		Width:  960,
		Height: 600,
		URL:    baseURL + "albums/medium/processing.png",
	}

	img.Large = &models.ImageSize{
		Width:  1440,
		Height: 900,
		URL:    baseURL + "albums/large/processing.png",
	}
}

func loadImageNotCached(srv ConfigProvider, tx *database.AutoTx, imageID int64) *models.Image {
	baseURL := srv.ConfigString("images.base_url")

	var authorID int64
	var path, extension, previewExt string
	var processing bool

	const loadQuery = `SELECT user_id, path, extension, preview_extension, processing FROM images WHERE id = $1`

	tx.Query(loadQuery, imageID).
		Scan(&authorID, &path, &extension, &previewExt, &processing)
	if authorID == 0 {
		return nil
	}

	img := &models.Image{
		ID: imageID,
		Author: &models.User{
			ID: authorID,
		},
		IsAnimated: previewExt != "" && !processing,
		Processing: processing,
	}

	if processing {
		setProcessingImage(baseURL, img)
		return img
	}

	filePath := path + "." + extension

	var previewPath string
	if img.IsAnimated {
		previewPath = path + "." + previewExt
	}

	var width, height int64
	var size string
	tx.Query(`
		SELECT width, height, (SELECT type FROM size WHERE size.id = image_sizes.size)
		FROM image_sizes
		WHERE image_sizes.image_id = $1
	`, imageID)

	for tx.Scan(&width, &height, &size) {
		switch size {
		case "thumbnail":
			img.Thumbnail = &models.ImageSize{
				Height: height,
				Width:  width,
				URL:    baseURL + "albums/thumbnails/" + filePath,
			}
			if img.IsAnimated {
				img.Thumbnail.Preview = baseURL + "albums/thumbnails/" + previewPath
			}
		case "small":
			img.Small = &models.ImageSize{
				Height: height,
				Width:  width,
				URL:    baseURL + "albums/small/" + filePath,
			}
			if img.IsAnimated {
				img.Small.Preview = baseURL + "albums/small/" + previewPath
			}
		case "medium":
			img.Medium = &models.ImageSize{
				Height: height,
				Width:  width,
				URL:    baseURL + "albums/medium/" + filePath,
			}
			if img.IsAnimated {
				img.Medium.Preview = baseURL + "albums/medium/" + previewPath
			}
		case "large":
			img.Large = &models.ImageSize{
				Height: height,
				Width:  width,
				URL:    baseURL + "albums/large/" + filePath,
			}
			if img.IsAnimated {
				img.Large.Preview = baseURL + "albums/large/" + previewPath
			}
		}
	}

	return img
}

// LoadImage loads an image by ID with caching support
func LoadImage(srv ConfigProvider, imgCache *cache.Cache, tx *database.AutoTx, imageID int64) *models.Image {
	var img models.Image
	idStr := strconv.FormatInt(imageID, 10)
	oldImg, found := imgCache.Get(idStr)
	if found {
		img = *oldImg.(*models.Image)
	} else {
		newImg := loadImageNotCached(srv, tx, imageID)
		if newImg == nil || newImg.Processing {
			return newImg
		}
		if !newImg.Processing {
			imgCache.Set(idStr, newImg, cache.DefaultExpiration)
		}

		img = *newImg
	}

	img.Author = userutil.LoadUser(tx, img.Author.ID)

	return &img
}
