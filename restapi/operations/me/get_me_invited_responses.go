// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// GetMeInvitedOKCode is the HTTP code returned for type GetMeInvitedOK
const GetMeInvitedOKCode int = 200

/*GetMeInvitedOK User list

swagger:response getMeInvitedOK
*/
type GetMeInvitedOK struct {

	/*
	  In: Body
	*/
	Payload *models.FriendList `json:"body,omitempty"`
}

// NewGetMeInvitedOK creates GetMeInvitedOK with default headers values
func NewGetMeInvitedOK() *GetMeInvitedOK {
	return &GetMeInvitedOK{}
}

// WithPayload adds the payload to the get me invited o k response
func (o *GetMeInvitedOK) WithPayload(payload *models.FriendList) *GetMeInvitedOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get me invited o k response
func (o *GetMeInvitedOK) SetPayload(payload *models.FriendList) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMeInvitedOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetMeInvitedForbiddenCode is the HTTP code returned for type GetMeInvitedForbidden
const GetMeInvitedForbiddenCode int = 403

/*GetMeInvitedForbidden access denied

swagger:response getMeInvitedForbidden
*/
type GetMeInvitedForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetMeInvitedForbidden creates GetMeInvitedForbidden with default headers values
func NewGetMeInvitedForbidden() *GetMeInvitedForbidden {
	return &GetMeInvitedForbidden{}
}

// WithPayload adds the payload to the get me invited forbidden response
func (o *GetMeInvitedForbidden) WithPayload(payload *models.Error) *GetMeInvitedForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get me invited forbidden response
func (o *GetMeInvitedForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMeInvitedForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
