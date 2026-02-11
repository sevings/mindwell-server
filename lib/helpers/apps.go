package helpers

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/sevings/mindwell-server/lib/auth"
	"github.com/sevings/mindwell-server/lib/database"
	"github.com/sevings/mindwell-server/lib/textutil"
)

type CreateAppParameters struct {
	DevName     string
	Type        string
	Flow        string
	RedirectUri string
	Name        string
	ShowName    string
	Platform    string
	Info        string
}

func CreateApp(tx *database.AutoTx, hash auth.TokenHash, app CreateAppParameters) (int64, string, error) {
	var secretHash []byte
	var secret string

	switch app.Type {
	case "private":
		secret = textutil.GenerateString(64)
		secretHash = hash.AppSecretHash(secret)
	case "public":
		secret = "(не сгенерирован)"
	default:
		return 0, "", fmt.Errorf("тип приложения может быть только public или private")
	}

	flow := 1
	switch app.Flow {
	case "code":
		flow += 2
	case "password":
		flow += 4
	default:
		return 0, "", fmt.Errorf("тип авторизации может быть только code или password")
	}

	uriRe := regexp.MustCompile(`^\w+://[^#]+$`)
	if !uriRe.MatchString(app.RedirectUri) {
		return 0, "", fmt.Errorf("redirect uri должен содержать схему и не содержать якорь")
	}

	app.Name = strings.TrimSpace(app.Name)
	nameRe := regexp.MustCompile(`^\w+$`)
	if !nameRe.MatchString(app.Name) {
		return 0, "", fmt.Errorf("название приложения может содержать только латинские буквы, цифры и знак подчеркивания")
	}

	app.DevName = strings.TrimSpace(app.DevName)
	dev, err := auth.LoadUserIDByName(tx, app.DevName)
	if err != nil {
		return 0, "", fmt.Errorf("пользователь %s не найден", app.DevName)
	}
	if dev.Ban.Account {
		return 0, "", auth.ErrUnauthorized
	}

	appID := int64(rand.Int31())

	const query = `
INSERT INTO apps(id, secret_hash, redirect_uri, developer_id, flow, name, show_name, platform, info)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

	tx.Exec(query, appID, secretHash, app.RedirectUri, dev.ID, flow,
		app.Name, app.ShowName, app.Platform, app.Info)
	if tx.Error() != nil {
		return 0, "", tx.Error()
	}

	return appID, secret, nil
}
