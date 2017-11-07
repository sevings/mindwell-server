// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// WatchingStatus watching status
// swagger:model WatchingStatus

type WatchingStatus struct {

	// id
	// Minimum: 1
	ID int64 `json:"id,omitempty"`

	// is watching
	IsWatching bool `json:"isWatching,omitempty"`
}

/* polymorph WatchingStatus id false */

/* polymorph WatchingStatus isWatching false */

// Validate validates this watching status
func (m *WatchingStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *WatchingStatus) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.MinimumInt("id", "body", int64(m.ID), 1, false); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *WatchingStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WatchingStatus) UnmarshalBinary(b []byte) error {
	var res WatchingStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
