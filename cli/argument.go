package cli

import "strings"

type Argument struct {
	Name string
	Type ParameterType
}

func NewArgument(name string, tp ParameterType) *Argument {
	if len(name) == 0 {
		return nil
	}
	if strings.HasPrefix(name, "--") {
		return nil
	}
	return &Argument{
		Name: name,
		Type: tp,
	}
}
