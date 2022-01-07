// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetUsersNameImagesHandlerFunc turns a function with the right signature into a get users name images handler
type GetUsersNameImagesHandlerFunc func(GetUsersNameImagesParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetUsersNameImagesHandlerFunc) Handle(params GetUsersNameImagesParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetUsersNameImagesHandler interface for that can handle valid get users name images params
type GetUsersNameImagesHandler interface {
	Handle(GetUsersNameImagesParams, *models.UserID) middleware.Responder
}

// NewGetUsersNameImages creates a new http.Handler for the get users name images operation
func NewGetUsersNameImages(ctx *middleware.Context, handler GetUsersNameImagesHandler) *GetUsersNameImages {
	return &GetUsersNameImages{Context: ctx, Handler: handler}
}

/* GetUsersNameImages swagger:route GET /users/{name}/images users getUsersNameImages

GetUsersNameImages get users name images API

*/
type GetUsersNameImages struct {
	Context *middleware.Context
	Handler GetUsersNameImagesHandler
}

func (o *GetUsersNameImages) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetUsersNameImagesParams()
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
