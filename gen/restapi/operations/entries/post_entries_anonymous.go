// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/sevings/yummy-server/gen/models"
)

// PostEntriesAnonymousHandlerFunc turns a function with the right signature into a post entries anonymous handler
type PostEntriesAnonymousHandlerFunc func(PostEntriesAnonymousParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn PostEntriesAnonymousHandlerFunc) Handle(params PostEntriesAnonymousParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// PostEntriesAnonymousHandler interface for that can handle valid post entries anonymous params
type PostEntriesAnonymousHandler interface {
	Handle(PostEntriesAnonymousParams, *models.UserID) middleware.Responder
}

// NewPostEntriesAnonymous creates a new http.Handler for the post entries anonymous operation
func NewPostEntriesAnonymous(ctx *middleware.Context, handler PostEntriesAnonymousHandler) *PostEntriesAnonymous {
	return &PostEntriesAnonymous{Context: ctx, Handler: handler}
}

/*PostEntriesAnonymous swagger:route POST /entries/anonymous entries postEntriesAnonymous

PostEntriesAnonymous post entries anonymous API

*/
type PostEntriesAnonymous struct {
	Context *middleware.Context
	Handler PostEntriesAnonymousHandler
}

func (o *PostEntriesAnonymous) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostEntriesAnonymousParams()

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
