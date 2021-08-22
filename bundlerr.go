// Package bundlerr let's you bundle multiple errors into one
package bundlerr

import (
	"errors"
)

type Bundle struct {
	errors    []error
	formatter Formatter
}

func New(errs ...error) *Bundle {
	return NewWithFormatter(defaultFormatFn, errs...)
}

func NewWithFormatter(formatter Formatter, errs ...error) *Bundle {
	return &Bundle{
		errors:    errs,
		formatter: formatter,
	}
}

// Append appends new error to the bundle. If the new error is nil, it is ignore
// if the new error is another bundle, their errors are merged into this bundle creating flat structure
func (b *Bundle) Append(e error) {
	if e == nil {
		return
	}
	if b == e {
		return
	}

	if bundle, ok := e.(*Bundle); ok {
		b.errors = append(b.errors, bundle.errors...)
		return
	} else if bundle, ok := e.(Bundle); ok {
		b.errors = append(b.errors, bundle.errors...)
		return
	}

	b.errors = append(b.errors, e)
}

// Evaluate evaluates the bundle to nil if no errors are bundled or returns the bundle
func (b *Bundle) Evaluate() error {
	if b == nil || len(b.errors) == 0 {
		return nil
	}

	return b
}

// Errors returns all errors in this bundle
func (b *Bundle) Errors() []error {
	if b == nil {
		return []error{}
	}

	return b.errors
}

func (b Bundle) Error() string {
	return b.formatter(b)
}

// Is returns true if any of the bundled errors is the target error - same as errors.Is
func (b *Bundle) Is(target error) bool {
	if b == nil {
		return false
	}

	for _, err := range b.errors {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

// As returns true if any of the bundled errors are the target error - same as errors.As
func (b *Bundle) As(target interface{}) bool {
	if b == nil {
		return false
	}

	for _, err := range b.errors {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}
