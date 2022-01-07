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

// GetAccountNameNameHandlerFunc turns a function with the right signature into a get account name name handler
type GetAccountNameNameHandlerFunc func(GetAccountNameNameParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAccountNameNameHandlerFunc) Handle(params GetAccountNameNameParams) middleware.Responder {
	return fn(params)
}

// GetAccountNameNameHandler interface for that can handle valid get account name name params
type GetAccountNameNameHandler interface {
	Handle(GetAccountNameNameParams) middleware.Responder
}

// NewGetAccountNameName creates a new http.Handler for the get account name name operation
func NewGetAccountNameName(ctx *middleware.Context, handler GetAccountNameNameHandler) *GetAccountNameName {
	return &GetAccountNameName{Context: ctx, Handler: handler}
}

/* GetAccountNameName swagger:route GET /account/name/{name} account getAccountNameName

check if name is used

*/
type GetAccountNameName struct {
	Context *middleware.Context
	Handler GetAccountNameNameHandler
}

func (o *GetAccountNameName) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetAccountNameNameParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetAccountNameNameOKBody get account name name o k body
// Example: {"isFree":false,"name":"example"}
//
// swagger:model GetAccountNameNameOKBody
type GetAccountNameNameOKBody struct {

	// is free
	IsFree bool `json:"isFree,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this get account name name o k body
func (o *GetAccountNameNameOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get account name name o k body based on context it is used
func (o *GetAccountNameNameOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetAccountNameNameOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetAccountNameNameOKBody) UnmarshalBinary(b []byte) error {
	var res GetAccountNameNameOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
