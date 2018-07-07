// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// GetUsersNameFollowersOKCode is the HTTP code returned for type GetUsersNameFollowersOK
const GetUsersNameFollowersOKCode int = 200

/*GetUsersNameFollowersOK User list

swagger:response getUsersNameFollowersOK
*/
type GetUsersNameFollowersOK struct {

	/*
	  In: Body
	*/
	Payload *models.FriendList `json:"body,omitempty"`
}

// NewGetUsersNameFollowersOK creates GetUsersNameFollowersOK with default headers values
func NewGetUsersNameFollowersOK() *GetUsersNameFollowersOK {
	return &GetUsersNameFollowersOK{}
}

// WithPayload adds the payload to the get users name followers o k response
func (o *GetUsersNameFollowersOK) WithPayload(payload *models.FriendList) *GetUsersNameFollowersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users name followers o k response
func (o *GetUsersNameFollowersOK) SetPayload(payload *models.FriendList) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersNameFollowersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersNameFollowersForbiddenCode is the HTTP code returned for type GetUsersNameFollowersForbidden
const GetUsersNameFollowersForbiddenCode int = 403

/*GetUsersNameFollowersForbidden access denied

swagger:response getUsersNameFollowersForbidden
*/
type GetUsersNameFollowersForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersNameFollowersForbidden creates GetUsersNameFollowersForbidden with default headers values
func NewGetUsersNameFollowersForbidden() *GetUsersNameFollowersForbidden {
	return &GetUsersNameFollowersForbidden{}
}

// WithPayload adds the payload to the get users name followers forbidden response
func (o *GetUsersNameFollowersForbidden) WithPayload(payload *models.Error) *GetUsersNameFollowersForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users name followers forbidden response
func (o *GetUsersNameFollowersForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersNameFollowersForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersNameFollowersNotFoundCode is the HTTP code returned for type GetUsersNameFollowersNotFound
const GetUsersNameFollowersNotFoundCode int = 404

/*GetUsersNameFollowersNotFound User not found

swagger:response getUsersNameFollowersNotFound
*/
type GetUsersNameFollowersNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersNameFollowersNotFound creates GetUsersNameFollowersNotFound with default headers values
func NewGetUsersNameFollowersNotFound() *GetUsersNameFollowersNotFound {
	return &GetUsersNameFollowersNotFound{}
}

// WithPayload adds the payload to the get users name followers not found response
func (o *GetUsersNameFollowersNotFound) WithPayload(payload *models.Error) *GetUsersNameFollowersNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users name followers not found response
func (o *GetUsersNameFollowersNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersNameFollowersNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
