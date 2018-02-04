// Code generated by go-swagger; DO NOT EDIT.

package relations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/sevings/yummy-server/models"
)

// DeleteRelationsFromIDHandlerFunc turns a function with the right signature into a delete relations from ID handler
type DeleteRelationsFromIDHandlerFunc func(DeleteRelationsFromIDParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteRelationsFromIDHandlerFunc) Handle(params DeleteRelationsFromIDParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// DeleteRelationsFromIDHandler interface for that can handle valid delete relations from ID params
type DeleteRelationsFromIDHandler interface {
	Handle(DeleteRelationsFromIDParams, *models.UserID) middleware.Responder
}

// NewDeleteRelationsFromID creates a new http.Handler for the delete relations from ID operation
func NewDeleteRelationsFromID(ctx *middleware.Context, handler DeleteRelationsFromIDHandler) *DeleteRelationsFromID {
	return &DeleteRelationsFromID{Context: ctx, Handler: handler}
}

/*DeleteRelationsFromID swagger:route DELETE /relations/from/{id} relations deleteRelationsFromId

cancel following request or unsubscribe the user

*/
type DeleteRelationsFromID struct {
	Context *middleware.Context
	Handler DeleteRelationsFromIDHandler
}

func (o *DeleteRelationsFromID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteRelationsFromIDParams()

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