package middleware

import (
	"net/http"

	"github.com/wlMalk/gapi/request"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *request.Request)
}

type HandlerFunc func(http.ResponseWriter, *request.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *request.Request) {
	f(w, r)
}

type Middleware interface {
	Run(Handler) Handler
}

type MiddlewareFunc func(Handler) Handler

func (f MiddlewareFunc) Run(h Handler) Handler {
	return f(h)
}
