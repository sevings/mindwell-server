// Code generated by go-swagger; DO NOT EDIT.

package notifications

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"
	models "github.com/sevings/mindwell-server/models"
)

// PutNotificationsReadHandlerFunc turns a function with the right signature into a put notifications read handler
type PutNotificationsReadHandlerFunc func(PutNotificationsReadParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn PutNotificationsReadHandlerFunc) Handle(params PutNotificationsReadParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// PutNotificationsReadHandler interface for that can handle valid put notifications read params
type PutNotificationsReadHandler interface {
	Handle(PutNotificationsReadParams, *models.UserID) middleware.Responder
}

// NewPutNotificationsRead creates a new http.Handler for the put notifications read operation
func NewPutNotificationsRead(ctx *middleware.Context, handler PutNotificationsReadHandler) *PutNotificationsRead {
	return &PutNotificationsRead{Context: ctx, Handler: handler}
}

/*PutNotificationsRead swagger:route PUT /notifications/read notifications putNotificationsRead

PutNotificationsRead put notifications read API

*/
type PutNotificationsRead struct {
	Context *middleware.Context
	Handler PutNotificationsReadHandler
}

func (o *PutNotificationsRead) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPutNotificationsReadParams()

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

// PutNotificationsReadOKBody put notifications read o k body
// swagger:model PutNotificationsReadOKBody
type PutNotificationsReadOKBody struct {

	// unread
	Unread int64 `json:"unread,omitempty"`
}

// Validate validates this put notifications read o k body
func (o *PutNotificationsReadOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutNotificationsReadOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutNotificationsReadOKBody) UnmarshalBinary(b []byte) error {
	var res PutNotificationsReadOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
