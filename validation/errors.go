package validation

import (
	"fmt"
)

type Err string

func (err Err) Err(vals ...interface{}) error {
	return ValidationErr(fmt.Sprintf(string(err), vals...))
}

type Msg string

func (m Msg) Msg(vals ...interface{}) string {
	return fmt.Sprintf(string(m), vals...)
}

type ValidationErr string

func (err ValidationErr) Error() string {
	return string(err)
}

const (
	ErrEq     = Err("%v must equal %v") // prefix "Validation Error: " will be added
	ErrRegexp = Err("%v must match %v")
)

const (
	MsgEq     = Msg("must equal %v") // prefix "Validation Error: " will be added
	MsgRegexp = Msg("must match %v")
)
