// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// FavoriteStatus favorite status
//
// swagger:model FavoriteStatus
type FavoriteStatus struct {

	// count
	Count int64 `json:"count,omitempty"`

	// id
	// Minimum: 1
	ID int64 `json:"id,omitempty"`

	// is favorited
	IsFavorited bool `json:"isFavorited,omitempty"`
}

// Validate validates this favorite status
func (m *FavoriteStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *FavoriteStatus) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.MinimumInt("id", "body", m.ID, 1, false); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this favorite status based on context it is used
func (m *FavoriteStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *FavoriteStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FavoriteStatus) UnmarshalBinary(b []byte) error {
	var res FavoriteStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
