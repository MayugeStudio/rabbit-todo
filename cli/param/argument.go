package param

import (
	"fmt"
	"strings"
)

type Argument struct {
	Name string
	Type Type
}

func NewArgument(name string, tp Type) (*Argument, error) {
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
