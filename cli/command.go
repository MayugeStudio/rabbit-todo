package cli

import (
	"fmt"
	"strings"
)

type Action func(args []string, opts map[string]OptionValue) (string, error)

type Command struct {
	Name      string
	Arguments []*Argument
	Options   []*Option
	Action    Action
	Usage     string
}

// Execute method execute action
func (c *Command) Execute(inputParams []string) (string, error) {
	var (
		args []string
		opts = make(map[string]OptionValue)
	)

	for i := 0; i < len(inputParams); i++ {
		param := inputParams[i]
		if isOption(param) {
			optName := param
			optValue := ""
			optType := ParameterType(-100)
			i++

			for _, eOpt := range c.Options {
				if eOpt.Name == optName {
					optType = eOpt.Type
					break
				}
			}

			if optType == -100 {
				return "", fmt.Errorf("invalid option %s", optName)
			}

			if i < len(inputParams) {
				optValue = inputParams[i]
			}

			ov, err := convertToOptionValue(optValue, optType)
			if err != nil {
				return "", fmt.Errorf("invalid option \"%s\": %w", optName, err)
			}

			optName = strings.TrimPrefix(param, "--")
			opts[optName] = ov

		} else {
			args = append(args, param)
		}
	}

	if len(args) < len(c.Arguments) {
		return "", fmt.Errorf("not enough arguments: actual %d, expected %d", len(args), len(c.Arguments))
	} else if len(args) > len(c.Arguments) {
		return "", fmt.Errorf("too many arguments: actual %d, expected %d", len(args), len(c.Arguments))
	}

	return c.Action(args, opts)
}

func NewCommand(name string, args []*Argument, opts []*Option, action Action) Command {
	return Command{
		Name:      name,
		Arguments: args,
		Options:   opts,
		Action:    action,
		Usage:     createUsageString(name, args, opts),
	}
}

func isOption(param string) bool {
	if strings.HasPrefix(param, "--") {
		return true
	}
	return false
}

func createUsageString(commandName string, args []*Argument, opts []*Option) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Usage: %s", commandName))
	if len(args) > 0 {
		builder.WriteString(" [arguments]")
	}
	if len(opts) > 0 {
		builder.WriteString(" [options]")
	}
	return builder.String()
}
