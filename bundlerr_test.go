package bundlerr

import (
	"errors"
	"reflect"
	"testing"
)

func TestErrorsBundleEvaluate(t *testing.T) {
	tests := []struct {
		name          string
		bundle        *Bundle
		appends       []error
		wantErr       bool
		wantErrString string
	}{
		{
			name:    "nil",
			bundle:  New(),
			wantErr: false,
		},
		{
			name:          "single error",
			bundle:        New(errors.New("error1")),
			wantErr:       true,
			wantErrString: "error1",
		},
		{
			name:          "multiple errors",
			bundle:        New(errors.New("error1"), errors.New("error2")),
			wantErr:       true,
			wantErrString: "error1 █ error2",
		},
		{
			name:    "empty with nil append",
			bundle:  New(),
			appends: []error{nil},
			wantErr: false,
		},
		{
			name:          "empty with errorous append",
			bundle:        New(),
			appends:       []error{errors.New("error1")},
			wantErr:       true,
			wantErrString: "error1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.appends {
				tt.bundle.Append(e)
			}

			err := tt.bundle.Evaluate()
			if (err != nil) != tt.wantErr {
				t.Errorf("errorsBundle.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got := tt.bundle.Error(); !reflect.DeepEqual(got, tt.wantErrString) {
				t.Errorf("NewErrorsBundle() = %v, want %v", got, tt.wantErrString)
			}

		})
	}
}
