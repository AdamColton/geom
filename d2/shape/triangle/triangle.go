package triangle

import (
	"math"
	"strings"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

type Triangle [3]d2.Pt

func (t Triangle) Contains(pt d2.Pt) bool {
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
func (t Triangle) SignedArea() float64 {
	v1 := t[0].Subtract(t[1])
	v2 := t[0].Subtract(t[2])
	return 0.5 * v1.Cross(v2)
}

// Area of the triangle
func (t Triangle) Area() float64 {
	return math.Abs(t.SignedArea())
}

// Perimeter of the triangle
func (t Triangle) Perimeter() float64 {
	return t[0].Distance(t[1]) + t[1].Distance(t[2]) + t[2].Distance(t[0])
}

// Centroid returns the center of mass of the triangle
func (t Triangle) Centroid() d2.Pt {
	return d2.Pt{
		X: (t[0].X + t[1].X + t[2].X) / 3.0,
		Y: (t[0].Y + t[1].Y + t[2].Y) / 3.0,
	}
}

func (t Triangle) String() string {
	return strings.Join([]string{
		"Triangle[ ",
		t[0].String(), ",",
		t[1].String(), ",",
		t[2].String(), " ]",
	}, "")
}

func (t Triangle) Pt1c0() d2.Pt1 {
	return line.Segments(append(t[:], t[0]))
}

func (t Triangle) Pt1(t0 float64) d2.Pt {
	return t.Pt1c0().Pt1(t0)
}

func (t Triangle) Pt2c1(t0 float64) d2.Pt1 {
	m := line.New(t[0], t[1]).Pt1(0.5)
	p0 := line.New(t[0], m).Pt1(t0)
	p1 := line.New(t[2], t[1]).Pt1(t0)
	return line.New(p0, p1)
}

func (t Triangle) Pt2(t0, t1 float64) d2.Pt {
	return t.Pt2c1(t0).Pt1(t1)
}

func (t Triangle) L(ts, c int) d2.Limit {
	if (ts == 1 && c == 1) ||
		(ts == 2 && c == 2) ||
		(ts == 1 && c == 0) ||
		(ts == 2 && c == 1) {
		return d2.LimitBounded
	}
	return d2.LimitUndefined
}

func (t Triangle) Intersections(l line.Line) []float64 {
	var out []float64
	prev := t[2]
	for _, cur := range t {
		l2 := line.New(prev, cur)
		prev = cur
		tt, ok := l.LineIntersection(l2)
		if tt >= 0 && tt < 1 && ok {
			t0, _ := l2.LineIntersection(l)
			out = append(out, t0)
		}
	}
	return out
}

func (t Triangle) BoundingBox() (d2.Pt, d2.Pt) {
	return d2.MinMax(t[:]...)
}

func (t Triangle) CircumCenter() d2.Pt {
	l01 := line.Bisect(t[0], t[1])
	l02 := line.Bisect(t[0], t[2])
	t0, ok := l01.LineIntersection(l02)
	if ok {
		return l02.Pt1(t0)
	}
	if t[0] == t[1] || t[1] == t[2] {
		return line.New(t[0], t[2]).Pt1(0.5)
	}
	return line.New(t[0], t[1]).Pt1(0.5)
}
