// Code generated by go-swagger; DO NOT EDIT.

package oauth2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostOauth2TokenHandlerFunc turns a function with the right signature into a post oauth2 token handler
type PostOauth2TokenHandlerFunc func(PostOauth2TokenParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostOauth2TokenHandlerFunc) Handle(params PostOauth2TokenParams) middleware.Responder {
	return fn(params)
}

// PostOauth2TokenHandler interface for that can handle valid post oauth2 token params
type PostOauth2TokenHandler interface {
	Handle(PostOauth2TokenParams) middleware.Responder
}

// NewPostOauth2Token creates a new http.Handler for the post oauth2 token operation
func NewPostOauth2Token(ctx *middleware.Context, handler PostOauth2TokenHandler) *PostOauth2Token {
	return &PostOauth2Token{Context: ctx, Handler: handler}
}

/* PostOauth2Token swagger:route POST /oauth2/token oauth2 postOauth2Token

PostOauth2Token post oauth2 token API

*/
type PostOauth2Token struct {
	Context *middleware.Context
	Handler PostOauth2TokenHandler
}

func (o *PostOauth2Token) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostOauth2TokenParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostOauth2TokenOKBody post oauth2 token o k body
//
// swagger:model PostOauth2TokenOKBody
type PostOauth2TokenOKBody struct {

	// access token
	AccessToken string `json:"access_token,omitempty"`

	// expires in
	ExpiresIn int64 `json:"expires_in,omitempty"`

	// refresh token
	RefreshToken string `json:"refresh_token,omitempty"`

	// scope
	Scope []string `json:"scope"`

	// token type
	// Enum: [bearer]
	TokenType string `json:"token_type,omitempty"`
}

// Validate validates this post oauth2 token o k body
func (o *PostOauth2TokenOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateTokenType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var postOauth2TokenOKBodyTypeTokenTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["bearer"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		postOauth2TokenOKBodyTypeTokenTypePropEnum = append(postOauth2TokenOKBodyTypeTokenTypePropEnum, v)
	}
}

const (

	// PostOauth2TokenOKBodyTokenTypeBearer captures enum value "bearer"
	PostOauth2TokenOKBodyTokenTypeBearer string = "bearer"
)

// prop value enum
func (o *PostOauth2TokenOKBody) validateTokenTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, postOauth2TokenOKBodyTypeTokenTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *PostOauth2TokenOKBody) validateTokenType(formats strfmt.Registry) error {
	if swag.IsZero(o.TokenType) { // not required
		return nil
	}

	// value enum
	if err := o.validateTokenTypeEnum("postOauth2TokenOK"+"."+"token_type", "body", o.TokenType); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this post oauth2 token o k body based on context it is used
func (o *PostOauth2TokenOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostOauth2TokenOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostOauth2TokenOKBody) UnmarshalBinary(b []byte) error {
	var res PostOauth2TokenOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
