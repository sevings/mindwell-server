// Code generated by go-swagger; DO NOT EDIT.

package relations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// DeleteRelationsToNameHandlerFunc turns a function with the right signature into a delete relations to name handler
type DeleteRelationsToNameHandlerFunc func(DeleteRelationsToNameParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteRelationsToNameHandlerFunc) Handle(params DeleteRelationsToNameParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// DeleteRelationsToNameHandler interface for that can handle valid delete relations to name params
type DeleteRelationsToNameHandler interface {
	Handle(DeleteRelationsToNameParams, *models.UserID) middleware.Responder
}

// NewDeleteRelationsToName creates a new http.Handler for the delete relations to name operation
func NewDeleteRelationsToName(ctx *middleware.Context, handler DeleteRelationsToNameHandler) *DeleteRelationsToName {
	return &DeleteRelationsToName{Context: ctx, Handler: handler}
}

/* DeleteRelationsToName swagger:route DELETE /relations/to/{name} relations deleteRelationsToName

DeleteRelationsToName delete relations to name API

*/
type DeleteRelationsToName struct {
	Context *middleware.Context
	Handler DeleteRelationsToNameHandler
}

func (o *DeleteRelationsToName) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteRelationsToNameParams()
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
