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

// Entry entry
// swagger:model Entry
type Entry struct {

	// author
	Author *User `json:"author,omitempty"`

	// comment count
	CommentCount int64 `json:"commentCount,omitempty"`

	// comments
	Comments EntryComments `json:"comments"`

	// content
	Content string `json:"content,omitempty"`

	// created at
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// id
	// Minimum: 1
	ID int64 `json:"id,omitempty"`

	// is favorited
	IsFavorited bool `json:"isFavorited,omitempty"`

	// is votable
	IsVotable bool `json:"isVotable,omitempty"`

	// is watching
	IsWatching bool `json:"isWatching,omitempty"`

	// privacy
	Privacy string `json:"privacy,omitempty"`

	// rating
	Rating int64 `json:"rating,omitempty"`

	// title
	Title string `json:"title,omitempty"`

	// visible for
	VisibleFor EntryVisibleFor `json:"visibleFor"`

	// vote
	Vote string `json:"vote,omitempty"`

	// word count
	WordCount int64 `json:"wordCount,omitempty"`
}

// Validate validates this entry
func (m *Entry) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAuthor(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validatePrivacy(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateVote(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Entry) validateAuthor(formats strfmt.Registry) error {

	if swag.IsZero(m.Author) { // not required
		return nil
	}

	if m.Author != nil {

		if err := m.Author.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("author")
			}
			return err
		}
	}

	return nil
}

func (m *Entry) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.MinimumInt("id", "body", int64(m.ID), 1, false); err != nil {
		return err
	}

	return nil
}

var entryTypePrivacyPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["all","some","me","anonymous"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		entryTypePrivacyPropEnum = append(entryTypePrivacyPropEnum, v)
	}
}

const (
	// EntryPrivacyAll captures enum value "all"
	EntryPrivacyAll string = "all"
	// EntryPrivacySome captures enum value "some"
	EntryPrivacySome string = "some"
	// EntryPrivacyMe captures enum value "me"
	EntryPrivacyMe string = "me"
	// EntryPrivacyAnonymous captures enum value "anonymous"
	EntryPrivacyAnonymous string = "anonymous"
)

// prop value enum
func (m *Entry) validatePrivacyEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, entryTypePrivacyPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Entry) validatePrivacy(formats strfmt.Registry) error {

	if swag.IsZero(m.Privacy) { // not required
		return nil
	}

	// value enum
	if err := m.validatePrivacyEnum("privacy", "body", m.Privacy); err != nil {
		return err
	}

	return nil
}

var entryTypeVotePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["not","pos","neg","ban"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		entryTypeVotePropEnum = append(entryTypeVotePropEnum, v)
	}
}

const (
	// EntryVoteNot captures enum value "not"
	EntryVoteNot string = "not"
	// EntryVotePos captures enum value "pos"
	EntryVotePos string = "pos"
	// EntryVoteNeg captures enum value "neg"
	EntryVoteNeg string = "neg"
	// EntryVoteBan captures enum value "ban"
	EntryVoteBan string = "ban"
)

// prop value enum
func (m *Entry) validateVoteEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, entryTypeVotePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Entry) validateVote(formats strfmt.Registry) error {

	if swag.IsZero(m.Vote) { // not required
		return nil
	}

	// value enum
	if err := m.validateVoteEnum("vote", "body", m.Vote); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Entry) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Entry) UnmarshalBinary(b []byte) error {
	var res Entry
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
