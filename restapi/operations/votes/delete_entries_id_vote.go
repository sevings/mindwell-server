// Code generated by go-swagger; DO NOT EDIT.

package votes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// DeleteEntriesIDVoteHandlerFunc turns a function with the right signature into a delete entries ID vote handler
type DeleteEntriesIDVoteHandlerFunc func(DeleteEntriesIDVoteParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteEntriesIDVoteHandlerFunc) Handle(params DeleteEntriesIDVoteParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// DeleteEntriesIDVoteHandler interface for that can handle valid delete entries ID vote params
type DeleteEntriesIDVoteHandler interface {
	Handle(DeleteEntriesIDVoteParams, *models.UserID) middleware.Responder
}

// NewDeleteEntriesIDVote creates a new http.Handler for the delete entries ID vote operation
func NewDeleteEntriesIDVote(ctx *middleware.Context, handler DeleteEntriesIDVoteHandler) *DeleteEntriesIDVote {
	return &DeleteEntriesIDVote{Context: ctx, Handler: handler}
}

/* DeleteEntriesIDVote swagger:route DELETE /entries/{id}/vote votes deleteEntriesIdVote

DeleteEntriesIDVote delete entries ID vote API

*/
type DeleteEntriesIDVote struct {
	Context *middleware.Context
	Handler DeleteEntriesIDVoteHandler
}

func (o *DeleteEntriesIDVote) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteEntriesIDVoteParams()
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
