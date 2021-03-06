// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// GetUsersNameInvitedOKCode is the HTTP code returned for type GetUsersNameInvitedOK
const GetUsersNameInvitedOKCode int = 200

/*GetUsersNameInvitedOK User list

swagger:response getUsersNameInvitedOK
*/
type GetUsersNameInvitedOK struct {

	/*
	  In: Body
	*/
	Payload *models.FriendList `json:"body,omitempty"`
}

// NewGetUsersNameInvitedOK creates GetUsersNameInvitedOK with default headers values
func NewGetUsersNameInvitedOK() *GetUsersNameInvitedOK {

	return &GetUsersNameInvitedOK{}
}

// WithPayload adds the payload to the get users name invited o k response
func (o *GetUsersNameInvitedOK) WithPayload(payload *models.FriendList) *GetUsersNameInvitedOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users name invited o k response
func (o *GetUsersNameInvitedOK) SetPayload(payload *models.FriendList) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersNameInvitedOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersNameInvitedForbiddenCode is the HTTP code returned for type GetUsersNameInvitedForbidden
const GetUsersNameInvitedForbiddenCode int = 403

/*GetUsersNameInvitedForbidden access denied

swagger:response getUsersNameInvitedForbidden
*/
type GetUsersNameInvitedForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersNameInvitedForbidden creates GetUsersNameInvitedForbidden with default headers values
func NewGetUsersNameInvitedForbidden() *GetUsersNameInvitedForbidden {

	return &GetUsersNameInvitedForbidden{}
}

// WithPayload adds the payload to the get users name invited forbidden response
func (o *GetUsersNameInvitedForbidden) WithPayload(payload *models.Error) *GetUsersNameInvitedForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users name invited forbidden response
func (o *GetUsersNameInvitedForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersNameInvitedForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersNameInvitedNotFoundCode is the HTTP code returned for type GetUsersNameInvitedNotFound
const GetUsersNameInvitedNotFoundCode int = 404

/*GetUsersNameInvitedNotFound User not found

swagger:response getUsersNameInvitedNotFound
*/
type GetUsersNameInvitedNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersNameInvitedNotFound creates GetUsersNameInvitedNotFound with default headers values
func NewGetUsersNameInvitedNotFound() *GetUsersNameInvitedNotFound {

	return &GetUsersNameInvitedNotFound{}
}

// WithPayload adds the payload to the get users name invited not found response
func (o *GetUsersNameInvitedNotFound) WithPayload(payload *models.Error) *GetUsersNameInvitedNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users name invited not found response
func (o *GetUsersNameInvitedNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersNameInvitedNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
