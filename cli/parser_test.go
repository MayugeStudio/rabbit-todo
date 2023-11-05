package cli

import (
	"fmt"
	"rabbit-todo/cli/param"
	"reflect"
	"testing"
)

func TestNewParser(t *testing.T) {
	t.Run("NewParser", func(t *testing.T) {
		NewParser()
	})

}

func TestParser_Execute(t *testing.T) {
	type inputType struct {
		commands []Command
		args     []string
	}
	type testCase struct {
		testName   string
		input      inputType
		want       string
		wantErr    bool
		wantErrStr string
	}

	actionGen := func(name string) Action {
		return func(args map[string]param.Value, opts map[string]param.Value) (string, error) {
			msg1 := args["msg-1"].StringVal
			msg2 := args["msg-2"].StringVal
			to := opts["to"].StringVal
			from := opts["from"].StringVal
			str := "from:"
			str += from
			str += " -> "
			str += "\""
			str += fmt.Sprintf("%s %s", msg1, msg2)
			str += "\""
			str += " -> "
			str += "to:"
			str += to
			return str, nil
		}
	}
	testAction1 := actionGen("John")
	testAction2 := actionGen("Mike")

	toOption, _ := param.NewOption("--to", param.STRING)
	fromOption, _ := param.NewOption("--from", param.STRING)
	msgArg1, _ := param.NewArgument("msg-1", param.STRING)
	msgArg2, _ := param.NewArgument("msg-2", param.STRING)

	command1 := NewCommand("command-1", testAction1)
	_ = command1.AddArgument(msgArg1)
	_ = command1.AddArgument(msgArg2)
	_ = command1.AddOption(toOption)
	_ = command1.AddOption(fromOption)
	command2 := NewCommand("command-2", testAction2)
	_ = command2.AddArgument(msgArg1)
	_ = command2.AddArgument(msgArg2)
	_ = command2.AddOption(toOption)
	_ = command2.AddOption(fromOption)

	commands := []Command{command1, command2}

	tests := []testCase{
		{
			testName: "Ok-ExecuteSuccessfully",
			input: inputType{
				commands: commands,
				args:     []string{"command-1", "Hello", "!!", "--to", "John", "--from", "Mike"},
			},
			want:       "from:Mike -> \"Hello !!\" -> to:John",
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Err-NoCommandProvided",
			input: inputType{
				commands: commands,
				args:     []string{},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "no command provided",
		},
		{
			testName: "Error-UnknownCommand",
			input: inputType{
				commands: commands,
				args:     []string{"command-100", "Hello", "!!", "--to", "John", "--from", "Mike"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "unknown command command-100",
		},
		{
			testName: "Error-MissingArguments",
			input: inputType{
				commands: commands,
				args:     []string{"command-1", "Hello", "--to", "John", "--from", "Mike"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "not enough arguments: actual 1, expected 2",
		},
		{
			testName: "Error-MissingOptions'--to'",
			input: inputType{
				commands: commands,
				args:     []string{"command-1", "Hello", "!!", "--to", "--from", "Mike"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "\"--to\" option requires a \"string\" type argument",
		},
		{
			testName: "Error-MissingOptions'--from'",
			input: inputType{
				commands: commands,
				args:     []string{"command-1", "Hello", "!!", "--to", "John", "--from"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "\"--from\" option requires a \"string\" type argument",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			parser := Parser{tc.input.commands}
			got, err := parser.Execute(tc.input.args)
			isErr := err != nil
			if isErr != tc.wantErr {
				t.Fatalf("Parser.Execute() error = %v, wantError %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if err.Error() != tc.wantErrStr {
					t.Errorf("Parser.Execute() error = %q, wantErrStr %q", err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Parser.Execute() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParser_AddCommand(t *testing.T) {
	type inputType struct {
		commands []Command
		command  Command
	}
	type testCase struct {
		testName   string
		input      inputType
		wantErr    bool
		wantErrStr string
	}

	tests := []testCase{
		{
			testName: "Ok-AddCommandSuccessfully",
			input: inputType{
				commands: []Command{
					{Name: "command-1"},
				},
				command: Command{Name: "command-2"},
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Error-DuplicateCommandName",
			input: inputType{
				commands: []Command{
					{
						Name: "command-1",
					},
				},
				command: Command{Name: "command-1"},
			},
			wantErr:    true,
			wantErrStr: "duplicate command name command-1",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			p := &Parser{
				commands: tc.input.commands,
			}
			err := p.AddCommand(tc.input.command)
			if (err != nil) != tc.wantErr {
				t.Errorf("AddCommand() error = %v, wantErr %v", err, tc.wantErr)
			} else if err != nil && err.Error() != tc.wantErrStr {
				t.Errorf("AddCommand() gotErrStr = %v, wantErrStr %v", err, tc.wantErrStr)
			}
		})
	}
}
