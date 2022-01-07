// Code generated by go-swagger; DO NOT EDIT.

package votes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetEntriesIDVoteHandlerFunc turns a function with the right signature into a get entries ID vote handler
type GetEntriesIDVoteHandlerFunc func(GetEntriesIDVoteParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetEntriesIDVoteHandlerFunc) Handle(params GetEntriesIDVoteParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetEntriesIDVoteHandler interface for that can handle valid get entries ID vote params
type GetEntriesIDVoteHandler interface {
	Handle(GetEntriesIDVoteParams, *models.UserID) middleware.Responder
}

// NewGetEntriesIDVote creates a new http.Handler for the get entries ID vote operation
func NewGetEntriesIDVote(ctx *middleware.Context, handler GetEntriesIDVoteHandler) *GetEntriesIDVote {
	return &GetEntriesIDVote{Context: ctx, Handler: handler}
}

/* GetEntriesIDVote swagger:route GET /entries/{id}/vote votes getEntriesIdVote

GetEntriesIDVote get entries ID vote API

*/
type GetEntriesIDVote struct {
	Context *middleware.Context
	Handler GetEntriesIDVoteHandler
}

func (o *GetEntriesIDVote) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetEntriesIDVoteParams()
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
