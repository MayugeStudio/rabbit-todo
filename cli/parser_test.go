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
		return func(args []string, opts map[string]param.Value) (string, error) {
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

	toOption, _ := param.NewOption("--to", param.STRING)
	fromOption, _ := param.NewOption("--from", param.STRING)
	msgArg1, _ := param.NewArgument("msg-1", param.STRING)
	msgArg2, _ := param.NewArgument("msg-2", param.STRING)

	command1 := NewCommand("command-1", []*param.Argument{msgArg1, msgArg2}, []*param.Option{toOption, fromOption}, testAction1)
	command2 := NewCommand("command-2", []*param.Argument{msgArg1, msgArg2}, []*param.Option{toOption, fromOption}, testAction2)

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
