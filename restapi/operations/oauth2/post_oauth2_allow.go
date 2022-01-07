// Code generated by go-swagger; DO NOT EDIT.

package oauth2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/sevings/mindwell-server/models"
)

// PostOauth2AllowHandlerFunc turns a function with the right signature into a post oauth2 allow handler
type PostOauth2AllowHandlerFunc func(PostOauth2AllowParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn PostOauth2AllowHandlerFunc) Handle(params PostOauth2AllowParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// PostOauth2AllowHandler interface for that can handle valid post oauth2 allow params
type PostOauth2AllowHandler interface {
	Handle(PostOauth2AllowParams, *models.UserID) middleware.Responder
}

// NewPostOauth2Allow creates a new http.Handler for the post oauth2 allow operation
func NewPostOauth2Allow(ctx *middleware.Context, handler PostOauth2AllowHandler) *PostOauth2Allow {
	return &PostOauth2Allow{Context: ctx, Handler: handler}
}

/* PostOauth2Allow swagger:route POST /oauth2/allow oauth2 postOauth2Allow

only for internal usage

*/
type PostOauth2Allow struct {
	Context *middleware.Context
	Handler PostOauth2AllowHandler
}

func (o *PostOauth2Allow) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostOauth2AllowParams()
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

// PostOauth2AllowOKBody post oauth2 allow o k body
//
// swagger:model PostOauth2AllowOKBody
type PostOauth2AllowOKBody struct {

	// code
	Code string `json:"code,omitempty"`

	// state
	State string `json:"state,omitempty"`
}

// Validate validates this post oauth2 allow o k body
func (o *PostOauth2AllowOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post oauth2 allow o k body based on context it is used
func (o *PostOauth2AllowOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostOauth2AllowOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostOauth2AllowOKBody) UnmarshalBinary(b []byte) error {
	var res PostOauth2AllowOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
