package geomerr

import (
	"fmt"
	"strings"
)

// SliceErrRecord represents an error comparing two slices.
type SliceErrRecord struct {
	// Index of the reference slice
	Index int
	Err   error
}

// SliceErrs collects SliceErrRecord and treats them as a single error.
type SliceErrs []SliceErrRecord

// MaxSliceErrs limits the number of errors that will be reported by
// SliceErrs.Error
var MaxSliceErrs = 10

// Error fulfills the error interface.
func (e SliceErrs) Error() string {
	var out []string
	ln := len(e)
	if ln > MaxSliceErrs {
		out = make([]string, MaxSliceErrs+1)
		out[MaxSliceErrs] = fmt.Sprintf("Omitting %d more", ln-MaxSliceErrs)
		ln = MaxSliceErrs
	} else {
		out = make([]string, ln)
	}

	for i := 0; i < ln; i++ {
		r := e[i]
		out[i] = fmt.Sprintf("\t%d: %s", r.Index, r.Err.Error())
	}

	return strings.Join(out, "\n")
}
