// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// PutMeAvatarOKCode is the HTTP code returned for type PutMeAvatarOK
const PutMeAvatarOKCode int = 200

/*PutMeAvatarOK Avatar

swagger:response putMeAvatarOK
*/
type PutMeAvatarOK struct {

	/*
	  In: Body
	*/
	Payload *models.Avatar `json:"body,omitempty"`
}

// NewPutMeAvatarOK creates PutMeAvatarOK with default headers values
func NewPutMeAvatarOK() *PutMeAvatarOK {
	return &PutMeAvatarOK{}
}

// WithPayload adds the payload to the put me avatar o k response
func (o *PutMeAvatarOK) WithPayload(payload *models.Avatar) *PutMeAvatarOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put me avatar o k response
func (o *PutMeAvatarOK) SetPayload(payload *models.Avatar) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutMeAvatarOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutMeAvatarBadRequestCode is the HTTP code returned for type PutMeAvatarBadRequest
const PutMeAvatarBadRequestCode int = 400

/*PutMeAvatarBadRequest bad request

swagger:response putMeAvatarBadRequest
*/
type PutMeAvatarBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutMeAvatarBadRequest creates PutMeAvatarBadRequest with default headers values
func NewPutMeAvatarBadRequest() *PutMeAvatarBadRequest {
	return &PutMeAvatarBadRequest{}
}

// WithPayload adds the payload to the put me avatar bad request response
func (o *PutMeAvatarBadRequest) WithPayload(payload *models.Error) *PutMeAvatarBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put me avatar bad request response
func (o *PutMeAvatarBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutMeAvatarBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
