package poly

import (
	poly1d "github.com/adamcolton/geom/calc/poly"
	"github.com/adamcolton/geom/d2"
)

// Poly is a 2D polynomial curve.
type Poly struct {
	Coefficients
}

// New polynomial curve
func New(pts ...d2.V) Poly {
	return Poly{Slice(pts)}
}

// Pt1 returns the point on the curve at t0.
func (p Poly) Pt1(t0 float64) d2.Pt {
	return d2.Pt{
		X: p.X().F(t0),
		Y: p.Y().F(t0),
	}
}

// X returns the 1D polynomial formed by the X values.
func (p Poly) X() poly1d.Poly {
	return poly1d.Poly{X(p)}
}

// X returns the 1D polynomial formed by the Y values.
func (p Poly) Y() poly1d.Poly {
	return poly1d.Poly{Y(p)}
}

// Add creates a new polynomial by summinging p with p2.
func (p Poly) Add(p2 Poly) Poly {
	return Poly{Sum{p, p2}}
}
