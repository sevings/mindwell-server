// Code generated by go-swagger; DO NOT EDIT.

package account

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// GetAccountInvitesOKCode is the HTTP code returned for type GetAccountInvitesOK
const GetAccountInvitesOKCode int = 200

/*GetAccountInvitesOK invite list

swagger:response getAccountInvitesOK
*/
type GetAccountInvitesOK struct {

	/*
	  In: Body
	*/
	Payload *models.GetAccountInvitesOKBody `json:"body,omitempty"`
}

// NewGetAccountInvitesOK creates GetAccountInvitesOK with default headers values
func NewGetAccountInvitesOK() *GetAccountInvitesOK {
	return &GetAccountInvitesOK{}
}

// WithPayload adds the payload to the get account invites o k response
func (o *GetAccountInvitesOK) WithPayload(payload *models.GetAccountInvitesOKBody) *GetAccountInvitesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get account invites o k response
func (o *GetAccountInvitesOK) SetPayload(payload *models.GetAccountInvitesOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAccountInvitesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
