// Code generated by go-swagger; DO NOT EDIT.

package watchings

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy-server/gen/models"
)

// PutEntriesIDWatchingOKCode is the HTTP code returned for type PutEntriesIDWatchingOK
const PutEntriesIDWatchingOKCode int = 200

/*PutEntriesIDWatchingOK watching status

swagger:response putEntriesIdWatchingOK
*/
type PutEntriesIDWatchingOK struct {

	/*
	  In: Body
	*/
	Payload *models.WatchingStatus `json:"body,omitempty"`
}

// NewPutEntriesIDWatchingOK creates PutEntriesIDWatchingOK with default headers values
func NewPutEntriesIDWatchingOK() *PutEntriesIDWatchingOK {
	return &PutEntriesIDWatchingOK{}
}

// WithPayload adds the payload to the put entries Id watching o k response
func (o *PutEntriesIDWatchingOK) WithPayload(payload *models.WatchingStatus) *PutEntriesIDWatchingOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put entries Id watching o k response
func (o *PutEntriesIDWatchingOK) SetPayload(payload *models.WatchingStatus) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutEntriesIDWatchingOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutEntriesIDWatchingForbiddenCode is the HTTP code returned for type PutEntriesIDWatchingForbidden
const PutEntriesIDWatchingForbiddenCode int = 403

/*PutEntriesIDWatchingForbidden access denied

swagger:response putEntriesIdWatchingForbidden
*/
type PutEntriesIDWatchingForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutEntriesIDWatchingForbidden creates PutEntriesIDWatchingForbidden with default headers values
func NewPutEntriesIDWatchingForbidden() *PutEntriesIDWatchingForbidden {
	return &PutEntriesIDWatchingForbidden{}
}

// WithPayload adds the payload to the put entries Id watching forbidden response
func (o *PutEntriesIDWatchingForbidden) WithPayload(payload *models.Error) *PutEntriesIDWatchingForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put entries Id watching forbidden response
func (o *PutEntriesIDWatchingForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutEntriesIDWatchingForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutEntriesIDWatchingNotFoundCode is the HTTP code returned for type PutEntriesIDWatchingNotFound
const PutEntriesIDWatchingNotFoundCode int = 404

/*PutEntriesIDWatchingNotFound Entry not found

swagger:response putEntriesIdWatchingNotFound
*/
type PutEntriesIDWatchingNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutEntriesIDWatchingNotFound creates PutEntriesIDWatchingNotFound with default headers values
func NewPutEntriesIDWatchingNotFound() *PutEntriesIDWatchingNotFound {
	return &PutEntriesIDWatchingNotFound{}
}

// WithPayload adds the payload to the put entries Id watching not found response
func (o *PutEntriesIDWatchingNotFound) WithPayload(payload *models.Error) *PutEntriesIDWatchingNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put entries Id watching not found response
func (o *PutEntriesIDWatchingNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutEntriesIDWatchingNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
