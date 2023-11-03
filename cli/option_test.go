package cli

import (
	"reflect"
	"testing"
)

func TestNewOption(t *testing.T) {
	type args struct {
		optName string
		optType ParameterType
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
			},
			want: nil,
		},
		{
			testName: "Test-Fail-Invalid-OptName-NotStartWithDoubleDash",
			args: args{
				optName: "invalid-opt-name",
				optType: BOOL,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := NewOption(tt.args.optName, tt.args.optType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFlagOption(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		testName string
		args     args
		want     *Option
	}{
		{
			testName: "Test-Success-CreateFlagOption",
			args:     args{name: "--simple-option"},
			want: &Option{
				Name:   "--simple-option",
				Type:   BOOL,
				IsFlag: true,
			},
		},
		{
			testName: "Test-Fail-InvalidOptName-ZeroLength",
			args:     args{name: ""},
			want:     nil,
		},
		{
			testName: "Test-Fail-InvalidOptName-NotStartWithDoubleDash",
			args:     args{name: "invalid-opt-name"},
			want:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := NewFlagOption(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFlagOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
