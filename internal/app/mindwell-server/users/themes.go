package users

import (
	"fmt"
	"github.com/aofei/cameron"
	"github.com/go-openapi/runtime/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/themes"
	"github.com/sevings/mindwell-server/utils"
	"image/jpeg"
	"os"
	"strconv"
	"strings"
)

func loadTopThemes(srv *utils.MindwellServer, tx *utils.AutoTx, top string) []*models.Friend {
	query := usersQuerySelect + ", 0 "

	if top == "rank" {
		query += "FROM users, gender, user_privacy WHERE creator_id IS NOT NULL" + usersQueryJoins + "ORDER BY rank ASC"
		query += " LIMIT 50"
		tx.Query(query)
	} else if top == "new" {
		query += "FROM users, gender, user_privacy WHERE creator_id IS NOT NULL" + usersQueryJoins + " ORDER BY created_at DESC"
		query += " LIMIT 50"
		tx.Query(query)
	} else {
		fmt.Printf("Unknown themes top: %s\n", top)
		return nil
	}

	list, _, _ := loadUserList(srv, tx, false)
	return list
}

func newThemesLoader(srv *utils.MindwellServer) func(themes.GetThemesParams, *models.UserID) middleware.Responder {
	return func(params themes.GetThemesParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			body := &themes.GetThemesOKBody{}

			if params.Query != nil {
				body.Themes = searchUsers(srv, tx, "creator_id IS NOT NULL", *params.Query)
				body.Query = *params.Query
			} else {
				body.Themes = loadTopThemes(srv, tx, *params.Top)
				body.Top = *params.Top
			}

			return themes.NewGetThemesOK().WithPayload(body)
		})
	}
}

func newThemeLoader(srv *utils.MindwellServer) func(themes.GetThemesNameParams, *models.UserID) middleware.Responder {
	return func(params themes.GetThemesNameParams, userID *models.UserID) middleware.Responder {
		const query = profileQuery + "WHERE lower(users.name) = lower($1) AND users.creator_id IS NOT NULL"

		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			theme := loadUserProfile(srv, tx, query, userID, params.Name)
			if theme.ID == 0 {
				return themes.NewGetThemesNameNotFound()
			}

			return themes.NewGetThemesNameOK().WithPayload(theme)
		})
	}
}

func removeInvite(tx *utils.AutoTx, userID int64) bool {
	const idQ = `
		SELECT id
		FROM invites
		WHERE referrer_id = $1
		LIMIT 1
`

	inviteID := tx.QueryInt64(idQ, userID)
	if inviteID == 0 {
		return false
	}

	const q = `
		DELETE FROM invites
		WHERE id = $1
		`

	tx.Exec(q, inviteID)

	return tx.RowsAffected() == 1
}

func saveAvatar(srv *utils.MindwellServer, size, blockSize int, name, folder, fileName string) {
	path := srv.ImagesFolder() + "avatars/" + strconv.Itoa(size) + "/" + folder
	err := os.MkdirAll(path, 0777)
	if err != nil {
		srv.LogApi().Error(err.Error())
	}

	f, err := os.Create(path + fileName)
	if err != nil {
		srv.LogApi().Error(err.Error())
	}
	defer f.Close()

	img := cameron.Identicon([]byte(name), size, blockSize)
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 85})
	if err != nil {
		srv.LogApi().Error(err.Error())
	}
}

func generateAvatar(srv *utils.MindwellServer, name string) string {
	folder := name[:1] + "/"
	fileName := utils.GenerateString(5) + ".jpg"

	saveAvatar(srv, 124, 15, name, folder, fileName)
	saveAvatar(srv, 92, 11, name, folder, fileName)
	saveAvatar(srv, 42, 5, name, folder, fileName)

	return folder + fileName
}

func createTheme(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, name, showName string) int64 {
	const rankQ = "SELECT COUNT(*) + 1 FROM users WHERE creator_id IS NOT NULL AND karma >= 0"
	rank := tx.QueryInt64(rankQ)

	avatar := generateAvatar(srv, name)

	const q = `
		INSERT INTO users 
		(name, show_name, email, password_hash, avatar, creator_id, rank)
		values($1, $2, $1, '', $3, $4, $5)
		RETURNING id`

	user := tx.QueryInt64(q, name, showName, avatar, userID.ID, rank)
	return user
}

func newThemeCreator(srv *utils.MindwellServer) func(themes.PostThemesParams, *models.UserID) middleware.Responder {
	return func(params themes.PostThemesParams, userID *models.UserID) middleware.Responder {
		if userID.Ban.Invite {
			err := srv.NewError(&i18n.Message{ID: "cant_create_theme", Other: "You are not allowed to create themes."})
			return themes.NewPostThemesForbidden().WithPayload(err)
		}

		name := strings.TrimSpace(params.Name)
		showName := strings.TrimSpace(params.ShowName)

		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			const idQ = "SELECT id FROM users WHERE lower(name) = lower($1)"
			oldID := tx.QueryInt64(idQ, name)
			if oldID > 0 {
				err := srv.NewError(&i18n.Message{ID: "theme_name_is_not_free", Other: "Theme name is not free."})
				return themes.NewPostThemesBadRequest().WithPayload(err)
			}

			if userID.Authority != models.UserIDAuthorityAdmin {
				if ok := removeInvite(tx, userID.ID); !ok {
					err := srv.StandardError("invalid_invite")
					return themes.NewPostThemesBadRequest().WithPayload(err)
				}
			}

			const query = profileQuery + "WHERE users.id = $1 AND users.creator_id IS NOT NULL"
			themeID := createTheme(srv, tx, userID, name, showName)
			theme := loadUserProfile(srv, tx, query, userID, themeID)
			if theme.ID == 0 {
				err := srv.NewError(nil)
				return themes.NewPostThemesBadRequest().WithPayload(err)
			}

			return themes.NewPostThemesOK().WithPayload(theme)
		})
	}
}

func checkThemeAdmin(tx *utils.AutoTx, userID *models.UserID, name string) bool {
	const q = `SELECT creator_id = $1 FROM users WHERE lower(name) = lower($2)`

	return tx.QueryBool(q, userID.ID, name)
}

func editThemeProfile(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, params themes.PutThemesNameParams) *models.Profile {
	name := params.Name

	if params.IsDaylog != nil {
		const q = "update users set is_daylog = $2 where lower(name) = lower($1)"
		tx.Exec(q, name, *params.IsDaylog)
	}

	const q = "update users set privacy = (select id from user_privacy where type = $2), show_name = $3 where lower(name) = lower($1)"
	showName := strings.TrimSpace(params.ShowName)
	tx.Exec(q, name, params.Privacy, showName)

	if params.ShowInTops != nil {
		const q = "update users set show_in_tops = $2 where lower(name) = lower($1)"
		tx.Exec(q, name, *params.ShowInTops)
	}

	if params.Title != nil {
		const q = "update users set title = $2 where lower(name) = lower($1)"
		title := strings.TrimSpace(*params.Title)
		tx.Exec(q, name, title)
	}

	const loadQuery = profileQuery + "WHERE lower(users.name) = lower($1)"
	return loadUserProfile(srv, tx, loadQuery, userID, name)
}

func newThemeEditor(srv *utils.MindwellServer) func(themes.PutThemesNameParams, *models.UserID) middleware.Responder {
	return func(params themes.PutThemesNameParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			if !checkThemeAdmin(tx, userID, params.Name) {
				err := srv.StandardError("no_theme")
				return themes.NewPutThemesNameForbidden().WithPayload(err)
			}

			theme := editThemeProfile(srv, tx, userID, params)
			if theme.ID == 0 {
				err := srv.StandardError("no_theme")
				return themes.NewPutThemesNameForbidden().WithPayload(err)
			}

			return themes.NewPutThemesNameOK().WithPayload(theme)
		})
	}
}

func newThemeFollowersLoader(srv *utils.MindwellServer) func(themes.GetThemesNameFollowersParams, *models.UserID) middleware.Responder {
	return func(params themes.GetThemesNameFollowersParams, userID *models.UserID) middleware.Responder {
		list := loadRelatedUsers(srv, userID, usersQueryToName, usersToNameQueryWhere,
			models.RelationshipRelationFollowed, params.Name, models.FriendListRelationFollowers,
			*params.After, *params.Before, *params.Limit)

		if list == nil {
			err := srv.StandardError("no_theme")
			return themes.NewGetThemesNameFollowersNotFound().WithPayload(err)
		}

		return themes.NewGetThemesNameFollowersOK().WithPayload(list)
	}
}
