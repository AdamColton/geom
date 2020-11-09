// Package geomtest provides helpers for testing the geom packages.
package geomtest

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	Small = 1e-10
	Big   = 1.0 / Small
)

func IsSmall(f float64) bool {
	return f > -Small && f < Small
}

type AssertEqualizer interface {
	AssertEqual(to interface{}) error
}

type ErrWrapper struct {
	Code             int
	expected, actual interface{}
}

const (
	NotEqualCode int = iota + 1
	TypeMismatchCode
	LenMismatchCode
)

func (e ErrWrapper) Error() string {
	switch e.Code {
	case TypeMismatchCode:
		et := reflect.TypeOf(e.expected)
		at := reflect.TypeOf(e.actual)
		return fmt.Sprintf(`Types do not match: expected "%s", got "%s"`, et, at)
	case NotEqualCode:
		return fmt.Sprintf("Expected %s got %s", e.expected, e.actual)
	case LenMismatchCode:
		eLn := reflect.ValueOf(e.expected).Len()
		aLn := reflect.ValueOf(e.actual).Len()
		return fmt.Sprintf("Lengths do not match: Expected %d got %d", eLn, aLn)
	}
	return "Unsupported Error Code"
}

func NotEqual(expected, actual interface{}) error {
	return ErrWrapper{
		expected: expected,
		actual:   actual,
		Code:     NotEqualCode,
	}
}

func TypeMismatch(expected, actual interface{}) error {
	return ErrWrapper{
		expected: expected,
		actual:   actual,
		Code:     TypeMismatchCode,
	}
}

func LenMismatch(expected, actual interface{}) error {
	return ErrWrapper{
		expected: expected,
		actual:   actual,
		Code:     LenMismatchCode,
	}
}

type SliceErrRecord struct {
	Index int
	Err   error
}

type SliceErrs []SliceErrRecord

var MaxSliceErrs = 10

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
