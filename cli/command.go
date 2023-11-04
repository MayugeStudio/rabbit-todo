package cli

import (
	"fmt"
	"rabbit-todo/cli/param"
	"strings"
)

const (
	optionPrefix = "--"
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
	args := make([]string, 0, len(c.Arguments))
	opts := c.initializeOptions()
	flagOpts := c.flagOptions()
	optionNow := false

	for i := 0; i < len(inputParams); i++ {
		p := inputParams[i]
		if !optionNow && isArgument(p) {
			args = append(args, p)
		} else {
			optionNow = true
			optName, v, err := c.parseOption(p, inputParams, &i, flagOpts)
			if err != nil {
				return nil, nil, err
			}
			opts[optName] = *v
		}
	}

	if err := c.validateArguments(args); err != nil {
		return nil, nil, err
	}
	return args, opts, nil
}

func (c *Command) parseOption(optParam string, inputParams []string, idxPtr *int, flagOpts []*param.Option) (string, *param.Value, error) {
	optType, isFlag, err := c.getOptionTypeAndFlag(optParam)
	if err != nil {
		return "", nil, err
	}

	if isFlag {
		return c.processFlagOption(optParam, inputParams, idxPtr, flagOpts)
	} else {
		return c.processRegularOption(optParam, optType, inputParams, idxPtr)
	}
}

// validateArguments ensures the correct number of arguments have been provided.
func (c *Command) validateArguments(args []string) error {
	if len(args) < len(c.Arguments) {
		return fmt.Errorf("not enough arguments: actual %d, expected %d", len(args), len(c.Arguments))
	} else if len(args) > len(c.Arguments) {
		return fmt.Errorf("too many arguments: actual %d, expected %d", len(args), len(c.Arguments))
	}
	return nil
}

// initializeOptions initializes the options map with default values for flags.
func (c *Command) initializeOptions() map[string]param.Value {
	options := make(map[string]param.Value)
	for _, option := range c.Options {
		if option.IsFlag {
			options[strings.TrimPrefix(option.Name, optionPrefix)] = *param.NewBoolParameterPtr(false)
		}
	}
	return options
}

func (c *Command) getOptionTypeAndFlag(optionName string) (param.Type, bool, error) {
	for _, option := range c.Options {
		if option.Name == optionName {
			return option.Type, option.IsFlag, nil
		}
	}
	return -1, false, fmt.Errorf("invalid option %s", optionName)
}

func (c *Command) processFlagOption(optionName string, inputParams []string, idxPtr *int, flagOpts []*param.Option) (string, *param.Value, error) {
	// Make sure the next parameter is not an argument not starting with `--`
	// Flag Option cannot have value
	if *idxPtr+1 < len(inputParams) && isArgument(inputParams[*idxPtr+1]) {
		return "", nil, fmt.Errorf("flag-option %s cannot have value", optionName)
	}

	// Delete FlagOption from FlagOption slice
	c.removeFlagOption(optionName, flagOpts)
	optionName = strings.TrimPrefix(optionName, optionPrefix)
	return optionName, param.NewBoolParameterPtr(true), nil
}

func (c *Command) processRegularOption(optionName string, optType param.Type, inputParams []string, idxPtr *int) (string, *param.Value, error) {
	// Make sure the  next parameter is not an option starting with `--`
	// Whether normal option is last parameter or not
	// Normal Option always accepts one argument
	if *idxPtr+1 >= len(inputParams) || !isArgument(inputParams[*idxPtr+1]) {
		typeStr := param.ParameterTypeToString(optType)
		return "", nil, fmt.Errorf("\"%s\" option requires a \"%s\" type argument", optionName, typeStr)
	}

	// Increase idx by 1
	*idxPtr++
	value := inputParams[*idxPtr]
	paramValue, err := param.ToParameterValue(value, optType)
	if err != nil {
		return "", nil, fmt.Errorf("invalid option \"%s\": %w", optionName, err)
	}

	optionName = strings.TrimPrefix(optionName, optionPrefix)
	return optionName, paramValue, nil
}

func (c *Command) removeFlagOption(name string, flagOpts []*param.Option) {
	for i, opt := range flagOpts {
		if opt.Name == name {
			flagOpts[i] = flagOpts[len(flagOpts)-1]
			flagOpts = flagOpts[:len(flagOpts)-1]
			break
		}
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
