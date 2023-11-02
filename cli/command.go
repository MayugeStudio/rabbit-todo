package cli

type Action func(args []string, opts []string) (string, error)
type UsageFunc func() string

type Command struct {
	Name      string
	Arguments []string
	Options   []string
	Action    Action
	Usage     UsageFunc
}

func (c *Command) Execute() (string, error) {
	return c.Action(c.Arguments, c.Options)
}
