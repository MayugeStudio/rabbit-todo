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
		args        []string
		opts        = make(map[string]OptionValue)
		flagOpts    = c.flagOptions()
		startOption = false
	)

	for i := 0; i < len(inputParams); i++ {
		param := inputParams[i]
		if isArgument(param) && !startOption {
			args = append(args, param)
		} else {
			startOption = true
			optName, ov, err := c.parseOption(param, inputParams, &i, &flagOpts)
			if err != nil {
				return "", err
			}
			opts[optName] = *ov
		}
	}

	if len(args) < len(c.Arguments) {
		return "", fmt.Errorf("not enough arguments: actual %d, expected %d", len(args), len(c.Arguments))
	} else if len(args) > len(c.Arguments) {
		return "", fmt.Errorf("too many arguments: actual %d, expected %d", len(args), len(c.Arguments))
	}

	// FlagOptions not passed are initialized to false
	for _, flagOpt := range flagOpts {
		opts[strings.TrimPrefix(flagOpt.Name, "-")] = *getBoolOptionPtr(false)
	}

	return c.Action(args, opts)
}

func (c *Command) parseOption(param string, inputParams []string, idx *int, flagOptsPtr *[]*Option) (string, *OptionValue, error) {
	optName := param
	optValue := ""
	optType := ParameterType(-100)
	isFlag := false
	flagOpts := *flagOptsPtr

	for _, eOpt := range c.Options {
		if eOpt.Name == optName {
			optType = eOpt.Type
			isFlag = eOpt.IsFlag
			break
		}
	}

	if optType == -100 {
		return "", nil, fmt.Errorf("invalid option %s", optName)
	}

	if isFlag {
		// Check Next Parameter
		if *idx+1 < len(inputParams) {
			nextParam := inputParams[*idx+1]
			// Flag Option cannot have value
			if isArgument(nextParam) {
				return "", nil, fmt.Errorf("flag-option %s cannot have value", optName)
			}
		}

		// Delete FlagOption from FlagOption slice
		for i, fOpt := range flagOpts {
			if fOpt.Name == optName {
				flagOpts[i] = flagOpts[len(flagOpts)-1]
				flagOpts = flagOpts[:len(flagOpts)-1]
			}
		}
		optName = strings.TrimPrefix(optName, "--")
		return optName, getBoolOptionPtr(true), nil
	} else {
		// Not Flag Option
		// Check Next Parameter
		*idx++
		if *idx < len(inputParams) {
			optValue = inputParams[*idx]
		}

		ov, err := convertToOptionValue(optValue, optType)
		if err != nil {
			return "", nil, fmt.Errorf("invalid option \"%s\": %w", optName, err)
		}

		optName = strings.TrimPrefix(param, "--")
		return optName, ov, nil
	}
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

func isArgument(param string) bool {
	if strings.HasPrefix(param, "--") {
		return false
	}
	return true
}

func (c *Command) flagOptions() []*Option {
	var flagOpts []*Option

	for _, option := range c.Options {
		if option.IsFlag {
			flagOpts = append(flagOpts, option)
		}
	}
	return flagOpts
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
