package d2

import (
	"strconv"
	"strings"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
)

// Pt represets a two dimensional point.
type Pt D2

// Pt is defined on Pt to fulfill Point
func (pt Pt) Pt() Pt { return pt }

// V converts Pt to V
func (pt Pt) V() V { return V(pt) }

// Polar converts Pt to Polar
func (pt Pt) Polar() Polar { return D2(pt).Polar() }

// Angle returns the angle of the point relative to the origin
func (pt Pt) Angle() angle.Rad { return D2(pt).Angle() }

// Mag2 returns the square of the magnitude. For comparisions this can be more
// efficient as it avoids a sqrt call.
func (pt Pt) Mag2() float64 { return D2(pt).Mag2() }

// Mag returns the magnitude of the point relative to the origin
func (pt Pt) Mag() float64 { return D2(pt).Mag() }

// Subtract returns the difference between two points as V
func (pt Pt) Subtract(pt2 Pt) V {
	return D2{
		pt.X - pt2.X,
		pt.Y - pt2.Y,
	}.V()
}

// Add a V to a Pt
func (pt Pt) Add(v V) Pt {
	return D2{
		pt.X + v.X,
		pt.Y + v.Y,
	}.Pt()
}

// Distance returns the distance between to points
func (pt Pt) Distance(pt2 Pt) float64 {
	return pt.Subtract(pt2).Mag()
}

// Multiply performs a scalar multiplication on the Pt
func (pt Pt) Multiply(scale float64) Pt {
	return D2{pt.X * scale, pt.Y * scale}.Pt()
}

// Product of two points
func (pt Pt) Product(pt2 Pt) Pt {
	return D2{pt.X * pt2.X, pt.Y * pt2.Y}.Pt()
}

// Prec is the precision for the String method on F
var Prec = 4

// String fulfills Stringer, returns the vector as "(X, Y)"
func (pt Pt) String() string {
	return strings.Join([]string{
		"Pt(",
		strconv.FormatFloat(pt.X, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(pt.Y, 'f', Prec, 64),
		")",
	}, "")
}

// Min returns a Pt with the lowest X and the lowest Y.
func Min(pts ...Pt) Pt {
	if len(pts) == 0 {
		return Pt{}
	}
	m := pts[0]
	for _, pt := range pts[1:] {
		if pt.X < m.X {
			m.X = pt.X
		}
		if pt.Y < m.Y {
			m.Y = pt.Y
		}
	}
	return m
}

// Max returns a Pt with the highest X and highest Y.
func Max(pts ...Pt) Pt {
	if len(pts) == 0 {
		return Pt{}
	}
	m := pts[0]
	for _, pt := range pts[1:] {
		if pt.X > m.X {
			m.X = pt.X
		}
		if pt.Y > m.Y {
			m.Y = pt.Y
		}
	}
	return m
}

// MinMax takes any number of points and returns a min point with the lowest X
// and the lowest Y in the entire set and a max point with the highest X and
// highest Y in the set.
func MinMax(pts ...Pt) (Pt, Pt) {
	if len(pts) == 0 {
		return Pt{}, Pt{}
	}
	min, max := pts[0], pts[0]
	for _, pt := range pts[1:] {
		min, max = Min(min, pt), Max(max, pt)
	}
	return min, max
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (pt Pt) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	pt2, ok := actual.(Pt)
	if !ok {
		return geomerr.TypeMismatch(pt, actual)
	}
	v := pt.Subtract(pt2)
	if !t.Zero(v.X) || !t.Zero(v.Y) {
		return geomerr.NotEqual(pt, pt2)
	}
	return nil
}
