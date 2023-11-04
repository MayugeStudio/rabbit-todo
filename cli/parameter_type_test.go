package cli

import "testing"

func TestParameterTypeToString(t *testing.T) {
	type inputType struct {
		paramType ParameterType
	}
	type testCase struct {
		testName string
		input    inputType
		want     string
	}
	tests := []testCase{
		{
			testName: "Ok-ConvertToString",
			input: inputType{
				paramType: STRING,
			},
			want: "string",
		},
		{
			testName: "Ok-ConvertToInteger",
			input: inputType{
				paramType: INT,
			},
			want: "int",
		},
		{
			testName: "Ok-ConvertToBoolean",
			input: inputType{
				paramType: BOOL,
			},
			want: "bool",
		},
		{
			testName: "Ok-UnknownType",
			input: inputType{
				paramType: -1,
			},
			want: "unknownType",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := ParameterTypeToString(tt.input.paramType); got != tt.want {
				t.Errorf("ParameterTypeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
