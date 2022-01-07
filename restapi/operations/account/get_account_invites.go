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

	"github.com/sevings/mindwell-server/models"
)

// GetAccountInvitesHandlerFunc turns a function with the right signature into a get account invites handler
type GetAccountInvitesHandlerFunc func(GetAccountInvitesParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAccountInvitesHandlerFunc) Handle(params GetAccountInvitesParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// GetAccountInvitesHandler interface for that can handle valid get account invites params
type GetAccountInvitesHandler interface {
	Handle(GetAccountInvitesParams, *models.UserID) middleware.Responder
}

// NewGetAccountInvites creates a new http.Handler for the get account invites operation
func NewGetAccountInvites(ctx *middleware.Context, handler GetAccountInvitesHandler) *GetAccountInvites {
	return &GetAccountInvites{Context: ctx, Handler: handler}
}

/* GetAccountInvites swagger:route GET /account/invites account getAccountInvites

GetAccountInvites get account invites API

*/
type GetAccountInvites struct {
	Context *middleware.Context
	Handler GetAccountInvitesHandler
}

func (o *GetAccountInvites) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetAccountInvitesParams()
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

// GetAccountInvitesOKBody get account invites o k body
//
// swagger:model GetAccountInvitesOKBody
type GetAccountInvitesOKBody struct {

	// invites
	Invites []string `json:"invites"`
}

// Validate validates this get account invites o k body
func (o *GetAccountInvitesOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get account invites o k body based on context it is used
func (o *GetAccountInvitesOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetAccountInvitesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetAccountInvitesOKBody) UnmarshalBinary(b []byte) error {
	var res GetAccountInvitesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
