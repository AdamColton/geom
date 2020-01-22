package shape

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

// Subtract defines a shape by subtracting the second shape from the first.
type Subtract [2]Shape

// Contains returns true if the point is inside the first shape but not the
// second.
func (s Subtract) Contains(pt d2.Pt) bool {
	return s[0].Contains(pt) && !s[1].Contains(pt)
}

// LineIntersections fulfills line.LineIntersector returning the intesection
// points on the perimeter of the Subtraction of the shapes.
func (s Subtract) LineIntersections(l line.Line) []float64 {
	i0, i1 := s[0].LineIntersections(l), s[1].LineIntersections(l)
	out := make([]float64, 0, len(i0)+len(i1))
	for _, t := range i0 {
		pt := l.Pt1(t)
		if !s[1].Contains(pt) {
			out = append(out, t)
		}
	}
	for _, t := range i1 {
		pt := l.Pt1(t)
		if s[0].Contains(pt) {
			out = append(out, t)
		}
	}
	return out
}

func (s Subtract) BoundingBox() (d2.Pt, d2.Pt) {
	// Bounding box may be tighter, but that's not easy to determine.
	return s[0].BoundingBox()
}
