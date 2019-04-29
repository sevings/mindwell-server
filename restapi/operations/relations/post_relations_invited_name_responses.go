// Code generated by go-swagger; DO NOT EDIT.

package relations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/sevings/mindwell-server/models"
)

// PostRelationsInvitedNameNoContentCode is the HTTP code returned for type PostRelationsInvitedNameNoContent
const PostRelationsInvitedNameNoContentCode int = 204

/*PostRelationsInvitedNameNoContent invited

swagger:response postRelationsInvitedNameNoContent
*/
type PostRelationsInvitedNameNoContent struct {
}

// NewPostRelationsInvitedNameNoContent creates PostRelationsInvitedNameNoContent with default headers values
func NewPostRelationsInvitedNameNoContent() *PostRelationsInvitedNameNoContent {

	return &PostRelationsInvitedNameNoContent{}
}

// WriteResponse to the client
func (o *PostRelationsInvitedNameNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PostRelationsInvitedNameForbiddenCode is the HTTP code returned for type PostRelationsInvitedNameForbidden
const PostRelationsInvitedNameForbiddenCode int = 403

/*PostRelationsInvitedNameForbidden invalid invite or the user is invited already

swagger:response postRelationsInvitedNameForbidden
*/
type PostRelationsInvitedNameForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostRelationsInvitedNameForbidden creates PostRelationsInvitedNameForbidden with default headers values
func NewPostRelationsInvitedNameForbidden() *PostRelationsInvitedNameForbidden {

	return &PostRelationsInvitedNameForbidden{}
}

// WithPayload adds the payload to the post relations invited name forbidden response
func (o *PostRelationsInvitedNameForbidden) WithPayload(payload *models.Error) *PostRelationsInvitedNameForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post relations invited name forbidden response
func (o *PostRelationsInvitedNameForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostRelationsInvitedNameForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostRelationsInvitedNameNotFoundCode is the HTTP code returned for type PostRelationsInvitedNameNotFound
const PostRelationsInvitedNameNotFoundCode int = 404

/*PostRelationsInvitedNameNotFound User not found

swagger:response postRelationsInvitedNameNotFound
*/
type PostRelationsInvitedNameNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostRelationsInvitedNameNotFound creates PostRelationsInvitedNameNotFound with default headers values
func NewPostRelationsInvitedNameNotFound() *PostRelationsInvitedNameNotFound {

	return &PostRelationsInvitedNameNotFound{}
}

// WithPayload adds the payload to the post relations invited name not found response
func (o *PostRelationsInvitedNameNotFound) WithPayload(payload *models.Error) *PostRelationsInvitedNameNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post relations invited name not found response
func (o *PostRelationsInvitedNameNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostRelationsInvitedNameNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}