// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/sevings/mindwell-server/models"
)

// GetMeCalendarHandlerFunc turns a function with the right signature into a get me calendar handler
type GetMeCalendarHandlerFunc func(GetMeCalendarParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetMeCalendarHandlerFunc) Handle(params GetMeCalendarParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetMeCalendarHandler interface for that can handle valid get me calendar params
type GetMeCalendarHandler interface {
	Handle(GetMeCalendarParams, *models.UserID) middleware.Responder
}

// NewGetMeCalendar creates a new http.Handler for the get me calendar operation
func NewGetMeCalendar(ctx *middleware.Context, handler GetMeCalendarHandler) *GetMeCalendar {
	return &GetMeCalendar{Context: ctx, Handler: handler}
}

/*GetMeCalendar swagger:route GET /me/calendar me getMeCalendar

GetMeCalendar get me calendar API

*/
type GetMeCalendar struct {
	Context *middleware.Context
	Handler GetMeCalendarHandler
}

func (o *GetMeCalendar) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetMeCalendarParams()

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