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

// ProfileAllOf1 profile all of1
// swagger:model profileAllOf1
type ProfileAllOf1 struct {

	// age lower bound
	AgeLowerBound int64 `json:"ageLowerBound,omitempty"`

	// age upper bound
	AgeUpperBound int64 `json:"ageUpperBound,omitempty"`

	// city
	// Max Length: 50
	City string `json:"city,omitempty"`

	// country
	// Max Length: 50
	Country string `json:"country,omitempty"`

	// counts
	Counts *ProfileAllOf1Counts `json:"counts,omitempty"`

	// created at
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// design
	Design *Design `json:"design,omitempty"`

	// gender
	Gender string `json:"gender,omitempty"`

	// invited by
	InvitedBy *User `json:"invitedBy,omitempty"`

	// is daylog
	IsDaylog bool `json:"isDaylog,omitempty"`

	// karma
	Karma float32 `json:"karma,omitempty"`

	// last seen at
	LastSeenAt strfmt.DateTime `json:"lastSeenAt,omitempty"`

	// privacy
	Privacy string `json:"privacy,omitempty"`

	// relations
	Relations *ProfileAllOf1Relations `json:"relations,omitempty"`

	// title
	// Max Length: 260
	Title string `json:"title,omitempty"`
}

// Validate validates this profile all of1
func (m *ProfileAllOf1) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCity(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateCountry(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateCounts(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateDesign(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateGender(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateInvitedBy(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validatePrivacy(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateRelations(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTitle(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProfileAllOf1) validateCity(formats strfmt.Registry) error {

	if swag.IsZero(m.City) { // not required
		return nil
	}

	if err := validate.MaxLength("city", "body", string(m.City), 50); err != nil {
		return err
	}

	return nil
}

func (m *ProfileAllOf1) validateCountry(formats strfmt.Registry) error {

	if swag.IsZero(m.Country) { // not required
		return nil
	}

	if err := validate.MaxLength("country", "body", string(m.Country), 50); err != nil {
		return err
	}

	return nil
}

func (m *ProfileAllOf1) validateCounts(formats strfmt.Registry) error {

	if swag.IsZero(m.Counts) { // not required
		return nil
	}

	if m.Counts != nil {

		if err := m.Counts.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("counts")
			}
			return err
		}
	}

	return nil
}

func (m *ProfileAllOf1) validateDesign(formats strfmt.Registry) error {

	if swag.IsZero(m.Design) { // not required
		return nil
	}

	if m.Design != nil {

		if err := m.Design.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("design")
			}
			return err
		}
	}

	return nil
}

var profileAllOf1TypeGenderPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["male","female","not set"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		profileAllOf1TypeGenderPropEnum = append(profileAllOf1TypeGenderPropEnum, v)
	}
}

const (
	// ProfileAllOf1GenderMale captures enum value "male"
	ProfileAllOf1GenderMale string = "male"
	// ProfileAllOf1GenderFemale captures enum value "female"
	ProfileAllOf1GenderFemale string = "female"
	// ProfileAllOf1GenderNotSet captures enum value "not set"
	ProfileAllOf1GenderNotSet string = "not set"
)

// prop value enum
func (m *ProfileAllOf1) validateGenderEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, profileAllOf1TypeGenderPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ProfileAllOf1) validateGender(formats strfmt.Registry) error {

	if swag.IsZero(m.Gender) { // not required
		return nil
	}

	// value enum
	if err := m.validateGenderEnum("gender", "body", m.Gender); err != nil {
		return err
	}

	return nil
}

func (m *ProfileAllOf1) validateInvitedBy(formats strfmt.Registry) error {

	if swag.IsZero(m.InvitedBy) { // not required
		return nil
	}

	if m.InvitedBy != nil {

		if err := m.InvitedBy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("invitedBy")
			}
			return err
		}
	}

	return nil
}

var profileAllOf1TypePrivacyPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["all","followers"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		profileAllOf1TypePrivacyPropEnum = append(profileAllOf1TypePrivacyPropEnum, v)
	}
}

const (
	// ProfileAllOf1PrivacyAll captures enum value "all"
	ProfileAllOf1PrivacyAll string = "all"
	// ProfileAllOf1PrivacyFollowers captures enum value "followers"
	ProfileAllOf1PrivacyFollowers string = "followers"
)

// prop value enum
func (m *ProfileAllOf1) validatePrivacyEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, profileAllOf1TypePrivacyPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ProfileAllOf1) validatePrivacy(formats strfmt.Registry) error {

	if swag.IsZero(m.Privacy) { // not required
		return nil
	}

	// value enum
	if err := m.validatePrivacyEnum("privacy", "body", m.Privacy); err != nil {
		return err
	}

	return nil
}

func (m *ProfileAllOf1) validateRelations(formats strfmt.Registry) error {

	if swag.IsZero(m.Relations) { // not required
		return nil
	}

	if m.Relations != nil {

		if err := m.Relations.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("relations")
			}
			return err
		}
	}

	return nil
}

func (m *ProfileAllOf1) validateTitle(formats strfmt.Registry) error {

	if swag.IsZero(m.Title) { // not required
		return nil
	}

	if err := validate.MaxLength("title", "body", string(m.Title), 260); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProfileAllOf1) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProfileAllOf1) UnmarshalBinary(b []byte) error {
	var res ProfileAllOf1
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
