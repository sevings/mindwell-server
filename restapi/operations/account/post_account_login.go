// Code generated by go-swagger; DO NOT EDIT.

package account

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PostAccountLoginHandlerFunc turns a function with the right signature into a post account login handler
type PostAccountLoginHandlerFunc func(PostAccountLoginParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostAccountLoginHandlerFunc) Handle(params PostAccountLoginParams) middleware.Responder {
	return fn(params)
}

// PostAccountLoginHandler interface for that can handle valid post account login params
type PostAccountLoginHandler interface {
	Handle(PostAccountLoginParams) middleware.Responder
}

// NewPostAccountLogin creates a new http.Handler for the post account login operation
func NewPostAccountLogin(ctx *middleware.Context, handler PostAccountLoginHandler) *PostAccountLogin {
	return &PostAccountLogin{Context: ctx, Handler: handler}
}

/* PostAccountLogin swagger:route POST /account/login account postAccountLogin

PostAccountLogin post account login API

*/
type PostAccountLogin struct {
	Context *middleware.Context
	Handler PostAccountLoginHandler
}

func (o *PostAccountLogin) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostAccountLoginParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
