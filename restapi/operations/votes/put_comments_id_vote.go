// Code generated by go-swagger; DO NOT EDIT.

package votes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/sevings/mindwell-server/models"
)

// PutCommentsIDVoteHandlerFunc turns a function with the right signature into a put comments ID vote handler
type PutCommentsIDVoteHandlerFunc func(PutCommentsIDVoteParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn PutCommentsIDVoteHandlerFunc) Handle(params PutCommentsIDVoteParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// PutCommentsIDVoteHandler interface for that can handle valid put comments ID vote params
type PutCommentsIDVoteHandler interface {
	Handle(PutCommentsIDVoteParams, *models.UserID) middleware.Responder
}

// NewPutCommentsIDVote creates a new http.Handler for the put comments ID vote operation
func NewPutCommentsIDVote(ctx *middleware.Context, handler PutCommentsIDVoteHandler) *PutCommentsIDVote {
	return &PutCommentsIDVote{Context: ctx, Handler: handler}
}

/*PutCommentsIDVote swagger:route PUT /comments/{id}/vote votes putCommentsIdVote

PutCommentsIDVote put comments ID vote API

*/
type PutCommentsIDVote struct {
	Context *middleware.Context
	Handler PutCommentsIDVoteHandler
}

func (o *PutCommentsIDVote) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPutCommentsIDVoteParams()

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
