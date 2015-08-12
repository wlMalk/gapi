package validation

import (
	"regexp"
)

func Eq(a interface{}) Validator {
	return NewFuncValidator(func(v *Value, r Requester) error {
		if a != v.Value() {
			return ErrEq.Err([]interface{}{v.Name(), a}...)
		}
		return nil
	}, MsgEq.Msg(a))
}

func Regexp(p *regexp.Regexp) Validator {
	return NewFuncValidator(func(v *Value, r Requester) error {
		if !p.MatchString(v.String()) {
			return ErrRegexp.Err([]interface{}{v.Name(), p.String()}...)
		}
		return nil
	}, MsgRegexp.Msg(p.String()))
}

func If(f func(*Value, Requester) bool, validators ...Validator) Validator {
	return ValidatorFunc(func(v *Value, r Requester) error {
		if f(v, r) {
			for _, va := range validators {
				err := va.Validate(v, r)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func IfElse(f func(*Value, Requester) bool, tValidators []Validator, fValidators []Validator) Validator {
	return ValidatorFunc(func(v *Value, r Requester) error {
		if f(v, r) {
			for _, va := range tValidators {
				err := va.Validate(v, r)
				if err != nil {
					return err
				}
			}
		} else {
			for _, va := range fValidators {
				err := va.Validate(v, r)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
