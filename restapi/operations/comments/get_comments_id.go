// Code generated by go-swagger; DO NOT EDIT.

package comments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/sevings/mindwell-server/models"
)

// GetCommentsIDHandlerFunc turns a function with the right signature into a get comments ID handler
type GetCommentsIDHandlerFunc func(GetCommentsIDParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetCommentsIDHandlerFunc) Handle(params GetCommentsIDParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetCommentsIDHandler interface for that can handle valid get comments ID params
type GetCommentsIDHandler interface {
	Handle(GetCommentsIDParams, *models.UserID) middleware.Responder
}

// NewGetCommentsID creates a new http.Handler for the get comments ID operation
func NewGetCommentsID(ctx *middleware.Context, handler GetCommentsIDHandler) *GetCommentsID {
	return &GetCommentsID{Context: ctx, Handler: handler}
}

/*GetCommentsID swagger:route GET /comments/{id} comments getCommentsId

GetCommentsID get comments ID API

*/
type GetCommentsID struct {
	Context *middleware.Context
	Handler GetCommentsIDHandler
}

func (o *GetCommentsID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetCommentsIDParams()

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
