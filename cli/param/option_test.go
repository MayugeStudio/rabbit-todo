package param

import (
	"reflect"
	"testing"
)

func TestNewOption(t *testing.T) {
	type inputType struct {
		optName string
		optType Type
	}
	type testCase struct {
		testName   string
		input      inputType
		want       *Option
		wantErr    bool
		wantErrStr string
	}
	tests := []testCase{
		{
			testName: "Ok-IntTypeOpt",
			input: inputType{
				optName: "--int-type-opt",
				optType: INT,
			},
			want: &Option{
				Name:   "--int-type-opt",
				Type:   INT,
				IsFlag: false,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-StrTypeOpt",
			input: inputType{
				optName: "--str-type-opt",
				optType: STRING,
			},
			want: &Option{
				Name:   "--str-type-opt",
				Type:   STRING,
				IsFlag: false,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-BoolTypeOpt",
			input: inputType{
				optName: "--bool-type-opt",
				optType: BOOL,
			},
			want: &Option{
				Name:   "--bool-type-opt",
				Type:   BOOL,
				IsFlag: false,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Error-InvalidOptName-Empty",
			input: inputType{
				optName: "",
				optType: STRING,
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "name must not be empty",
		},
		{
			testName: "Error-InvalidOptName-NotStartsWithDoubleDash",
			input: inputType{
				optName: "invalid-opt-name",
				optType: INT,
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "name must start with '--'",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := NewOption(tc.input.optName, tc.input.optType)
			isErr := err != nil
			if isErr != tc.wantErr {
				t.Fatalf("NewOption() error = %v, wantError %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if err.Error() != tc.wantErrStr {
					t.Errorf("NewOption() error = %v, wantErrStr %v", err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("NewOption() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNewFlagOption(t *testing.T) {
	type testCase struct {
		testName   string
		input      string
		want       *Option
		wantErr    bool
		wantErrStr string
	}
	tests := []testCase{
		{
			testName: "Ok-CreateFlagOption",
			input:    "--simple-option",
			want: &Option{
				Name:   "--simple-option",
				Type:   BOOL,
				IsFlag: true,
			},
		},
		{
			testName:   "Error-InvalidOptName-Empty",
			input:      "",
			want:       nil,
			wantErr:    true,
			wantErrStr: "name must not be empty",
		},
		{
			testName:   "Error-InvalidOptName-NotStartWithDoubleDash",
			input:      "invalid-opt-name",
			want:       nil,
			wantErr:    true,
			wantErrStr: "name must start with '--'",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := NewFlagOption(tc.input)
			isErr := err != nil
			if isErr != tc.wantErr {
				t.Fatalf("NewFlagOption() error = %v, wantError %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if err.Error() != tc.wantErrStr {
					t.Errorf("NewFlagOption() error = %q, wantErrStr %q", err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("NewFlagOption() = %v, want %v", got, tc.want)
			}
		})
	}
}
