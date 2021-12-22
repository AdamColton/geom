package d3

import (
	"math"
	"strconv"
	"strings"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
)

type Pt D3

func (pt Pt) Mag() float64  { return D3(pt).Mag() }
func (pt Pt) Mag2() float64 { return D3(pt).Mag2() }

func (pt Pt) Subtract(pt2 Pt) V {
	return V{
		pt.X - pt2.X,
		pt.Y - pt2.Y,
		pt.Z - pt2.Z,
	}
}

func (pt Pt) Pt() Pt {
	return pt
}

func (pt Pt) Add(v V) Pt {
	return Pt{
		pt.X + v.X,
		pt.Y + v.Y,
		pt.Z + v.Z,
	}
}

// Distance returns the distance between to points
func (pt Pt) Distance(pt2 Pt) float64 {
	return pt.Subtract(pt2).Mag()
}

func (pt Pt) Multiply(scale float64) Pt {
	return Pt{
		pt.X * scale,
		pt.Y * scale,
		pt.Z * scale,
	}
}

func (pt Pt) Round() Pt {
	return Pt{
		X: math.Round(pt.X),
		Y: math.Round(pt.Y),
		Z: math.Round(pt.Z),
	}
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
		", ",
		strconv.FormatFloat(pt.Z, 'f', Prec, 64),
		")",
	}, "")
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (pt Pt) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	pt2, ok := actual.(Pt)
	if !ok {
		return geomerr.TypeMismatch(pt, actual)
	}
	d := pt.Subtract(pt2)
	if !t.Zero(d.X) || !t.Zero(d.Y) || !t.Zero(d.Z) {
		return geomerr.NotEqual(pt, pt2)
	}
	return nil
}

// Min returns a Pt with the lowest X, lowest Z and the lowest Y.
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
		if pt.Z < m.Z {
			m.Z = pt.Z
		}
	}
	return m
}

// Max returns a Pt with the highest X, highest Y and highest Z.
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
		if pt.Z > m.Z {
			m.Z = pt.Z
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
