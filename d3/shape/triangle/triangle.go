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
	if f, b := t.Intersection(l); b {
		return []float64{f}
	}
	return nil
}

var Small = 1e-6

// Intersection returns the intersection point if there is one and bool
// indicating if there was an intersection.
func (t *Triangle) Intersection(l line.Line) (float64, bool) {
	// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
	v1 := t[1].Subtract(t[0])
	v2 := t[2].Subtract(t[0])
	h := l.D.Cross(v2)
	a := v1.Dot(h)
	if a > -epsilon && a < epsilon {
		// line is parallel to triangle
		return 0, false
	}
	f := 1.0 / a
	s := l.T0.Subtract(t[0])

	// u and v are Barycentric Coordinates
	u := f * s.Dot(h)
	if u < -Small || u > 1+Small {
		// point on plane is outside triangle
		return 0, false
	}
	q := s.Cross(v1)
	v := f * l.D.Dot(q)
	if v < -Small || u+v > 1+Small {
		// point on plane is outside triangle
		return 0, false
	}
	return f * v2.Dot(q), true
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
