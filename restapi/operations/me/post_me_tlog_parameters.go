// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostMeTlogMaxParseMemory sets the maximum size in bytes for
// the multipart form parser for this operation.
//
// The default value is 32 MB.
// The multipart parser stores up to this + 10MB.
var PostMeTlogMaxParseMemory int64 = 32 << 20

// NewPostMeTlogParams creates a new PostMeTlogParams object
// with the default values initialized.
func NewPostMeTlogParams() PostMeTlogParams {

	var (
		// initialize parameters with default values

		inLiveDefault    = bool(false)
		isVotableDefault = bool(false)

		titleDefault = string("")
	)

	return PostMeTlogParams{
		InLive: &inLiveDefault,

		IsVotable: &isVotableDefault,

		Title: &titleDefault,
	}
}

// PostMeTlogParams contains all the bound params for the post me tlog operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostMeTlog
type PostMeTlogParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  Max Length: 30000
	  Min Length: 1
	  Pattern: \s*\S+.*
	  In: formData
	*/
	Content string
	/*
	  Max Items: 10
	  Unique: true
	  In: formData
	*/
	Images []int64
	/*
	  In: formData
	  Default: false
	*/
	InLive *bool
	/*
	  In: formData
	  Default: false
	*/
	IsVotable *bool
	/*
	  Required: true
	  In: formData
	*/
	Privacy string
	/*
	  Max Items: 5
	  Unique: true
	  In: formData
	*/
	Tags []string
	/*
	  Max Length: 500
	  In: formData
	  Default: ""
	*/
	Title *string
	/*
	  In: formData
	*/
	VisibleFor []int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostMeTlogParams() beforehand.
func (o *PostMeTlogParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := r.ParseMultipartForm(PostMeTlogMaxParseMemory); err != nil {
		if err != http.ErrNotMultipart {
			return errors.New(400, "%v", err)
		} else if err := r.ParseForm(); err != nil {
			return errors.New(400, "%v", err)
		}
	}
	fds := runtime.Values(r.Form)

	fdContent, fdhkContent, _ := fds.GetOK("content")
	if err := o.bindContent(fdContent, fdhkContent, route.Formats); err != nil {
		res = append(res, err)
	}

	fdImages, fdhkImages, _ := fds.GetOK("images")
	if err := o.bindImages(fdImages, fdhkImages, route.Formats); err != nil {
		res = append(res, err)
	}

	fdInLive, fdhkInLive, _ := fds.GetOK("inLive")
	if err := o.bindInLive(fdInLive, fdhkInLive, route.Formats); err != nil {
		res = append(res, err)
	}

	fdIsVotable, fdhkIsVotable, _ := fds.GetOK("isVotable")
	if err := o.bindIsVotable(fdIsVotable, fdhkIsVotable, route.Formats); err != nil {
		res = append(res, err)
	}

	fdPrivacy, fdhkPrivacy, _ := fds.GetOK("privacy")
	if err := o.bindPrivacy(fdPrivacy, fdhkPrivacy, route.Formats); err != nil {
		res = append(res, err)
	}

	fdTags, fdhkTags, _ := fds.GetOK("tags")
	if err := o.bindTags(fdTags, fdhkTags, route.Formats); err != nil {
		res = append(res, err)
	}

	fdTitle, fdhkTitle, _ := fds.GetOK("title")
	if err := o.bindTitle(fdTitle, fdhkTitle, route.Formats); err != nil {
		res = append(res, err)
	}

	fdVisibleFor, fdhkVisibleFor, _ := fds.GetOK("visibleFor")
	if err := o.bindVisibleFor(fdVisibleFor, fdhkVisibleFor, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindContent binds and validates parameter Content from formData.
func (o *PostMeTlogParams) bindContent(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("content", "formData", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("content", "formData", raw); err != nil {
		return err
	}
	o.Content = raw

	if err := o.validateContent(formats); err != nil {
		return err
	}

	return nil
}

// validateContent carries on validations for parameter Content
func (o *PostMeTlogParams) validateContent(formats strfmt.Registry) error {

	if err := validate.MinLength("content", "formData", o.Content, 1); err != nil {
		return err
	}

	if err := validate.MaxLength("content", "formData", o.Content, 30000); err != nil {
		return err
	}

	if err := validate.Pattern("content", "formData", o.Content, `\s*\S+.*`); err != nil {
		return err
	}

	return nil
}

// bindImages binds and validates array parameter Images from formData.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *PostMeTlogParams) bindImages(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var qvImages string
	if len(rawData) > 0 {
		qvImages = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	imagesIC := swag.SplitByFormat(qvImages, "")
	if len(imagesIC) == 0 {
		return nil
	}

	var imagesIR []int64
	for i, imagesIV := range imagesIC {
		// items.Format: "int64"
		imagesI, err := swag.ConvertInt64(imagesIV)
		if err != nil {
			return errors.InvalidType(fmt.Sprintf("%s.%v", "images", i), "formData", "int64", imagesI)
		}

		if err := validate.MinimumInt(fmt.Sprintf("%s.%v", "images", i), "formData", imagesI, 1, false); err != nil {
			return err
		}

		imagesIR = append(imagesIR, imagesI)
	}

	o.Images = imagesIR
	if err := o.validateImages(formats); err != nil {
		return err
	}

	return nil
}

// validateImages carries on validations for parameter Images
func (o *PostMeTlogParams) validateImages(formats strfmt.Registry) error {

	imagesSize := int64(len(o.Images))

	// maxItems: 10
	if err := validate.MaxItems("images", "formData", imagesSize, 10); err != nil {
		return err
	}

	// uniqueItems: true
	if err := validate.UniqueItems("images", "formData", o.Images); err != nil {
		return err
	}
	return nil
}

// bindInLive binds and validates parameter InLive from formData.
func (o *PostMeTlogParams) bindInLive(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewPostMeTlogParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("inLive", "formData", "bool", raw)
	}
	o.InLive = &value

	return nil
}

// bindIsVotable binds and validates parameter IsVotable from formData.
func (o *PostMeTlogParams) bindIsVotable(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewPostMeTlogParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("isVotable", "formData", "bool", raw)
	}
	o.IsVotable = &value

	return nil
}

// bindPrivacy binds and validates parameter Privacy from formData.
func (o *PostMeTlogParams) bindPrivacy(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("privacy", "formData", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("privacy", "formData", raw); err != nil {
		return err
	}
	o.Privacy = raw

	if err := o.validatePrivacy(formats); err != nil {
		return err
	}

	return nil
}

// validatePrivacy carries on validations for parameter Privacy
func (o *PostMeTlogParams) validatePrivacy(formats strfmt.Registry) error {

	if err := validate.EnumCase("privacy", "formData", o.Privacy, []interface{}{"all", "registered", "invited", "followers", "some", "me"}, true); err != nil {
		return err
	}

	return nil
}

// bindTags binds and validates array parameter Tags from formData.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *PostMeTlogParams) bindTags(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var qvTags string
	if len(rawData) > 0 {
		qvTags = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	tagsIC := swag.SplitByFormat(qvTags, "")
	if len(tagsIC) == 0 {
		return nil
	}

	var tagsIR []string
	for i, tagsIV := range tagsIC {
		tagsI := tagsIV

		if err := validate.MinLength(fmt.Sprintf("%s.%v", "tags", i), "formData", tagsI, 1); err != nil {
			return err
		}

		if err := validate.MaxLength(fmt.Sprintf("%s.%v", "tags", i), "formData", tagsI, 50); err != nil {
			return err
		}

		tagsIR = append(tagsIR, tagsI)
	}

	o.Tags = tagsIR
	if err := o.validateTags(formats); err != nil {
		return err
	}

	return nil
}

// validateTags carries on validations for parameter Tags
func (o *PostMeTlogParams) validateTags(formats strfmt.Registry) error {

	tagsSize := int64(len(o.Tags))

	// maxItems: 5
	if err := validate.MaxItems("tags", "formData", tagsSize, 5); err != nil {
		return err
	}

	// uniqueItems: true
	if err := validate.UniqueItems("tags", "formData", o.Tags); err != nil {
		return err
	}
	return nil
}

// bindTitle binds and validates parameter Title from formData.
func (o *PostMeTlogParams) bindTitle(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewPostMeTlogParams()
		return nil
	}
	o.Title = &raw

	if err := o.validateTitle(formats); err != nil {
		return err
	}

	return nil
}

// validateTitle carries on validations for parameter Title
func (o *PostMeTlogParams) validateTitle(formats strfmt.Registry) error {

	if err := validate.MaxLength("title", "formData", *o.Title, 500); err != nil {
		return err
	}

	return nil
}

// bindVisibleFor binds and validates array parameter VisibleFor from formData.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *PostMeTlogParams) bindVisibleFor(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var qvVisibleFor string
	if len(rawData) > 0 {
		qvVisibleFor = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	visibleForIC := swag.SplitByFormat(qvVisibleFor, "")
	if len(visibleForIC) == 0 {
		return nil
	}

	var visibleForIR []int64
	for i, visibleForIV := range visibleForIC {
		// items.Format: "int64"
		visibleForI, err := swag.ConvertInt64(visibleForIV)
		if err != nil {
			return errors.InvalidType(fmt.Sprintf("%s.%v", "visibleFor", i), "formData", "int64", visibleForI)
		}

		if err := validate.MinimumInt(fmt.Sprintf("%s.%v", "visibleFor", i), "formData", visibleForI, 1, false); err != nil {
			return err
		}

		visibleForIR = append(visibleForIR, visibleForI)
	}

	o.VisibleFor = visibleForIR

	return nil
}
