// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetUsersNameFollowersHandlerFunc turns a function with the right signature into a get users name followers handler
type GetUsersNameFollowersHandlerFunc func(GetUsersNameFollowersParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetUsersNameFollowersHandlerFunc) Handle(params GetUsersNameFollowersParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetUsersNameFollowersHandler interface for that can handle valid get users name followers params
type GetUsersNameFollowersHandler interface {
	Handle(GetUsersNameFollowersParams, *models.UserID) middleware.Responder
}

// NewGetUsersNameFollowers creates a new http.Handler for the get users name followers operation
func NewGetUsersNameFollowers(ctx *middleware.Context, handler GetUsersNameFollowersHandler) *GetUsersNameFollowers {
	return &GetUsersNameFollowers{Context: ctx, Handler: handler}
}

/* GetUsersNameFollowers swagger:route GET /users/{name}/followers users getUsersNameFollowers

GetUsersNameFollowers get users name followers API

*/
type GetUsersNameFollowers struct {
	Context *middleware.Context
	Handler GetUsersNameFollowersHandler
}

func (o *GetUsersNameFollowers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetUsersNameFollowersParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal *models.UserID
	if uprinc != nil {
		principal = uprinc.(*models.UserID) // this is really a models.UserID, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
