package cli

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewOption(t *testing.T) {
	type args struct {
		optName    string
		optType    ParameterType
		wantErr    bool
		wantErrStr string
	}
	tests := []struct {
		testName string
		args     args
		want     *Option
	}{
		{
			testName: "Test-Ok-IntTypeOpt",
			args: args{
				optName: "--int-type-opt",
				optType: INT,
				wantErr: false,
			},
			want: &Option{
				Name:   "--int-type-opt",
				Type:   INT,
				IsFlag: false,
			},
		},
		{
			testName: "Test-Ok-StrTypeOpt",
			args: args{
				optName: "--str-type-opt",
				optType: STRING,
				wantErr: false,
			},
			want: &Option{
				Name:   "--str-type-opt",
				Type:   STRING,
				IsFlag: false,
			},
		},
		{
			testName: "Test-Ok-BoolTypeOpt",
			args: args{
				optName: "--bool-type-opt",
				optType: BOOL,
				wantErr: false,
			},
			want: &Option{
				Name:   "--bool-type-opt",
				Type:   BOOL,
				IsFlag: false,
			},
		},
		{
			testName: "Test-Fail-InvalidOptName-ZeroLength",
			args: args{
				optName:    "",
				optType:    BOOL,
				wantErr:    true,
				wantErrStr: "name must not be empty",
			},
			want: nil,
		},
		{
			testName: "Test-Fail-InvalidOptName-NotStartWithDoubleDash",
			args: args{
				optName:    "invalid-opt-name",
				optType:    BOOL,
				wantErr:    true,
				wantErrStr: "name must start with '--'",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, err := NewOption(tt.args.optName, tt.args.optType)
			if (err != nil) != tt.args.wantErr {
				t.Fatalf("NewOption() error = %v, wantError %v", err, tt.args.wantErr)
			}
			if tt.args.wantErr {
				if err.Error() != tt.args.wantErrStr {
					t.Errorf("NewOption() error = %v, wantErrStr %v", err, tt.args.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFlagOption(t *testing.T) {
	type args struct {
		optName    string
		wantErr    bool
		wantErrStr string
	}
	tests := []struct {
		testName string
		args     args
		want     *Option
	}{
		{
			testName: "Test-Success-CreateFlagOption",
			args: args{
				optName: "--simple-option",
				wantErr: false,
			},
			want: &Option{
				Name:   "--simple-option",
				Type:   BOOL,
				IsFlag: true,
			},
		},
		{
			testName: "Test-Fail-InvalidOptName-ZeroLength",
			args: args{
				optName:    "",
				wantErr:    true,
				wantErrStr: "name must not be empty",
			},
			want: nil,
		},
		{
			testName: "Test-Fail-InvalidOptName-NotStartWithDoubleDash",
			args: args{
				optName:    "invalid-opt-name",
				wantErr:    true,
				wantErrStr: "name must start with '--'",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got, err := NewFlagOption(tt.args.optName); (err != nil) != tt.args.wantErr ||
				!reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFlagOption() = %v, error = %v, want %v, wantError %v", got, err, tt.want, tt.args.wantErr)
				if tt.args.wantErr && !strings.Contains(err.Error(), tt.args.wantErrStr) {
					t.Errorf("NewFlagOption() error = %v, wantErrStr %v", err, tt.args.wantErrStr)
				}
			}
		})
	}
}
