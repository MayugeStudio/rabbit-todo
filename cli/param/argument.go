package param

import (
	"fmt"
	"strings"
)

type Argument struct {
	Name string
	Type ParameterType
}

func NewArgument(name string, tp ParameterType) (*Argument, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("name must not be empty")
	}

	if strings.HasPrefix(name, "--") {
		return nil, fmt.Errorf("name must not start with '--'")
	}
	return &Argument{
		Name: name,
		Type: tp,
	}, nil
}
