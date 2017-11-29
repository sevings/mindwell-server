// Code generated by go-swagger; DO NOT EDIT.

package me

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

// NewGetUsersMeEntriesParams creates a new GetUsersMeEntriesParams object
// with the default values initialized.
func NewGetUsersMeEntriesParams() GetUsersMeEntriesParams {
	var (
		limitDefault = int64(50)
		skipDefault  = int64(0)
	)
	return GetUsersMeEntriesParams{
		Limit: &limitDefault,

		Skip: &skipDefault,
	}
}

// GetUsersMeEntriesParams contains all the bound params for the get users me entries operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetUsersMeEntries
type GetUsersMeEntriesParams struct {

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
	  Maximum: 100
	  Minimum: 1
	  In: query
	  Default: 50
	*/
	Limit *int64
	/*
	  In: query
	  Default: 0
	*/
	Skip *int64
	/*
	  Max Length: 50
	  In: query
	*/
	Tag *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetUsersMeEntriesParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	if err := o.bindXUserKey(r.Header[http.CanonicalHeaderKey("X-User-Key")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	qSkip, qhkSkip, _ := qs.GetOK("skip")
	if err := o.bindSkip(qSkip, qhkSkip, route.Formats); err != nil {
		res = append(res, err)
	}

	qTag, qhkTag, _ := qs.GetOK("tag")
	if err := o.bindTag(qTag, qhkTag, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetUsersMeEntriesParams) bindXUserKey(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

func (o *GetUsersMeEntriesParams) validateXUserKey(formats strfmt.Registry) error {

	if err := validate.MinLength("X-User-Key", "header", o.XUserKey, 32); err != nil {
		return err
	}

	if err := validate.MaxLength("X-User-Key", "header", o.XUserKey, 32); err != nil {
		return err
	}

	return nil
}

func (o *GetUsersMeEntriesParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		var limitDefault int64 = int64(50)
		o.Limit = &limitDefault
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("limit", "query", "int64", raw)
	}
	o.Limit = &value

	if err := o.validateLimit(formats); err != nil {
		return err
	}

	return nil
}

func (o *GetUsersMeEntriesParams) validateLimit(formats strfmt.Registry) error {

	if err := validate.MinimumInt("limit", "query", int64(*o.Limit), 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("limit", "query", int64(*o.Limit), 100, false); err != nil {
		return err
	}

	return nil
}

func (o *GetUsersMeEntriesParams) bindSkip(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		var skipDefault int64 = int64(0)
		o.Skip = &skipDefault
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("skip", "query", "int64", raw)
	}
	o.Skip = &value

	return nil
}

func (o *GetUsersMeEntriesParams) bindTag(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Tag = &raw

	if err := o.validateTag(formats); err != nil {
		return err
	}

	return nil
}

func (o *GetUsersMeEntriesParams) validateTag(formats strfmt.Registry) error {

	if err := validate.MaxLength("tag", "query", (*o.Tag), 50); err != nil {
		return err
	}

	return nil
}
