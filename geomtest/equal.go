package geomtest

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
)

// Equal calls AssertEqual with the default value of Small. If there is an error
// it is passed into t.Error. The return bool will be true if the values were
// equal.
func Equal(t *testing.T, expected, actual interface{}) bool {
	return EqualInDelta(t, expected, actual, Small)
}

// EqualInDelta calls AssertEqual. If there is an error it is passed into
// t.Error. The return bool will be true if the values were equal.
func EqualInDelta(t *testing.T, expected, actual interface{}, delta cmpr.Tolerance) bool {
	err := AssertEqual(expected, actual, delta)
	if err == nil {
		return true
	}
	t.Error(err)
	return false
}

// AssertEqual can compare anything that implements geomtest.AssertEqualizer.
// There is also logic to handle comparing float64 values Any two slices whose
// elements can be compared with Equal can be compared. The provided delta value
// will be passed to anything that implements AssertEqualizer. If the equality
// check fails, an error is returned.
func AssertEqual(expected, actual interface{}, delta cmpr.Tolerance) error {
	ev := reflect.ValueOf(expected)
	if ev.Kind() == reflect.Slice {
		av := reflect.ValueOf(actual)
		if av.Kind() != reflect.Slice || av.Type().Elem() != ev.Type().Elem() {
			return geomerr.TypeMismatch(expected, actual)
		}
		ln := ev.Len()
		if aln := av.Len(); ln != aln {
			return geomerr.LenMismatch(ln, aln)
		}
		var errs geomerr.SliceErrs
		for i := 0; i < ln; i++ {
			err := AssertEqual(ev.Index(i).Interface(), av.Index(i).Interface(), delta)
			if err != nil {
				errs = append(errs, geomerr.SliceErrRecord{
					Err:   err,
					Index: i,
				})
			}
		}
		if len(errs) == 0 {
			return nil
		}
		return errs
	}

	if eq, ok := expected.(AssertEqualizer); ok {
		return eq.AssertEqual(actual, delta)
	} else if f0, ok := expected.(float64); ok {
		if f1, ok := actual.(float64); ok {
			if delta.Equal(f0, f1) {
				return nil
			}
			return geomerr.NotEqual(f0, f1)
		}
	}

	return fmt.Errorf("unsupported_type")
}
