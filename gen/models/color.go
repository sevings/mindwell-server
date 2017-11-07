// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// Color color in rgb
// swagger:model Color

type Color string

// Validate validates this color
func (m Color) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validate.Pattern("", "body", string(m), `#[0-9a-d]{6}`); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
