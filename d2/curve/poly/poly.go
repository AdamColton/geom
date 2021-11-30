package poly

import (
	poly1d "github.com/adamcolton/geom/calc/poly"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

// Poly is a 2D polynomial curve.
type Poly struct {
	Coefficients
}

// New polynomial curve
func New(pts ...d2.V) Poly {
	return Poly{Slice(pts)}
}

// Copy the coefficients into an instance of Slice. The provided buffer will
// be used if it has sufficient capacity.
func (p Poly) Copy(buf []d2.V) Poly {
	ln := p.Len()
	out := Buf(ln, buf)
	for i := range out {
		out[i] = p.Coefficient(i)
	}
	return Poly{out}
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

// Copy the coefficients into an instance of Slice. The provided buffer will
// be used if it has sufficient capacity.
func (v V) Copy(buf []d2.V) V {
	return V{v.Poly.Copy(buf)}
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
// Note that this is still not buffered, so for repeated calls, make a copy to
// reduce duplicated work.
func (p Poly) V1c0() d2.V1 {
	return p.V()
}

// V1 takes the derivate of p at t0.
func (p Poly) V1(t0 float64) d2.V {
	return p.V1c0().V1(t0)
}

// PolyLineIntersections returns the intersection points relative to the
// Polynomial curve.
func (p Poly) PolyLineIntersections(l line.Line, buf []float64) []float64 {
	if l.D.X == 0 {
		d0 := poly1d.New(-l.T0.X)
		return p.X().Add(d0).Roots(buf)
	}
	if l.D.Y == 0 {
		d0 := poly1d.New(-l.T0.Y)
		return p.Y().Add(d0).Roots(buf)
	}
	m := l.M()
	p2 := p.Y().Add(p.X().Scale(-m))
	d0 := poly1d.New(m*l.T0.X - l.T0.Y)
	p2 = p2.Add(d0)

	return p2.Roots(buf)
}

// LineIntersections fulfills line.Intersector and returns the intersections
// relative to the line.
func (p Poly) LineIntersections(l line.Line, buf []float64) []float64 {
	ln := len(buf)
	ts := p.PolyLineIntersections(l, buf)
	if lnTs := len(ts); lnTs > 0 {
		if ln == 0 || lnTs < ln {
			ln = lnTs
		}
		var toCoord, atCoord func(float64) float64
		if l.D.X == 0 {
			toCoord, atCoord = p.Y().F, l.AtY
		} else {
			toCoord, atCoord = p.X().F, l.AtX
		}
		for i, t := range ts[:ln] {
			ts[i] = atCoord(toCoord(t))
		}
	}
	return ts
}
