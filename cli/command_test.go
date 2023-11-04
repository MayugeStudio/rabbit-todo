package cli

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewCommand(t *testing.T) {
	type inputType struct {
		commandName string
		arguments   []*Argument
		options     []*Option
	}
	type testCase struct {
		testName string
		input    inputType
		want     string
	}
	tests := []testCase{
		{
			testName: "Test-Ok-WithOneArgAndOneOpt",
			input: inputType{
				commandName: "test-command",
				arguments:   []*Argument{{Name: "Hello", Type: STRING}},
				options:     []*Option{{Name: "--hello", Type: STRING}},
			},
			want: "Usage: test-command [arguments] [options]",
		},
		{
			testName: "Test-Ok-WithTwoArgAndTwoOpt",
			input: inputType{
				commandName: "test-command",
				arguments: []*Argument{
					{
						Name: "Hello",
						Type: STRING,
					},
					{
						Name: "World",
						Type: STRING,
					},
				},
				options: []*Option{
					{
						Name:   "--hello",
						Type:   STRING,
						IsFlag: false,
					},
					{
						Name:   "--world",
						Type:   STRING,
						IsFlag: false,
					},
				},
			},
			want: "Usage: test-command [arguments] [options]",
		},
		{
			testName: "Test-Ok-WithOneArgAndZeroOpt",
			input: inputType{
				commandName: "test-command",
				arguments: []*Argument{
					{
						Name: "OneArg",
						Type: STRING,
					},
				},
				options: []*Option{},
			},
			want: "Usage: test-command [arguments]",
		},
		{
			testName: "Test-Ok-WithZeroArgAndOneOpt",
			input: inputType{
				commandName: "test-command",
				arguments:   []*Argument{},
				options: []*Option{
					{
						Name:   "--one-arg",
						Type:   STRING,
						IsFlag: false,
					},
				},
			},
			want: "Usage: test-command [options]",
		},
		{
			testName: "Test-Ok-WithZeroArgAndZeroOpt",
			input: inputType{
				commandName: "test-command",
				arguments:   []*Argument{},
				options:     []*Option{},
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
	testAction := func(args []string, opts map[string]OptionValue) (string, error) { return strings.Join(args, ""), nil }

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
			testName: "Test-Ok-HelloWorld",
			input: inputType{
				command: Command{
					Name: "return-HelloWorld-command",
					Arguments: []*Argument{
						{
							Name: "a",
							Type: STRING,
						},
						{
							Name: "b",
							Type: STRING,
						},
					},
					Options: []*Option{},
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
			testName: "Test-Fail-NotEnoughArguments",
			input: inputType{
				command: Command{
					Name: "fail-command",
					Arguments: []*Argument{
						{
							Name: "arg1",
							Type: INT,
						},
						{
							Name: "arg2",
							Type: STRING,
						},
					},
					Options: []*Option{},
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
			testName: "Test-Fail-TooManyArguments",
			input: inputType{
				command: Command{
					Name: "fail-command",
					Arguments: []*Argument{
						{
							Name: "arg1",
							Type: INT,
						},
					},
					Options: []*Option{},
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
	testActionAddStrings := func(args []string, opts map[string]OptionValue) (string, error) {
		var opt OptionValue
		opt = opts["opt1"]
		arg1 := opt.StringVal

		opt = opts["opt2"]
		arg2 := opt.StringVal
		return arg1 + arg2, nil
	}
	testActionAddIntegers := func(args []string, opts map[string]OptionValue) (string, error) {
		var opt OptionValue
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
			testName: "Test-Ok-WithTwoStringOptionCommand",
			input: inputType{
				command: NewCommand(
					"two-option-command",
					[]*Argument{},
					[]*Option{
						{Name: "--opt1", Type: STRING, IsFlag: false},
						{Name: "--opt2", Type: STRING, IsFlag: false},
					},
					testActionAddStrings),
				inputParams: []string{"--opt1", "Hello", "--opt2", "World"},
			},
			want:       "HelloWorld",
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Test-Ok-WithTwoIntegerOptionCommand",
			input: inputType{
				command: NewCommand(
					"two-option-command",
					[]*Argument{},
					[]*Option{
						{Name: "--opt1", Type: INT, IsFlag: false},
						{Name: "--opt2", Type: INT, IsFlag: false},
					},
					testActionAddIntegers),
				inputParams: []string{"--opt1", "1", "--opt2", "2"},
			},
			want:       "3",
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Test-Err-InvalidOptionType",
			input: inputType{
				command: NewCommand(
					"two-option-command",
					[]*Argument{},
					[]*Option{
						{Name: "--opt1", Type: -1, IsFlag: false},
						{Name: "--opt2", Type: INT, IsFlag: false},
					},
					testActionAddIntegers),
				inputParams: []string{"--opt1", "1", "--opt2", "2"},
			},
			want:       "3",
			wantErr:    true,
			wantErrStr: "invalid option \"--opt1\": unknown parameter type -1",
		},
		{
			testName: "Test-Err-InvalidOption",
			input: inputType{
				command: Command{
					Name:      "fail-command",
					Arguments: []*Argument{},
					Options: []*Option{
						{Name: "--opt-1", Type: INT, IsFlag: false},
						{Name: "--opt-2", Type: INT, IsFlag: false},
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
