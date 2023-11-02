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
		testName    string
		commandName string
		args        []string
		opts        []string
		want        string
	}{
		{
			testName:    "With 1 arg and 1 opt",
			commandName: "test-command",
			args:        []string{"Hello"},
			opts:        []string{"--hello"},
			want:        "Usage: test-command [arguments] [options]",
		},
		{
			testName:    "With 2 arg and 2 opt",
			commandName: "test-command",
			args:        []string{"Hello", "World"},
			opts:        []string{"--hello", "--world"},
			want:        "Usage: test-command [arguments] [options]",
		},
		{
			testName:    "With 1 arg and 0 opt",
			commandName: "test-command",
			args:        []string{"OneArg"},
			opts:        []string{},
			want:        "Usage: test-command [arguments]",
		},
		{
			testName:    "With 0 arg and 1 opt",
			commandName: "test-command",
			args:        []string{},
			opts:        []string{"--one-opt"},
			want:        "Usage: test-command [options]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := NewCommand(tt.commandName, tt.args, tt.opts, nil)

			if cmd.Name != tt.commandName {
				t.Errorf("got: name = %v, want %v", cmd.Name, tt.commandName)
			}

			if !reflect.DeepEqual(cmd.Arguments, tt.args) {
				t.Errorf("got: arguments = %v, want %v", cmd.Arguments, tt.args)
			}

			if !reflect.DeepEqual(cmd.Options, tt.opts) {
				t.Errorf("got: options = %v, want %v", cmd.Options, tt.opts)
			}

			if cmd.Usage != tt.want {
				t.Errorf("got: %v, want: %v", cmd.Usage, tt.want)
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
		wantError   bool
	}{
		{
			testName:    "WantSuccess: Expect HelloWorld string",
			commandName: "return-HelloWorld-command",
			parameters:  Parameters{args: []string{"a", "b"}, opts: []string{}},
			input:       Input{args: []string{"Hello", "World"}, opts: []string{}},
			action:      testAction,
			want:        "HelloWorld",
			wantError:   false,
		},
		{
			testName:    "WantError: Not-Enough-Parameters",
			commandName: "fail-command",
			parameters:  Parameters{args: []string{"a", "b"}, opts: []string{}},
			input:       Input{args: []string{"one-arg"}, opts: []string{}},
			action:      testAction,
			want:        "error: not enough arguments, parameters: 2",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := NewCommand(tt.commandName, tt.parameters.args, tt.parameters.opts, tt.action)
			got, err := cmd.Execute(tt.input.args, tt.input.opts)
			checkError(t, tt.testName, err, tt.want, tt.wantError)
			if !tt.wantError && got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
