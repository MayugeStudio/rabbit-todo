package cli

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		testName     string
		commandName  string
		args         []*Argument
		opts         []*Option
		wantUsageStr string
	}{
		{
			testName:     "With 1 arg and 1 opt",
			commandName:  "test-command",
			args:         []*Argument{{Name: "Hello", Type: STRING}},
			opts:         []*Option{{Name: "--hello", Type: STRING}},
			wantUsageStr: "Usage: test-command [arguments] [options]",
		},
		{
			testName:    "With 2 arg and 2 opt",
			commandName: "test-command",
			args: []*Argument{
				{
					Name: "Hello",
					Type: STRING,
				},
				{
					Name: "World",
					Type: STRING,
				},
			},
			opts: []*Option{
				NewOption("--hello", STRING),
				NewOption("--world", STRING),
			},
			wantUsageStr: "Usage: test-command [arguments] [options]",
		},
		{
			testName:    "With 1 arg and 0 opt",
			commandName: "test-command",
			args: []*Argument{
				{
					Name: "OneArg",
					Type: STRING,
				},
			},
			opts:         []*Option{},
			wantUsageStr: "Usage: test-command [arguments]",
		},
		{
			testName:    "With 0 arg and 1 opt",
			commandName: "test-command",
			args:        []*Argument{},
			opts: []*Option{
				NewOption("--one-arg", STRING),
			},
			wantUsageStr: "Usage: test-command [options]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := NewCommand(tt.commandName, tt.args, tt.opts, nil)

			if got.Name != tt.commandName {
				t.Errorf("got: name = %v, want %v", got.Name, tt.commandName)
			}

			if !reflect.DeepEqual(got.Arguments, tt.args) {
				t.Errorf("got: arguments = %v, want %v", got.Arguments, tt.args)
			}

			if !reflect.DeepEqual(got.Options, tt.opts) {
				t.Errorf("got: options = %v, want %v", got.Options, tt.opts)
			}

			if got.Usage != tt.wantUsageStr {
				t.Errorf("got: %v, want: %v", got.Usage, tt.wantUsageStr)
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
		args []*Argument
		opts []*Option
	}

	type Input struct {
		args []string
		opts []string
	}

	tests := []struct {
		testName    string
		commandName string
		parameters  Parameters
		input       Input
		action      Action
		wantResult  string
		wantErr     bool
		wantErrStr  string
	}{
		{
			testName:    "WantSuccess: Expect HelloWorld string",
			commandName: "return-HelloWorld-command",
			parameters: Parameters{
				args: []*Argument{
					{
						Name: "a",
						Type: STRING,
					},
					{
						Name: "b",
						Type: STRING,
					},
				},
				opts: []*Option{},
			},
			input: Input{
				args: []string{"Hello", "World"},
				opts: []string{},
			},
			action:     testAction,
			wantResult: "HelloWorld",
			wantErr:    false,
		},
		{
			testName:    "WantError: Not-Enough-Arguments",
			commandName: "fail-command",
			parameters: Parameters{
				args: []*Argument{
					{
						Name: "arg1",
						Type: STRING,
					},
					{
						Name: "arg2",
						Type: STRING,
					},
				},
				opts: []*Option{},
			},
			input: Input{
				args: []string{"one-arg"},
				opts: []string{},
			},
			action:     testAction,
			wantErr:    true,
			wantErrStr: "not enough arguments, expected: 2, got: 1",
		},
		{
			testName:    "WantError: Too-Many-Options",
			commandName: "fail-command",
			parameters: Parameters{
				args: []*Argument{
					{
						Name: "a",
						Type: STRING,
					},
					{
						Name: "b",
						Type: STRING,
					},
				},
				opts: []*Option{
					NewOption("--option-1", STRING),
				},
			},
			input: Input{
				args: []string{"Hello", "World"},
				opts: []string{"--input-option-1", "--input-option-2"},
			},
			action:     testAction,
			wantErr:    true,
			wantErrStr: "too many options",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := NewCommand(tt.commandName, tt.parameters.args, tt.parameters.opts, tt.action)
			if cmd == nil {
				t.Fatalf("unexpected issue ocurred, got nil from NewCommand()")
			}
			got, err := cmd.Execute(tt.input.args, tt.input.opts)
			if tt.wantErr {
				// Expected error
				if err == nil {
					t.Errorf("got nil, want %q", tt.wantErrStr)
				}
				// Error message check
				if !strings.Contains(err.Error(), tt.wantErrStr) {
					t.Errorf("got %q, want %q", err, tt.wantErrStr)
				}
			} else {
				if got != tt.wantResult {
					t.Errorf("got: %v, want: %v", got, tt.wantResult)
				}
			}
		})
	}
}
