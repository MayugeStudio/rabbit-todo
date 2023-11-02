package cli

import (
	"testing"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		testName    string
		commandName string
		args        []string
		opts        []string
	}{
		{
			testName:    "With two arguments and two options",
			commandName: "test-command",
			args:        []string{"arg1", "arg2"},
			opts:        []string{"--opt1", "--opt2"},
		},
		{
			testName:    "With no arguments and options",
			commandName: "test-command",
			args:        []string{},
			opts:        []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := NewCommand(tt.commandName, tt.args, tt.opts, nil)

			if cmd.Name != tt.commandName {
				t.Errorf("got: name = %v, want %v", cmd.Name, tt.commandName)
			}
			if len(cmd.Arguments) != len(tt.args) {
				t.Errorf("got: Argument length = %v, want %v", len(cmd.Arguments), len(tt.args))
			}
			for i, arg := range cmd.Arguments {
				if arg != tt.args[i] {
					t.Errorf("got: Argument[%d] = %v, want %v", i, arg, tt.args[i])
				}
			}
			if len(cmd.Options) != len(tt.opts) {
				t.Errorf("got: Options length = %v, want %v", len(cmd.Options), len(tt.opts))
			}
			for i, opt := range cmd.Options {
				if opt != tt.opts[i] {
					t.Errorf("got: Options[%d] = %v, want %v", i, opt, tt.opts[i])
				}
			}
		})
	}
}

func TestCommand_Usage(t *testing.T) {
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
			opts:        []string{"OneOpt"},
			want:        "Usage: test-command [options]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := NewCommand(tt.commandName, tt.args, tt.opts, nil)
			if cmd.Usage != tt.want {
				t.Errorf("got: %v, want: %v", cmd.Usage, tt.want)
			}
		})
	}
}

func TestCommand_Execute(t *testing.T) {
	var testFunction Action
	testFunction = func(args []string, opts []string) (string, error) {
		result := ""

		for _, arg := range args {
			result += arg
		}
		return result, nil
	}

	tests := []struct {
		testName string
		args     []string
		opts     []string
		function Action
		want     string
	}{
		{
			testName: "3+1",
			args:     []string{"3", "1"},
			opts:     []string{},
			function: testFunction,
			want:     "31",
		},
		{
			testName: "Hello+World",
			args:     []string{"Hello", "World"},
			opts:     []string{},
			function: testFunction,
			want:     "HelloWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := &Command{
				Action: tt.function,
			}
			got, err := cmd.Execute(tt.args, tt.opts)
			if err != nil {
				t.Errorf("got: error %v", err)
			}
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
