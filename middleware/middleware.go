package middleware

import (
	"github.com/wlMalk/gapi/request"
	"github.com/wlMalk/gapi/response"
)

type Handler interface {
	ServeHTTP(*response.Response, *request.Request)
}

type HandlerFunc func(*response.Response, *request.Request)

func (f HandlerFunc) ServeHTTP(w *response.Response, r *request.Request) {
	f(w, r)
}

type Middleware interface {
	Run(Handler) Handler
}

type MiddlewareFunc func(Handler) Handler

func (f MiddlewareFunc) Run(h Handler) Handler {
	return f(h)
}
