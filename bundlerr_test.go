package bundlerr

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestEvaluate(t *testing.T) {
	var b = New(errors.New("error1"))

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
			name:          "empty with errornous append",
			bundle:        New(),
			appends:       []error{errors.New("error1")},
			wantErr:       true,
			wantErrString: "error1",
		},
		{
			name:          "empty with bundle append",
			bundle:        New(),
			appends:       []error{New(errors.New("foobar"))},
			wantErr:       true,
			wantErrString: "foobar",
		},
		{
			name:          "append itself",
			bundle:        b,
			appends:       []error{b},
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

type ErrCustom struct{ i int }

func (e ErrCustom) Error() string {
	return fmt.Sprintf("custom error: %d", e.i)
}

func TestIs(t *testing.T) {
	var (
		ErrFooBar = errors.New("foobar")
		ErrLipsum = errors.New("lipsum")
	)

	tests := []struct {
		name   string
		bundle *Bundle
		isErr  error
		is     bool
	}{
		{
			name:   "nil",
			bundle: New(),
			isErr:  ErrFooBar,
			is:     false,
		},
		{
			name:   "simple",
			bundle: New(ErrFooBar),
			isErr:  ErrFooBar,
			is:     true,
		},
		{
			name:   "multiple",
			bundle: New(ErrFooBar, ErrLipsum),
			isErr:  ErrFooBar,
			is:     true,
		},
		{
			name:   "custom type",
			bundle: New(ErrCustom{1}),
			isErr:  ErrCustom{1},
			is:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := errors.Is(tt.bundle, tt.isErr)
			if res != tt.is {
				t.Errorf("errorsBundle.Is() error = `%v`, wantErr = `%v`, want `%v`", tt.bundle, tt.isErr, tt.is)
			}
		})
	}
}

func TestAs(t *testing.T) {
	var (
		ErrFooBar = errors.New("foobar")
		// ErrLipsum = errors.New("lipsum")
	)

	tests := []struct {
		name   string
		bundle *Bundle
		asErr  error
		as     bool
	}{
		{
			name:   "simple",
			bundle: New(ErrFooBar),
			asErr:  ErrFooBar,
			as:     true,
		},
		// NOTE: This will work because the error type is the same. It doesn't check for value
		// {
		//     name:   "simple different",
		//     bundle: New(ErrFooBar),
		//     isErr:  ErrLipsum,
		//     is:     false,
		// },
		{
			name:   "custom type - same value",
			bundle: New(ErrCustom{1}),
			asErr:  ErrCustom{1},
			as:     true,
		},
		{
			name:   "custom type - different value",
			bundle: New(ErrCustom{1}),
			asErr:  ErrCustom{},
			as:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.bundle.As(&tt.asErr)
			if res != tt.as {
				t.Errorf("errorsBundle.As() error = `%v`, wantErr = `%v`, want `%v`", tt.bundle, tt.asErr, tt.as)
			}
		})
	}
}
