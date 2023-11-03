package cli

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewArgument(t *testing.T) {
	type args struct {
		argName string
		argType ParameterType
	}
	tests := []struct {
		testName   string
		args       args
		want       *Argument
		wantErr    bool
		wantErrStr string
	}{
		{
			testName: "Test-Ok-IntTypeArg",
			args: args{
				argName: "intTypeArg",
				argType: INT,
			},
			want: &Argument{
				Name: "intTypeArg",
				Type: INT,
			},
			wantErr: false,
		},
		{
			testName: "Test-Ok-StrTypeArg",
			args: args{
				argName: "strTypeArg",
				argType: STRING,
			},
			want: &Argument{
				Name: "strTypeArg",
				Type: STRING,
			},
			wantErr: false,
		},
		{
			testName: "Test-Ok-BoolTypeArg",
			args: args{
				argName: "boolTypeArg",
				argType: BOOL,
			},
			want: &Argument{
				Name: "boolTypeArg",
				Type: BOOL,
			},
			wantErr: false,
		},
		{
			testName: "Test-Fail-InvalidArgName-ZeroLength",
			args: args{
				argName: "",
				argType: BOOL,
			},
			wantErr:    true,
			wantErrStr: "name must not be empty",
		},
		{
			testName: "Test-Fail-InvalidArgName-DoubleDashPrefix",
			args: args{
				argName: "--invalid-arg-name",
				argType: BOOL,
			},
			wantErr:    true,
			wantErrStr: "name must not start with '--'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, err := NewArgument(tt.args.argName, tt.args.argType)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewArgument() = %v, %q, want nil, %q", got, err, tt.wantErrStr)
				}
				if !strings.Contains(err.Error(), tt.wantErrStr) {
					t.Errorf("NewArgument() = %v, %q, want nil, %q", got, err, tt.wantErrStr)
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewArgument() = %v, nil, want %v, nil", got, tt.want)
				}
			}
		})
	}
}
