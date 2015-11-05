package operation

import (
	"fmt"

	"github.com/wlMalk/gapi/defaults"
	"github.com/wlMalk/gapi/middleware"
	"github.com/wlMalk/gapi/param"
	"github.com/wlMalk/gapi/request"
	"github.com/wlMalk/gapi/response"
)

type Operations struct {
	operations []*Operation
	index      map[string]int
	urlIndex   map[string]map[string]int
	isLocked   bool
}

func NewOperations() *Operations {
	return &Operations{
		index:    map[string]int{},
		urlIndex: map[string]map[string]int{},
	}
}

func (ops *Operations) Append(operations ...*Operation) {
	if !ops.isLocked {
		le := len(ops.operations)
		for i := 0; i < len(operations); i++ {
			ops.operations = append(ops.operations, operations[i])
			name := operations[i].name
			if name == "" {
				continue
			}
			if _, ok := ops.index[name]; ok {
				panic(fmt.Sprintf("Detected 2 operations with the same name: \"%s\".", name))
			}
			ops.index[name] = le + i
		}
	} else {
		panic("Can not edit Operations while it's locked.")
	}
}

func (ops *Operations) Set(operations ...*Operation) {
	if !ops.isLocked {
		ops.operations = nil
		ops.index = map[string]int{}
		ops.urlIndex = map[string]map[string]int{}
		ops.Append(operations...)
	} else {
		panic("Can not edit Operations while it's locked.")
	}
}

func (ops *Operations) Get() []*Operation {
	return ops.operations
}

func (ops *Operations) GetByIndex(i int) *Operation {
	return ops.operations[i]
}

func (ops *Operations) GetByName(name string) *Operation {
	return ops.operations[ops.index[name]]
}

func (ops *Operations) Len() int {
	return len(ops.operations)
}

func (ops *Operations) Prepare() {
	if !ops.isLocked {
		for i := 0; i < len(ops.operations); i++ {
			ops.operations[i].Prepare()
		}
	} else {
		panic("Can not edit Operations while it's locked.")
	}
}

func (ops *Operations) SetBasePath(basePath string) {
	if !ops.isLocked {
		for i := 0; i < len(ops.operations); i++ {
			ops.operations[i].setBasePath(basePath)
		}
	} else {
		panic("Can not edit Operations while it's locked.")
	}
}

func (ops *Operations) Lock() {
	ops.isLocked = true
}

type Operation struct {
	name     string
	handler  middleware.Handler
	path     string
	method   string
	params   *param.Params
	schemes  []string
	consumes []string
	produces []string
}

func New(path string) *Operation {
	return &Operation{
		path:     path,
		method:   "GET",
		params:   &param.Params{},
		schemes:  defaults.Schemes,
		consumes: defaults.Consumes,
		produces: defaults.Produces,
	}
}

func GET(path string) *Operation {
	return New(path).Method("GET")
}

func POST(path string) *Operation {
	return New(path).Method("POST")
}

func PUT(path string) *Operation {
	return New(path).Method("PUT")
}

func DELETE(path string) *Operation {
	return New(path).Method("DELETE")
}

func (o *Operation) Uses(handler middleware.Handler) *Operation {
	o.handler = handler
	return o
}

func (o *Operation) UsesFunc(f func(*response.Response, *request.Request)) *Operation {
	o.handler = middleware.HandlerFunc(f)
	return o
}

func getHandler(h middleware.Handler, mw []middleware.Middleware) middleware.Handler {
	final := h
	for i := len(mw) - 1; i >= 0; i-- {
		final = mw[i].Run(final)
	}
	return final
}

func (o *Operation) With(mw ...middleware.Middleware) *Operation {
	o.handler = getHandler(o.handler, mw)
	return o
}

func (o *Operation) Apply(temps ...func(*Operation)) *Operation {
	for i := 0; i < len(temps); i++ {
		temps[i](o)
	}
	return o
}

func getHandlerFunc(h middleware.Handler, mw []func(middleware.Handler) middleware.Handler) middleware.Handler {
	final := h
	for i := len(mw) - 1; i >= 0; i-- {
		final = mw[i](final)
	}
	return final
}

func (o *Operation) WithFunc(fmw ...func(middleware.Handler) middleware.Handler) *Operation {
	o.handler = getHandlerFunc(o.handler, fmw)
	return o
}

func (o *Operation) Name(name string) *Operation {
	o.name = name
	return o
}

func (o *Operation) GetName() string {
	return o.name
}

func (o *Operation) Method(method string) *Operation {
	o.method = method
	return o
}

func (o *Operation) Params(params ...*param.Param) *Operation {
	o.params.Set(params...)
	return o
}

func (o *Operation) Schemes(schemes ...string) *Operation {
	o.schemes = schemes
	return o
}

func (o *Operation) Consumes(consumes ...string) *Operation {
	o.consumes = consumes
	return o
}

func (o *Operation) Produces(produces ...string) *Operation {
	o.produces = produces
	return o
}

func (o *Operation) GetMethod() string {
	return o.method
}

func (o *Operation) GetPath() string {
	return o.path
}

func (o *Operation) setBasePath(basePath string) {
	o.path = basePath + o.path
}

func (o *Operation) Prepare() {
	o.handler = middleware.CheckSchemes(o.schemes)(middleware.ValidateParams(o.params)(o.handler))
	o.params.Lock()
}

func (o *Operation) ServeHTTP(w *response.Response, r *request.Request) {
	o.handler.ServeHTTP(w, r)
}

type middlewareContainer struct {
	middleware []middleware.Middleware
}

func With(mw ...middleware.Middleware) *middlewareContainer {
	return &middlewareContainer{middleware: mw}
}

func (c middlewareContainer) Handle(operations ...*Operation) []*Operation {
	for i := 0; i < len(operations); i++ {
		operations[i].With(c.middleware...)
	}
	return operations
}

type templateContainer struct {
	templates []func(*Operation)
}

func Apply(temps ...func(*Operation)) *templateContainer {
	return &templateContainer{templates: temps}
}

func (c templateContainer) To(operations ...*Operation) []*Operation {
	for i := 0; i < len(operations); i++ {
		operations[i].Apply(c.templates...)
	}
	return operations
}
