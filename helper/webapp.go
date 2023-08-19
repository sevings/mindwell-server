package helper

import (
	"github.com/sevings/mindwell-server/utils"
	"github.com/zpatrick/go-config"
	"log"
	"strconv"
	"strings"
)

type hashConfig struct {
	cfg *config.Config
}

func (hc hashConfig) ConfigString(key string) string {
	value, err := hc.cfg.String(key)
	if err != nil {
		log.Println(err)
	}

	return value
}

func CreateWebApp(tx *utils.AutoTx, cfg *config.Config) {
	hc := hashConfig{cfg: cfg}
	baseURL := hc.ConfigString("server.base_url")

	app := utils.CreateAppParameters{
		DevName:     "Mindwell",
		Type:        "private",
		Flow:        "password",
		RedirectUri: strings.Replace(baseURL, "://", "://auth.", 1) + "blank",
		Name:        "MindwellWeb",
		ShowName:    "MindwellWeb",
		Platform:    "Web",
		Info:        "The official Mindwell web app.",
	}

	hash := utils.NewTokenHash(hashConfig{cfg: cfg})
	appID, secret, err := utils.CreateApp(tx, hash, app)
	if err != nil {
		log.Println(err)
		return
	}

	text := "Приложение создано."
	text += "\nclient id: " + strconv.FormatInt(appID, 10)
	text += "\nclient secret: " + secret
	text += "\napp name: " + app.Name
	text += "\napp show name: " + app.ShowName
	text += "\ndeveloper name: " + app.DevName
	text += "\nredirect uri: " + app.RedirectUri
	text += "\nflow: " + app.Flow
	text += "\nplatform: " + app.Platform
	text += "\ninfo: " + app.Info

	log.Println(text)
}
