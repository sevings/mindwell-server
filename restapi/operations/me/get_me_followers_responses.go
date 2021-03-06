// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// GetMeFollowersOKCode is the HTTP code returned for type GetMeFollowersOK
const GetMeFollowersOKCode int = 200

/*GetMeFollowersOK User list

swagger:response getMeFollowersOK
*/
type GetMeFollowersOK struct {

	/*
	  In: Body
	*/
	Payload *models.FriendList `json:"body,omitempty"`
}

// NewGetMeFollowersOK creates GetMeFollowersOK with default headers values
func NewGetMeFollowersOK() *GetMeFollowersOK {

	return &GetMeFollowersOK{}
}

// WithPayload adds the payload to the get me followers o k response
func (o *GetMeFollowersOK) WithPayload(payload *models.FriendList) *GetMeFollowersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get me followers o k response
func (o *GetMeFollowersOK) SetPayload(payload *models.FriendList) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMeFollowersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
