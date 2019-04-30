// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Notification notification
// swagger:model Notification
type Notification struct {

	// comment
	Comment *Comment `json:"comment,omitempty"`

	// created at
	CreatedAt float64 `json:"createdAt,omitempty"`

	// entry
	Entry *Entry `json:"entry,omitempty"`

	// id
	// Minimum: 1
	ID int64 `json:"id,omitempty"`

	// read
	Read bool `json:"read,omitempty"`

	// type
	// Enum: [comment follower request accept invite welcome invited]
	Type string `json:"type,omitempty"`

	// user
	User *User `json:"user,omitempty"`
}

// Validate validates this notification
func (m *Notification) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateComment(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEntry(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUser(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Notification) validateComment(formats strfmt.Registry) error {

	if swag.IsZero(m.Comment) { // not required
		return nil
	}

	if m.Comment != nil {
		if err := m.Comment.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("comment")
			}
			return err
		}
	}

	return nil
}

func (m *Notification) validateEntry(formats strfmt.Registry) error {

	if swag.IsZero(m.Entry) { // not required
		return nil
	}

	if m.Entry != nil {
		if err := m.Entry.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("entry")
			}
			return err
		}
	}

	return nil
}

func (m *Notification) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.MinimumInt("id", "body", int64(m.ID), 1, false); err != nil {
		return err
	}

	return nil
}

var notificationTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["comment","follower","request","accept","invite","welcome","invited"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		notificationTypeTypePropEnum = append(notificationTypeTypePropEnum, v)
	}
}

const (

	// NotificationTypeComment captures enum value "comment"
	NotificationTypeComment string = "comment"

	// NotificationTypeFollower captures enum value "follower"
	NotificationTypeFollower string = "follower"

	// NotificationTypeRequest captures enum value "request"
	NotificationTypeRequest string = "request"

	// NotificationTypeAccept captures enum value "accept"
	NotificationTypeAccept string = "accept"

	// NotificationTypeInvite captures enum value "invite"
	NotificationTypeInvite string = "invite"

	// NotificationTypeWelcome captures enum value "welcome"
	NotificationTypeWelcome string = "welcome"

	// NotificationTypeInvited captures enum value "invited"
	NotificationTypeInvited string = "invited"
)

// prop value enum
func (m *Notification) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, notificationTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Notification) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

func (m *Notification) validateUser(formats strfmt.Registry) error {

	if swag.IsZero(m.User) { // not required
		return nil
	}

	if m.User != nil {
		if err := m.User.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("user")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Notification) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Notification) UnmarshalBinary(b []byte) error {
	var res Notification
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
