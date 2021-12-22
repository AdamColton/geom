// Package geomtest provides helpers for testing the geom packages.
package geomtest

import (
	"fmt"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/stretchr/testify/assert"
)

const (
	// Small is the value that will be passed into AssertEqualizer
	Small cmpr.Tolerance = 1e-10
)

// AssertEqualizer allows a type to define an equality test.
//
// Note that when fulfilling an interface Go will coerce a pointer type to it's
// base type, but not the other way. So if AssertEqual is on the base type
// is on the base type and a pointer to that type is passed into geomtest.Equal
// it will be cast to the base type.
type AssertEqualizer interface {
	AssertEqual(to interface{}, t cmpr.Tolerance) error
}

type GeomAssert struct {
	*assert.Assertions
	assert.TestingT
}

func New(t assert.TestingT) *GeomAssert {
	return &GeomAssert{
		Assertions: assert.New(t),
		TestingT:   t,
	}
}

func (g *GeomAssert) Equal(expected, actual interface{}, msg ...interface{}) bool {
	if _, isAssert := expected.(AssertEqualizer); isAssert {
		return EqualInDelta(g.TestingT, expected, actual, Small, msg...)
	}
	if _, isFloat := expected.(float64); isFloat {
		return EqualInDelta(g.TestingT, expected, actual, Small, msg...)
	}
	return g.Assertions.Equal(expected, actual, msg...)
}

func Message(msg ...interface{}) string {
	ln := len(msg)
	if ln == 0 {
		return ""
	}
	if ln == 1 {
		if s, ok := msg[0].(string); ok {
			return s
		}
		return fmt.Sprint(msg[0])
	}
	if f, ok := msg[0].(string); ok {
		return fmt.Sprintf(f, msg[1:]...)
	}
	return fmt.Sprint(msg...)
}
