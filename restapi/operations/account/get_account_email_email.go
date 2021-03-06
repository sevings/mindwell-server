// Code generated by go-swagger; DO NOT EDIT.

package account

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetAccountEmailEmailHandlerFunc turns a function with the right signature into a get account email email handler
type GetAccountEmailEmailHandlerFunc func(GetAccountEmailEmailParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAccountEmailEmailHandlerFunc) Handle(params GetAccountEmailEmailParams) middleware.Responder {
	return fn(params)
}

// GetAccountEmailEmailHandler interface for that can handle valid get account email email params
type GetAccountEmailEmailHandler interface {
	Handle(GetAccountEmailEmailParams) middleware.Responder
}

// NewGetAccountEmailEmail creates a new http.Handler for the get account email email operation
func NewGetAccountEmailEmail(ctx *middleware.Context, handler GetAccountEmailEmailHandler) *GetAccountEmailEmail {
	return &GetAccountEmailEmail{Context: ctx, Handler: handler}
}

/* GetAccountEmailEmail swagger:route GET /account/email/{email} account getAccountEmailEmail

check if email is used

*/
type GetAccountEmailEmail struct {
	Context *middleware.Context
	Handler GetAccountEmailEmailHandler
}

func (o *GetAccountEmailEmail) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetAccountEmailEmailParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetAccountEmailEmailOKBody get account email email o k body
// Example: {"email":"mail@example.com","isFree":true}
//
// swagger:model GetAccountEmailEmailOKBody
type GetAccountEmailEmailOKBody struct {

	// email
	Email string `json:"email,omitempty"`

	// is free
	IsFree bool `json:"isFree,omitempty"`
}

// Validate validates this get account email email o k body
func (o *GetAccountEmailEmailOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get account email email o k body based on context it is used
func (o *GetAccountEmailEmailOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetAccountEmailEmailOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetAccountEmailEmailOKBody) UnmarshalBinary(b []byte) error {
	var res GetAccountEmailEmailOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
