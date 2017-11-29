// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy-server/gen/models"
)

// GetUsersMeIgnoredOKCode is the HTTP code returned for type GetUsersMeIgnoredOK
const GetUsersMeIgnoredOKCode int = 200

/*GetUsersMeIgnoredOK User list

swagger:response getUsersMeIgnoredOK
*/
type GetUsersMeIgnoredOK struct {

	/*
	  In: Body
	*/
	Payload *models.UserList `json:"body,omitempty"`
}

// NewGetUsersMeIgnoredOK creates GetUsersMeIgnoredOK with default headers values
func NewGetUsersMeIgnoredOK() *GetUsersMeIgnoredOK {
	return &GetUsersMeIgnoredOK{}
}

// WithPayload adds the payload to the get users me ignored o k response
func (o *GetUsersMeIgnoredOK) WithPayload(payload *models.UserList) *GetUsersMeIgnoredOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users me ignored o k response
func (o *GetUsersMeIgnoredOK) SetPayload(payload *models.UserList) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersMeIgnoredOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersMeIgnoredForbiddenCode is the HTTP code returned for type GetUsersMeIgnoredForbidden
const GetUsersMeIgnoredForbiddenCode int = 403

/*GetUsersMeIgnoredForbidden access denied

swagger:response getUsersMeIgnoredForbidden
*/
type GetUsersMeIgnoredForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersMeIgnoredForbidden creates GetUsersMeIgnoredForbidden with default headers values
func NewGetUsersMeIgnoredForbidden() *GetUsersMeIgnoredForbidden {
	return &GetUsersMeIgnoredForbidden{}
}

// WithPayload adds the payload to the get users me ignored forbidden response
func (o *GetUsersMeIgnoredForbidden) WithPayload(payload *models.Error) *GetUsersMeIgnoredForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users me ignored forbidden response
func (o *GetUsersMeIgnoredForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersMeIgnoredForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
