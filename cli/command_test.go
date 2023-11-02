package cli

import (
	"reflect"
	"strings"
	"testing"
)

func checkError(t *testing.T, testName string, gotErr error, wantErrStr string, wantErr bool) {
	t.Helper()
	if wantErr {
		if gotErr == nil {
			t.Fatalf("expected an error for test case %s but got nil", testName)
		}
		if !strings.Contains(gotErr.Error(), wantErrStr) {
			t.Errorf("got error %q, want %q", gotErr, wantErrStr)
		}
	} else {
		if gotErr != nil {
			t.Fatalf("expected no error for test case %s but got %v", testName, gotErr)
		}
	}
}

func TestNewCommand(t *testing.T) {
	tests := []struct {
		testName     string
		commandName  string
		args         []string
		opts         []string
		wantUsageStr string
		wantErrStr   string
		wantErr      bool
	}{
		{
			testName:     "With 1 arg and 1 opt",
			commandName:  "test-command",
			args:         []string{"Hello"},
			opts:         []string{"--hello"},
			wantUsageStr: "Usage: test-command [arguments] [options]",
			wantErr:      false,
		},
		{
			testName:     "With 2 arg and 2 opt",
			commandName:  "test-command",
			args:         []string{"Hello", "World"},
			opts:         []string{"--hello", "--world"},
			wantUsageStr: "Usage: test-command [arguments] [options]",
			wantErr:      false,
		},
		{
			testName:     "With 1 arg and 0 opt",
			commandName:  "test-command",
			args:         []string{"OneArg"},
			opts:         []string{},
			wantUsageStr: "Usage: test-command [arguments]",
			wantErr:      false,
		},
		{
			testName:     "With 0 arg and 1 opt",
			commandName:  "test-command",
			args:         []string{},
			opts:         []string{"--one-opt"},
			wantUsageStr: "Usage: test-command [options]",
			wantErr:      false,
		},
		{
			testName:    "Error: 0 character argument",
			commandName: "test-command",
			args:        []string{""},
			opts:        []string{},
			wantErrStr:  "error: argument must be at least 1 character",
			wantErr:     true,
		},
		{
			testName:    "Error: 0 character option",
			commandName: "test-command",
			args:        []string{},
			opts:        []string{""},
			wantErrStr:  "error: option must be at least 1 character",
			wantErr:     true,
		},
		{
			testName:    "Error: option must be start with `--`",
			commandName: "test-command",
			args:        []string{},
			opts:        []string{"no-dash-option"},
			wantErrStr:  "error: option must be start with `--`",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd, err := NewCommand(tt.commandName, tt.args, tt.opts, nil)

			checkError(t, tt.testName, err, tt.wantErrStr, tt.wantErr)
			if !tt.wantErr {
				if cmd.Name != tt.commandName {
					t.Errorf("got: name = %v, want %v", cmd.Name, tt.commandName)
				}

				if !reflect.DeepEqual(cmd.Arguments, tt.args) {
					t.Errorf("got: arguments = %v, want %v", cmd.Arguments, tt.args)
				}

				if !reflect.DeepEqual(cmd.Options, tt.opts) {
					t.Errorf("got: options = %v, want %v", cmd.Options, tt.opts)
				}

				if cmd.Usage != tt.wantUsageStr {
					t.Errorf("got: %v, want: %v", cmd.Usage, tt.wantUsageStr)
				}
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
		args []string
		opts []string
	}

	type Input Parameters

	tests := []struct {
		testName    string
		commandName string
		parameters  Parameters
		input       Input
		action      Action
		want        string
		wantErr     bool
	}{
		{
			testName:    "WantSuccess: Expect HelloWorld string",
			commandName: "return-HelloWorld-command",
			parameters:  Parameters{args: []string{"a", "b"}, opts: []string{}},
			input:       Input{args: []string{"Hello", "World"}, opts: []string{}},
			action:      testAction,
			want:        "HelloWorld",
			wantErr:     false,
		},
		{
			testName:    "WantError: Not-Enough-Arguments",
			commandName: "fail-command",
			parameters:  Parameters{args: []string{"a", "b"}, opts: []string{}},
			input:       Input{args: []string{"one-arg"}, opts: []string{}},
			action:      testAction,
			want:        "error: not enough arguments, expected: 2, got: 1",
			wantErr:     true,
		},
		{
			testName:    "WantError: Too-Many-Options",
			commandName: "fail-command",
			parameters:  Parameters{args: []string{"a", "b"}, opts: []string{"--option-1"}},
			input:       Input{args: []string{"Hello", "World"}, opts: []string{"--input-option-1", "--input-option-2"}},
			action:      testAction,
			want:        "error: too many options",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd, _ := NewCommand(tt.commandName, tt.parameters.args, tt.parameters.opts, tt.action)
			if cmd == nil {
				t.Fatalf("unexpected issue ocurred, got nil from NewCommand()")
			}
			got, err := cmd.Execute(tt.input.args, tt.input.opts)
			checkError(t, tt.testName, err, tt.want, tt.wantErr)
			if !tt.wantErr && got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
