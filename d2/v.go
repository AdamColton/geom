package d2

import (
	"strconv"
	"strings"

	"github.com/adamcolton/geom/angle"
)

type V D2

func (v V) Pt() Pt           { return Pt(v) }
func (v V) V() V             { return v }
func (v V) Polar() Polar     { return D2(v).Polar() }
func (v V) Angle() angle.Rad { return D2(v).Angle() }
func (v V) Mag2() float64    { return D2(v).Mag2() }
func (v V) Mag() float64     { return D2(v).Mag() }

// Cross returns the cross product of the two vectors
func (v V) Cross(v2 V) float64 {
	return v.X*v2.Y - v2.X*v.Y
}

// Dot returns the dot product of two vectors
func (v V) Dot(v2 V) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

func (v V) Multiply(scale float64) V {
	return V{v.X * scale, v.Y * scale}
}

func (v V) Product(v2 V) V {
	return V{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
	}
}

// String fulfills Stringer, returns the vector as "(X, Y)"
func (v V) String() string {
	return strings.Join([]string{
		"V(",
		strconv.FormatFloat(v.X, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(v.Y, 'f', Prec, 64),
		")",
	}, "")
}

func (v V) Add(v2 V) V {
	return V{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v V) Subtract(v2 V) V {
	return V{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

func (v V) Abs() V {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	return v
}
