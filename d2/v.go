package d2

import (
	"math"
	"strconv"
	"strings"
)

type Vector interface {
	V() V
}

type V Pt

// Angle returns the angle in radians
func (v V) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v V) V() V { return v }

// Mag2 Returns the sqaure of the magnitude of the vector.
func (v V) Mag2() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Mag returns the magnitude (distance to origin) of the vector
func (v V) Mag() float64 {
	return math.Sqrt(v.Mag2())
}

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

func Polar(radius, angle float64) V {
	s, c := math.Sincos(angle)
	return V{c * radius, s * radius}
}
