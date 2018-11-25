// Code generated by go-swagger; DO NOT EDIT.

package adm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
	models "github.com/sevings/mindwell-server/models"
)

// GetAdmGrandsonStatusOKCode is the HTTP code returned for type GetAdmGrandsonStatusOK
const GetAdmGrandsonStatusOKCode int = 200

/*GetAdmGrandsonStatusOK status of your gifts

swagger:response getAdmGrandsonStatusOK
*/
type GetAdmGrandsonStatusOK struct {

	/*
	  In: Body
	*/
	Payload *GetAdmGrandsonStatusOKBody `json:"body,omitempty"`
}

// NewGetAdmGrandsonStatusOK creates GetAdmGrandsonStatusOK with default headers values
func NewGetAdmGrandsonStatusOK() *GetAdmGrandsonStatusOK {

	return &GetAdmGrandsonStatusOK{}
}

// WithPayload adds the payload to the get adm grandson status o k response
func (o *GetAdmGrandsonStatusOK) WithPayload(payload *GetAdmGrandsonStatusOKBody) *GetAdmGrandsonStatusOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get adm grandson status o k response
func (o *GetAdmGrandsonStatusOK) SetPayload(payload *GetAdmGrandsonStatusOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAdmGrandsonStatusOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetAdmGrandsonStatusForbiddenCode is the HTTP code returned for type GetAdmGrandsonStatusForbidden
const GetAdmGrandsonStatusForbiddenCode int = 403

/*GetAdmGrandsonStatusForbidden you're not registered in adm

swagger:response getAdmGrandsonStatusForbidden
*/
type GetAdmGrandsonStatusForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAdmGrandsonStatusForbidden creates GetAdmGrandsonStatusForbidden with default headers values
func NewGetAdmGrandsonStatusForbidden() *GetAdmGrandsonStatusForbidden {

	return &GetAdmGrandsonStatusForbidden{}
}

// WithPayload adds the payload to the get adm grandson status forbidden response
func (o *GetAdmGrandsonStatusForbidden) WithPayload(payload *models.Error) *GetAdmGrandsonStatusForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get adm grandson status forbidden response
func (o *GetAdmGrandsonStatusForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAdmGrandsonStatusForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
