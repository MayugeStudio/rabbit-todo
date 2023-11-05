package param

import (
	"fmt"
	"strconv"
)

type Value struct {
	StringVal string
	IntVal    int
	BoolVal   bool
	Type      Type
}

func (v *Value) Value() interface{} {
	switch v.Type {
	case STRING:
		return v.StringVal
	case INT:
		return v.IntVal
	case BOOL:
		return v.BoolVal
	default:
		return nil
	}
}

func ToParameterValue(value string, paramType Type) (*Value, error) {
	switch paramType {
	case STRING:
		return NewStringParameterPtr(value), nil
	case INT:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to Integer", value)
		}
		return NewIntegerParameterPtr(intValue), nil
	case BOOL:
		if value == "" {
			return NewBoolParameterPtr(true), nil
		}
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to Boolean", value)
		}
		return NewBoolParameterPtr(boolValue), nil
	default:
		return nil, fmt.Errorf("unknown parameter type %v", paramType)
	}
}

func NewStringParameterPtr(value string) *Value {
	return &Value{
		StringVal: value,
		IntVal:    0,
		BoolVal:   false,
		Type:      STRING,
	}
}

func NewIntegerParameterPtr(value int) *Value {
	return &Value{
		StringVal: "",
		IntVal:    value,
		BoolVal:   false,
		Type:      INT,
	}
}

func NewBoolParameterPtr(value bool) *Value {
	return &Value{
		StringVal: "",
		IntVal:    0,
		BoolVal:   value,
		Type:      BOOL,
	}
}
