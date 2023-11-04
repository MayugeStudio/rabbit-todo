package cli

import (
	"fmt"
	"strconv"
)

type OptionValue struct {
	StringVal string
	IntVal    int
	BoolVal   bool
	Type      ParameterType
}

func (ov *OptionValue) Value() interface{} {
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

func convertToOptionValue(value string, paramType ParameterType) (*OptionValue, error) {
	switch paramType {
	case STRING:
		return getStringOptionPtr(value), nil
	case INT:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to Integer", value)
		}
		return getIntegerOptionPtr(intValue), nil
	case BOOL:
		if value == "" {
			return getBoolOptionPtr(true), nil
		}
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to Boolean", value)
		}
		return getBoolOptionPtr(boolValue), nil
	default:
		return nil, fmt.Errorf("unknown parameter type %v", paramType)
	}
}

func getStringOptionPtr(value string) *OptionValue {
	return &OptionValue{
		StringVal: value,
		IntVal:    0,
		BoolVal:   false,
		Type:      STRING,
	}
}

func getIntegerOptionPtr(value int) *OptionValue {
	return &OptionValue{
		StringVal: "",
		IntVal:    value,
		BoolVal:   false,
		Type:      INT,
	}
}

func getBoolOptionPtr(value bool) *OptionValue {
	return &OptionValue{
		StringVal: "",
		IntVal:    0,
		BoolVal:   value,
		Type:      BOOL,
	}
}
