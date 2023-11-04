package cli

import (
	"reflect"
	"testing"
)

func TestNewArgument(t *testing.T) {
	type inputType struct {
		argName string
		argType ParameterType
	}
	type testCase struct {
		testName   string
		input      inputType
		want       *Argument
		wantErr    bool
		wantErrStr string
	}
	tests := []testCase{
		{
			testName: "Ok-IntTypeArg",
			input: inputType{
				argName: "intTypeArg",
				argType: INT,
			},
			want: &Argument{
				Name: "intTypeArg",
				Type: INT,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-StrTypeArg",
			input: inputType{
				argName: "strTypeArg",
				argType: STRING,
			},
			want: &Argument{
				Name: "strTypeArg",
				Type: STRING,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-BoolTypeArg",
			input: inputType{
				argName: "boolTypeArg",
				argType: BOOL,
			},
			want: &Argument{
				Name: "boolTypeArg",
				Type: BOOL,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Error-InvalidArgName-Empty",
			input: inputType{
				argName: "",
				argType: BOOL,
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "name must not be empty",
		},
		{
			testName: "Error-InvalidArgName-StartWithDoubleDash",
			input: inputType{
				argName: "--invalid-arg-name",
				argType: BOOL,
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "name must not start with '--'",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := NewArgument(tc.input.argName, tc.input.argType)
			isErr := err != nil
			if isErr != tc.wantErr {
				t.Fatalf("NewArgument() error = %v, wantError %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if err.Error() != tc.wantErrStr {
					t.Errorf("NewArgument() error = %q, wantErrStr %q", err, tc.wantErrStr)
				}
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("NewArgument() = %v, want %v", got, tc.want)
			}
		})
	}
}
