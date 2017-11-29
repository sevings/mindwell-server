// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteEntriesIDParams creates a new DeleteEntriesIDParams object
// with the default values initialized.
func NewDeleteEntriesIDParams() DeleteEntriesIDParams {
	var ()
	return DeleteEntriesIDParams{}
}

// DeleteEntriesIDParams contains all the bound params for the delete entries ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters DeleteEntriesID
type DeleteEntriesIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  Max Length: 32
	  Min Length: 32
	  In: header
	*/
	XUserKey string
	/*
	  Required: true
	  Minimum: 1
	  In: path
	*/
	ID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *DeleteEntriesIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if err := o.bindXUserKey(r.Header[http.CanonicalHeaderKey("X-User-Key")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	rID, rhkID, _ := route.Params.GetOK("id")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DeleteEntriesIDParams) bindXUserKey(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-User-Key", "header")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if err := validate.RequiredString("X-User-Key", "header", raw); err != nil {
		return err
	}

	o.XUserKey = raw

	if err := o.validateXUserKey(formats); err != nil {
		return err
	}

	return nil
}

func (o *DeleteEntriesIDParams) validateXUserKey(formats strfmt.Registry) error {

	if err := validate.MinLength("X-User-Key", "header", o.XUserKey, 32); err != nil {
		return err
	}

	if err := validate.MaxLength("X-User-Key", "header", o.XUserKey, 32); err != nil {
		return err
	}

	return nil
}

func (o *DeleteEntriesIDParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("id", "path", "int64", raw)
	}
	o.ID = value

	if err := o.validateID(formats); err != nil {
		return err
	}

	return nil
}

func (o *DeleteEntriesIDParams) validateID(formats strfmt.Registry) error {

	if err := validate.MinimumInt("id", "path", int64(o.ID), 1, false); err != nil {
		return err
	}

	return nil
}
