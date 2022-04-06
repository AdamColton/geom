package d2

import (
	"strconv"
	"strings"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
)

// V represents a Vector, the difference between two points.
type V D2

// Pt converts a V to Pt
func (v V) Pt() Pt { return Pt(v) }

// V fulfills the Vector interface
func (v V) V() V { return v }

// Polar converts V to Polar
func (v V) Polar() Polar { return D2(v).Polar() }

// Angle of the vector
func (v V) Angle() angle.Rad { return D2(v).Angle() }

// Mag2 returns the square magnitude of the vector. This can be useful for comparisons
// to avoid the additional cost of a Sqrt call.
func (v V) Mag2() float64 { return D2(v).Mag2() }

// Mag returns the magnitude of the vector
func (v V) Mag() float64 { return D2(v).Mag() }

// Cross returns the cross product of the two vectors
func (v V) Cross(v2 V) float64 {
	return v.X*v2.Y - v2.X*v.Y
}

// Dot returns the dot product of two vectors
func (v V) Dot(v2 V) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

// Multiply returns the scalar product of a vector.
func (v V) Multiply(scale float64) V {
	return V{v.X * scale, v.Y * scale}
}

// Product of two vectors
func (v V) Product(v2 V) V {
	return V{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
	}
}

// Sincos returns the sin and cos of the angle of the vector. This fulfills
// Sincoser.
func (v V) Sincos() (float64, float64) {
	h := v.Mag()
	return v.Y / h, v.X / h
}

// String fulfills Stringer, returns the vector as "V(X, Y)"
func (v V) String() string {
	return strings.Join([]string{
		"V(",
		strconv.FormatFloat(v.X, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(v.Y, 'f', Prec, 64),
		")",
	}, "")
}

// Add two vectors
func (v V) Add(v2 V) V {
	return V{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

// Subtract two vectors
func (v V) Subtract(v2 V) V {
	return V{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

// Abs of vector for both X and Y
func (v V) Abs() V {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	return v
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (v V) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	v2, ok := actual.(V)
	if !ok {
		return geomerr.TypeMismatch(v, actual)
	}
	d := v.Subtract(v2)
	if !t.Zero(d.X) || !t.Zero(d.Y) {
		return geomerr.NotEqual(v, v2)
	}
	return nil
}
