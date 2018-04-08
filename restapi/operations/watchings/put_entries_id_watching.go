// Code generated by go-swagger; DO NOT EDIT.

package watchings

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// PutEntriesIDWatchingHandlerFunc turns a function with the right signature into a put entries ID watching handler
type PutEntriesIDWatchingHandlerFunc func(PutEntriesIDWatchingParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn PutEntriesIDWatchingHandlerFunc) Handle(params PutEntriesIDWatchingParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// PutEntriesIDWatchingHandler interface for that can handle valid put entries ID watching params
type PutEntriesIDWatchingHandler interface {
	Handle(PutEntriesIDWatchingParams, *models.UserID) middleware.Responder
}

// NewPutEntriesIDWatching creates a new http.Handler for the put entries ID watching operation
func NewPutEntriesIDWatching(ctx *middleware.Context, handler PutEntriesIDWatchingHandler) *PutEntriesIDWatching {
	return &PutEntriesIDWatching{Context: ctx, Handler: handler}
}

/*PutEntriesIDWatching swagger:route PUT /entries/{id}/watching watchings putEntriesIdWatching

PutEntriesIDWatching put entries ID watching API

*/
type PutEntriesIDWatching struct {
	Context *middleware.Context
	Handler PutEntriesIDWatchingHandler
}

func (o *PutEntriesIDWatching) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPutEntriesIDWatchingParams()

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
