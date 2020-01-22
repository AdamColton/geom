package shape

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

// Intersection takes two Shapes and creates a new Shape that is their
// intersection.
type Intersection [2]Shape

// Contains returns true if the point is inside both shapes.
func (a Intersection) Contains(pt d2.Pt) bool {
	return a[0].Contains(pt) && a[1].Contains(pt)
}

// LineIntersections fulfills line.LineIntersector returning the intesection
// points on the perimeter of the Intersection of the shapes.
func (a Intersection) LineIntersections(l line.Line) []float64 {
	i0, i1 := a[0].LineIntersections(l), a[1].LineIntersections(l)
	out := make([]float64, 0, len(i0)+len(i1))
	for _, t := range i0 {
		pt := l.Pt1(t)
		if a[1].Contains(pt) {
			out = append(out, t)
		}
	}
	for _, t := range i1 {
		pt := l.Pt1(t)
		if a[0].Contains(pt) {
			out = append(out, t)
		}
	}
	return out
}

// BoundingBox of the Intersection
func (a Intersection) BoundingBox() (d2.Pt, d2.Pt) {
	m, M := a[0].BoundingBox()
	m1, M1 := a[1].BoundingBox()
	return d2.Max(m, m1), d2.Min(M, M1)
}
