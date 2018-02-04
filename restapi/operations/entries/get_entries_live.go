// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/sevings/yummy-server/models"
)

// GetEntriesLiveHandlerFunc turns a function with the right signature into a get entries live handler
type GetEntriesLiveHandlerFunc func(GetEntriesLiveParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetEntriesLiveHandlerFunc) Handle(params GetEntriesLiveParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetEntriesLiveHandler interface for that can handle valid get entries live params
type GetEntriesLiveHandler interface {
	Handle(GetEntriesLiveParams, *models.UserID) middleware.Responder
}

// NewGetEntriesLive creates a new http.Handler for the get entries live operation
func NewGetEntriesLive(ctx *middleware.Context, handler GetEntriesLiveHandler) *GetEntriesLive {
	return &GetEntriesLive{Context: ctx, Handler: handler}
}

/*GetEntriesLive swagger:route GET /entries/live entries getEntriesLive

GetEntriesLive get entries live API

*/
type GetEntriesLive struct {
	Context *middleware.Context
	Handler GetEntriesLiveHandler
}

func (o *GetEntriesLive) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetEntriesLiveParams()

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