// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy-server/gen/models"
)

// GetEntriesUsersMeOKCode is the HTTP code returned for type GetEntriesUsersMeOK
const GetEntriesUsersMeOKCode int = 200

/*GetEntriesUsersMeOK Entry list

swagger:response getEntriesUsersMeOK
*/
type GetEntriesUsersMeOK struct {

	/*
	  In: Body
	*/
	Payload *models.Feed `json:"body,omitempty"`
}

// NewGetEntriesUsersMeOK creates GetEntriesUsersMeOK with default headers values
func NewGetEntriesUsersMeOK() *GetEntriesUsersMeOK {
	return &GetEntriesUsersMeOK{}
}

// WithPayload adds the payload to the get entries users me o k response
func (o *GetEntriesUsersMeOK) WithPayload(payload *models.Feed) *GetEntriesUsersMeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entries users me o k response
func (o *GetEntriesUsersMeOK) SetPayload(payload *models.Feed) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntriesUsersMeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetEntriesUsersMeForbiddenCode is the HTTP code returned for type GetEntriesUsersMeForbidden
const GetEntriesUsersMeForbiddenCode int = 403

/*GetEntriesUsersMeForbidden access denied

swagger:response getEntriesUsersMeForbidden
*/
type GetEntriesUsersMeForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetEntriesUsersMeForbidden creates GetEntriesUsersMeForbidden with default headers values
func NewGetEntriesUsersMeForbidden() *GetEntriesUsersMeForbidden {
	return &GetEntriesUsersMeForbidden{}
}

// WithPayload adds the payload to the get entries users me forbidden response
func (o *GetEntriesUsersMeForbidden) WithPayload(payload *models.Error) *GetEntriesUsersMeForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entries users me forbidden response
func (o *GetEntriesUsersMeForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntriesUsersMeForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}