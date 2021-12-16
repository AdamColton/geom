// Package geomtest provides helpers for testing the geom packages.
package geomtest

import "github.com/adamcolton/geom/calc/cmpr"

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
