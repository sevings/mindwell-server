// Code generated by go-swagger; DO NOT EDIT.

package votes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/sevings/mindwell-server/models"
)

// PutEntriesIDVoteHandlerFunc turns a function with the right signature into a put entries ID vote handler
type PutEntriesIDVoteHandlerFunc func(PutEntriesIDVoteParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn PutEntriesIDVoteHandlerFunc) Handle(params PutEntriesIDVoteParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// PutEntriesIDVoteHandler interface for that can handle valid put entries ID vote params
type PutEntriesIDVoteHandler interface {
	Handle(PutEntriesIDVoteParams, *models.UserID) middleware.Responder
}

// NewPutEntriesIDVote creates a new http.Handler for the put entries ID vote operation
func NewPutEntriesIDVote(ctx *middleware.Context, handler PutEntriesIDVoteHandler) *PutEntriesIDVote {
	return &PutEntriesIDVote{Context: ctx, Handler: handler}
}

/*PutEntriesIDVote swagger:route PUT /entries/{id}/vote votes putEntriesIdVote

PutEntriesIDVote put entries ID vote API

*/
type PutEntriesIDVote struct {
	Context *middleware.Context
	Handler PutEntriesIDVoteHandler
}

func (o *PutEntriesIDVote) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPutEntriesIDVoteParams()

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
