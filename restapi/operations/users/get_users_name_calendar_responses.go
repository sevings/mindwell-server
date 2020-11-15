// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/sevings/mindwell-server/models"
)

// GetUsersNameCalendarOKCode is the HTTP code returned for type GetUsersNameCalendarOK
const GetUsersNameCalendarOKCode int = 200

/*GetUsersNameCalendarOK Entry list

swagger:response getUsersNameCalendarOK
*/
type GetUsersNameCalendarOK struct {

	/*
	  In: Body
	*/
	Payload *models.Calendar `json:"body,omitempty"`
}

// NewGetUsersNameCalendarOK creates GetUsersNameCalendarOK with default headers values
func NewGetUsersNameCalendarOK() *GetUsersNameCalendarOK {

	return &GetUsersNameCalendarOK{}
}

// WithPayload adds the payload to the get users name calendar o k response
func (o *GetUsersNameCalendarOK) WithPayload(payload *models.Calendar) *GetUsersNameCalendarOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users name calendar o k response
func (o *GetUsersNameCalendarOK) SetPayload(payload *models.Calendar) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersNameCalendarOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersNameCalendarNotFoundCode is the HTTP code returned for type GetUsersNameCalendarNotFound
const GetUsersNameCalendarNotFoundCode int = 404

/*GetUsersNameCalendarNotFound User not found

swagger:response getUsersNameCalendarNotFound
*/
type GetUsersNameCalendarNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersNameCalendarNotFound creates GetUsersNameCalendarNotFound with default headers values
func NewGetUsersNameCalendarNotFound() *GetUsersNameCalendarNotFound {

	return &GetUsersNameCalendarNotFound{}
}

// WithPayload adds the payload to the get users name calendar not found response
func (o *GetUsersNameCalendarNotFound) WithPayload(payload *models.Error) *GetUsersNameCalendarNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users name calendar not found response
func (o *GetUsersNameCalendarNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersNameCalendarNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}