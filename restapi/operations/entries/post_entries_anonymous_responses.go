// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy-server/models"
)

// PostEntriesAnonymousCreatedCode is the HTTP code returned for type PostEntriesAnonymousCreated
const PostEntriesAnonymousCreatedCode int = 201

/*PostEntriesAnonymousCreated Entry data

swagger:response postEntriesAnonymousCreated
*/
type PostEntriesAnonymousCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Entry `json:"body,omitempty"`
}

// NewPostEntriesAnonymousCreated creates PostEntriesAnonymousCreated with default headers values
func NewPostEntriesAnonymousCreated() *PostEntriesAnonymousCreated {
	return &PostEntriesAnonymousCreated{}
}

// WithPayload adds the payload to the post entries anonymous created response
func (o *PostEntriesAnonymousCreated) WithPayload(payload *models.Entry) *PostEntriesAnonymousCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post entries anonymous created response
func (o *PostEntriesAnonymousCreated) SetPayload(payload *models.Entry) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEntriesAnonymousCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostEntriesAnonymousForbiddenCode is the HTTP code returned for type PostEntriesAnonymousForbidden
const PostEntriesAnonymousForbiddenCode int = 403

/*PostEntriesAnonymousForbidden access denied

swagger:response postEntriesAnonymousForbidden
*/
type PostEntriesAnonymousForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostEntriesAnonymousForbidden creates PostEntriesAnonymousForbidden with default headers values
func NewPostEntriesAnonymousForbidden() *PostEntriesAnonymousForbidden {
	return &PostEntriesAnonymousForbidden{}
}

// WithPayload adds the payload to the post entries anonymous forbidden response
func (o *PostEntriesAnonymousForbidden) WithPayload(payload *models.Error) *PostEntriesAnonymousForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post entries anonymous forbidden response
func (o *PostEntriesAnonymousForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEntriesAnonymousForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
