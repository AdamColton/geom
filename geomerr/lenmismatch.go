package geomerr

import (
	"fmt"
)

// ErrLenMismatch represents a mis-matched length.
type ErrLenMismatch struct {
	Expected, Actual int
}

// LenMismatch returns an ErrLenMismatch.
func LenMismatch(expected, actual int) error {
	return ErrLenMismatch{
		Expected: expected,
		Actual:   actual,
	}
}

// Error fulfills the error interface.
func (e ErrLenMismatch) Error() string {
	return fmt.Sprintf("Lengths do not match: Expected %v got %v", e.Expected, e.Actual)
}
