// Code generated by go-swagger; DO NOT EDIT.

package votes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/mindwell-server/models"
)

// PutEntriesIDVoteOKCode is the HTTP code returned for type PutEntriesIDVoteOK
const PutEntriesIDVoteOKCode int = 200

/*PutEntriesIDVoteOK vote status

swagger:response putEntriesIdVoteOK
*/
type PutEntriesIDVoteOK struct {

	/*
	  In: Body
	*/
	Payload *models.VoteStatus `json:"body,omitempty"`
}

// NewPutEntriesIDVoteOK creates PutEntriesIDVoteOK with default headers values
func NewPutEntriesIDVoteOK() *PutEntriesIDVoteOK {
	return &PutEntriesIDVoteOK{}
}

// WithPayload adds the payload to the put entries Id vote o k response
func (o *PutEntriesIDVoteOK) WithPayload(payload *models.VoteStatus) *PutEntriesIDVoteOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put entries Id vote o k response
func (o *PutEntriesIDVoteOK) SetPayload(payload *models.VoteStatus) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutEntriesIDVoteOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutEntriesIDVoteForbiddenCode is the HTTP code returned for type PutEntriesIDVoteForbidden
const PutEntriesIDVoteForbiddenCode int = 403

/*PutEntriesIDVoteForbidden access denied

swagger:response putEntriesIdVoteForbidden
*/
type PutEntriesIDVoteForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutEntriesIDVoteForbidden creates PutEntriesIDVoteForbidden with default headers values
func NewPutEntriesIDVoteForbidden() *PutEntriesIDVoteForbidden {
	return &PutEntriesIDVoteForbidden{}
}

// WithPayload adds the payload to the put entries Id vote forbidden response
func (o *PutEntriesIDVoteForbidden) WithPayload(payload *models.Error) *PutEntriesIDVoteForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put entries Id vote forbidden response
func (o *PutEntriesIDVoteForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutEntriesIDVoteForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutEntriesIDVoteNotFoundCode is the HTTP code returned for type PutEntriesIDVoteNotFound
const PutEntriesIDVoteNotFoundCode int = 404

/*PutEntriesIDVoteNotFound Entry not found

swagger:response putEntriesIdVoteNotFound
*/
type PutEntriesIDVoteNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutEntriesIDVoteNotFound creates PutEntriesIDVoteNotFound with default headers values
func NewPutEntriesIDVoteNotFound() *PutEntriesIDVoteNotFound {
	return &PutEntriesIDVoteNotFound{}
}

// WithPayload adds the payload to the put entries Id vote not found response
func (o *PutEntriesIDVoteNotFound) WithPayload(payload *models.Error) *PutEntriesIDVoteNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put entries Id vote not found response
func (o *PutEntriesIDVoteNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutEntriesIDVoteNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
