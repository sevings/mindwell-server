package users

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/users"
)

func newUserLoader(db *sql.DB) func(users.GetUsersIDParams, *models.UserID) middleware.Responder {
	return func(params users.GetUsersIDParams, userID *models.UserID) middleware.Responder {
		const q = profileQuery + "WHERE long_users.id = $1"
		return loadProfile(db, q, userID, params.ID)
	}
}

const privacyQueryID = privacyQueryStart + "users.id = $1"
const usersQueryToID = usersQueryStart + "relations.to_id = $1 AND relations.from_id = short_users.id" + usersQueryEnd
const usersQueryFromID = usersQueryStart + "relations.from_id = $1 AND relations.to_id = short_users.id" + usersQueryEnd

func loadUsersRelatedToID(db *sql.DB, usersQuery, relation string,
	userID *models.UserID, args ...interface{}) middleware.Responder {
	return loadUsers(db, usersQuery, privacyQueryID, relationToIDQuery, loadUserQueryID, relation,
		userID, args...)
}

func newFollowersLoader(db *sql.DB) func(users.GetUsersIDFollowersParams, *models.UserID) middleware.Responder {
	return func(params users.GetUsersIDFollowersParams, userID *models.UserID) middleware.Responder {
		return loadUsersRelatedToID(db, usersQueryToID, models.UserListRelationFollowers,
			userID,
			params.ID, models.RelationshipRelationFollowed, *params.Limit, *params.Skip)
	}
}

func newFollowingsLoader(db *sql.DB) func(users.GetUsersIDFollowingsParams, *models.UserID) middleware.Responder {
	return func(params users.GetUsersIDFollowingsParams, userID *models.UserID) middleware.Responder {
		return loadUsersRelatedToID(db, usersQueryFromID, models.UserListRelationFollowings,
			userID,
			params.ID, models.RelationshipRelationFollowed, *params.Limit, *params.Skip)
	}
}

func newInvitedLoader(db *sql.DB) func(users.GetUsersIDInvitedParams, *models.UserID) middleware.Responder {
	return func(params users.GetUsersIDInvitedParams, userID *models.UserID) middleware.Responder {
		return loadUsersRelatedToID(db, invitedUsersQuery, models.UserListRelationInvited,
			userID, params.ID, *params.Limit, *params.Skip)
	}
}
