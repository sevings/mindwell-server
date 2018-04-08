// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/sevings/mindwell-server/models"
)

// GetEntriesUsersMeHandlerFunc turns a function with the right signature into a get entries users me handler
type GetEntriesUsersMeHandlerFunc func(GetEntriesUsersMeParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetEntriesUsersMeHandlerFunc) Handle(params GetEntriesUsersMeParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetEntriesUsersMeHandler interface for that can handle valid get entries users me params
type GetEntriesUsersMeHandler interface {
	Handle(GetEntriesUsersMeParams, *models.UserID) middleware.Responder
}

// NewGetEntriesUsersMe creates a new http.Handler for the get entries users me operation
func NewGetEntriesUsersMe(ctx *middleware.Context, handler GetEntriesUsersMeHandler) *GetEntriesUsersMe {
	return &GetEntriesUsersMe{Context: ctx, Handler: handler}
}

/*GetEntriesUsersMe swagger:route GET /entries/users/me entries getEntriesUsersMe

GetEntriesUsersMe get entries users me API

*/
type GetEntriesUsersMe struct {
	Context *middleware.Context
	Handler GetEntriesUsersMeHandler
}

func (o *GetEntriesUsersMe) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetEntriesUsersMeParams()

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
