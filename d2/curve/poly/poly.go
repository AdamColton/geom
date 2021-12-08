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

// Multiply creates a new polynomial by taking the produce of p with p2.
func (p Poly) Multiply(p2 Poly) Poly {
	return Poly{Product{p, p2}}
}

// V represents the derivative of a Polynomial and will return d2.V instead of
// d2.Pt.
type V struct {
	Poly
}

// V returns and instace of V that holds the derivative of p.
func (p Poly) V() V {
	return V{Poly{Derivative{p}}}
}

// V1 returns V at t0.
func (v V) V1(t0 float64) d2.V {
	return v.Pt1(t0).V()
}

// V1c0 returns and instance of V fulfilling d2.V1 and caching the derivative.
func (p Poly) V1c0() d2.V1 {
	return p.V()
}

// V1 takes the derivate of p at t0.
func (p Poly) V1(t0 float64) d2.V {
	return p.V1c0().V1(t0)
}
