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
	testAction = func(args []string, opts []string) (string, error) {
		result := ""

		for _, arg := range args {
			result += arg
		}
		return result, nil
	}

	type Parameters struct {
		arguments []*Argument
		options   []*Option
	}

	type Input struct {
		arguments []string
		options   []string
	}

	type args struct {
		commandName string
		parameters  Parameters
		input       Input
		action      Action
		wantErr     bool
		wantErrStr  string
	}

	tests := []struct {
		testName string
		args     args
		want     string
	}{
		{
			testName: "WantSuccess: Expect HelloWorld string",
			args: args{
				commandName: "return-HelloWorld-command",
				parameters: Parameters{
					arguments: []*Argument{
						{
							Name: "a",
							Type: STRING,
						},
						{
							Name: "b",
							Type: STRING,
						},
					},
					options: []*Option{},
				},
				input: Input{
					arguments: []string{"Hello", "World"},
					options:   []string{},
				},
				action:  testAction,
				wantErr: false,
			},
			want: "HelloWorld",
		},
		{
			testName: "WantError: Not-Enough-Arguments",
			args: args{
				commandName: "fail-command",
				parameters: Parameters{
					arguments: []*Argument{
						{
							Name: "arg1",
							Type: STRING,
						},
						{
							Name: "arg2",
							Type: STRING,
						},
					},
					options: []*Option{},
				},
				input: Input{
					arguments: []string{"one-arg"},
					options:   []string{},
				},
				action:     testAction,
				wantErr:    true,
				wantErrStr: "not enough arguments, expected: 2, got: 1",
			},
		},
		{
			testName: "WantError: Too-Many-Options",
			args: args{
				commandName: "fail-command",
				parameters: Parameters{
					arguments: []*Argument{
						{
							Name: "a",
							Type: STRING,
						},
						{
							Name: "b",
							Type: STRING,
						},
					},
					options: []*Option{
						{
							Name:   "--option-1",
							Type:   STRING,
							IsFlag: false,
						},
					},
				},
				input: Input{
					arguments: []string{"Hello", "World"},
					options:   []string{"--input-option-1", "--input-option-2"},
				},
				action:     testAction,
				wantErr:    true,
				wantErrStr: "too many options",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := NewCommand(tt.args.commandName, tt.args.parameters.arguments, tt.args.parameters.options, tt.args.action)
			if cmd == nil {
				t.Fatalf("unexpected issue ocurred, got nil from NewCommand()")
			}
			got, err := cmd.Execute(tt.args.input.arguments, tt.args.input.options)
			if tt.args.wantErr {
				// Expected error
				if err == nil {
					t.Errorf("Command.Execute() = (%v, nil), want (nil, %q)", got, tt.args.wantErrStr)
				}
				// Error message check
				if !strings.Contains(err.Error(), tt.args.wantErrStr) {
					t.Errorf("Command.Execute() = (nil, %q), want (nil, %q)", err, tt.args.wantErrStr)
				}
			} else {
				if got != tt.want {
					t.Errorf("Command.Execute() = (%v, nil), want (%v, nil)", got, tt.want)
				}
			}
		})
	}
}
