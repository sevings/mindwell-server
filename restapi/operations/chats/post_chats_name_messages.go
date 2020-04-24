// Code generated by go-swagger; DO NOT EDIT.

package chats

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/sevings/mindwell-server/models"
)

// PostChatsNameMessagesHandlerFunc turns a function with the right signature into a post chats name messages handler
type PostChatsNameMessagesHandlerFunc func(PostChatsNameMessagesParams, *models.UserID) middleware.Responder

// Handle executing the request and returning a response
func (fn PostChatsNameMessagesHandlerFunc) Handle(params PostChatsNameMessagesParams, principal *models.UserID) middleware.Responder {
	return fn(params, principal)
}

// PostChatsNameMessagesHandler interface for that can handle valid post chats name messages params
type PostChatsNameMessagesHandler interface {
	Handle(PostChatsNameMessagesParams, *models.UserID) middleware.Responder
}

// NewPostChatsNameMessages creates a new http.Handler for the post chats name messages operation
func NewPostChatsNameMessages(ctx *middleware.Context, handler PostChatsNameMessagesHandler) *PostChatsNameMessages {
	return &PostChatsNameMessages{Context: ctx, Handler: handler}
}

/*PostChatsNameMessages swagger:route POST /chats/{name}/messages chats postChatsNameMessages

PostChatsNameMessages post chats name messages API

*/
type PostChatsNameMessages struct {
	Context *middleware.Context
	Handler PostChatsNameMessagesHandler
}

func (o *PostChatsNameMessages) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostChatsNameMessagesParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.UserID
	if uprinc != nil {
		principal = uprinc.(*models.UserID) // this is really a models.UserID, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}