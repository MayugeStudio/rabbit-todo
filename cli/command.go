package cli

import (
	"fmt"
	"strings"
)

type Action func(args []string, opts []string) (string, error)

type Command struct {
	Name      string
	Arguments []*Argument
	Options   []*Option
	Action    Action
	Usage     string
}

func NewCommand(name string, args []*Argument, opts []*Option, action Action) *Command {
	return &Command{
		Name:      name,
		Arguments: args,
		Options:   opts,
		Action:    action,
		Usage:     createUsageString(name, args, opts),
	}
}

func (c *Command) Execute(inputParams []string) (string, error) {
	var (
		actualArgs []string
		actualOpts []string
	)

	for _, param := range inputParams {
		if strings.HasPrefix(param, "--") {
			actualOpts = append(actualOpts, param)
		} else {
			actualArgs = append(actualArgs, param)
		}
	}

	if len(actualArgs) < len(c.Arguments) {
		return "", fmt.Errorf("not enough arguments")
	} else if len(actualArgs) > len(c.Arguments) {
		return "", fmt.Errorf("too many arguments")
	}

	if len(actualOpts) > len(c.Options) {
		return "", fmt.Errorf("too many options")
	}

	return c.Action(actualArgs, actualOpts)
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
