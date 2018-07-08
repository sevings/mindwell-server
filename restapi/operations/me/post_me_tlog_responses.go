// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// PostMeTlogCreatedCode is the HTTP code returned for type PostMeTlogCreated
const PostMeTlogCreatedCode int = 201

/*PostMeTlogCreated Entry data

swagger:response postMeTlogCreated
*/
type PostMeTlogCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Entry `json:"body,omitempty"`
}

// NewPostMeTlogCreated creates PostMeTlogCreated with default headers values
func NewPostMeTlogCreated() *PostMeTlogCreated {
	return &PostMeTlogCreated{}
}

// WithPayload adds the payload to the post me tlog created response
func (o *PostMeTlogCreated) WithPayload(payload *models.Entry) *PostMeTlogCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post me tlog created response
func (o *PostMeTlogCreated) SetPayload(payload *models.Entry) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMeTlogCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostMeTlogForbiddenCode is the HTTP code returned for type PostMeTlogForbidden
const PostMeTlogForbiddenCode int = 403

/*PostMeTlogForbidden access denied

swagger:response postMeTlogForbidden
*/
type PostMeTlogForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostMeTlogForbidden creates PostMeTlogForbidden with default headers values
func NewPostMeTlogForbidden() *PostMeTlogForbidden {
	return &PostMeTlogForbidden{}
}

// WithPayload adds the payload to the post me tlog forbidden response
func (o *PostMeTlogForbidden) WithPayload(payload *models.Error) *PostMeTlogForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post me tlog forbidden response
func (o *PostMeTlogForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMeTlogForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
