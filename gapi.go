package gapi

import (
	"net/http"

	"github.com/wlMalk/gapi/operation"
)

type Gapi struct {
	wrapper    wrapper
	Operations *operation.Operations
	basePath   string
}

type wrapper interface {
	http.Handler
	Set(*operation.Operations)
	URL(*operation.Operation, ...string) string
}

func New(w wrapper) *Gapi {
	return &Gapi{
		wrapper:    w,
		Operations: operation.NewOperations(),
	}
}

func (g *Gapi) URL(o *operation.Operation, vars ...string) string {
	return g.wrapper.URL(o, vars...)
}

func (g *Gapi) SetBasePath(basePath string) *Gapi {
	g.basePath = basePath
	return g
}

func (g *Gapi) Prepare() {
	if g.basePath != "" {
		g.Operations.SetBasePath(g.basePath)
	}
	g.Operations.Prepare()
	g.wrapper.Set(g.Operations)
	g.Operations.Lock()
}

func (g *Gapi) Start(addr string) error {
	g.Prepare()
	return http.ListenAndServe(addr, g)
}

func (g *Gapi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.wrapper.ServeHTTP(w, r)
}
