// Code generated by go-swagger; DO NOT EDIT.

package chats

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/sevings/mindwell-server/models"
)

// DeleteMessagesIDHandlerFunc turns a function with the right signature into a delete messages ID handler
type DeleteMessagesIDHandlerFunc func(DeleteMessagesIDParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteMessagesIDHandlerFunc) Handle(params DeleteMessagesIDParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// DeleteMessagesIDHandler interface for that can handle valid delete messages ID params
type DeleteMessagesIDHandler interface {
	Handle(DeleteMessagesIDParams, *models.UserID) middleware.Responder
}

// NewDeleteMessagesID creates a new http.Handler for the delete messages ID operation
func NewDeleteMessagesID(ctx *middleware.Context, handler DeleteMessagesIDHandler) *DeleteMessagesID {
	return &DeleteMessagesID{Context: ctx, Handler: handler}
}

/*DeleteMessagesID swagger:route DELETE /messages/{id} chats deleteMessagesId

DeleteMessagesID delete messages ID API

*/
type DeleteMessagesID struct {
	Context *middleware.Context
	Handler DeleteMessagesIDHandler
}

func (o *DeleteMessagesID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteMessagesIDParams()

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