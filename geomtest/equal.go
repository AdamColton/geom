package geomtest

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
	"github.com/stretchr/testify/assert"
)

type tHelper interface {
	Helper()
}

// Equal calls AssertEqual with the default value of Small. If there is an error
// it is passed into t.Error. The return bool will be true if the values were
// equal.
func Equal(t assert.TestingT, expected, actual interface{}, msg ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return EqualInDelta(t, expected, actual, Small, msg...)
}

var equalType = reflect.TypeOf((*AssertEqualizer)(nil)).Elem()

// EqualInDelta calls AssertEqual. If there is an error it is passed into
// t.Error. The return bool will be true if the values were equal.
func EqualInDelta(t assert.TestingT, expected, actual interface{}, delta cmpr.Tolerance, msg ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	err := AssertEqual(expected, actual, delta)
	if err == nil {
		return true
	}
	if len(msg) > 0 {
		t.Errorf("%s: %s", err.Error(), Message(msg...))
	} else {
		t.Errorf("%s", err.Error())
	}
	return false
}

// AssertEqual can compare anything that implements geomtest.AssertEqualizer.
// There is also logic to handle comparing float64 values Any two slices whose
// elements can be compared with Equal can be compared. The provided delta value
// will be passed to anything that implements AssertEqualizer. If the equality
// check fails, an error is returned.
func AssertEqual(expected, actual interface{}, delta cmpr.Tolerance) error {
	ev := reflect.ValueOf(expected)
	if eq, ok := expected.(AssertEqualizer); ok {
		return eq.AssertEqual(actual, delta)
	}

	if ev.Kind() == reflect.Slice {
		av := reflect.ValueOf(actual)
		if av.Kind() != reflect.Slice {
			return geomerr.TypeMismatch(expected, actual)
		}
		return geomerr.NewSliceErrs(ev.Len(), av.Len(), func(i int) error {
			return AssertEqual(ev.Index(i).Interface(), av.Index(i).Interface(), delta)
		})
	}

	if ef, ok := expected.(float64); ok {
		if err := geomerr.NewTypeMismatch(expected, actual); err != nil {
			return err
		}
		af := actual.(float64)
		return geomerr.NewNotEqual(delta.Equal(ef, af), ef, af)
	}

	format := "unsupported_type: %s"
	t := ev.Type()
	if t.Kind() != reflect.Ptr {
		if p := reflect.PtrTo(t); p.Implements(equalType) {
			format = fmt.Sprintf("%s (%s fulfills AssertEqualizer)", format, p.String())
		}
	}

	return fmt.Errorf(format, t.String())
}
