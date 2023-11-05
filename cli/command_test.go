package cli

import (
	"fmt"
	"rabbit-todo/cli/param"
	"reflect"
	"strings"
	"testing"
)

func TestNewCommand(t *testing.T) {
	type inputType struct{ commandName string }
	type testCase struct {
		testName string
		input    inputType
		want     string
	}
	tests := []testCase{
		{
			testName: "Ok-CreateSuccessfully",
			input: inputType{
				commandName: "test-command",
			},
			want: "test-command",
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			got := NewCommand(tc.input.commandName, nil)

			if got.Name != tc.input.commandName {
				t.Errorf("Command.Name = %v, want %v", got.Name, tc.input.commandName)
			}
		})
	}
}

func TestCommand_Usage(t *testing.T) {
	type inputType struct {
		command   Command
		arguments []*param.Argument
		options   []*param.Option
	}
	type testCase struct {
		testName string
		input    inputType
		want     string
	}
	tests := []testCase{
		{
			testName: "Ok-TwoArgAndTwoOpt",
			input: inputType{
				command: NewCommand("test-command", nil),
				arguments: []*param.Argument{
					{Name: "arg1", Type: param.STRING},
					{Name: "arg2", Type: param.STRING},
				},
				options: []*param.Option{
					{Name: "opt1", Type: param.INT},
					{Name: "opt2", Type: param.INT},
				},
			},
			want: "Usage: test-command [arguments] [options]",
		},
		{
			testName: "Ok-OneArgAndZeroOpt",
			input: inputType{
				command:   NewCommand("test-command", nil),
				arguments: []*param.Argument{{Name: "arg1", Type: param.STRING}},
				options:   nil,
			},
			want: "Usage: test-command [arguments]",
		},
		{
			testName: "Ok-ZeroArgAndOneOpt",
			input: inputType{
				command:   NewCommand("test-command", nil),
				arguments: nil,
				options:   []*param.Option{{Name: "opt1", Type: param.INT}},
			},
			want: "Usage: test-command [options]",
		},
		{
			testName: "Ok-ZeroArgAndZeroOpt",
			input: inputType{
				command:   NewCommand("test-command", nil),
				arguments: nil,
				options:   nil,
			},
			want: "Usage: test-command",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			for _, argument := range tc.input.arguments {
				err := tc.input.command.AddArgument(argument)
				if err != nil {
					t.Fatalf("Command.AddArgument() error = %v", err)
				}
			}
			for _, option := range tc.input.options {
				err := tc.input.command.AddOption(option)
				if err != nil {
					t.Fatalf("Command.AddOption() error = %v", err)
				}
			}
			got := tc.input.command.Usage()
			if got != tc.want {
				t.Errorf("Command.Usage() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCommand_AddArgument(t *testing.T) {
	type inputType struct {
		arguments []*param.Argument
	}
	type testCase struct {
		testName   string
		input      inputType
		want       []*param.Argument
		wantErr    bool
		wantErrStr string
	}
	tests := []testCase{
		{
			testName: "Ok-AddArgumentSuccessfully",
			input: inputType{
				arguments: []*param.Argument{{Name: "arg1", Type: param.STRING}},
			},
			want: []*param.Argument{{Name: "arg1", Type: param.STRING}},
		},
		{
			testName: "Error-DuplicateArgumentName",
			input: inputType{
				arguments: []*param.Argument{
					{Name: "arg1", Type: param.STRING},
					{Name: "arg2", Type: param.STRING},
					{Name: "arg2", Type: param.STRING},
				},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "duplicate argument name arg2",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			var Err error = nil
			cmd := NewCommand("test-command", nil)
			for _, argument := range tc.input.arguments {
				err := cmd.AddArgument(argument)
				if err != nil {
					Err = err
				}
			}
			isErr := Err != nil
			if isErr != tc.wantErr {
				t.Fatalf("Command.AddArgument() error = %v, wantError %v", Err, tc.wantErr)
			}
			if tc.wantErr {
				if Err.Error() != tc.wantErrStr {
					t.Errorf("Command.AddArgument() error = %q, wantErrStr %q", Err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(cmd.arguments, tc.want) {
				t.Errorf("Command.Execute() = %v, want %v", cmd.arguments, tc.want)
			}
		})
	}
}

func TestCommand_AddOption(t *testing.T) {
	type inputType struct {
		options []*param.Option
	}
	type testCase struct {
		testName   string
		input      inputType
		want       []*param.Option
		wantErr    bool
		wantErrStr string
	}
	tests := []testCase{
		{
			testName: "Ok-AddOptionSuccessfully",
			input: inputType{
				options: []*param.Option{{Name: "opt1", Type: param.INT}},
			},
			want: []*param.Option{{Name: "opt1", Type: param.INT}},
		},
		{
			testName: "Error-DuplicateOptionName",
			input: inputType{
				options: []*param.Option{
					{Name: "opt1", Type: param.INT},
					{Name: "opt2", Type: param.INT},
					{Name: "opt2", Type: param.INT},
				},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "duplicate option name opt2",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			cmd := NewCommand("test-command", nil)
			var Err error = nil
			for _, option := range tc.input.options {
				err := cmd.AddOption(option)
				if err != nil {
					Err = err
				}
			}

			isErr := Err != nil
			if isErr != tc.wantErr {
				t.Fatalf("Command.AddOption() error = %v, wantError %v", Err, tc.wantErr)
			}
			if tc.wantErr {
				if Err.Error() != tc.wantErrStr {
					t.Errorf("Command.AddOption() error = %q, wantErrStr %q", Err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(cmd.options, tc.want) {
				t.Errorf("Command.AddOption() = %v, want %v", cmd.options, tc.want)
			}
		})
	}
}

func TestCommand_Execute_With_Arguments(t *testing.T) {
	testAction := func(args []string, opts map[string]param.Value) (string, error) {
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
					arguments: []*param.Argument{
						{
							Name: "a",
							Type: param.STRING,
						},
						{
							Name: "b",
							Type: param.STRING,
						},
					},
					options: []*param.Option{},
					action:  testAction,
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
					arguments: []*param.Argument{
						{
							Name: "arg1",
							Type: param.INT,
						},
						{
							Name: "arg2",
							Type: param.STRING,
						},
					},
					options: []*param.Option{},
					action:  testAction,
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
					arguments: []*param.Argument{
						{
							Name: "arg1",
							Type: param.INT,
						},
					},
					options: []*param.Option{},
					action:  testAction,
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
	testActionAddStrings := func(args []string, opts map[string]param.Value) (string, error) {
		opt1 := opts["opt1"].StringVal
		opt2 := opts["opt2"].StringVal
		return opt1 + opt2, nil
	}
	testActionAddIntegers := func(args []string, opts map[string]param.Value) (string, error) {
		opt1 := opts["opt1"].IntVal
		opt2 := opts["opt2"].IntVal
		return fmt.Sprintf("%d", opt1+opt2), nil
	}

	type inputType struct {
		command     Command
		arguments   []*param.Argument
		options     []*param.Option
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
					testActionAddStrings),
				arguments: []*param.Argument{},
				options: []*param.Option{
					{Name: "--opt1", Type: param.STRING, IsFlag: false},
					{Name: "--opt2", Type: param.STRING, IsFlag: false},
				},
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
					testActionAddIntegers),
				arguments: []*param.Argument{},
				options: []*param.Option{
					{Name: "--opt1", Type: param.INT, IsFlag: false},
					{Name: "--opt2", Type: param.INT, IsFlag: false},
				},
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
					testActionAddIntegers),
				arguments: []*param.Argument{},
				options: []*param.Option{
					{Name: "--opt1", Type: -1, IsFlag: false},
					{Name: "--opt2", Type: param.INT, IsFlag: false},
				},
				inputParams: []string{"--opt1", "1", "--opt2", "2"},
			},
			want:       "3",
			wantErr:    true,
			wantErrStr: "invalid option \"--opt1\": unknown parameter type -1",
		},
		{
			testName: "Error-InvalidOption",
			input: inputType{
				command:   NewCommand("fail-command", testActionAddIntegers),
				arguments: []*param.Argument{},
				options: []*param.Option{
					{Name: "--opt-1", Type: param.INT, IsFlag: false},
					{Name: "--opt-2", Type: param.INT, IsFlag: false},
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
			for _, argument := range tc.input.arguments {
				err := tc.input.command.AddArgument(argument)
				if err != nil {
					t.Fatalf("Command.AddArgument() error = %v", err)
				}
			}
			for _, option := range tc.input.options {
				err := tc.input.command.AddOption(option)
				if err != nil {
					t.Fatalf("Command.AddOption() error = %v", err)
				}
			}
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
	testAction := func(args []string, opts map[string]param.Value) (string, error) {
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
					arguments: nil,
					options: []*param.Option{
						{
							Name:   "--opt1",
							Type:   param.BOOL,
							IsFlag: true,
						},
						{
							Name:   "--opt2",
							Type:   param.STRING,
							IsFlag: false,
						},
					},
					action: testAction,
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
					arguments: nil,
					options: []*param.Option{
						{
							Name:   "--opt1",
							Type:   param.BOOL,
							IsFlag: true,
						},
						{
							Name:   "--opt2",
							Type:   param.STRING,
							IsFlag: false,
						},
					},
					action: testAction,
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
					arguments: nil,
					options: []*param.Option{
						{
							Name:   "--opt1",
							Type:   param.BOOL,
							IsFlag: true,
						},
						{
							Name:   "--opt2",
							Type:   param.STRING,
							IsFlag: false,
						},
					},
					action: testAction,
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

	testAction1 := func(args []string, opts map[string]param.Value) (string, error) {
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

	toOption, _ := param.NewOption("--to", param.STRING)
	fromOption, _ := param.NewOption("--from", param.STRING)
	msgArg1, _ := param.NewArgument("msg-1", param.STRING)
	msgArg2, _ := param.NewArgument("msg-2", param.STRING)

	command1 := NewCommand("command-1", testAction1)
	_ = command1.AddArgument(msgArg1)
	_ = command1.AddArgument(msgArg2)
	_ = command1.AddOption(toOption)
	_ = command1.AddOption(fromOption)

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
			wantErrStr: "\"--to\" option requires a \"string\" type argument",
		},
		{
			testName: "Error-MissingArgumentOfOption'--from'",
			input: inputType{
				command: command1,
				args:    []string{"Hello", "!!", "--to", "John", "--from"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "\"--from\" option requires a \"string\" type argument",
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
