package cli

import (
	"reflect"
	"testing"
)

func TestNewOption(t *testing.T) {
	type args struct {
		optName string
		optType ParameterType
		isFlag  bool
	}
	tests := []struct {
		testName string
		args     args
		want     *Option
	}{
		{
			testName: "Test-IntTypeOpt",
			args: args{
				optName: "--int-type-opt",
				optType: INT,
				isFlag:  false,
			},
			want: &Option{
				Name:   "--int-type-opt",
				Type:   INT,
				IsFlag: false,
			},
		},
		{
			testName: "Test-StrTypeOpt",
			args: args{
				optName: "--str-type-opt",
				optType: STRING,
				isFlag:  false,
			},
			want: &Option{
				Name:   "--str-type-opt",
				Type:   STRING,
				IsFlag: false,
			},
		},
		{
			testName: "Test-BoolTypeOpt",
			args: args{
				optName: "--bool-type-opt",
				optType: BOOL,
				isFlag:  false,
			},
			want: &Option{
				Name:   "--bool-type-opt",
				Type:   BOOL,
				IsFlag: false,
			},
		},
		{
			testName: "Test-Fail-Invalid-OptName-ZeroLength",
			args: args{
				optName: "",
				optType: BOOL,
				isFlag:  false,
			},
			want: nil,
		},
		{
			testName: "Test-Fail-Invalid-OptName-NotStartWithDoubleDash",
			args: args{
				optName: "invalid-opt-name",
				optType: BOOL,
				isFlag:  false,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := NewOption(tt.args.optName, tt.args.optType, tt.args.isFlag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
