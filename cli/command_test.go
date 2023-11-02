package cli

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCommand_Execute(t *testing.T) {
	type args struct {
		argument []string
		option   []string
	}
	tests := []struct {
		name     string
		args     args
		function Action
		want     string
	}{
		{
			name: "3+1",
			args: args{
				argument: []string{"3", "1"},
				option:   nil,
			},
			function: func(args []string, _ []string) (string, error) {
				a, err := strconv.Atoi(args[0])
				if err != nil {
					return "", err
				}
				b, err := strconv.Atoi(args[1])
				if err != nil {
					return "", err
				}
				return strconv.Itoa(a + b), nil
			},
			want: "4",
		},
		{
			name: "Hello+World",
			args: args{
				argument: []string{"Hello", "World"},
				option:   nil,
			},
			function: func(args []string, _ []string) (string, error) {
				a := args[0]
				b := args[1]
				return fmt.Sprintf("%s %s", a, b), nil
			},
			want: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &Command{
				Arguments: tt.args.argument,
				Options:   tt.args.option,
				Action:    tt.function,
			}
			got, err := cmd.Execute()
			if err != nil {
				t.Errorf("got error %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCommand(t *testing.T) {
	tests := []struct {
		testName string
		args     []string
		opts     []string
		wantName string
	}{
		{
			testName: "simple command",
			args:     []string{"arg1", "arg2"},
			opts:     []string{"--opt1", "--opt2"},
			wantName: "test-command",
		},
		{
			testName: "simple with no args and opts",
			args:     []string{},
			opts:     []string{},
			wantName: "simple-command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := NewCommand(tt.wantName, tt.args, tt.opts, nil)

			if cmd.Name != tt.wantName {
				t.Errorf("NewCommand() name = %v, want %v", cmd.Name, tt.wantName)
			}
			if len(cmd.Arguments) != len(tt.args) {
				t.Errorf("NewCommand() Argument length = %v, want %v", len(cmd.Arguments), len(tt.args))
			}
			for i, arg := range cmd.Arguments {
				if arg != tt.args[i] {
					t.Errorf("NewCommand() Argument[%d] = %v, want %v", i, arg, tt.args[i])
				}
			}
			if len(cmd.Options) != len(tt.opts) {
				t.Errorf("NewCommand() Options length = %v, want %v", len(cmd.Options), len(tt.opts))
			}
			for i, opt := range cmd.Options {
				if opt != tt.opts[i] {
					t.Errorf("NewCommand() Options[%d] = %v, want %v", i, opt, tt.opts[i])
				}
			}
		})
	}
}
