// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/sevings/yummy-server/gen/models"
)

// DeleteEntriesIDHandlerFunc turns a function with the right signature into a delete entries ID handler
type DeleteEntriesIDHandlerFunc func(DeleteEntriesIDParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteEntriesIDHandlerFunc) Handle(params DeleteEntriesIDParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// DeleteEntriesIDHandler interface for that can handle valid delete entries ID params
type DeleteEntriesIDHandler interface {
	Handle(DeleteEntriesIDParams, *models.UserID) middleware.Responder
}

// NewDeleteEntriesID creates a new http.Handler for the delete entries ID operation
func NewDeleteEntriesID(ctx *middleware.Context, handler DeleteEntriesIDHandler) *DeleteEntriesID {
	return &DeleteEntriesID{Context: ctx, Handler: handler}
}

/*DeleteEntriesID swagger:route DELETE /entries/{id} entries deleteEntriesId

DeleteEntriesID delete entries ID API

*/
type DeleteEntriesID struct {
	Context *middleware.Context
	Handler DeleteEntriesIDHandler
}

func (o *DeleteEntriesID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteEntriesIDParams()

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
