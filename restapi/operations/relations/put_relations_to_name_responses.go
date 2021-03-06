// Code generated by go-swagger; DO NOT EDIT.

package relations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// PutRelationsToNameOKCode is the HTTP code returned for type PutRelationsToNameOK
const PutRelationsToNameOKCode int = 200

/*PutRelationsToNameOK your relationship with the user

swagger:response putRelationsToNameOK
*/
type PutRelationsToNameOK struct {

	/*
	  In: Body
	*/
	Payload *models.Relationship `json:"body,omitempty"`
}

// NewPutRelationsToNameOK creates PutRelationsToNameOK with default headers values
func NewPutRelationsToNameOK() *PutRelationsToNameOK {

	return &PutRelationsToNameOK{}
}

// WithPayload adds the payload to the put relations to name o k response
func (o *PutRelationsToNameOK) WithPayload(payload *models.Relationship) *PutRelationsToNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put relations to name o k response
func (o *PutRelationsToNameOK) SetPayload(payload *models.Relationship) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutRelationsToNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutRelationsToNameForbiddenCode is the HTTP code returned for type PutRelationsToNameForbidden
const PutRelationsToNameForbiddenCode int = 403

/*PutRelationsToNameForbidden access denied

swagger:response putRelationsToNameForbidden
*/
type PutRelationsToNameForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutRelationsToNameForbidden creates PutRelationsToNameForbidden with default headers values
func NewPutRelationsToNameForbidden() *PutRelationsToNameForbidden {

	return &PutRelationsToNameForbidden{}
}

// WithPayload adds the payload to the put relations to name forbidden response
func (o *PutRelationsToNameForbidden) WithPayload(payload *models.Error) *PutRelationsToNameForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put relations to name forbidden response
func (o *PutRelationsToNameForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutRelationsToNameForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutRelationsToNameNotFoundCode is the HTTP code returned for type PutRelationsToNameNotFound
const PutRelationsToNameNotFoundCode int = 404

/*PutRelationsToNameNotFound User not found

swagger:response putRelationsToNameNotFound
*/
type PutRelationsToNameNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutRelationsToNameNotFound creates PutRelationsToNameNotFound with default headers values
func NewPutRelationsToNameNotFound() *PutRelationsToNameNotFound {

	return &PutRelationsToNameNotFound{}
}

// WithPayload adds the payload to the put relations to name not found response
func (o *PutRelationsToNameNotFound) WithPayload(payload *models.Error) *PutRelationsToNameNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put relations to name not found response
func (o *PutRelationsToNameNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutRelationsToNameNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
