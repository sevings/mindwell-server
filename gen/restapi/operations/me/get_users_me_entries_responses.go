// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy/gen/models"
)

// GetUsersMeEntriesOKCode is the HTTP code returned for type GetUsersMeEntriesOK
const GetUsersMeEntriesOKCode int = 200

/*GetUsersMeEntriesOK Entry list

swagger:response getUsersMeEntriesOK
*/
type GetUsersMeEntriesOK struct {

	/*
	  In: Body
	*/
	Payload *models.Feed `json:"body,omitempty"`
}

// NewGetUsersMeEntriesOK creates GetUsersMeEntriesOK with default headers values
func NewGetUsersMeEntriesOK() *GetUsersMeEntriesOK {
	return &GetUsersMeEntriesOK{}
}

// WithPayload adds the payload to the get users me entries o k response
func (o *GetUsersMeEntriesOK) WithPayload(payload *models.Feed) *GetUsersMeEntriesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users me entries o k response
func (o *GetUsersMeEntriesOK) SetPayload(payload *models.Feed) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersMeEntriesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersMeEntriesForbiddenCode is the HTTP code returned for type GetUsersMeEntriesForbidden
const GetUsersMeEntriesForbiddenCode int = 403

/*GetUsersMeEntriesForbidden access denied

swagger:response getUsersMeEntriesForbidden
*/
type GetUsersMeEntriesForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersMeEntriesForbidden creates GetUsersMeEntriesForbidden with default headers values
func NewGetUsersMeEntriesForbidden() *GetUsersMeEntriesForbidden {
	return &GetUsersMeEntriesForbidden{}
}

// WithPayload adds the payload to the get users me entries forbidden response
func (o *GetUsersMeEntriesForbidden) WithPayload(payload *models.Error) *GetUsersMeEntriesForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users me entries forbidden response
func (o *GetUsersMeEntriesForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersMeEntriesForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
