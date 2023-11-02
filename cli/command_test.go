package cli

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCommand(t *testing.T) {
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
				a, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return "", err
				}
				b, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					return "", err
				}
				return fmt.Sprintf("%d", a+b), nil
			},
			want: "4",
		},
		{
			name: "1+1",
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
