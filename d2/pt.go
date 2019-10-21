package d2

import (
	"strconv"
	"strings"

	"github.com/adamcolton/geom/angle"
)

type Pt D2

func (pt Pt) Pt() Pt           { return pt }
func (pt Pt) V() V             { return V(pt) }
func (pt Pt) Polar() Polar     { return D2(pt).Polar() }
func (pt Pt) Angle() angle.Rad { return D2(pt).Angle() }
func (pt Pt) Mag2() float64    { return D2(pt).Mag2() }
func (pt Pt) Mag() float64     { return D2(pt).Mag() }

func (pt Pt) Subtract(pt2 Pt) V {
	return D2{
		pt.X - pt2.X,
		pt.Y - pt2.Y,
	}.V()
}

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

func (pt Pt) Multiply(scale float64) Pt {
	return D2{pt.X * scale, pt.Y * scale}.Pt()
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

func Min(pt1, pt2 Pt) Pt {
	if pt2.X < pt1.X {
		pt1.X = pt2.X
	}
	if pt2.Y < pt1.Y {
		pt1.Y = pt2.Y
	}
	return pt1
}

func Max(pt1, pt2 Pt) Pt {
	if pt2.X > pt1.X {
		pt1.X = pt2.X
	}
	if pt2.Y > pt1.Y {
		pt1.Y = pt2.Y
	}
	return pt1
}

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
