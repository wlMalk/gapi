package httprouter

import (
	"net/http"

	"github.com/wlMalk/gapi/middleware"
	"github.com/wlMalk/gapi/operation"
	"github.com/wlMalk/gapi/request"
	"github.com/wlMalk/gapi/wrapper"

	"github.com/julienschmidt/httprouter"
)

type Wrapper struct {
	router *httprouter.Router
}

func Wrap(router *httprouter.Router) *Wrapper {
	return &Wrapper{router}
}

func (w *Wrapper) Set(operations *operation.Operations) {
	for _, o := range operations.Get() {
		w.router.Handle(o.GetMethod(), wrapper.CurlyToColon(o.GetPath()), getHandle(o))
	}
}

func (w *Wrapper) URL(o *operation.Operation, vars ...string) string {
	return ""
}

func (w *Wrapper) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w.router.ServeHTTP(rw, r)
}

func getHandle(h middleware.Handler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req := request.New(r)
		params := map[string]string{}
		for _, p := range ps {
			params[p.Key] = p.Value
		}
		req.Attributes.Set(wrapper.PathParamsKey, params)
		h.ServeHTTP(w, req)
	})
}
