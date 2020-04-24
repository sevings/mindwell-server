// Code generated by go-swagger; DO NOT EDIT.

package chats

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/sevings/mindwell-server/models"
)

// DeleteMessagesIDOKCode is the HTTP code returned for type DeleteMessagesIDOK
const DeleteMessagesIDOKCode int = 200

/*DeleteMessagesIDOK OK

swagger:response deleteMessagesIdOK
*/
type DeleteMessagesIDOK struct {
}

// NewDeleteMessagesIDOK creates DeleteMessagesIDOK with default headers values
func NewDeleteMessagesIDOK() *DeleteMessagesIDOK {

	return &DeleteMessagesIDOK{}
}

// WriteResponse to the client
func (o *DeleteMessagesIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteMessagesIDForbiddenCode is the HTTP code returned for type DeleteMessagesIDForbidden
const DeleteMessagesIDForbiddenCode int = 403

/*DeleteMessagesIDForbidden access denied

swagger:response deleteMessagesIdForbidden
*/
type DeleteMessagesIDForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteMessagesIDForbidden creates DeleteMessagesIDForbidden with default headers values
func NewDeleteMessagesIDForbidden() *DeleteMessagesIDForbidden {

	return &DeleteMessagesIDForbidden{}
}

// WithPayload adds the payload to the delete messages Id forbidden response
func (o *DeleteMessagesIDForbidden) WithPayload(payload *models.Error) *DeleteMessagesIDForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete messages Id forbidden response
func (o *DeleteMessagesIDForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteMessagesIDForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteMessagesIDNotFoundCode is the HTTP code returned for type DeleteMessagesIDNotFound
const DeleteMessagesIDNotFoundCode int = 404

/*DeleteMessagesIDNotFound Message not found

swagger:response deleteMessagesIdNotFound
*/
type DeleteMessagesIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteMessagesIDNotFound creates DeleteMessagesIDNotFound with default headers values
func NewDeleteMessagesIDNotFound() *DeleteMessagesIDNotFound {

	return &DeleteMessagesIDNotFound{}
}

// WithPayload adds the payload to the delete messages Id not found response
func (o *DeleteMessagesIDNotFound) WithPayload(payload *models.Error) *DeleteMessagesIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete messages Id not found response
func (o *DeleteMessagesIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteMessagesIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}