package cli

import (
	"fmt"
	"rabbit-todo/cli/param"
	"strings"
)

const (
	optionPrefix = "--"
)

// Command represents a CLI command, including its name, expected arguments,
// options, the action to execute, and usage information.
// The struct is used to define commands in a CLI application and provides
// methods to execute and validate command input.
type Command struct {
	Name      string
	arguments []*param.Argument
	options   []*param.Option
	action    Action
}

// Action defines the function signature for actions that commands execute.
// It takes a slice of arguments and a map of options, which are parsed from the CLI input.
// The function returns a string that usually contains the result of the command execution or
// a message to the user, and an error if the execution fails.
type Action func(args []string, opts map[string]param.Value) (string, error)

// NewCommand constructs a new Command object with the given name and action,
// then It Generates the param.Argument and param.Option pointer slice.
func NewCommand(name string, action Action) Command {
	return Command{
		Name:      name,
		arguments: make([]*param.Argument, 0),
		options:   make([]*param.Option, 0),
		action:    action,
	}
}

// Usage generates a usage string for the command that includes its name,
// a placeholder for arguments, and a placeholder for options if the command has any.
// The generated string is intended to be shown to users to demonstrate how to use the command.
func (c *Command) Usage() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Usage: %s", c.Name))
	if len(c.arguments) > 0 {
		builder.WriteString(" [arguments]")
	}
	if len(c.options) > 0 {
		builder.WriteString(" [options]")
	}
	return builder.String()
}

// AddArgument check whether given arg name is duplicate or not,
// then add it to param.Argument slice.
func (c *Command) AddArgument(arg *param.Argument) error {
	for _, argument := range c.arguments {
		if arg.Name == argument.Name {
			return fmt.Errorf("duplicate argument name %s", arg.Name)
		}
	}
	c.arguments = append(c.arguments, arg)
	return nil
}

// AddOption check whether given opt name is duplicate or not,
// then add it to param.Option slice.
func (c *Command) AddOption(opt *param.Option) error {
	for _, option := range c.options {
		if opt.Name == option.Name {
			return fmt.Errorf("duplicate option name %s", opt.Name)
		}
	}
	c.options = append(c.options, opt)
	return nil
}

// Execute runs the command with the provided input parameters.
// It separates the input parameters into arguments and options,
// validates them, and then calls the commands' Action function.
// It returns the result-string of the Action function or an error encountered during validation or execution.
func (c *Command) Execute(inputParams []string) (string, error) {
	args, opts, err := c.validate(inputParams)
	if err != nil {
		return "", err
	}

	return c.action(args, opts)
}

// validate parses and validates the input parameters for the command.
// It separates the input parameters into arguments and options,
// checks them against the command's requirements, and returns
// a slice of arguments and a map of option if they are valid.
// It returns an error if there are too few or too many arguments,
// or if an invalid option is provided.
func (c *Command) validate(inputParams []string) ([]string, map[string]param.Value, error) {
	args := make([]string, 0, len(c.arguments))
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

// parseOption takes an option parameter and the current index of inputParams,
// determines whether it's a flag or a regular option, and then processes it accordingly.
// It returns the name of the option, its value, and an error if the option is invalid
// or if there's a problem processing the value.
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

// validateArguments checks if the correct number of the arguments has been provided for the command.
// It returns an error if the number of provided arguments is less than required or greater than allowed.
func (c *Command) validateArguments(args []string) error {
	if len(args) < len(c.arguments) {
		return fmt.Errorf("not enough arguments: actual %d, expected %d", len(args), len(c.arguments))
	} else if len(args) > len(c.arguments) {
		return fmt.Errorf("too many arguments: actual %d, expected %d", len(args), len(c.arguments))
	}
	return nil
}

// initializeOptions creates a map with default values for all options defined in the command.
// It initializes flags with a default boolean value of false.
// The map is used to store the values of options parses from the input parameters.
func (c *Command) initializeOptions() map[string]param.Value {
	options := make(map[string]param.Value)
	for _, option := range c.options {
		if option.IsFlag {
			options[strings.TrimPrefix(option.Name, optionPrefix)] = *param.NewBoolParameterPtr(false)
		}
	}
	return options
}

// getOptionTypeAndFlag  retrieves the type and flag status of the option with the given name.
// It returns the type of the option, a boolean indicating if it's a flag,
// if the option does not exist in the command's options list.
func (c *Command) getOptionTypeAndFlag(optionName string) (param.Type, bool, error) {
	for _, option := range c.options {
		if option.Name == optionName {
			return option.Type, option.IsFlag, nil
		}
	}
	return -1, false, fmt.Errorf("invalid option %s", optionName)
}

// processFlagOption processes a flag option from the input parameters.
// A flag option does not take a value; it is simply present or absent.
// The method updates the flagOpts slice to remove processed options and
// returns the name of the option and its boolean value.
func (c *Command) processFlagOption(optionName string, inputParams []string, idxPtr *int, flagOpts []*param.Option) (string, *param.Value, error) {
	// Make sure the next parameter is not an argument not starting with `--`
	// Flag Option cannot have value
	if *idxPtr+1 < len(inputParams) && isArgument(inputParams[*idxPtr+1]) {
		return "", nil, fmt.Errorf("flag-option %s cannot have value", optionName)
	}

	c.removeFlagOption(optionName, flagOpts)
	optionName = strings.TrimPrefix(optionName, optionPrefix)
	return optionName, param.NewBoolParameterPtr(true), nil
}

// processRegularOption processes a regular (non-flag) option.
// It checks if the next parameter is a valid value for the option.
// parses it, and returns the option's name and its parsed value.
func (c *Command) processRegularOption(optionName string, optType param.Type, inputParams []string, idxPtr *int) (string, *param.Value, error) {
	// Make sure the  next parameter is not an option starting with `--`
	// Whether normal option is last parameter or not
	// Normal Option always accepts one argument
	if *idxPtr+1 >= len(inputParams) || !isArgument(inputParams[*idxPtr+1]) {
		typeStr := param.ParameterTypeToString(optType)
		return "", nil, fmt.Errorf("\"%s\" option requires a \"%s\" type argument", optionName, typeStr)
	}

	*idxPtr++
	value := inputParams[*idxPtr]
	paramValue, err := param.ToParameterValue(value, optType)
	if err != nil {
		return "", nil, fmt.Errorf("invalid option \"%s\": %w", optionName, err)
	}

	optionName = strings.TrimPrefix(optionName, optionPrefix)
	return optionName, paramValue, nil
}

// removeFlagOption removes a flag option from the slice of flag options once it has been processed.
// This ensures that each flag  option is only processed once.
func (c *Command) removeFlagOption(name string, flagOpts []*param.Option) {
	for i, opt := range flagOpts {
		if opt.Name == name {
			flagOpts[i] = flagOpts[len(flagOpts)-1]
			flagOpts = flagOpts[:len(flagOpts)-1]
			break
		}
	}
}

// isArgument checks if the provided string is an argument, i.e., does not start with '--'.
// It returns true if the string is an argument, false otherwise.
func isArgument(p string) bool {
	if strings.HasPrefix(p, "--") {
		return false
	}
	return true
}

// flagOptions constructs a lice of all flag options from the command's options.
// Flag options are options that do not take an argument and act as boolean toggles.
func (c *Command) flagOptions() []*param.Option {
	var flagOpts []*param.Option

	for _, option := range c.options {
		if option.IsFlag {
			flagOpts = append(flagOpts, option)
		}
	}
	return flagOpts
}
