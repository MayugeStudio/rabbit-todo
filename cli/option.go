package cli

import "strings"

type Option struct {
	Name   string
	Type   ParameterType
	IsFlag bool
}

func NewOption(name string, tp ParameterType) *Option {
	if !isValidOption(name) {
		return nil
	}
	return &Option{
		Name:   name,
		Type:   tp,
		IsFlag: false,
	}
}

func NewFlagOption(name string) *Option {
	if !isValidOption(name) {
		return nil
	}
	return &Option{
		Name:   name,
		Type:   BOOL,
		IsFlag: true,
	}
}

func isValidOption(name string) bool {
	if len(name) == 0 {
		return false
	}
	if !strings.HasPrefix(name, "--") {
		return false
	}
	return true
}
