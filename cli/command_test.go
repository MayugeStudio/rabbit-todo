package cli

import (
	"fmt"
	"rabbit-todo/cli/param"
	"reflect"
	"strings"
	"testing"
)

func TestNewCommand(t *testing.T) {
	type inputType struct {
		commandName string
		arguments   []*parameter.Argument
		options     []*parameter.Option
	}
	type testCase struct {
		testName string
		input    inputType
		want     string
	}
	tests := []testCase{
		{
			testName: "Ok-WithOneArgAndOneOpt",
			input: inputType{
				commandName: "test-command",
				arguments:   []*parameter.Argument{{Name: "Hello", Type: parameter.STRING}},
				options:     []*parameter.Option{{Name: "--hello", Type: parameter.STRING}},
			},
			want: "Usage: test-command [arguments] [options]",
		},
		{
			testName: "Ok-WithTwoArgAndTwoOpt",
			input: inputType{
				commandName: "test-command",
				arguments: []*parameter.Argument{
					{
						Name: "Hello",
						Type: parameter.STRING,
					},
					{
						Name: "World",
						Type: parameter.STRING,
					},
				},
				options: []*parameter.Option{
					{
						Name:   "--hello",
						Type:   parameter.STRING,
						IsFlag: false,
					},
					{
						Name:   "--world",
						Type:   parameter.STRING,
						IsFlag: false,
					},
				},
			},
			want: "Usage: test-command [arguments] [options]",
		},
		{
			testName: "Ok-WithOneArgAndZeroOpt",
			input: inputType{
				commandName: "test-command",
				arguments: []*parameter.Argument{
					{
						Name: "OneArg",
						Type: parameter.STRING,
					},
				},
				options: []*parameter.Option{},
			},
			want: "Usage: test-command [arguments]",
		},
		{
			testName: "Ok-WithZeroArgAndOneOpt",
			input: inputType{
				commandName: "test-command",
				arguments:   []*parameter.Argument{},
				options: []*parameter.Option{
					{
						Name:   "--one-arg",
						Type:   parameter.STRING,
						IsFlag: false,
					},
				},
			},
			want: "Usage: test-command [options]",
		},
		{
			testName: "Ok-WithZeroArgAndZeroOpt",
			input: inputType{
				commandName: "test-command",
				arguments:   []*parameter.Argument{},
				options:     []*parameter.Option{},
			},
			want: "Usage: test-command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := NewCommand(tt.input.commandName, tt.input.arguments, tt.input.options, nil)

			if got.Name != tt.input.commandName {
				t.Errorf("Command.Name = %v, want %v", got.Name, tt.input.commandName)
			}

			if !reflect.DeepEqual(got.Arguments, tt.input.arguments) {
				t.Errorf("Command.Arguments = %v, want %v", got.Arguments, tt.input.arguments)
			}

			if !reflect.DeepEqual(got.Options, tt.input.options) {
				t.Errorf("Command.Options = %v, want %v", got.Options, tt.input.options)
			}

			if got.Usage != tt.want {
				t.Errorf("Command.Usage = %v, want %v", got.Usage, tt.want)
			}
		})
	}
}

func TestCommand_Execute_With_Arguments(t *testing.T) {
	testAction := func(args []string, opts map[string]parameter.ParamValue) (string, error) {
		return strings.Join(args, ""), nil
	}

	type inputType struct {
		command     Command
		inputParams []string
	}

	type testCase struct {
		testName   string
		input      inputType
		want       string
		wantErr    bool
		wantErrStr string
	}

	tests := []testCase{
		{
			testName: "Ok-HelloWorld",
			input: inputType{
				command: Command{
					Name: "return-HelloWorld-command",
					Arguments: []*parameter.Argument{
						{
							Name: "a",
							Type: parameter.STRING,
						},
						{
							Name: "b",
							Type: parameter.STRING,
						},
					},
					Options: []*parameter.Option{},
					Action:  testAction,
					Usage:   "",
				},
				inputParams: []string{"Hello", "World"},
			},
			want:       "HelloWorld",
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Error-NotEnoughArguments",
			input: inputType{
				command: Command{
					Name: "fail-command",
					Arguments: []*parameter.Argument{
						{
							Name: "arg1",
							Type: parameter.INT,
						},
						{
							Name: "arg2",
							Type: parameter.STRING,
						},
					},
					Options: []*parameter.Option{},
					Action:  testAction,
					Usage:   "",
				},
				inputParams: []string{"one-arg"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "not enough arguments: actual 1, expected 2",
		},
		{
			testName: "Error-TooManyArguments",
			input: inputType{
				command: Command{
					Name: "fail-command",
					Arguments: []*parameter.Argument{
						{
							Name: "arg1",
							Type: parameter.INT,
						},
					},
					Options: []*parameter.Option{},
					Action:  testAction,
					Usage:   "",
				},
				inputParams: []string{"arg1", "arg2", "arg3"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "too many arguments: actual 3, expected 1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := tc.input.command.Execute(tc.input.inputParams)
			isErr := err != nil
			if isErr != tc.wantErr {
				t.Fatalf("Command.Execute() error = %v, wantError %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if err.Error() != tc.wantErrStr {
					t.Errorf("Command.Execute() error = %q, wantErrStr %q", err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Command.Execute() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCommand_Execute_With_Options(t *testing.T) {
	testActionAddStrings := func(args []string, opts map[string]parameter.ParamValue) (string, error) {
		var opt parameter.ParamValue
		opt = opts["opt1"]
		arg1 := opt.StringVal

		opt = opts["opt2"]
		arg2 := opt.StringVal
		return arg1 + arg2, nil
	}
	testActionAddIntegers := func(args []string, opts map[string]parameter.ParamValue) (string, error) {
		var opt parameter.ParamValue
		opt = opts["opt1"]
		arg1 := opt.IntVal

		opt = opts["opt2"]
		arg2 := opt.IntVal
		return fmt.Sprintf("%d", arg1+arg2), nil
	}

	type inputType struct {
		command     Command
		inputParams []string
	}

	type testCase struct {
		testName   string
		input      inputType
		want       string
		wantErr    bool
		wantErrStr string
	}

	tests := []testCase{
		{
			testName: "Ok-WithTwoStringOptionCommand",
			input: inputType{
				command: NewCommand(
					"two-option-command",
					[]*parameter.Argument{},
					[]*parameter.Option{
						{Name: "--opt1", Type: parameter.STRING, IsFlag: false},
						{Name: "--opt2", Type: parameter.STRING, IsFlag: false},
					},
					testActionAddStrings),
				inputParams: []string{"--opt1", "Hello", "--opt2", "World"},
			},
			want:       "HelloWorld",
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-WithTwoIntegerOptionCommand",
			input: inputType{
				command: NewCommand(
					"two-option-command",
					[]*parameter.Argument{},
					[]*parameter.Option{
						{Name: "--opt1", Type: parameter.INT, IsFlag: false},
						{Name: "--opt2", Type: parameter.INT, IsFlag: false},
					},
					testActionAddIntegers),
				inputParams: []string{"--opt1", "1", "--opt2", "2"},
			},
			want:       "3",
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Error-InvalidOptionType",
			input: inputType{
				command: NewCommand(
					"two-option-command",
					[]*parameter.Argument{},
					[]*parameter.Option{
						{Name: "--opt1", Type: -1, IsFlag: false},
						{Name: "--opt2", Type: parameter.INT, IsFlag: false},
					},
					testActionAddIntegers),
				inputParams: []string{"--opt1", "1", "--opt2", "2"},
			},
			want:       "3",
			wantErr:    true,
			wantErrStr: "invalid option \"--opt1\": unknown parameter type -1",
		},
		{
			testName: "Error-InvalidOption",
			input: inputType{
				command: Command{
					Name:      "fail-command",
					Arguments: []*parameter.Argument{},
					Options: []*parameter.Option{
						{Name: "--opt-1", Type: parameter.INT, IsFlag: false},
						{Name: "--opt-2", Type: parameter.INT, IsFlag: false},
					},
					Action: testActionAddIntegers,
					Usage:  "",
				},
				inputParams: []string{"--opt-1", "1", "--opt-2", "2", "--opt-3", "3"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "invalid option --opt-3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := tc.input.command.Execute(tc.input.inputParams)
			isErr := err != nil
			if isErr != tc.wantErr {
				t.Fatalf("Command.Execute() error = %v, wantError %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if err.Error() != tc.wantErrStr {
					t.Errorf("Command.Execute() error = %q, wantErrStr %q", err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Command.Execute() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCommand_Execute_With_FlagOptions(t *testing.T) {
	testAction := func(args []string, opts map[string]parameter.ParamValue) (string, error) {
		var opt1 bool
		var opt2 string

		opt1 = opts["opt1"].BoolVal
		opt2 = opts["opt2"].StringVal

		if opt1 {
			return opt2 + " Hello!", nil
		}

		return opt2 + " Bye!", nil
	}

	type inputType struct {
		command     Command
		inputParams []string
	}

	type testCase struct {
		testName   string
		input      inputType
		want       string
		wantErr    bool
		wantErrStr string
	}

	tests := []testCase{
		{
			testName: "Ok-TrueFlag",
			input: inputType{
				command: Command{
					Name:      "flag-command",
					Arguments: nil,
					Options: []*parameter.Option{
						{
							Name:   "--opt1",
							Type:   parameter.BOOL,
							IsFlag: true,
						},
						{
							Name:   "--opt2",
							Type:   parameter.STRING,
							IsFlag: false,
						},
					},
					Action: testAction,
					Usage:  "Usage: flag-command [options]",
				},
				inputParams: []string{"--opt1", "--opt2", "John"},
			},
			want:       "John Hello!",
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-FalseFlag",
			input: inputType{
				command: Command{
					Name:      "flag-command",
					Arguments: nil,
					Options: []*parameter.Option{
						{
							Name:   "--opt1",
							Type:   parameter.BOOL,
							IsFlag: true,
						},
						{
							Name:   "--opt2",
							Type:   parameter.STRING,
							IsFlag: false,
						},
					},
					Action: testAction,
					Usage:  "Usage: flag-command [options]",
				},
				inputParams: []string{"--opt2", "John"},
			},
			want:       "John Bye!",
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Error-FlagHasValue",
			input: inputType{
				command: Command{
					Name:      "flag-command",
					Arguments: nil,
					Options: []*parameter.Option{
						{
							Name:   "--opt1",
							Type:   parameter.BOOL,
							IsFlag: true,
						},
						{
							Name:   "--opt2",
							Type:   parameter.STRING,
							IsFlag: false,
						},
					},
					Action: testAction,
					Usage:  "Usage: flag-command [options]",
				},
				inputParams: []string{"--opt1", "Mike", "--opt2", "John"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "flag-option --opt1 cannot have value",
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := tc.input.command.Execute(tc.input.inputParams)
			isErr := err != nil
			if isErr != tc.wantErr {
				t.Fatalf("Command.Execute() error = %v, wantError %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if err.Error() != tc.wantErrStr {
					t.Errorf("Command.Execute() error = %q, wantErrStr %q", err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Command.Execute() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCommand_Execute_Integration(t *testing.T) {
	type inputType struct {
		command Command
		args    []string
	}
	type testCase struct {
		testName   string
		input      inputType
		want       string
		wantErr    bool
		wantErrStr string
	}

	testAction1 := func(args []string, opts map[string]parameter.ParamValue) (string, error) {
		to := opts["to"].StringVal
		from := opts["from"].StringVal
		str := "from:"
		str += from
		str += " -> "
		str += "\""
		str += strings.Join(args, " ")
		str += "\""
		str += " -> "
		str += "to:"
		str += to
		return str, nil
	}

	toOption, _ := parameter.NewOption("--to", parameter.STRING)
	fromOption, _ := parameter.NewOption("--from", parameter.STRING)
	msgArg1, _ := parameter.NewArgument("msg-1", parameter.STRING)
	msgArg2, _ := parameter.NewArgument("msg-2", parameter.STRING)

	command1 := NewCommand("command-1", []*parameter.Argument{msgArg1, msgArg2}, []*parameter.Option{toOption, fromOption}, testAction1)

	tests := []testCase{
		{
			testName: "Ok-ExecuteSuccessfully",
			input: inputType{
				command: command1,
				args:    []string{"Hello", "!!", "--to", "John", "--from", "Mike"},
			},
			want:       "from:Mike -> \"Hello !!\" -> to:John",
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Error-MissingArguments",
			input: inputType{
				command: command1,
				args:    []string{"Hello", "--to", "John", "--from", "Mike"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "not enough arguments: actual 1, expected 2",
		},
		{
			testName: "Error-MissingArgumentOfOption'--to'",
			input: inputType{
				command: command1,
				args:    []string{"Hello", "!!", "--to", "--from", "Mike"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "\"--to\" option require one \"string\" type argument",
		},
		{
			testName: "Error-MissingArgumentOfOption'--from'",
			input: inputType{
				command: command1,
				args:    []string{"Hello", "!!", "--to", "John", "--from"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "\"--from\" option require one \"string\" type argument",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := tc.input.command.Execute(tc.input.args)
			isErr := err != nil
			if isErr != tc.wantErr {
				t.Fatalf("Command.Execute() error = %v, wantError %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if err.Error() != tc.wantErrStr {
					t.Errorf("Command.Execute() error = %q, wantErrStr %q", err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Command.Execute() = %v, want %v", got, tc.want)
			}
		})
	}
}
