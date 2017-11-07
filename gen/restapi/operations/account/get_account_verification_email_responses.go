// Code generated by go-swagger; DO NOT EDIT.

package account

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy-server/gen/models"
)

// GetAccountVerificationEmailOKCode is the HTTP code returned for type GetAccountVerificationEmailOK
const GetAccountVerificationEmailOKCode int = 200

/*GetAccountVerificationEmailOK verified

swagger:response getAccountVerificationEmailOK
*/
type GetAccountVerificationEmailOK struct {
}

// NewGetAccountVerificationEmailOK creates GetAccountVerificationEmailOK with default headers values
func NewGetAccountVerificationEmailOK() *GetAccountVerificationEmailOK {
	return &GetAccountVerificationEmailOK{}
}

// WriteResponse to the client
func (o *GetAccountVerificationEmailOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
}

// GetAccountVerificationEmailBadRequestCode is the HTTP code returned for type GetAccountVerificationEmailBadRequest
const GetAccountVerificationEmailBadRequestCode int = 400

/*GetAccountVerificationEmailBadRequest code or email is not valid

swagger:response getAccountVerificationEmailBadRequest
*/
type GetAccountVerificationEmailBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAccountVerificationEmailBadRequest creates GetAccountVerificationEmailBadRequest with default headers values
func NewGetAccountVerificationEmailBadRequest() *GetAccountVerificationEmailBadRequest {
	return &GetAccountVerificationEmailBadRequest{}
}

// WithPayload adds the payload to the get account verification email bad request response
func (o *GetAccountVerificationEmailBadRequest) WithPayload(payload *models.Error) *GetAccountVerificationEmailBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get account verification email bad request response
func (o *GetAccountVerificationEmailBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAccountVerificationEmailBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
