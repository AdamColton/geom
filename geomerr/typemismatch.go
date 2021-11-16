package geomerr

import (
	"fmt"
	"reflect"
)

// ErrTypeMismatch indicates that two types that were expected to be equal were
// not.
type ErrTypeMismatch struct {
	Expected, Actual reflect.Type
}

// TypeMismatch creates an instance of ErrTypeMismatch.
func TypeMismatch(expected, actual interface{}) error {
	return ErrTypeMismatch{
		Expected: reflect.TypeOf(expected),
		Actual:   reflect.TypeOf(actual),
	}
}

// Error fulfills the error interface.
func (e ErrTypeMismatch) Error() string {
	return fmt.Sprintf(`Types do not match: expected "%v", got "%v"`, e.Expected, e.Actual)
}
