// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/sevings/mindwell-server/models"
)

// PutMeOKCode is the HTTP code returned for type PutMeOK
const PutMeOKCode int = 200

/*PutMeOK your data

swagger:response putMeOK
*/
type PutMeOK struct {

	/*
	  In: Body
	*/
	Payload *models.Profile `json:"body,omitempty"`
}

// NewPutMeOK creates PutMeOK with default headers values
func NewPutMeOK() *PutMeOK {

	return &PutMeOK{}
}

// WithPayload adds the payload to the put me o k response
func (o *PutMeOK) WithPayload(payload *models.Profile) *PutMeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put me o k response
func (o *PutMeOK) SetPayload(payload *models.Profile) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutMeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
