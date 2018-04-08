// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// GetUsersIDFollowingsOKCode is the HTTP code returned for type GetUsersIDFollowingsOK
const GetUsersIDFollowingsOKCode int = 200

/*GetUsersIDFollowingsOK User list

swagger:response getUsersIdFollowingsOK
*/
type GetUsersIDFollowingsOK struct {

	/*
	  In: Body
	*/
	Payload *models.UserList `json:"body,omitempty"`
}

// NewGetUsersIDFollowingsOK creates GetUsersIDFollowingsOK with default headers values
func NewGetUsersIDFollowingsOK() *GetUsersIDFollowingsOK {
	return &GetUsersIDFollowingsOK{}
}

// WithPayload adds the payload to the get users Id followings o k response
func (o *GetUsersIDFollowingsOK) WithPayload(payload *models.UserList) *GetUsersIDFollowingsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users Id followings o k response
func (o *GetUsersIDFollowingsOK) SetPayload(payload *models.UserList) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersIDFollowingsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersIDFollowingsForbiddenCode is the HTTP code returned for type GetUsersIDFollowingsForbidden
const GetUsersIDFollowingsForbiddenCode int = 403

/*GetUsersIDFollowingsForbidden access denied

swagger:response getUsersIdFollowingsForbidden
*/
type GetUsersIDFollowingsForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersIDFollowingsForbidden creates GetUsersIDFollowingsForbidden with default headers values
func NewGetUsersIDFollowingsForbidden() *GetUsersIDFollowingsForbidden {
	return &GetUsersIDFollowingsForbidden{}
}

// WithPayload adds the payload to the get users Id followings forbidden response
func (o *GetUsersIDFollowingsForbidden) WithPayload(payload *models.Error) *GetUsersIDFollowingsForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users Id followings forbidden response
func (o *GetUsersIDFollowingsForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersIDFollowingsForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersIDFollowingsNotFoundCode is the HTTP code returned for type GetUsersIDFollowingsNotFound
const GetUsersIDFollowingsNotFoundCode int = 404

/*GetUsersIDFollowingsNotFound User not found

swagger:response getUsersIdFollowingsNotFound
*/
type GetUsersIDFollowingsNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersIDFollowingsNotFound creates GetUsersIDFollowingsNotFound with default headers values
func NewGetUsersIDFollowingsNotFound() *GetUsersIDFollowingsNotFound {
	return &GetUsersIDFollowingsNotFound{}
}

// WithPayload adds the payload to the get users Id followings not found response
func (o *GetUsersIDFollowingsNotFound) WithPayload(payload *models.Error) *GetUsersIDFollowingsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users Id followings not found response
func (o *GetUsersIDFollowingsNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersIDFollowingsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
