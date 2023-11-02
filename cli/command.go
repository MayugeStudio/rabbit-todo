package cli

import (
	"fmt"
	"strings"
)

type Action func(args []string, opts []string) (string, error)
type UsageFunc func() string

type Command struct {
	Name      string
	Arguments []string
	Options   []string
	Action    Action
	Usage     UsageFunc
}

func NewCommand(name string, args []string, opts []string, action Action) *Command {
	return &Command{
		Name:      name,
		Arguments: args,
		Options:   opts,
		Action:    action,
		Usage: func() string {
			var builder strings.Builder
			builder.WriteString(fmt.Sprintf("Usage: %s", name))
			if len(args) > 0 {
				builder.WriteString(" [arguments]")
			}
			if len(args) > 0 {
				builder.WriteString(" [options]")
			}
			return builder.String()
		},
	}
}

func (c *Command) Execute() (string, error) {
	return c.Action(c.Arguments, c.Options)
}
