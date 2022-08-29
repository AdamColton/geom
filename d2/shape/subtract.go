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
func (s Subtract) LineIntersections(l line.Line, buf []float64) []float64 {
	max := len(buf)
	buf = s[0].LineIntersections(l, buf[:0])
	ln := len(buf)
	for i := 0; i < ln; i++ {
		t := buf[i]
		pt := l.Pt1(t)
		if s[1].Contains(pt) {
			ln--
			buf[i] = buf[ln]
			i--
		}
	}
	buf = buf[:ln]
	if max > 0 && len(buf) >= max {
		return buf[:max]
	}

	i1 := s[1].LineIntersections(l, buf[ln:])
	for _, t := range i1 {
		pt := l.Pt1(t)
		if s[0].Contains(pt) {
			buf = append(buf, t)
		}
	}
	if max > 0 && len(buf) >= max {
		return buf[:max]
	}
	return buf
}

// BoundingBox fulfills shape.Shape, it returns a box that contains the shape.
func (s Subtract) BoundingBox() (d2.Pt, d2.Pt) {
	// Bounding box may be tighter, but that's not easy to determine.
	return s[0].BoundingBox()
}

// ConvexHull fulfills shape.ConvexHuller. It returns the convex hull of the
// Intersection. It just returns the convex hull of the shape being subtracted
// from, so the result may not be tight.
func (s Subtract) ConvexHull() []d2.Pt {
	ch := s[0].(ConvexHuller)
	return ch.ConvexHull()
}
