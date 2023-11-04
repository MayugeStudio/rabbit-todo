package cli

import (
	"rabbit-todo/cli/param"
	"reflect"
	"strings"
	"testing"
)

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
		return func(args []string, opts map[string]parameter.ParamValue) (string, error) {
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
	}
	testAction1 := actionGen("John")
	testAction2 := actionGen("Mike")

	toOption, _ := parameter.NewOption("--to", parameter.STRING)
	fromOption, _ := parameter.NewOption("--from", parameter.STRING)
	msgArg1, _ := parameter.NewArgument("msg-1", parameter.STRING)
	msgArg2, _ := parameter.NewArgument("msg-2", parameter.STRING)

	command1 := NewCommand("command-1", []*parameter.Argument{msgArg1, msgArg2}, []*parameter.Option{toOption, fromOption}, testAction1)
	command2 := NewCommand("command-2", []*parameter.Argument{msgArg1, msgArg2}, []*parameter.Option{toOption, fromOption}, testAction2)

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
			wantErrStr: "\"--to\" option require one \"string\" type argument",
		},
		{
			testName: "Error-MissingOptions'--from'",
			input: inputType{
				commands: commands,
				args:     []string{"command-1", "Hello", "!!", "--to", "John", "--from"},
			},
			want:       "",
			wantErr:    true,
			wantErrStr: "\"--from\" option require one \"string\" type argument",
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
