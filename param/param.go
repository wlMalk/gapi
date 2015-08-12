package param

import (
	"fmt"

	"github.com/wlMalk/gapi/constants"
	"github.com/wlMalk/gapi/validation"
)

type Params struct {
	params             map[string]*Param
	containsFiles      bool
	containsBodyParams bool
	isLocked           bool
}

func NewParams() *Params {
	return &Params{
		params: map[string]*Param{},
	}
}

func (p *Params) Get() map[string]*Param {
	return p.params
}

func (p *Params) ContainsFiles() bool {
	return p.containsFiles
}

func (p *Params) ContainsBodyParams() bool {
	return p.containsBodyParams
}

func (p *Params) Append(params ...*Param) {
	if !p.isLocked {
		for i := 0; i < len(params); i++ {
			name := params[i].name
			if name == "" {
				panic("Detected a param without a name.")
			}
			if _, ok := p.params[name]; ok {
				panic(fmt.Sprintf("Detected 2 params with the same name: \"%s\".", name))
			}
			p.params[name] = params[i]
		}
	} else {
		panic("Can not edit Params while it's locked.")
	}
}

func (p *Params) Set(params ...*Param) {
	if !p.isLocked {
		p.params = map[string]*Param{}
		p.Append(params...)
	} else {
		panic("Can not edit Params while it's locked.")
	}
}

func (p *Params) Len() int {
	return len(p.params)
}

func (p *Params) Lock() {
	p.isLocked = true
}

// Param
type Param struct {
	name       string
	validators []validation.Validator
	def        interface{}
	as         int
	strSep     string
	isRequired bool
	isMultiple bool
	isFile     bool
	isInPath   bool
	isInQuery  bool
	isInHeader bool
	isInBody   bool
}

func New(name string) *Param {
	return &Param{
		name: name,
		as:   constants.TYPE_STRING,
	}
}

func PathParam(name string) *Param {
	return New(name).In(constants.IN_PATH)
}

func QueryParam(name string) *Param {
	return New(name).In(constants.IN_QUERY)
}

func HeaderParam(name string) *Param {
	return New(name).In(constants.IN_HEADER)
}

func BodyParam(name string) *Param {
	return New(name).In(constants.IN_BODY)
}

func (p *Param) Name(name string) *Param {
	p.name = name
	return p
}

func (p *Param) GetName() string {
	return p.name
}

func (p *Param) Required() *Param {
	p.isRequired = true
	return p
}

func (p *Param) IsRequired() bool {
	return p.isRequired
}

func (p *Param) File() *Param {
	p.isFile = true
	return p
}

func (p *Param) IsFile() bool {
	return p.isFile
}

func (p *Param) Multiple() *Param {
	p.isMultiple = true
	return p
}

func (p *Param) IsMultiple() bool {
	return p.isMultiple
}

func (p *Param) As(as int) *Param {
	if as == constants.TYPE_STRING ||
		as == constants.TYPE_INT ||
		as == constants.TYPE_INT64 ||
		as == constants.TYPE_FLOAT ||
		as == constants.TYPE_FLOAT64 ||
		as == constants.TYPE_BOOL {

		p.as = as
	}
	return p
}

func (p *Param) GetAs() int {
	return p.as
}

// If a Param is in path then it is required.
func (p *Param) In(in ...int) *Param {
	for _, i := range in {
		switch i {
		case constants.IN_PATH:
			p.isInPath = true
			p.isRequired = true
		case constants.IN_QUERY:
			p.isInQuery = true
		case constants.IN_HEADER:
			p.isInHeader = true
		case constants.IN_BODY:
			p.isInBody = true
		}
	}
	return p
}

// If a Param is in path then it is required.
func (p *Param) IsIn(in ...int) bool {
	for _, i := range in {
		switch i {
		case constants.IN_PATH:
			if !p.isInPath {
				return false
			}
		case constants.IN_QUERY:
			if !p.isInQuery {
				return false
			}
		case constants.IN_HEADER:
			if !p.isInHeader {
				return false
			}
		case constants.IN_BODY:
			if !p.isInBody {
				return false
			}
		}
	}
	return true
}

// Must sets the validators to use.
func (p *Param) Must(validators ...validation.Validator) *Param {
	p.validators = validators
	return p
}

// Validate returns the first error it encountered
func (p *Param) Validate(v *validation.Value, req validation.Requester) error {
	for _, va := range p.validators {
		err := va.Validate(v, req)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateAll returns all the errors it encountered
func (p *Param) ValidateAll(v *validation.Value, req validation.Requester) []error {
	var errs []error
	for _, va := range p.validators {
		err := va.Validate(v, req)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
