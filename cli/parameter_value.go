package cli

import (
	"fmt"
	"strconv"
)

type ParameterValue struct {
	StringVal string
	IntVal    int
	BoolVal   bool
	Type      ParameterType
}

func (ov *ParameterValue) Value() interface{} {
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

func convertToParameterValue(value string, paramType ParameterType) (*ParameterValue, error) {
	switch paramType {
	case STRING:
		return getStringParameterPtr(value), nil
	case INT:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to Integer", value)
		}
		return getIntegerParameterPtr(intValue), nil
	case BOOL:
		if value == "" {
			return getBoolParameterPtr(true), nil
		}
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to Boolean", value)
		}
		return getBoolParameterPtr(boolValue), nil
	default:
		return nil, fmt.Errorf("unknown parameter type %v", paramType)
	}
}

func getStringParameterPtr(value string) *ParameterValue {
	return &ParameterValue{
		StringVal: value,
		IntVal:    0,
		BoolVal:   false,
		Type:      STRING,
	}
}

func getIntegerParameterPtr(value int) *ParameterValue {
	return &ParameterValue{
		StringVal: "",
		IntVal:    value,
		BoolVal:   false,
		Type:      INT,
	}
}

func getBoolParameterPtr(value bool) *ParameterValue {
	return &ParameterValue{
		StringVal: "",
		IntVal:    0,
		BoolVal:   value,
		Type:      BOOL,
	}
}
