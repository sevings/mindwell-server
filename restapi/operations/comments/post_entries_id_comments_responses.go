// Code generated by go-swagger; DO NOT EDIT.

package comments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// PostEntriesIDCommentsCreatedCode is the HTTP code returned for type PostEntriesIDCommentsCreated
const PostEntriesIDCommentsCreatedCode int = 201

/*PostEntriesIDCommentsCreated Comment data

swagger:response postEntriesIdCommentsCreated
*/
type PostEntriesIDCommentsCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Comment `json:"body,omitempty"`
}

// NewPostEntriesIDCommentsCreated creates PostEntriesIDCommentsCreated with default headers values
func NewPostEntriesIDCommentsCreated() *PostEntriesIDCommentsCreated {
	return &PostEntriesIDCommentsCreated{}
}

// WithPayload adds the payload to the post entries Id comments created response
func (o *PostEntriesIDCommentsCreated) WithPayload(payload *models.Comment) *PostEntriesIDCommentsCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post entries Id comments created response
func (o *PostEntriesIDCommentsCreated) SetPayload(payload *models.Comment) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEntriesIDCommentsCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostEntriesIDCommentsForbiddenCode is the HTTP code returned for type PostEntriesIDCommentsForbidden
const PostEntriesIDCommentsForbiddenCode int = 403

/*PostEntriesIDCommentsForbidden access denied

swagger:response postEntriesIdCommentsForbidden
*/
type PostEntriesIDCommentsForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostEntriesIDCommentsForbidden creates PostEntriesIDCommentsForbidden with default headers values
func NewPostEntriesIDCommentsForbidden() *PostEntriesIDCommentsForbidden {
	return &PostEntriesIDCommentsForbidden{}
}

// WithPayload adds the payload to the post entries Id comments forbidden response
func (o *PostEntriesIDCommentsForbidden) WithPayload(payload *models.Error) *PostEntriesIDCommentsForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post entries Id comments forbidden response
func (o *PostEntriesIDCommentsForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEntriesIDCommentsForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostEntriesIDCommentsNotFoundCode is the HTTP code returned for type PostEntriesIDCommentsNotFound
const PostEntriesIDCommentsNotFoundCode int = 404

/*PostEntriesIDCommentsNotFound Entry not found

swagger:response postEntriesIdCommentsNotFound
*/
type PostEntriesIDCommentsNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostEntriesIDCommentsNotFound creates PostEntriesIDCommentsNotFound with default headers values
func NewPostEntriesIDCommentsNotFound() *PostEntriesIDCommentsNotFound {
	return &PostEntriesIDCommentsNotFound{}
}

// WithPayload adds the payload to the post entries Id comments not found response
func (o *PostEntriesIDCommentsNotFound) WithPayload(payload *models.Error) *PostEntriesIDCommentsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post entries Id comments not found response
func (o *PostEntriesIDCommentsNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEntriesIDCommentsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
