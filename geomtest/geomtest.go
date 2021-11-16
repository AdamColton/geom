// Package geomtest provides helpers for testing the geom packages.
package geomtest

import "github.com/adamcolton/geom/calc/cmpr"

const (
	// Small is the value that will be passed into AssertEqualizer
	Small cmpr.Tolerance = 1e-10
)

// AssertEqualizer allows a type to define an equality test.
type AssertEqualizer interface {
	AssertEqual(to interface{}, t cmpr.Tolerance) error
}
