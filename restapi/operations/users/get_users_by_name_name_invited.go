// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetUsersByNameNameInvitedHandlerFunc turns a function with the right signature into a get users by name name invited handler
type GetUsersByNameNameInvitedHandlerFunc func(GetUsersByNameNameInvitedParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetUsersByNameNameInvitedHandlerFunc) Handle(params GetUsersByNameNameInvitedParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetUsersByNameNameInvitedHandler interface for that can handle valid get users by name name invited params
type GetUsersByNameNameInvitedHandler interface {
	Handle(GetUsersByNameNameInvitedParams, *models.UserID) middleware.Responder
}

// NewGetUsersByNameNameInvited creates a new http.Handler for the get users by name name invited operation
func NewGetUsersByNameNameInvited(ctx *middleware.Context, handler GetUsersByNameNameInvitedHandler) *GetUsersByNameNameInvited {
	return &GetUsersByNameNameInvited{Context: ctx, Handler: handler}
}

/*GetUsersByNameNameInvited swagger:route GET /users/byName/{name}/invited users getUsersByNameNameInvited

GetUsersByNameNameInvited get users by name name invited API

*/
type GetUsersByNameNameInvited struct {
	Context *middleware.Context
	Handler GetUsersByNameNameInvitedHandler
}

func (o *GetUsersByNameNameInvited) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetUsersByNameNameInvitedParams()

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
