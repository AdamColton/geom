package triangle

import (
	"math"
	"strings"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/geomerr"
)

// Triangle is a 2D triangle
type Triangle [3]d2.Pt

// Contains checks if a point is inside the triangle.
func (t *Triangle) Contains(pt d2.Pt) bool {
	// If a point is inside the triangle, the sign of the cross product from the
	// point to each vertex will be the same. But a cross product of exactly 0
	// doesn't have a sign.
	s1 := t[0].Subtract(t[1])
	r1 := pt.Subtract(t[1])
	s2 := t[1].Subtract(t[2])
	r2 := pt.Subtract(t[2])

	c1 := s1.Cross(r1)
	c2 := s2.Cross(r2)
	if !(c1 >= 0 && c2 >= 0) && !(c1 <= 0 && c2 <= 0) {
		return false
	}

	s3 := t[2].Subtract(t[0])
	r3 := pt.Subtract(t[0])
	c3 := s3.Cross(r3)
	return (c2 >= 0 && c3 >= 0) || (c2 <= 0 && c3 <= 0)
}

// SignedArea of the triangle
func (t *Triangle) SignedArea() float64 {
	v1 := t[0].Subtract(t[1])
	v2 := t[0].Subtract(t[2])
	return 0.5 * v1.Cross(v2)
}

// Area of the triangle
func (t *Triangle) Area() float64 {
	return math.Abs(t.SignedArea())
}

// Perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t[0].Distance(t[1]) + t[1].Distance(t[2]) + t[2].Distance(t[0])
}

// Centroid returns the center of mass of the triangle
func (t *Triangle) Centroid() d2.Pt {
	return d2.Pt{
		X: (t[0].X + t[1].X + t[2].X) / 3.0,
		Y: (t[0].Y + t[1].Y + t[2].Y) / 3.0,
	}
}

// String fulfills Stringer
func (t *Triangle) String() string {
	return strings.Join([]string{
		"Triangle[ ",
		t[0].String(), ",",
		t[1].String(), ",",
		t[2].String(), " ]",
	}, "")
}

// Pt1c0 converts the triangle to line.Segments
func (t *Triangle) Pt1c0() d2.Pt1 {
	return line.Segments(append(t[:], t[0]))
}

// Pt1 returns a point along the perimeter is t0 is between 0 and 1.
func (t *Triangle) Pt1(t0 float64) d2.Pt {
	return t.Pt1c0().Pt1(t0)
}

// Pt2c1 finds the fill line for t0.
func (t *Triangle) Pt2c1(t0 float64) d2.Pt1 {
	m := line.New(t[0], t[1]).Pt1(0.5)
	p0 := line.New(t[0], m).Pt1(t0)
	p1 := line.New(t[2], t[1]).Pt1(t0)
	return line.New(p0, p1)
}

// Pt2 finds a point inside the triange if t0 and t1 are both between 0 and 1
// inclusive. Conforms to shape filling rules.
func (t *Triangle) Pt2(t0, t1 float64) d2.Pt {
	return t.Pt2c1(t0).Pt1(t1)
}

// L fulfills Limiter and describes the parametric methods on the triangle.
func (t *Triangle) L(ts, c int) d2.Limit {
	if (ts == 1 && c == 1) ||
		(ts == 2 && c == 2) ||
		(ts == 1 && c == 0) ||
		(ts == 2 && c == 1) {
		return d2.LimitBounded
	}
	return d2.LimitUndefined
}

// LineIntersections find the intersections of the given line with the triangle
// relative to the line
func (t *Triangle) LineIntersections(l line.Line, buf []float64) []float64 {
	max := len(buf)
	buf = buf[:0]
	prev := t[2]
	for _, cur := range t {
		l2 := line.New(prev, cur)
		prev = cur
		t0, tt, ok := l.Intersection(l2)
		if tt >= 0 && tt < 1 && ok {
			buf = append(buf, t0)
			if max > 0 && len(buf) == max {
				return buf
			}
		}
	}
	return buf
}

// BoundingBox returns a bounding box that contains the triangle
func (t *Triangle) BoundingBox() (d2.Pt, d2.Pt) {
	return d2.MinMax(t[:]...)
}

// CircumCenter find the point where the bisectors of the sides intersect.
func (t *Triangle) CircumCenter() d2.Pt {
	l01 := line.Bisect(t[0], t[1])
	l02 := line.Bisect(t[0], t[2])
	_, t0, ok := l01.Intersection(l02)
	if ok {
		return l02.Pt1(t0)
	}
	if t[0] == t[1] || t[1] == t[2] {
		return line.New(t[0], t[2]).Pt1(0.5)
	}
	return line.New(t[0], t[1]).Pt1(0.5)
}

// ConvexHull fulfills shape.ConvexHuller. It returns the triangle as a slice.
func (t *Triangle) ConvexHull() []d2.Pt {
	return t[:]
}

// T applies the transform and returns a new Triangle.
func (t *Triangle) T(transform *d2.T) *Triangle {
	return &Triangle{
		transform.Pt(t[0]),
		transform.Pt(t[1]),
		transform.Pt(t[2]),
	}
}

// TransformShape fulfills shape.TransformShaper
func (t *Triangle) TransformShape(transform *d2.T) shape.Shape {
	return t.T(transform)
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (t *Triangle) AssertEqual(actual interface{}, tol cmpr.Tolerance) error {
	t2, ok := actual.(*Triangle)
	if !ok {
		return geomerr.TypeMismatch(t, actual)
	}

	return geomerr.NewSliceErrs(3, 3, func(i int) error {
		return t[i].AssertEqual(t2[i], tol)
	})
}
