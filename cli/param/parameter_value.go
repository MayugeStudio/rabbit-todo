package param

import (
	"fmt"
	"strconv"
)

type ParamValue struct {
	StringVal string
	IntVal    int
	BoolVal   bool
	Type      ParameterType
}

func (ov *ParamValue) Value() interface{} {
	switch ov.Type {
	case STRING:
		return ov.StringVal
	case INT:
		return ov.IntVal
	case BOOL:
		return ov.BoolVal
	default:
		return nil
	}
}

func ToParameterValue(value string, paramType ParameterType) (*ParamValue, error) {
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

func NewStringParameterPtr(value string) *ParamValue {
	return &ParamValue{
		StringVal: value,
		IntVal:    0,
		BoolVal:   false,
		Type:      STRING,
	}
}

func NewIntegerParameterPtr(value int) *ParamValue {
	return &ParamValue{
		StringVal: "",
		IntVal:    value,
		BoolVal:   false,
		Type:      INT,
	}
}

func NewBoolParameterPtr(value bool) *ParamValue {
	return &ParamValue{
		StringVal: "",
		IntVal:    0,
		BoolVal:   value,
		Type:      BOOL,
	}
}
