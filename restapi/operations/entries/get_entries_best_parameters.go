// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetEntriesBestParams creates a new GetEntriesBestParams object
// with the default values initialized.
func NewGetEntriesBestParams() GetEntriesBestParams {

	var (
		// initialize parameters with default values

		categoryDefault = string("month")
		limitDefault    = int64(30)
		queryDefault    = string("")
		tagDefault      = string("")
	)

	return GetEntriesBestParams{
		Category: &categoryDefault,

		Limit: &limitDefault,

		Query: &queryDefault,

		Tag: &tagDefault,
	}
}

// GetEntriesBestParams contains all the bound params for the get entries best operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetEntriesBest
type GetEntriesBestParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	  Default: "month"
	*/
	Category *string
	/*
	  Maximum: 100
	  Minimum: 1
	  In: query
	  Default: 30
	*/
	Limit *int64
	/*
	  Max Length: 100
	  In: query
	  Default: ""
	*/
	Query *string
	/*
	  Max Length: 50
	  In: query
	  Default: ""
	*/
	Tag *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetEntriesBestParams() beforehand.
func (o *GetEntriesBestParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qCategory, qhkCategory, _ := qs.GetOK("category")
	if err := o.bindCategory(qCategory, qhkCategory, route.Formats); err != nil {
		res = append(res, err)
	}

	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	qQuery, qhkQuery, _ := qs.GetOK("query")
	if err := o.bindQuery(qQuery, qhkQuery, route.Formats); err != nil {
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

// bindCategory binds and validates parameter Category from query.
func (o *GetEntriesBestParams) bindCategory(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetEntriesBestParams()
		return nil
	}
	o.Category = &raw

	if err := o.validateCategory(formats); err != nil {
		return err
	}

	return nil
}

// validateCategory carries on validations for parameter Category
func (o *GetEntriesBestParams) validateCategory(formats strfmt.Registry) error {

	if err := validate.EnumCase("category", "query", *o.Category, []interface{}{"month", "week"}, true); err != nil {
		return err
	}

	return nil
}

// bindLimit binds and validates parameter Limit from query.
func (o *GetEntriesBestParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetEntriesBestParams()
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

// validateLimit carries on validations for parameter Limit
func (o *GetEntriesBestParams) validateLimit(formats strfmt.Registry) error {

	if err := validate.MinimumInt("limit", "query", *o.Limit, 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("limit", "query", *o.Limit, 100, false); err != nil {
		return err
	}

	return nil
}

// bindQuery binds and validates parameter Query from query.
func (o *GetEntriesBestParams) bindQuery(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetEntriesBestParams()
		return nil
	}
	o.Query = &raw

	if err := o.validateQuery(formats); err != nil {
		return err
	}

	return nil
}

// validateQuery carries on validations for parameter Query
func (o *GetEntriesBestParams) validateQuery(formats strfmt.Registry) error {

	if err := validate.MaxLength("query", "query", *o.Query, 100); err != nil {
		return err
	}

	return nil
}

// bindTag binds and validates parameter Tag from query.
func (o *GetEntriesBestParams) bindTag(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetEntriesBestParams()
		return nil
	}
	o.Tag = &raw

	if err := o.validateTag(formats); err != nil {
		return err
	}

	return nil
}

// validateTag carries on validations for parameter Tag
func (o *GetEntriesBestParams) validateTag(formats strfmt.Registry) error {

	if err := validate.MaxLength("tag", "query", *o.Tag, 50); err != nil {
		return err
	}

	return nil
}
