package cli

import "strings"

type Option struct {
	Name   string
	Type   ParameterType
	IsFlag bool
}

func NewOption(name string, tp ParameterType, isFlag bool) *Option {
	if len(name) == 0 {
		return nil
	}
	if !strings.HasPrefix(name, "--") {
		return nil
	}
	return &Option{
		Name:   name,
		Type:   tp,
		IsFlag: isFlag,
	}
}
