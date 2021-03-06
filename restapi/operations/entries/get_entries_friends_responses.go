// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// GetEntriesFriendsOKCode is the HTTP code returned for type GetEntriesFriendsOK
const GetEntriesFriendsOKCode int = 200

/*GetEntriesFriendsOK Entry list

swagger:response getEntriesFriendsOK
*/
type GetEntriesFriendsOK struct {

	/*
	  In: Body
	*/
	Payload *models.Feed `json:"body,omitempty"`
}

// NewGetEntriesFriendsOK creates GetEntriesFriendsOK with default headers values
func NewGetEntriesFriendsOK() *GetEntriesFriendsOK {

	return &GetEntriesFriendsOK{}
}

// WithPayload adds the payload to the get entries friends o k response
func (o *GetEntriesFriendsOK) WithPayload(payload *models.Feed) *GetEntriesFriendsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entries friends o k response
func (o *GetEntriesFriendsOK) SetPayload(payload *models.Feed) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntriesFriendsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
