// Code generated by go-swagger; DO NOT EDIT.

package account

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/sevings/yummy-server/gen/models"
)

// PostAccountPasswordHandlerFunc turns a function with the right signature into a post account password handler
type PostAccountPasswordHandlerFunc func(PostAccountPasswordParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn PostAccountPasswordHandlerFunc) Handle(params PostAccountPasswordParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// PostAccountPasswordHandler interface for that can handle valid post account password params
type PostAccountPasswordHandler interface {
	Handle(PostAccountPasswordParams, *models.UserID) middleware.Responder
}

// NewPostAccountPassword creates a new http.Handler for the post account password operation
func NewPostAccountPassword(ctx *middleware.Context, handler PostAccountPasswordHandler) *PostAccountPassword {
	return &PostAccountPassword{Context: ctx, Handler: handler}
}

/*PostAccountPassword swagger:route POST /account/password account postAccountPassword

set new password

*/
type PostAccountPassword struct {
	Context *middleware.Context
	Handler PostAccountPasswordHandler
}

func (o *PostAccountPassword) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostAccountPasswordParams()

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
