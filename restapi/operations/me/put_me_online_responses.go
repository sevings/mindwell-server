// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PutMeOnlineOKCode is the HTTP code returned for type PutMeOnlineOK
const PutMeOnlineOKCode int = 200

/*PutMeOnlineOK unread counts

swagger:response putMeOnlineOK
*/
type PutMeOnlineOK struct {

	/*
	  In: Body
	*/
	Payload *PutMeOnlineOKBody `json:"body,omitempty"`
}

// NewPutMeOnlineOK creates PutMeOnlineOK with default headers values
func NewPutMeOnlineOK() *PutMeOnlineOK {

	return &PutMeOnlineOK{}
}

// WithPayload adds the payload to the put me online o k response
func (o *PutMeOnlineOK) WithPayload(payload *PutMeOnlineOKBody) *PutMeOnlineOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put me online o k response
func (o *PutMeOnlineOK) SetPayload(payload *PutMeOnlineOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutMeOnlineOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
