// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sevings/yummy-server/gen/models"
)

// GetUsersIDFavoritesOKCode is the HTTP code returned for type GetUsersIDFavoritesOK
const GetUsersIDFavoritesOKCode int = 200

/*GetUsersIDFavoritesOK Entry list

swagger:response getUsersIdFavoritesOK
*/
type GetUsersIDFavoritesOK struct {

	/*
	  In: Body
	*/
	Payload *models.Feed `json:"body,omitempty"`
}

// NewGetUsersIDFavoritesOK creates GetUsersIDFavoritesOK with default headers values
func NewGetUsersIDFavoritesOK() *GetUsersIDFavoritesOK {
	return &GetUsersIDFavoritesOK{}
}

// WithPayload adds the payload to the get users Id favorites o k response
func (o *GetUsersIDFavoritesOK) WithPayload(payload *models.Feed) *GetUsersIDFavoritesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users Id favorites o k response
func (o *GetUsersIDFavoritesOK) SetPayload(payload *models.Feed) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersIDFavoritesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersIDFavoritesForbiddenCode is the HTTP code returned for type GetUsersIDFavoritesForbidden
const GetUsersIDFavoritesForbiddenCode int = 403

/*GetUsersIDFavoritesForbidden access denied

swagger:response getUsersIdFavoritesForbidden
*/
type GetUsersIDFavoritesForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersIDFavoritesForbidden creates GetUsersIDFavoritesForbidden with default headers values
func NewGetUsersIDFavoritesForbidden() *GetUsersIDFavoritesForbidden {
	return &GetUsersIDFavoritesForbidden{}
}

// WithPayload adds the payload to the get users Id favorites forbidden response
func (o *GetUsersIDFavoritesForbidden) WithPayload(payload *models.Error) *GetUsersIDFavoritesForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users Id favorites forbidden response
func (o *GetUsersIDFavoritesForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersIDFavoritesForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersIDFavoritesNotFoundCode is the HTTP code returned for type GetUsersIDFavoritesNotFound
const GetUsersIDFavoritesNotFoundCode int = 404

/*GetUsersIDFavoritesNotFound User not found

swagger:response getUsersIdFavoritesNotFound
*/
type GetUsersIDFavoritesNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUsersIDFavoritesNotFound creates GetUsersIDFavoritesNotFound with default headers values
func NewGetUsersIDFavoritesNotFound() *GetUsersIDFavoritesNotFound {
	return &GetUsersIDFavoritesNotFound{}
}

// WithPayload adds the payload to the get users Id favorites not found response
func (o *GetUsersIDFavoritesNotFound) WithPayload(payload *models.Error) *GetUsersIDFavoritesNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users Id favorites not found response
func (o *GetUsersIDFavoritesNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersIDFavoritesNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
