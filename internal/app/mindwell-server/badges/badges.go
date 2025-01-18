package badges

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/me"
	"github.com/sevings/mindwell-server/restapi/operations/users"
	"github.com/sevings/mindwell-server/utils"
)

var baseURL string

// ConfigureAPI creates operations handlers
func ConfigureAPI(srv *utils.MindwellServer) {
	srv.API.MeGetMeBadgesHandler = me.GetMeBadgesHandlerFunc(newMyBadgesLoader(srv))
	srv.API.UsersGetUsersNameBadgesHandler = users.GetUsersNameBadgesHandlerFunc(newUserBadgesLoader(srv))

	baseURL = srv.ConfigString("images.base_url") + "badges/"
}

func badgesQuery(limit int64) *sqlf.Stmt {
	return sqlf.Select("badges.code, badges.title, badges.description, badges.icon").
		Select("extract(epoch from user_badges.given_at) as given_at").
		From("badges").
		Join("user_badges", "badges.id = user_badges.badge_id").
		OrderBy("user_count ASC").
		Limit(limit)
}

func myBadgesQuery(userID *models.UserID, limit int64) *sqlf.Stmt {
	return badgesQuery(limit).
		Where("user_id = ?", userID.ID)
}

func tlogBadgesQuery(tlog string, limit int64) *sqlf.Stmt {
	return badgesQuery(limit).
		Join("users", "users.id = user_badges.user_id").
		Where("lower(users.name) = lower(?)", tlog)
}

func idBadgeQuery(id int64) *sqlf.Stmt {
	return badgesQuery(1).
		Where("badges.id = ?", id)
}

func loadBadges(tx *utils.AutoTx) *models.BadgeList {
	list := &models.BadgeList{}

	var code, title, desc, icon string
	var at float64
	for tx.Scan(
		&code, &title, &desc, &icon,
		&at) {

		icon = baseURL + icon

		rw := &models.Badge{
			Code:        code,
			Title:       title,
			Description: desc,
			Icon:        icon,
			GivenAt:     at,
		}
		list.Data = append(list.Data, rw)
	}

	return list
}

func LoadBadgeByID(tx *utils.AutoTx, id int64) *models.Badge {
	q := idBadgeQuery(id)
	tx.QueryStmt(q)
	list := loadBadges(tx)
	if len(list.Data) == 0 {
		return nil
	}

	return list.Data[0]
}

func newMyBadgesLoader(srv *utils.MindwellServer) func(me.GetMeBadgesParams, *models.UserID) middleware.Responder {
	return func(params me.GetMeBadgesParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			q := myBadgesQuery(userID, *params.Limit)
			tx.QueryStmt(q)
			if tx.HasQueryError() {
				return me.NewGetMeBadgesOK()
			}

			list := loadBadges(tx)
			return me.NewGetMeBadgesOK().WithPayload(list)
		})
	}
}

func newUserBadgesLoader(srv *utils.MindwellServer) func(users.GetUsersNameBadgesParams, *models.UserID) middleware.Responder {
	return func(params users.GetUsersNameBadgesParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			if !utils.CanViewTlogName(tx, userID, params.Name) {
				err := srv.StandardError("no_tlog")
				return users.NewGetUsersNameBadgesNotFound().WithPayload(err)
			}

			q := tlogBadgesQuery(params.Name, *params.Limit)
			tx.QueryStmt(q)
			if tx.HasQueryError() {
				return users.NewGetUsersNameBadgesNotFound()
			}

			list := loadBadges(tx)
			return users.NewGetUsersNameBadgesOK().WithPayload(list)
		})
	}
}
