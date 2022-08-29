package shape

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/polygon"
)

// Intersection takes two Shapes and creates a new Shape that is their
// intersection.
type Intersection [2]Shape

// Contains returns true if the point is inside both shapes.
func (a Intersection) Contains(pt d2.Pt) bool {
	return a[0].Contains(pt) && a[1].Contains(pt)
}

// LineIntersections fulfills line.Intersector returning the intesection
// points on the perimeter of the Intersection of the shapes.
func (a Intersection) LineIntersections(l line.Line, buf []float64) []float64 {
	max := len(buf)
	buf = a[0].LineIntersections(l, buf[:0])
	ln := len(buf)
	for i := 0; i < ln; i++ {
		t := buf[i]
		pt := l.Pt1(t)
		if !a[1].Contains(pt) {
			ln--
			buf[i] = buf[ln]
			i--
		}
	}
	buf = buf[:ln]
	if max > 0 && len(buf) >= max {
		return buf[:max]
	}

	i1 := a[1].LineIntersections(l, buf[ln:])
	for _, t := range i1 {
		pt := l.Pt1(t)
		if a[0].Contains(pt) {
			buf = append(buf, t)
		}
	}
	if max > 0 && len(buf) > max {
		return buf[:max]
	}
	return buf
}

//BoundingBox fulfills shape.Shape, it returns a box that contains the shape.
func (a Intersection) BoundingBox() (d2.Pt, d2.Pt) {
	m, M := a[0].BoundingBox()
	m1, M1 := a[1].BoundingBox()
	return d2.Max(m, m1), d2.Min(M, M1)
}

// ConvexHull fulfills shape.ConvexHuller. It returns the convex hull of the
// Intersection. If the underlying Convex hulls are tight, the result will also
// be tight.
func (a Intersection) ConvexHull() []d2.Pt {
	a0h, a1h := a[0].(ConvexHuller), a[1].(ConvexHuller)
	s0h, s1h := polygon.Polygon(a0h.ConvexHull()), polygon.Polygon(a1h.ConvexHull())
	out := s0h.PolygonCollisions(s1h).P(s0h)
	out = PointsInContainer(s0h, a[1], true, out)
	out = PointsInContainer(s1h, a[0], true, out)
	return polygon.ConvexHull(out...)
}
