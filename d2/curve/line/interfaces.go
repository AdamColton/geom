package line

import (
	"github.com/adamcolton/geom/d2"
)

// Intersector finds the points where the interface intersects a line.
// The returned values should be relative to the line passed in. The buffer both
// provides reuse to avoid generating garbage and allows for fine tuning. If
// the buffer has a length of 0, all the intersections will be appended. If the
// buffer has a non-zero length, that output will be limited. So if a buffer
// of length 1 is passed in, only the first intersection is returned. But there
// is no guarenteed order.
type Intersector interface {
	LineIntersections(l Line, buf []float64) []float64
}

// TransformIntersectorWrapper applies a transform to an Intersector.
type TransformIntersectorWrapper struct {
	P d2.Pair
	Intersector
}

// LineIntersections fulfills Intersector. It applies the transform to the
// underlying Intersector.
func (tiw TransformIntersectorWrapper) LineIntersections(l Line, buf []float64) []float64 {
	return tiw.Intersector.LineIntersections(l.T(tiw.P[1]), buf)
}
