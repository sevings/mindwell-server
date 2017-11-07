// Code generated by go-swagger; DO NOT EDIT.

package watchings

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy-server/gen/models"
)

// GetEntriesIDWatchingOKCode is the HTTP code returned for type GetEntriesIDWatchingOK
const GetEntriesIDWatchingOKCode int = 200

/*GetEntriesIDWatchingOK watching status

swagger:response getEntriesIdWatchingOK
*/
type GetEntriesIDWatchingOK struct {

	/*
	  In: Body
	*/
	Payload *models.WatchingStatus `json:"body,omitempty"`
}

// NewGetEntriesIDWatchingOK creates GetEntriesIDWatchingOK with default headers values
func NewGetEntriesIDWatchingOK() *GetEntriesIDWatchingOK {
	return &GetEntriesIDWatchingOK{}
}

// WithPayload adds the payload to the get entries Id watching o k response
func (o *GetEntriesIDWatchingOK) WithPayload(payload *models.WatchingStatus) *GetEntriesIDWatchingOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entries Id watching o k response
func (o *GetEntriesIDWatchingOK) SetPayload(payload *models.WatchingStatus) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntriesIDWatchingOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetEntriesIDWatchingForbiddenCode is the HTTP code returned for type GetEntriesIDWatchingForbidden
const GetEntriesIDWatchingForbiddenCode int = 403

/*GetEntriesIDWatchingForbidden access denied

swagger:response getEntriesIdWatchingForbidden
*/
type GetEntriesIDWatchingForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetEntriesIDWatchingForbidden creates GetEntriesIDWatchingForbidden with default headers values
func NewGetEntriesIDWatchingForbidden() *GetEntriesIDWatchingForbidden {
	return &GetEntriesIDWatchingForbidden{}
}

// WithPayload adds the payload to the get entries Id watching forbidden response
func (o *GetEntriesIDWatchingForbidden) WithPayload(payload *models.Error) *GetEntriesIDWatchingForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entries Id watching forbidden response
func (o *GetEntriesIDWatchingForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntriesIDWatchingForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetEntriesIDWatchingNotFoundCode is the HTTP code returned for type GetEntriesIDWatchingNotFound
const GetEntriesIDWatchingNotFoundCode int = 404

/*GetEntriesIDWatchingNotFound Entry not found

swagger:response getEntriesIdWatchingNotFound
*/
type GetEntriesIDWatchingNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetEntriesIDWatchingNotFound creates GetEntriesIDWatchingNotFound with default headers values
func NewGetEntriesIDWatchingNotFound() *GetEntriesIDWatchingNotFound {
	return &GetEntriesIDWatchingNotFound{}
}

// WithPayload adds the payload to the get entries Id watching not found response
func (o *GetEntriesIDWatchingNotFound) WithPayload(payload *models.Error) *GetEntriesIDWatchingNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entries Id watching not found response
func (o *GetEntriesIDWatchingNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntriesIDWatchingNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
