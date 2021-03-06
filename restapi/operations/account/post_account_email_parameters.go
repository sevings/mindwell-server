// Code generated by go-swagger; DO NOT EDIT.

package account

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// PostAccountEmailMaxParseMemory sets the maximum size in bytes for
// the multipart form parser for this operation.
//
// The default value is 32 MB.
// The multipart parser stores up to this + 10MB.
var PostAccountEmailMaxParseMemory int64 = 32 << 20

// NewPostAccountEmailParams creates a new PostAccountEmailParams object
//
// There are no default values defined in the spec.
func NewPostAccountEmailParams() PostAccountEmailParams {

	return PostAccountEmailParams{}
}

// PostAccountEmailParams contains all the bound params for the post account email operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostAccountEmail
type PostAccountEmailParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  Max Length: 500
	  Pattern: .+@.+
	  In: formData
	*/
	Email string
	/*
	  Required: true
	  Max Length: 100
	  Min Length: 6
	  In: formData
	*/
	Password string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostAccountEmailParams() beforehand.
func (o *PostAccountEmailParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := r.ParseMultipartForm(PostAccountEmailMaxParseMemory); err != nil {
		if err != http.ErrNotMultipart {
			return errors.New(400, "%v", err)
		} else if err := r.ParseForm(); err != nil {
			return errors.New(400, "%v", err)
		}
	}
	fds := runtime.Values(r.Form)

	fdEmail, fdhkEmail, _ := fds.GetOK("email")
	if err := o.bindEmail(fdEmail, fdhkEmail, route.Formats); err != nil {
		res = append(res, err)
	}

	fdPassword, fdhkPassword, _ := fds.GetOK("password")
	if err := o.bindPassword(fdPassword, fdhkPassword, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindEmail binds and validates parameter Email from formData.
func (o *PostAccountEmailParams) bindEmail(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("email", "formData", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("email", "formData", raw); err != nil {
		return err
	}
	o.Email = raw

	if err := o.validateEmail(formats); err != nil {
		return err
	}

	return nil
}

// validateEmail carries on validations for parameter Email
func (o *PostAccountEmailParams) validateEmail(formats strfmt.Registry) error {

	if err := validate.MaxLength("email", "formData", o.Email, 500); err != nil {
		return err
	}

	if err := validate.Pattern("email", "formData", o.Email, `.+@.+`); err != nil {
		return err
	}

	return nil
}

// bindPassword binds and validates parameter Password from formData.
func (o *PostAccountEmailParams) bindPassword(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("password", "formData", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("password", "formData", raw); err != nil {
		return err
	}
	o.Password = raw

	if err := o.validatePassword(formats); err != nil {
		return err
	}

	return nil
}

// validatePassword carries on validations for parameter Password
func (o *PostAccountEmailParams) validatePassword(formats strfmt.Registry) error {

	if err := validate.MinLength("password", "formData", o.Password, 6); err != nil {
		return err
	}

	if err := validate.MaxLength("password", "formData", o.Password, 100); err != nil {
		return err
	}

	return nil
}
