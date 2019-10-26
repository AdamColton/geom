package d2

import (
	"math"

	"github.com/adamcolton/geom/angle"
)

// D2 is an abstract two dimensional struct holding the methods shared by both
// Pt and V.
type D2 struct {
	X, Y float64
}

// Pt converts D2 to a Pt
func (d D2) Pt() Pt {
	return Pt(d)
}

// V converts D2 to a V
func (d D2) V() V {
	return V(d)
}

// Polar converts D2 to a Polar
func (d D2) Polar() Polar {
	return Polar{d.Mag(), d.Angle()}
}

// Angle returns the angle in radians
func (d D2) Angle() angle.Rad {
	return angle.Atan(d.Y, d.X)
}

// Mag2 Returns the sqaure of the magnitude of the vector.
func (d D2) Mag2() float64 {
	return d.X*d.X + d.Y*d.Y
}

// Mag returns the magnitude (distance to origin) of the vector
func (d D2) Mag() float64 {
	return math.Sqrt(d.Mag2())
}
