package cli

import (
	"fmt"
	"strings"
)

type Option struct {
	Name   string
	Type   ParameterType
	IsFlag bool
}

func NewOption(name string, tp ParameterType) (*Option, error) {
	err := isValidOption(name)
	if err != nil {
		return nil, err
	}
	return &Option{
		Name:   name,
		Type:   tp,
		IsFlag: false,
	}, nil
}

func NewFlagOption(name string) (*Option, error) {
	err := isValidOption(name)
	if err != nil {
		return nil, err
	}
	return &Option{
		Name:   name,
		Type:   BOOL,
		IsFlag: true,
	}, nil
}

func isValidOption(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("name must not be empty")
	}
	if !strings.HasPrefix(name, "--") {
		return fmt.Errorf("name must start with '--'")
	}
	return nil
}
