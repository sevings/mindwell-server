// Code generated by go-swagger; DO NOT EDIT.

package relations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewPutRelationsToNameParams creates a new PutRelationsToNameParams object
// no default values defined in spec.
func NewPutRelationsToNameParams() PutRelationsToNameParams {

	return PutRelationsToNameParams{}
}

// PutRelationsToNameParams contains all the bound params for the put relations to name operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutRelationsToName
type PutRelationsToNameParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  Max Length: 20
	  Min Length: 1
	  In: path
	*/
	Name string
	/*
	  Required: true
	  In: query
	*/
	R string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPutRelationsToNameParams() beforehand.
func (o *PutRelationsToNameParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rName, rhkName, _ := route.Params.GetOK("name")
	if err := o.bindName(rName, rhkName, route.Formats); err != nil {
		res = append(res, err)
	}

	qR, qhkR, _ := qs.GetOK("r")
	if err := o.bindR(qR, qhkR, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindName binds and validates parameter Name from path.
func (o *PutRelationsToNameParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *PutRelationsToNameParams) validateName(formats strfmt.Registry) error {

	if err := validate.MinLength("name", "path", o.Name, 1); err != nil {
		return err
	}

	if err := validate.MaxLength("name", "path", o.Name, 20); err != nil {
		return err
	}

	return nil
}

// bindR binds and validates parameter R from query.
func (o *PutRelationsToNameParams) bindR(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("r", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("r", "query", raw); err != nil {
		return err
	}

	o.R = raw

	if err := o.validateR(formats); err != nil {
		return err
	}

	return nil
}

// validateR carries on validations for parameter R
func (o *PutRelationsToNameParams) validateR(formats strfmt.Registry) error {

	if err := validate.Enum("r", "query", o.R, []interface{}{"followed", "ignored"}); err != nil {
		return err
	}

	return nil
}
