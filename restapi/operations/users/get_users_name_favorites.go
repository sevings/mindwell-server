// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetUsersNameFavoritesHandlerFunc turns a function with the right signature into a get users name favorites handler
type GetUsersNameFavoritesHandlerFunc func(GetUsersNameFavoritesParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetUsersNameFavoritesHandlerFunc) Handle(params GetUsersNameFavoritesParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetUsersNameFavoritesHandler interface for that can handle valid get users name favorites params
type GetUsersNameFavoritesHandler interface {
	Handle(GetUsersNameFavoritesParams, *models.UserID) middleware.Responder
}

// NewGetUsersNameFavorites creates a new http.Handler for the get users name favorites operation
func NewGetUsersNameFavorites(ctx *middleware.Context, handler GetUsersNameFavoritesHandler) *GetUsersNameFavorites {
	return &GetUsersNameFavorites{Context: ctx, Handler: handler}
}

/*GetUsersNameFavorites swagger:route GET /users/{name}/favorites users getUsersNameFavorites

GetUsersNameFavorites get users name favorites API

*/
type GetUsersNameFavorites struct {
	Context *middleware.Context
	Handler GetUsersNameFavoritesHandler
}

func (o *GetUsersNameFavorites) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetUsersNameFavoritesParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
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