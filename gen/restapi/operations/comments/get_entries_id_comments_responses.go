// Code generated by go-swagger; DO NOT EDIT.

package comments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy-server/gen/models"
)

// GetEntriesIDCommentsOKCode is the HTTP code returned for type GetEntriesIDCommentsOK
const GetEntriesIDCommentsOKCode int = 200

/*GetEntriesIDCommentsOK comments list

swagger:response getEntriesIdCommentsOK
*/
type GetEntriesIDCommentsOK struct {

	/*
	  In: Body
	*/
	Payload *models.CommentList `json:"body,omitempty"`
}

// NewGetEntriesIDCommentsOK creates GetEntriesIDCommentsOK with default headers values
func NewGetEntriesIDCommentsOK() *GetEntriesIDCommentsOK {
	return &GetEntriesIDCommentsOK{}
}

// WithPayload adds the payload to the get entries Id comments o k response
func (o *GetEntriesIDCommentsOK) WithPayload(payload *models.CommentList) *GetEntriesIDCommentsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entries Id comments o k response
func (o *GetEntriesIDCommentsOK) SetPayload(payload *models.CommentList) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntriesIDCommentsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetEntriesIDCommentsForbiddenCode is the HTTP code returned for type GetEntriesIDCommentsForbidden
const GetEntriesIDCommentsForbiddenCode int = 403

/*GetEntriesIDCommentsForbidden access denied

swagger:response getEntriesIdCommentsForbidden
*/
type GetEntriesIDCommentsForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetEntriesIDCommentsForbidden creates GetEntriesIDCommentsForbidden with default headers values
func NewGetEntriesIDCommentsForbidden() *GetEntriesIDCommentsForbidden {
	return &GetEntriesIDCommentsForbidden{}
}

// WithPayload adds the payload to the get entries Id comments forbidden response
func (o *GetEntriesIDCommentsForbidden) WithPayload(payload *models.Error) *GetEntriesIDCommentsForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entries Id comments forbidden response
func (o *GetEntriesIDCommentsForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntriesIDCommentsForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetEntriesIDCommentsNotFoundCode is the HTTP code returned for type GetEntriesIDCommentsNotFound
const GetEntriesIDCommentsNotFoundCode int = 404

/*GetEntriesIDCommentsNotFound Entry not found

swagger:response getEntriesIdCommentsNotFound
*/
type GetEntriesIDCommentsNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetEntriesIDCommentsNotFound creates GetEntriesIDCommentsNotFound with default headers values
func NewGetEntriesIDCommentsNotFound() *GetEntriesIDCommentsNotFound {
	return &GetEntriesIDCommentsNotFound{}
}

// WithPayload adds the payload to the get entries Id comments not found response
func (o *GetEntriesIDCommentsNotFound) WithPayload(payload *models.Error) *GetEntriesIDCommentsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entries Id comments not found response
func (o *GetEntriesIDCommentsNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntriesIDCommentsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
