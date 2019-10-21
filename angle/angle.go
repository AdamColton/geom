package angle

import (
	"math"
)

// Rad is angle in radians
type Rad float64

// Deg is angle in degrees
func Deg(d float64) Rad {
	return Rad(d * math.Pi / 180)
}

// Rot is angle as percent of a rotation
func Rot(r float64) Rad {
	return Rad(r * 2 * math.Pi)
}

// Rad returns the angle as float64 reprenting radians
func (r Rad) Rad() float64 {
	return float64(r)
}

// Deg returns the angle as float64 reprenting degrees
func (r Rad) Deg() float64 {
	return float64(r) * 180 / math.Pi
}

// Rot returns the angle as float64 reprenting percent rotations
func (r Rad) Rot() float64 {
	return float64(r) / (2 * math.Pi)
}

// Sin returns the sine of the angle
func (r Rad) Sin() float64 {
	return math.Sin(float64(r))
}

// Cos returns the cosine of the angle
func (r Rad) Cos() float64 {
	return math.Cos(float64(r))
}

// Sincos returns both the sine and cosine of the angle
func (r Rad) Sincos() (float64, float64) {
	return math.Sincos(float64(r))
}
