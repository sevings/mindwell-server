// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetEntriesRandomHandlerFunc turns a function with the right signature into a get entries random handler
type GetEntriesRandomHandlerFunc func(GetEntriesRandomParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetEntriesRandomHandlerFunc) Handle(params GetEntriesRandomParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetEntriesRandomHandler interface for that can handle valid get entries random params
type GetEntriesRandomHandler interface {
	Handle(GetEntriesRandomParams, *models.UserID) middleware.Responder
}

// NewGetEntriesRandom creates a new http.Handler for the get entries random operation
func NewGetEntriesRandom(ctx *middleware.Context, handler GetEntriesRandomHandler) *GetEntriesRandom {
	return &GetEntriesRandom{Context: ctx, Handler: handler}
}

/* GetEntriesRandom swagger:route GET /entries/random entries getEntriesRandom

GetEntriesRandom get entries random API

*/
type GetEntriesRandom struct {
	Context *middleware.Context
	Handler GetEntriesRandomHandler
}

func (o *GetEntriesRandom) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetEntriesRandomParams()
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
