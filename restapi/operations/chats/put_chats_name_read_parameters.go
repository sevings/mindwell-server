// Code generated by go-swagger; DO NOT EDIT.

package chats

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewPutChatsNameReadParams creates a new PutChatsNameReadParams object
// no default values defined in spec.
func NewPutChatsNameReadParams() PutChatsNameReadParams {

	return PutChatsNameReadParams{}
}

// PutChatsNameReadParams contains all the bound params for the put chats name read operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutChatsNameRead
type PutChatsNameReadParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: query
	*/
	Message int64
	/*
	  Required: true
	  Max Length: 20
	  Min Length: 1
	  Pattern: ^[0-9\-_]*[a-zA-Z][a-zA-Z0-9\-_]*$
	  In: path
	*/
	Name string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPutChatsNameReadParams() beforehand.
func (o *PutChatsNameReadParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qMessage, qhkMessage, _ := qs.GetOK("message")
	if err := o.bindMessage(qMessage, qhkMessage, route.Formats); err != nil {
		res = append(res, err)
	}

	rName, rhkName, _ := route.Params.GetOK("name")
	if err := o.bindName(rName, rhkName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindMessage binds and validates parameter Message from query.
func (o *PutChatsNameReadParams) bindMessage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("message", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("message", "query", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("message", "query", "int64", raw)
	}
	o.Message = value

	return nil
}

// bindName binds and validates parameter Name from path.
func (o *PutChatsNameReadParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Name = raw

	if err := o.validateName(formats); err != nil {
		return err
	}

	return nil
}

// validateName carries on validations for parameter Name
func (o *PutChatsNameReadParams) validateName(formats strfmt.Registry) error {

	if err := validate.MinLength("name", "path", o.Name, 1); err != nil {
		return err
	}

	if err := validate.MaxLength("name", "path", o.Name, 20); err != nil {
		return err
	}

	if err := validate.Pattern("name", "path", o.Name, `^[0-9\-_]*[a-zA-Z][a-zA-Z0-9\-_]*$`); err != nil {
		return err
	}

	return nil
}
