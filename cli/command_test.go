package cli

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewCommand(t *testing.T) {
	type args struct {
		commandName string
		arguments   []*Argument
		options     []*Option
	}
	tests := []struct {
		testName string
		args     args
		want     string
	}{
		{
			testName: "With 1 arg and 1 opt",
			args: args{
				commandName: "test-command",
				arguments:   []*Argument{{Name: "Hello", Type: STRING}},
				options:     []*Option{{Name: "--hello", Type: STRING}},
			},
			want: "Usage: test-command [arguments] [options]",
		},
		{
			testName: "With 2 arg and 2 opt",
			args: args{
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
			testName: "With 1 arg and 0 opt",
			args: args{
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
			testName: "With 0 arg and 1 opt",
			args: args{
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
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := NewCommand(tt.args.commandName, tt.args.arguments, tt.args.options, nil)

			if got.Name != tt.args.commandName {
				t.Errorf("NewCommand().Name = %v, want %v", got.Name, tt.args.commandName)
			}

			if !reflect.DeepEqual(got.Arguments, tt.args.arguments) {
				t.Errorf("NewCommand().Arguments = %v, want %v", got.Arguments, tt.args.arguments)
			}

			if !reflect.DeepEqual(got.Options, tt.args.options) {
				t.Errorf("NewCommand().Options = %v, want %v", got.Options, tt.args.options)
			}

			if got.Usage != tt.want {
				t.Errorf("NewCommand().Usage = %v, want %v", got.Usage, tt.want)
			}
		})
	}
}

func TestCommand_Execute(t *testing.T) {
	var testAction Action
	testAction = func(args []string, opts []string) (string, error) { return strings.Join(args, ""), nil }

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
			wantErrStr: "not enough arguments",
		},
		{
			testName: "Test-Fail-TooManyOptions",
			input: inputType{
				command: Command{
					Name:      "fail-command",
					Arguments: []*Argument{},
					Options: []*Option{
						{
							Name:   "--opt-1",
							Type:   INT,
							IsFlag: false,
						},
						{
							Name:   "--opt-2",
							Type:   INT,
							IsFlag: false,
						},
					},
					Action: testAction,
					Usage:  "",
				},
				inputParams: []string{
					"--opt-1",
					"--opt-2",
					"--opt-3",
				},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "too many options",
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
					t.Errorf("NewOption() error = %q, wantErrStr %q", err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Command.Execute() = %v, want %v", got, tc.want)
			}
		})
	}
}
