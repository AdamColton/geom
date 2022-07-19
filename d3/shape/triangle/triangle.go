package triangle

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
)

// Triangle in 3D
type Triangle [3]d3.Pt

const epsilon float64 = 1e-8

// Intersections returns the intersection as a []float64
func (t *Triangle) Intersections(l line.Line) []float64 {
	if f, b := t.LineIntersection(l); b {
		return []float64{f}
	}
	return nil
}

var Small = 1e-6

// Intersector creates a TriangleIntersector that caches some of the computation
// for finding the intersections between the triangle and a line.
func (t *Triangle) Intersector() *TriangleIntersector {
	return &TriangleIntersector{
		T0: t[0],
		V1: t[1].Subtract(t[0]),
		V2: t[2].Subtract(t[0]),
	}
}

// LineIntersection returns the intersection point if there is one and bool
// indicating if there was an intersection.
func (t *Triangle) LineIntersection(l line.Line) (float64, bool) {
	return t.Intersector().LineIntersection(l)
}

// Intersection returns the intersection point if there is one and bool
// indicating if there was an intersection.
func (t *Triangle) Intersection(l line.Line) Intersection {
	return t.Intersector().Intersection(l)
}

// ErrCoincidentVerticies indicates that a Triangle has at least 2 identical
// points.
type ErrCoincidentVerticies struct{}

// Error fulfils the error interface
func (ErrCoincidentVerticies) Error() string {
	return "At least 2 verticies have the same value"
}

// Validate that a triangle has 3 unique points.
func (t *Triangle) Validate() error {
	if t[0] == t[1] || t[0] == t[2] || t[1] == t[2] {
		return ErrCoincidentVerticies{}
	}
	return nil
}

// Normal returns a vector that is perpendicular to the plane of the triangle.
func (t *Triangle) Normal() d3.V {
	return t[1].Subtract(t[0]).Cross(t[2].Subtract(t[0]))
}

// TriangleIntersector expresses a triangle as a point and 2 vectors which is
// the first step of computing line intersections.
type TriangleIntersector struct {
	T0     d3.Pt
	V1, V2 d3.V
}

// LineIntersection checks if the line intersects the triangle. It can only
// intersect once so the return is a float and a bool. The float indicates where
// the intersection occured and the bool indicates if there was an intersection.
func (t *TriangleIntersector) LineIntersection(l line.Line) (float64, bool) {
	i := t.RawIntersection(l)
	return i.T, i.Does
}

// Intersection of a line an a triangle. The U and V values indicate the
// coordinate on the triangle and T indicates the point on the line. Does
// indicates if there was any intersection. If a line intersects a triangle
// and a transformation is applied to both, the Intersection will still give
// the correct values.
type Intersection struct {
	U, V, T float64
	Does    bool
}

// Intersection checks if the given line intersects the triangle.
func (t *TriangleIntersector) Intersection(l line.Line) Intersection {
	var out Intersection
	// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
	h := l.D.Cross(t.V2)
	a := t.V1.Dot(h)
	if a > -epsilon && a < epsilon {
		// line is parallel to triangle
		return out
	}
	f := 1.0 / a
	s := l.T0.Subtract(t.T0)

	// u and v are Barycentric Coordinates
	out.U = f * s.Dot(h)
	if out.U < -Small || out.U > 1+Small {
		// point on plane is outside triangle
		return out
	}
	q := s.Cross(t.V1)
	out.V = f * l.D.Dot(q)
	if out.V < -Small || out.V+out.U > 1+Small {
		// point on plane is outside triangle
		return out
	}
	out.T = f * t.V2.Dot(q)
	out.Does = true
	return out
}

// RawIntersection saves overhead by doing all the calculation in place rather
// than calling functions. For raytracing it was shown to be a significant
// improvement.
func (t *TriangleIntersector) RawIntersection(l line.Line) Intersection {
	var out Intersection
	// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
	h := d3.V{
		l.D.Y*t.V2.Z - l.D.Z*t.V2.Y,
		l.D.Z*t.V2.X - l.D.X*t.V2.Z,
		l.D.X*t.V2.Y - l.D.Y*t.V2.X,
	}
	a := t.V1.X*h.X + t.V1.Y*h.Y + t.V1.Z*h.Z
	if a > -epsilon && a < epsilon {
		// line is parallel to triangle
		return out
	}
	f := 1.0 / a
	s := d3.V{
		l.T0.X - t.T0.X,
		l.T0.Y - t.T0.Y,
		l.T0.Z - t.T0.Z,
	}

	// u and v are Barycentric Coordinates
	out.U = f * (s.X*h.X + s.Y*h.Y + s.Z*h.Z)
	if out.U < -Small || out.U > 1+Small {
		// point on plane is outside triangle
		return out
	}
	q := d3.V{
		s.Y*t.V1.Z - s.Z*t.V1.Y,
		s.Z*t.V1.X - s.X*t.V1.Z,
		s.X*t.V1.Y - s.Y*t.V1.X,
	}
	out.V = f * (l.D.X*q.X + l.D.Y*q.Y + l.D.Z*q.Z)
	if out.V < -Small || out.V+out.U > 1+Small {
		// point on plane is outside triangle
		return out
	}
	out.T = f * (t.V2.X*q.X + t.V2.Y*q.Y + t.V2.Z*q.Z)
	out.Does = true
	return out
}
