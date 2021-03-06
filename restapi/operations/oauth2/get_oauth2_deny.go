// Code generated by go-swagger; DO NOT EDIT.

package oauth2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetOauth2DenyHandlerFunc turns a function with the right signature into a get oauth2 deny handler
type GetOauth2DenyHandlerFunc func(GetOauth2DenyParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetOauth2DenyHandlerFunc) Handle(params GetOauth2DenyParams) middleware.Responder {
	return fn(params)
}

// GetOauth2DenyHandler interface for that can handle valid get oauth2 deny params
type GetOauth2DenyHandler interface {
	Handle(GetOauth2DenyParams) middleware.Responder
}

// NewGetOauth2Deny creates a new http.Handler for the get oauth2 deny operation
func NewGetOauth2Deny(ctx *middleware.Context, handler GetOauth2DenyHandler) *GetOauth2Deny {
	return &GetOauth2Deny{Context: ctx, Handler: handler}
}

/* GetOauth2Deny swagger:route GET /oauth2/deny oauth2 getOauth2Deny

only for internal usage

*/
type GetOauth2Deny struct {
	Context *middleware.Context
	Handler GetOauth2DenyHandler
}

func (o *GetOauth2Deny) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetOauth2DenyParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
