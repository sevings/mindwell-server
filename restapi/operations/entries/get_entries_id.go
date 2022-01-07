// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetEntriesIDHandlerFunc turns a function with the right signature into a get entries ID handler
type GetEntriesIDHandlerFunc func(GetEntriesIDParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetEntriesIDHandlerFunc) Handle(params GetEntriesIDParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetEntriesIDHandler interface for that can handle valid get entries ID params
type GetEntriesIDHandler interface {
	Handle(GetEntriesIDParams, *models.UserID) middleware.Responder
}

// NewGetEntriesID creates a new http.Handler for the get entries ID operation
func NewGetEntriesID(ctx *middleware.Context, handler GetEntriesIDHandler) *GetEntriesID {
	return &GetEntriesID{Context: ctx, Handler: handler}
}

/* GetEntriesID swagger:route GET /entries/{id} entries getEntriesId

GetEntriesID get entries ID API

*/
type GetEntriesID struct {
	Context *middleware.Context
	Handler GetEntriesIDHandler
}

func (o *GetEntriesID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetEntriesIDParams()
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
