package validation

import (
	"strconv"

	"github.com/wlMalk/gapi/constants"
)

type Value struct {
	name             string
	value            interface{}
	strValue         string
	as               int
	strMultipleValue []string
	multiple         bool
	from             string
}

func NewValue(name string, value string, from string, as int) *Value {
	return &Value{
		name:     name,
		strValue: value,
		from:     from,
		as:       as,
	}
}

func NewMultipleValue(name string, value []string, from string, as int) *Value {
	return &Value{
		name:             name,
		strMultipleValue: value,
		multiple:         true,
		from:             from,
		as:               as,
	}
}

func (v *Value) Name() string {
	return v.name
}

func (v *Value) As() int {
	return v.as
}

func (v *Value) Value() interface{} {
	if v.value == nil {
		switch v.as {
		case constants.TYPE_STRING:
			v.String()
		case constants.TYPE_INT:
			v.Int()
		case constants.TYPE_INT64:
			v.Int64()
		case constants.TYPE_FLOAT:
			v.Float()
		case constants.TYPE_FLOAT64:
			v.Float64()
		case constants.TYPE_BOOL:
			v.Bool()
		}
	}
	return v.value
}

func (v *Value) RawString() string {
	return v.strValue
}

func (v *Value) String() string {
	if v.value != nil {
		rv := v.value.(string)
		return rv
	}
	if v.as == constants.TYPE_STRING {
		v.value = v.strValue
		return v.strValue
	}
	return ""
}

func (v *Value) BoolE() (bool, error) {
	if v.value != nil {
		rv, ok := v.value.(bool)
		if ok {
			return rv, nil
		}
	}
	if v.as == constants.TYPE_BOOL {
		rv, err := strconv.ParseBool(v.strValue)
		if err == nil {
			v.value = rv
		}
		return rv, err
	}
	return false, nil
}

func (v *Value) Bool() bool {
	rv, _ := v.BoolE()
	return rv
}

func (v *Value) IntE() (int, error) {
	if v.value != nil {
		rv, ok := v.value.(int)
		if ok {
			return rv, nil
		}
	}
	if v.as == constants.TYPE_INT {
		rv, err := strconv.Atoi(v.strValue)
		if err == nil {
			v.value = rv
		}
		return rv, err
	}
	return 0, nil
}

func (v *Value) Int() int {
	rv, _ := v.IntE()
	return rv
}

func (v *Value) Int64E() (int64, error) {
	if v.value != nil {
		rv, ok := v.value.(int64)
		if ok {
			return rv, nil
		}
	}
	if v.as == constants.TYPE_INT64 {
		rv, err := strconv.ParseInt(v.strValue, 10, 64)
		if err == nil {
			v.value = rv
		}
		return rv, err
	}
	return 0, nil
}

func (v *Value) Int64() int64 {
	rv, _ := v.Int64E()
	return rv
}

func (v *Value) FloatE() (float32, error) {
	if v.value != nil {
		rv, ok := v.value.(float32)
		if ok {
			return rv, nil
		}
	}
	if v.as == constants.TYPE_FLOAT {
		rv, err := strconv.ParseFloat(v.strValue, 32)
		if err == nil {
			v.value = rv
		}
		return float32(rv), err
	}
	return 0, nil
}

func (v *Value) Float() float32 {
	rv, _ := v.FloatE()
	return rv
}

func (v *Value) Float64E() (float64, error) {
	if v.value != nil {
		rv, ok := v.value.(float64)
		if ok {
			return rv, nil
		}
	}
	if v.as == constants.TYPE_FLOAT64 {
		rv, err := strconv.ParseFloat(v.strValue, 64)
		if err == nil {
			v.value = rv
		}
		return rv, err
	}
	return 0, nil
}

func (v *Value) Float64() float64 {
	rv, _ := v.Float64E()
	return rv
}
