package validation

import (
	"net/http"
)

type Validator interface {
	Validate(*Value, Requester) error
	Description() string
}

// ValidatorFunc is a convenient way to turn a function into a Validator.
// However, it returns empty string as a description.
type ValidatorFunc func(*Value, Requester) error

func (f ValidatorFunc) Validate(v *Value, r Requester) error {
	return f(v, r)
}

func (f ValidatorFunc) Description() string {
	return ""
}

func NewFuncValidator(f func(*Value, Requester) error, d string) Validator {
	return &FuncValidator{
		function:    f,
		description: d,
	}
}

type FuncValidator struct {
	function    func(*Value, Requester) error
	description string
}

func (f FuncValidator) Validate(v *Value, r Requester) error {
	return f.function(v, r)
}

func (f FuncValidator) Description() string {
	return f.description
}

type Requester interface {
	Req() *http.Request
	Param(string) *Value
	ParamOk(string) (*Value, bool)
	Params(...string) map[string]*Value
}

// type Valuer interface {
// 	Name() string
// 	Value() interface{}
// 	As() int
// 	Int() int
// 	Int64() int64
// 	Float() float32
// 	Float64() float64
// 	String() string
// 	Bool() bool
// 	// IntE() (int, error)
// 	// Int64E() (int64, error)
// 	// FloatE() (float32, error)
// 	// Float64E() (float64, error)
// 	// BoolE() (bool, error)
// }
