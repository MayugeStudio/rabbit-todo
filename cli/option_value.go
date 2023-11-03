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

func convertToOptionValue(value string, paramType ParameterType) (OptionValue, error) {
	switch paramType {
	case STRING:
		return OptionValue{StringVal: value, Type: STRING}, nil
	case INT:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return OptionValue{}, fmt.Errorf("cannot convert %s to Integer", value)
		}
		return OptionValue{IntVal: intValue, Type: INT}, nil
	case BOOL:
		if value == "" {
			return OptionValue{BoolVal: true, Type: BOOL}, nil
		}
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return OptionValue{}, fmt.Errorf("cannot convert %s to Boolean", value)
		}
		return OptionValue{BoolVal: boolValue, Type: BOOL}, nil
	default:
		return OptionValue{}, fmt.Errorf("unknown parameter type %v", paramType)
	}
}
