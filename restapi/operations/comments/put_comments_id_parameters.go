// Code generated by go-swagger; DO NOT EDIT.

package comments

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

// NewPutCommentsIDParams creates a new PutCommentsIDParams object
// with the default values initialized.
func NewPutCommentsIDParams() PutCommentsIDParams {
	var ()
	return PutCommentsIDParams{}
}

// PutCommentsIDParams contains all the bound params for the put comments ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutCommentsID
type PutCommentsIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  Max Length: 1000
	  Min Length: 1
	  In: formData
	*/
	Content string
	/*
	  Required: true
	  Minimum: 1
	  In: path
	*/
	ID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *PutCommentsIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		if err != http.ErrNotMultipart {
			return err
		} else if err := r.ParseForm(); err != nil {
			return err
		}
	}
	fds := runtime.Values(r.Form)

	fdContent, fdhkContent, _ := fds.GetOK("content")
	if err := o.bindContent(fdContent, fdhkContent, route.Formats); err != nil {
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

func (o *PutCommentsIDParams) bindContent(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("content", "formData")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if err := validate.RequiredString("content", "formData", raw); err != nil {
		return err
	}

	o.Content = raw

	if err := o.validateContent(formats); err != nil {
		return err
	}

	return nil
}

func (o *PutCommentsIDParams) validateContent(formats strfmt.Registry) error {

	if err := validate.MinLength("content", "formData", o.Content, 1); err != nil {
		return err
	}

	if err := validate.MaxLength("content", "formData", o.Content, 1000); err != nil {
		return err
	}

	return nil
}

func (o *PutCommentsIDParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

func (o *PutCommentsIDParams) validateID(formats strfmt.Registry) error {

	if err := validate.MinimumInt("id", "path", int64(o.ID), 1, false); err != nil {
		return err
	}

	return nil
}
