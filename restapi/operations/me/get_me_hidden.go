// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetMeHiddenHandlerFunc turns a function with the right signature into a get me hidden handler
type GetMeHiddenHandlerFunc func(GetMeHiddenParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetMeHiddenHandlerFunc) Handle(params GetMeHiddenParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetMeHiddenHandler interface for that can handle valid get me hidden params
type GetMeHiddenHandler interface {
	Handle(GetMeHiddenParams, *models.UserID) middleware.Responder
}

// NewGetMeHidden creates a new http.Handler for the get me hidden operation
func NewGetMeHidden(ctx *middleware.Context, handler GetMeHiddenHandler) *GetMeHidden {
	return &GetMeHidden{Context: ctx, Handler: handler}
}

/* GetMeHidden swagger:route GET /me/hidden me getMeHidden

GetMeHidden get me hidden API

*/
type GetMeHidden struct {
	Context *middleware.Context
	Handler GetMeHiddenHandler
}

func (o *GetMeHidden) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetMeHiddenParams()
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
