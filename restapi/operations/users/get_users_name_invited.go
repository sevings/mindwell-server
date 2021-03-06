// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetUsersNameInvitedHandlerFunc turns a function with the right signature into a get users name invited handler
type GetUsersNameInvitedHandlerFunc func(GetUsersNameInvitedParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetUsersNameInvitedHandlerFunc) Handle(params GetUsersNameInvitedParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetUsersNameInvitedHandler interface for that can handle valid get users name invited params
type GetUsersNameInvitedHandler interface {
	Handle(GetUsersNameInvitedParams, *models.UserID) middleware.Responder
}

// NewGetUsersNameInvited creates a new http.Handler for the get users name invited operation
func NewGetUsersNameInvited(ctx *middleware.Context, handler GetUsersNameInvitedHandler) *GetUsersNameInvited {
	return &GetUsersNameInvited{Context: ctx, Handler: handler}
}

/* GetUsersNameInvited swagger:route GET /users/{name}/invited users getUsersNameInvited

GetUsersNameInvited get users name invited API

*/
type GetUsersNameInvited struct {
	Context *middleware.Context
	Handler GetUsersNameInvitedHandler
}

func (o *GetUsersNameInvited) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetUsersNameInvitedParams()
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
