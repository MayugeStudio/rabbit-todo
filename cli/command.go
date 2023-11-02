package cli

import (
	"fmt"
	"strings"
)

type Action func(args []string, opts []string) (string, error)

type Command struct {
	Name      string
	Arguments []string
	Options   []string
	Action    Action
	Usage     string
}

func NewCommand(name string, args []string, opts []string, action Action) *Command {
	return &Command{
		Name:      name,
		Arguments: args,
		Options:   opts,
		Action:    action,
		Usage:     createUsageString(name, args, opts),
	}
}

func (c *Command) Execute(inputArgs []string, inputOpts []string) (string, error) {
	if len(inputArgs) != len(c.Arguments) {
		return "", fmt.Errorf("error: not enough arguments, expected: %d, got: %d", len(c.Arguments), len(inputArgs))
	}
	if len(inputOpts) > len(c.Options) {
		return "", fmt.Errorf("error: too many options")
	}
	return c.Action(inputArgs, inputOpts)
}

func createUsageString(commandName string, args []string, opts []string) string {
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
