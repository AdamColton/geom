package line

import (
	"strings"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomerr"
)

// Line in 2D space invoked parametrically
type Line struct {
	T0 d2.Pt
	D  d2.V
}

// Pt1 returns a Pt on the line
func (l Line) Pt1(t float64) d2.Pt {
	return l.T0.Add(l.D.Multiply(t))
}

// V1 always returns l.D, the slope of the line
func (l Line) V1(t float64) d2.V {
	return l.D
}

// AtX Returns the value of t at x. May return Inf.
func (l Line) AtX(x float64) float64 {
	return (x - l.T0.X) / l.D.X
}

// AtY Returns the value of t at x. May return Inf.
func (l Line) AtY(y float64) float64 {
	return (y - l.T0.Y) / l.D.Y
}

// B from the form y = mx + b, this will panic if l.D.X is zero
func (l Line) B() float64 {
	return l.Pt1(l.AtX(0)).Y
}

// M from the form y = mx + b, this will panic if l.D.X is zero
func (l Line) M() float64 {
	return l.D.Y / l.D.X
}

// LineIntersections returns the points at which the lines intersect. It fulls
// the Intersections interface. If the lines are parallel, nil is returned.
// Otherwise a slice with a single value is returned indicating the parametric
// point along l2 at which the intersection occures.
func (l Line) LineIntersections(l2 Line, buf []float64) []float64 {
	_, t, does := l.Intersection(l2)
	if !does {
		return buf[:0]
	}
	return append(buf[:0], t)
}

// Small is the value that will be used to compare against 0.
var Small = cmpr.Tolerance(1e-12)

// Intersection returns the parametric values of the intersection point on the
// line passed in as an argument and a bool indicating if there was an
// intersection.
func (l Line) Intersection(l2 Line) (tl float64, tl2 float64, does bool) {
	d := l.D.Cross(l2.D)
	does = !Small.Zero(d)
	if does {
		v := l.T0.Subtract(l2.T0)
		tl2 = (l.D.X*v.Y - l.D.Y*v.X) / d
		tl = (l2.D.Y*(v.X) + l2.D.X*(-v.Y)) / -d
	}
	return
}

// Closest returns the point on the line closest to pt
func (l Line) Closest(pt d2.Pt) d2.Pt {
	l2 := Line{
		T0: pt,
		D:  d2.V{-l.D.Y, l.D.X},
	}
	_, t0, hit := l2.Intersection(l)
	if !hit {
		return l.T0
	}
	return l.Pt1(t0)
}

// String fulfills Stringer
func (l Line) String() string {
	return strings.Join([]string{
		"Line( ",
		l.D.String(),
		"t + ",
		l.T0.String(),
		" )",
	}, "")
}

// New line from start to end so that l.Pt1(0)==start and l.Pt1(1)==end.
func New(start, end d2.Pt) Line {
	return Line{
		T0: start,
		D:  end.Subtract(start),
	}
}

// Bisect returns a line that bisects points a and b. All points on the line are
// equadistant from both a and b. At t=0, the mid-point is returned. At t=1, the
// point is the same distance from t=0 as the two definition points.
func Bisect(a, b d2.Point) Line {
	m, n := a.Pt(), b.Pt()
	return Line{
		T0: d2.Pt{(m.X + n.X) / 2.0, (m.Y + n.Y) / 2.0},
		D:  d2.V{(m.Y - n.Y) / 2.0, (n.X - m.X) / 2.0},
	}
}

// TangentLine takes a Pt1V1 and a parametric t0 and returns a line on
// the curve at that point, tangent to that point.
func TangentLine(c d2.Pt1V1, t0 float64) Line {
	return Line{
		T0: c.Pt1(t0),
		D:  c.V1(t0),
	}
}

// L fulfills d2.Limiter
func (Line) L(t, c int) d2.Limit {
	if t == 1 && c == 1 {
		return d2.LimitUnbounded
	}
	return d2.LimitUndefined
}

// VL fulfills d2.VLimiter
func (Line) VL(t, c int) d2.Limit {
	if t == 1 && c == 1 {
		return d2.LimitUnbounded
	}
	return d2.LimitUndefined
}

// T applies a transform to the line returning a new line.
func (l Line) T(t *d2.T) Line {
	return Line{
		T0: t.Pt(l.T0),
		D:  t.V(l.D),
	}
}

// Centroid point on the line
func (l Line) Centroid() d2.Pt {
	return l.Pt1(0.5)
}

// Cross product of the vector of the line with the vector from T0 to pt
func (l Line) Cross(pt d2.Pt) float64 {
	return l.D.Cross(pt.Subtract(l.T0))
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (l Line) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	l2, ok := actual.(Line)
	if !ok {
		return geomerr.TypeMismatch(l, actual)
	}
	if l.T0.AssertEqual(l2.T0, t) != nil || l.D.AssertEqual(l2.D, t) != nil {
		return geomerr.NotEqual(l, l2)
	}
	return nil
}
