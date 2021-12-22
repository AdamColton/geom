package d3

import (
	"math"
	"strconv"
	"strings"
)

// D3 is an abstract two dimensional struct holding the methods shared by both
// Pt and V.
type D3 struct {
	X, Y, Z float64
}

// Pt converts D2 to a Pt.
func (d D3) Pt() Pt { return Pt(d) }

// Pt converts D2 to a V.
func (d D3) V() V { return V(d) }

// Mag2 returns the sqaure of the magnitude of the D3.
func (d D3) Mag2() float64 {
	return d.X*d.X + d.Y*d.Y + d.Z*d.Z
}

// Mag returns the magnitude of the D3.
func (d D3) Mag() float64 {
	return math.Sqrt(d.Mag2())
}

// String fulfills Stringer representing the D3 as a String.
func (d D3) String() string {
	return strings.Join([]string{
		"D3(",
		strconv.FormatFloat(d.X, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(d.Y, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(d.Z, 'f', Prec, 64),
		")",
	}, "")
}
