package cli

import "strings"

type Option struct {
	Name   string
	Type   ParameterType
	IsFlag bool
}

func NewOption(name string, tp ParameterType) *Option {
	if len(name) == 0 {
		return nil
	}
	if !strings.HasPrefix(name, "--") {
		return nil
	}
	return &Option{
		Name:   name,
		Type:   tp,
		IsFlag: false,
	}
}

func NewFlagOption(name string) *Option {
	if len(name) == 0 {
		return nil
	}
	if !strings.HasPrefix(name, "--") {
		return nil
	}
	return &Option{
		Name:   name,
		Type:   BOOL,
		IsFlag: true,
	}
}
