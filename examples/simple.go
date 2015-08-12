package main

import (
	"fmt"
	"net/http"

	"github.com/wlMalk/gapi"
	"github.com/wlMalk/gapi/constants"
	"github.com/wlMalk/gapi/operation"
	"github.com/wlMalk/gapi/param"
	"github.com/wlMalk/gapi/request"

	"github.com/julienschmidt/httprouter"
	wrapper "github.com/wlMalk/gapi/wrapper/julienschmidt/httprouter"
)

func main() {
	w := wrapper.Wrap(httprouter.New())
	g := gapi.New(w).SetBasePath("/api")
	g.Operations.Set(
		operation.GET("/users/:id").
			UsesFunc(getUser).
			Schemes(constants.SCHEME_HTTP).
			Params(
			param.PathParam("id").As(constants.TYPE_INT64),
		),
	)
	g.Start(":8080")
}

func getUser(w http.ResponseWriter, r *request.Request) {
	p := r.Input["id"]
	id := p.Int64()
	w.Write([]byte(fmt.Sprint(id)))
}
