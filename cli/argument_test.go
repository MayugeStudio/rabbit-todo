package cli

import (
	"reflect"
	"testing"
)

func TestNewArgument(t *testing.T) {
	type args struct {
		argName string
		argType ParameterType
	}
	tests := []struct {
		testName string
		args     args
		want     *Argument
	}{
		{
			testName: "Test-IntTypeArg",
			args: args{
				argName: "intTypeArg",
				argType: INT,
			},
			want: &Argument{
				Name: "intTypeArg",
				Type: INT,
			},
		},
		{
			testName: "Test-StrTypeArg",
			args: args{
				argName: "strTypeArg",
				argType: STRING,
			},
			want: &Argument{
				Name: "strTypeArg",
				Type: STRING,
			},
		},
		{
			testName: "Test-BoolTypeArg",
			args: args{
				argName: "boolTypeArg",
				argType: BOOL,
			},
			want: &Argument{
				Name: "boolTypeArg",
				Type: BOOL,
			},
		},
		{
			testName: "Test-Fail-InvalidArgName-ZeroLength",
			args: args{
				argName: "",
				argType: BOOL,
			},
			want: nil,
		},
		{
			testName: "Test-Fail-InvalidArgName-DoubleDashPrefix",
			args: args{
				argName: "--invalid-arg-name",
				argType: BOOL,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := NewArgument(tt.args.argName, tt.args.argType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}
