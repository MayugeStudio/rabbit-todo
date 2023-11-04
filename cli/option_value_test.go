package cli

import (
	"reflect"
	"testing"
)

func TestOptionValue_Value(t *testing.T) {
	type fields struct {
		StringVal string
		IntVal    int
		BoolVal   bool
		Type      ParameterType
	}
	type testCase struct {
		testName string
		fields   fields
		want     interface{}
	}
	tests := []testCase{
		{
			testName: "Ok-StringValue",
			fields: fields{
				StringVal: "string",
				IntVal:    0,
				BoolVal:   false,
				Type:      STRING,
			},
			want: "string",
		},
		{
			testName: "Ok-IntegerValue",
			fields: fields{
				StringVal: "",
				IntVal:    10,
				BoolVal:   false,
				Type:      INT,
			},
			want: 10,
		},
		{
			testName: "Ok-BooleanValue",
			fields: fields{
				StringVal: "",
				IntVal:    0,
				BoolVal:   true,
				Type:      BOOL,
			},
			want: true,
		},
		{
			testName: "Ok-GetNil",
			fields: fields{
				StringVal: "",
				IntVal:    0,
				BoolVal:   false,
				Type:      -1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ov := &OptionValue{
				StringVal: tt.fields.StringVal,
				IntVal:    tt.fields.IntVal,
				BoolVal:   tt.fields.BoolVal,
				Type:      tt.fields.Type,
			}
			if got := ov.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToOptionValue(t *testing.T) {
	type inputType struct {
		value     string
		paramType ParameterType
	}
	type testCase struct {
		testName   string
		input      inputType
		want       OptionValue
		wantErr    bool
		wantErrStr string
	}
	tests := []testCase{
		{
			testName: "Ok-ConvertString",
			input: inputType{
				value:     "String-Value",
				paramType: STRING,
			},
			want: OptionValue{
				StringVal: "String-Value",
				IntVal:    0,
				BoolVal:   false,
				Type:      STRING,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-ConvertInteger",
			input: inputType{
				value:     "10",
				paramType: INT,
			},
			want: OptionValue{
				StringVal: "",
				IntVal:    10,
				BoolVal:   false,
				Type:      INT,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-ConvertBooleanEmpty",
			input: inputType{
				value:     "",
				paramType: BOOL,
			},
			want: OptionValue{
				StringVal: "",
				IntVal:    0,
				BoolVal:   true,
				Type:      BOOL,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-ConvertBooleanTrue",
			input: inputType{
				value:     "true",
				paramType: BOOL,
			},
			want: OptionValue{
				StringVal: "",
				IntVal:    0,
				BoolVal:   true,
				Type:      BOOL,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Ok-ConvertBooleanFalse",
			input: inputType{
				value:     "false",
				paramType: BOOL,
			},
			want: OptionValue{
				StringVal: "",
				IntVal:    0,
				BoolVal:   false,
				Type:      BOOL,
			},
			wantErr:    false,
			wantErrStr: "",
		},
		{
			testName: "Err-ConvertIntegerFail",
			input: inputType{
				value:     "notInteger",
				paramType: INT,
			},
			want:       OptionValue{},
			wantErr:    true,
			wantErrStr: "cannot convert notInteger to Integer",
		},
		{
			testName: "Err-ConvertBooleanFail",
			input: inputType{
				value:     "notBoolean",
				paramType: BOOL,
			},
			want:       OptionValue{},
			wantErr:    true,
			wantErrStr: "cannot convert notBoolean to Boolean",
		},
		{
			testName: "Err-UnknownParameter",
			input: inputType{
				value:     "unknown",
				paramType: -1,
			},
			want:       OptionValue{},
			wantErr:    true,
			wantErrStr: "unknown parameter type -1",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := convertToOptionValue(tc.input.value, tc.input.paramType)
			if (err != nil) != tc.wantErr {
				t.Errorf("convertToOptionValue() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if tc.wantErr {
				if err.Error() != tc.wantErrStr {
					t.Errorf("convertToOptionValue() error = %q, wantErrStr %q", err, tc.wantErrStr)
				}
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("convertToOptionValue() got = %v, want %v", got, tc.want)
			}
		},
		)
	}
}
