// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetMeFollowersHandlerFunc turns a function with the right signature into a get me followers handler
type GetMeFollowersHandlerFunc func(GetMeFollowersParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetMeFollowersHandlerFunc) Handle(params GetMeFollowersParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetMeFollowersHandler interface for that can handle valid get me followers params
type GetMeFollowersHandler interface {
	Handle(GetMeFollowersParams, *models.UserID) middleware.Responder
}

// NewGetMeFollowers creates a new http.Handler for the get me followers operation
func NewGetMeFollowers(ctx *middleware.Context, handler GetMeFollowersHandler) *GetMeFollowers {
	return &GetMeFollowers{Context: ctx, Handler: handler}
}

/* GetMeFollowers swagger:route GET /me/followers me getMeFollowers

GetMeFollowers get me followers API

*/
type GetMeFollowers struct {
	Context *middleware.Context
	Handler GetMeFollowersHandler
}

func (o *GetMeFollowers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetMeFollowersParams()
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
