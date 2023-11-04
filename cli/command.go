package cli

import (
	"fmt"
	"rabbit-todo/cli/param"
	"strings"
)

type Action func(args []string, opts map[string]param.Value) (string, error)

type Command struct {
	Name      string
	Arguments []*param.Argument
	Options   []*param.Option
	Action    Action
	Usage     string
}

// Execute method execute action
func (c *Command) Execute(inputParams []string) (string, error) {
	args, opts, err := c.Validate(inputParams)
	if err != nil {
		return "", err
	}

	return c.Action(args, opts)
}

func (c *Command) Validate(inputParams []string) ([]string, map[string]param.Value, error) {
	var (
		args        []string
		opts        = make(map[string]param.Value)
		flagOpts    = c.flagOptions()
		startOption = false
	)

	for i := 0; i < len(inputParams); i++ {
		p := inputParams[i]
		if isArgument(p) && !startOption {
			args = append(args, p)
		} else {
			startOption = true
			optName, ov, err := c.parseOption(p, inputParams, &i, flagOpts)
			if err != nil {
				return nil, nil, err
			}
			opts[optName] = *ov
		}
	}

	if len(args) < len(c.Arguments) {
		return nil, nil, fmt.Errorf("not enough arguments: actual %d, expected %d", len(args), len(c.Arguments))
	} else if len(args) > len(c.Arguments) {
		return nil, nil, fmt.Errorf("too many arguments: actual %d, expected %d", len(args), len(c.Arguments))
	}

	// FlagOptions not passed are initialized to false
	for _, flagOpt := range flagOpts {
		opts[strings.TrimPrefix(flagOpt.Name, "-")] = *param.NewBoolParameterPtr(false)
	}
	return args, opts, nil
}

func (c *Command) parseOption(p string, inputParams []string, idxPtr *int, flagOpts []*param.Option) (string, *param.Value, error) {
	name := p
	value := ""
	var optType param.Type
	isValid := false
	isFlag := false

	for _, eOpt := range c.Options {
		if eOpt.Name == name {
			optType = eOpt.Type
			isFlag = eOpt.IsFlag
			isValid = true
			break
		}
	}

	if !isValid {
		return "", nil, fmt.Errorf("invalid option %s", name)
	}

	if isFlag {
		// Make sure the next parameter is not an argument not starting with `--`
		if *idxPtr+1 < len(inputParams) {
			nextParam := inputParams[*idxPtr+1]
			// Flag Option cannot have value
			if isArgument(nextParam) {
				return "", nil, fmt.Errorf("flag-option %s cannot have value", name)
			}
		}

		// Delete FlagOption from FlagOption slice
		for i, fOpt := range flagOpts {
			if fOpt.Name == name {
				flagOpts[i] = flagOpts[len(flagOpts)-1]
				flagOpts = flagOpts[:len(flagOpts)-1]
			}
		}
		name = strings.TrimPrefix(name, "--")
		return name, param.NewBoolParameterPtr(true), nil
	} else {
		// Not Flag Option

		// Make sure the  next parameter is not an option starting with `--`
		if *idxPtr+1 < len(inputParams) {
			// Normal Option always accepts one argument
			if !isArgument(inputParams[*idxPtr+1]) {
				typeStr := param.ParameterTypeToString(optType)
				return "", nil, fmt.Errorf("\"%s\" option require one \"%s\" type argument", name, typeStr)
			}
		}

		// Whether normal option is last parameter or not
		if *idxPtr+1 == len(inputParams) {
			typeStr := param.ParameterTypeToString(optType)
			return "", nil, fmt.Errorf("\"%s\" option require one \"%s\" type argument", name, typeStr)
		}
		// Increase idx by one
		*idxPtr++
		value = inputParams[*idxPtr]

		ov, err := param.ToParameterValue(value, optType)
		if err != nil {
			return "", nil, fmt.Errorf("invalid option \"%s\": %w", name, err)
		}

		name = strings.TrimPrefix(p, "--")
		return name, ov, nil
	}
}

func NewCommand(name string, args []*param.Argument, opts []*param.Option, action Action) Command {
	return Command{
		Name:      name,
		Arguments: args,
		Options:   opts,
		Action:    action,
		Usage:     createUsageString(name, args, opts),
	}
}

func isArgument(p string) bool {
	if strings.HasPrefix(p, "--") {
		return false
	}
	return true
}

func (c *Command) flagOptions() []*param.Option {
	var flagOpts []*param.Option

	for _, option := range c.Options {
		if option.IsFlag {
			flagOpts = append(flagOpts, option)
		}
	}
	return flagOpts
}

func createUsageString(commandName string, args []*param.Argument, opts []*param.Option) string {
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
