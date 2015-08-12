package request

import (
	"net/http"

	"github.com/wlMalk/gapi/validation"
)

type Attributes map[interface{}]interface{}

func (a Attributes) Set(k interface{}, v interface{}) {
	a[k] = v
}

type Input map[string]*validation.Value

type Request struct {
	*http.Request
	Attributes Attributes
	Input      Input
}

func New(r *http.Request) *Request {
	return &Request{
		Request:    r,
		Input:      Input{},
		Attributes: Attributes{},
	}
}

// Req returns the request.
func (r *Request) Req() *http.Request {
	return r.Request
}

// Param returns the input parameter value by its name.
func (r *Request) Param(name string) *validation.Value {
	return r.Input[name]
}

// ParamOk returns the input parameter value by its name.
func (r *Request) ParamOk(name string) (*validation.Value, bool) {
	p, ok := r.Input[name]
	return p, ok
}

// Params returns a map of input parameters values by their names.
// If no names given then it returns r.Input
func (r *Request) Params(names ...string) map[string]*validation.Value {
	if len(names) == 0 {
		return r.Input
	}
	params := map[string]*validation.Value{}
	for _, n := range names {
		p, ok := r.Input[n]
		if !ok {
			continue
		}
		params[n] = p
	}
	return params
}
