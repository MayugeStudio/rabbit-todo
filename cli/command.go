package cli

import (
	"fmt"
	"rabbit-todo/cli/param"
	"strings"
)

type Action func(args []string, opts map[string]parameter.ParamValue) (string, error)

type Command struct {
	Name      string
	Arguments []*parameter.Argument
	Options   []*parameter.Option
	Action    Action
	Usage     string
}

// Execute method execute action
func (c *Command) Execute(inputParams []string) (string, error) {
	var (
		args        []string
		opts        = make(map[string]parameter.ParamValue)
		flagOpts    = c.flagOptions()
		startOption = false
	)

	for i := 0; i < len(inputParams); i++ {
		param := inputParams[i]
		if isArgument(param) && !startOption {
			args = append(args, param)
		} else {
			startOption = true
			optName, ov, err := c.parseOption(param, inputParams, &i, flagOpts)
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
		opts[strings.TrimPrefix(flagOpt.Name, "-")] = *parameter.NewBoolParameterPtr(false)
	}

	return c.Action(args, opts)
}

func (c *Command) parseOption(param string, inputParams []string, idxPtr *int, flagOpts []*parameter.Option) (string, *parameter.ParamValue, error) {
	name := param
	value := ""
	var oType parameter.ParameterType
	isValid := false
	isFlag := false

	for _, eOpt := range c.Options {
		if eOpt.Name == name {
			oType = eOpt.Type
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
		return name, parameter.NewBoolParameterPtr(true), nil
	} else {
		// Not Flag Option

		// Make sure the  next parameter is not an option starting with `--`
		if *idxPtr+1 < len(inputParams) {
			// Normal Option always accepts one argument
			if !isArgument(inputParams[*idxPtr+1]) {
				typeStr := parameter.ParameterTypeToString(oType)
				return "", nil, fmt.Errorf("\"%s\" option require one \"%s\" type argument", name, typeStr)
			}
		}

		// Whether normal option is last parameter or not
		if *idxPtr+1 == len(inputParams) {
			typeStr := parameter.ParameterTypeToString(oType)
			return "", nil, fmt.Errorf("\"%s\" option require one \"%s\" type argument", name, typeStr)
		}
		// Increase idx by one
		*idxPtr++
		value = inputParams[*idxPtr]

		ov, err := parameter.ToParameterValue(value, oType)
		if err != nil {
			return "", nil, fmt.Errorf("invalid option \"%s\": %w", name, err)
		}

		name = strings.TrimPrefix(param, "--")
		return name, ov, nil
	}
}

func NewCommand(name string, args []*parameter.Argument, opts []*parameter.Option, action Action) Command {
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

func (c *Command) flagOptions() []*parameter.Option {
	var flagOpts []*parameter.Option

	for _, option := range c.Options {
		if option.IsFlag {
			flagOpts = append(flagOpts, option)
		}
	}
	return flagOpts
}

func createUsageString(commandName string, args []*parameter.Argument, opts []*parameter.Option) string {
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
