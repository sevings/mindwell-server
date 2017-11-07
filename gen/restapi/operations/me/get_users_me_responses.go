// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy/gen/models"
)

// GetUsersMeOKCode is the HTTP code returned for type GetUsersMeOK
const GetUsersMeOKCode int = 200

/*GetUsersMeOK your data

swagger:response getUsersMeOK
*/
type GetUsersMeOK struct {

	/*
	  In: Body
	*/
	Payload *models.AuthProfile `json:"body,omitempty"`
}

// NewGetUsersMeOK creates GetUsersMeOK with default headers values
func NewGetUsersMeOK() *GetUsersMeOK {
	return &GetUsersMeOK{}
}

// WithPayload adds the payload to the get users me o k response
func (o *GetUsersMeOK) WithPayload(payload *models.AuthProfile) *GetUsersMeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users me o k response
func (o *GetUsersMeOK) SetPayload(payload *models.AuthProfile) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersMeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersMeForbiddenCode is the HTTP code returned for type GetUsersMeForbidden
const GetUsersMeForbiddenCode int = 403

/*GetUsersMeForbidden access denied

swagger:response getUsersMeForbidden
*/
type GetUsersMeForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersMeForbidden creates GetUsersMeForbidden with default headers values
func NewGetUsersMeForbidden() *GetUsersMeForbidden {
	return &GetUsersMeForbidden{}
}

// WithPayload adds the payload to the get users me forbidden response
func (o *GetUsersMeForbidden) WithPayload(payload *models.Error) *GetUsersMeForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users me forbidden response
func (o *GetUsersMeForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersMeForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
