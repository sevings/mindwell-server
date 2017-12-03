// Code generated by go-swagger; DO NOT EDIT.

package design

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

// NewPutDesignParams creates a new PutDesignParams object
// with the default values initialized.
func NewPutDesignParams() PutDesignParams {
	var ()
	return PutDesignParams{}
}

// PutDesignParams contains all the bound params for the put design operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutDesign
type PutDesignParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Pattern: #[0-9a-d]{6}
	  In: formData
	*/
	BackgroundColor *string
	/*
	  In: formData
	*/
	CSS *string
	/*
	  In: formData
	*/
	FontFamily *string
	/*
	  In: formData
	*/
	FontSize *int64
	/*
	  In: formData
	*/
	TextAlignment *string
	/*
	  Pattern: #[0-9a-d]{6}
	  In: formData
	*/
	TextColor *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *PutDesignParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
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

	fdBackgroundColor, fdhkBackgroundColor, _ := fds.GetOK("backgroundColor")
	if err := o.bindBackgroundColor(fdBackgroundColor, fdhkBackgroundColor, route.Formats); err != nil {
		res = append(res, err)
	}

	fdCSS, fdhkCSS, _ := fds.GetOK("css")
	if err := o.bindCSS(fdCSS, fdhkCSS, route.Formats); err != nil {
		res = append(res, err)
	}

	fdFontFamily, fdhkFontFamily, _ := fds.GetOK("fontFamily")
	if err := o.bindFontFamily(fdFontFamily, fdhkFontFamily, route.Formats); err != nil {
		res = append(res, err)
	}

	fdFontSize, fdhkFontSize, _ := fds.GetOK("fontSize")
	if err := o.bindFontSize(fdFontSize, fdhkFontSize, route.Formats); err != nil {
		res = append(res, err)
	}

	fdTextAlignment, fdhkTextAlignment, _ := fds.GetOK("textAlignment")
	if err := o.bindTextAlignment(fdTextAlignment, fdhkTextAlignment, route.Formats); err != nil {
		res = append(res, err)
	}

	fdTextColor, fdhkTextColor, _ := fds.GetOK("textColor")
	if err := o.bindTextColor(fdTextColor, fdhkTextColor, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutDesignParams) bindBackgroundColor(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.BackgroundColor = &raw

	if err := o.validateBackgroundColor(formats); err != nil {
		return err
	}

	return nil
}

func (o *PutDesignParams) validateBackgroundColor(formats strfmt.Registry) error {

	if err := validate.Pattern("backgroundColor", "formData", (*o.BackgroundColor), `#[0-9a-d]{6}`); err != nil {
		return err
	}

	return nil
}

func (o *PutDesignParams) bindCSS(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.CSS = &raw

	return nil
}

func (o *PutDesignParams) bindFontFamily(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.FontFamily = &raw

	return nil
}

func (o *PutDesignParams) bindFontSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("fontSize", "formData", "int64", raw)
	}
	o.FontSize = &value

	return nil
}

func (o *PutDesignParams) bindTextAlignment(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.TextAlignment = &raw

	if err := o.validateTextAlignment(formats); err != nil {
		return err
	}

	return nil
}

func (o *PutDesignParams) validateTextAlignment(formats strfmt.Registry) error {

	if err := validate.Enum("textAlignment", "formData", *o.TextAlignment, []interface{}{"left", "right", "center", "justify"}); err != nil {
		return err
	}

	return nil
}

func (o *PutDesignParams) bindTextColor(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.TextColor = &raw

	if err := o.validateTextColor(formats); err != nil {
		return err
	}

	return nil
}

func (o *PutDesignParams) validateTextColor(formats strfmt.Registry) error {

	if err := validate.Pattern("textColor", "formData", (*o.TextColor), `#[0-9a-d]{6}`); err != nil {
		return err
	}

	return nil
}
