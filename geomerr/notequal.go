package geomerr

import "fmt"

// ErrNotEqual is used to indicate two values that were expected to be equal.
// were not.
type ErrNotEqual struct {
	Expected, Actual interface{}
}

// NotEqual creates an instance of ErrNotEqual.
func NotEqual(expected, actual interface{}) error {
	return ErrNotEqual{
		Expected: expected,
		Actual:   actual,
	}
}

// Error fulfills the error interface.
func (e ErrNotEqual) Error() string {
	return fmt.Sprintf("Expected %v got %v", e.Expected, e.Actual)
}
