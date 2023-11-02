package cli

type Func func(args []string, opts []string) (string, error)

type Command struct {
	Name      string
	Arguments []string
	Options   []string
	Function  Func
}

func (c *Command) Execute() (string, error) {
	return c.Function(c.Arguments, c.Options)
}
