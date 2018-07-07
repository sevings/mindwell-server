// Code generated by go-swagger; DO NOT EDIT.

package relations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// DeleteRelationsToNameOKCode is the HTTP code returned for type DeleteRelationsToNameOK
const DeleteRelationsToNameOKCode int = 200

/*DeleteRelationsToNameOK your relationship with the user

swagger:response deleteRelationsToNameOK
*/
type DeleteRelationsToNameOK struct {

	/*
	  In: Body
	*/
	Payload *models.Relationship `json:"body,omitempty"`
}

// NewDeleteRelationsToNameOK creates DeleteRelationsToNameOK with default headers values
func NewDeleteRelationsToNameOK() *DeleteRelationsToNameOK {
	return &DeleteRelationsToNameOK{}
}

// WithPayload adds the payload to the delete relations to name o k response
func (o *DeleteRelationsToNameOK) WithPayload(payload *models.Relationship) *DeleteRelationsToNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete relations to name o k response
func (o *DeleteRelationsToNameOK) SetPayload(payload *models.Relationship) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteRelationsToNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteRelationsToNameNotFoundCode is the HTTP code returned for type DeleteRelationsToNameNotFound
const DeleteRelationsToNameNotFoundCode int = 404

/*DeleteRelationsToNameNotFound User not found

swagger:response deleteRelationsToNameNotFound
*/
type DeleteRelationsToNameNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteRelationsToNameNotFound creates DeleteRelationsToNameNotFound with default headers values
func NewDeleteRelationsToNameNotFound() *DeleteRelationsToNameNotFound {
	return &DeleteRelationsToNameNotFound{}
}

// WithPayload adds the payload to the delete relations to name not found response
func (o *DeleteRelationsToNameNotFound) WithPayload(payload *models.Error) *DeleteRelationsToNameNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete relations to name not found response
func (o *DeleteRelationsToNameNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteRelationsToNameNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
