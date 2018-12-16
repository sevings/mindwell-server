// Code generated by go-swagger; DO NOT EDIT.

package notifications

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PutNotificationsReadOKCode is the HTTP code returned for type PutNotificationsReadOK
const PutNotificationsReadOKCode int = 200

/*PutNotificationsReadOK unread count

swagger:response putNotificationsReadOK
*/
type PutNotificationsReadOK struct {

	/*
	  In: Body
	*/
	Payload *PutNotificationsReadOKBody `json:"body,omitempty"`
}

// NewPutNotificationsReadOK creates PutNotificationsReadOK with default headers values
func NewPutNotificationsReadOK() *PutNotificationsReadOK {

	return &PutNotificationsReadOK{}
}

// WithPayload adds the payload to the put notifications read o k response
func (o *PutNotificationsReadOK) WithPayload(payload *PutNotificationsReadOKBody) *PutNotificationsReadOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put notifications read o k response
func (o *PutNotificationsReadOK) SetPayload(payload *PutNotificationsReadOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutNotificationsReadOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
