// Package poly performs operations on polynomials.
package poly

import (
	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/calc/fbuf"
	"github.com/adamcolton/geom/geomerr"
)

// Poly is a 1D polynomial. The index corresponds power of X.
type Poly struct {
	Coefficients
}

// New 1D polynomial with the given coefficients.
func New(cs ...float64) Poly {
	ln := len(cs)
	if ln == 0 {
		return Poly{Empty{}}
	}
	if cs[ln-1] == 0 {
		return New(cs[:ln-1]...)
	}
	if ln == 1 {
		return Poly{D0(cs[0])}
	}
	if ln == 2 && cs[1] == 1 {
		return Poly{D1(cs[0])}
	}
	return Poly{Slice(cs)}
}

// Copy a Polynomial into a buffer.
func (p Poly) Copy(buf []float64) Poly {
	out := Slice(fbuf.Slice(p.Len(), buf))
	for i := range out {
		out[i] = p.Coefficient(i)
	}
	return Poly{out}
}

// Buf tries to get the Coefficients as a []float64. This is intended for
// recycling buffers.
func (p Poly) Buf() []float64 {
	buf, _ := p.Coefficients.(Slice)
	return buf
}

// F computes the value of p(x).
func (p Poly) F(x float64) float64 {
	idx := p.Len() - 1
	s := 0.0
	for ; idx >= 0; idx-- {
		s = p.Coefficient(idx) + s*x
	}
	return s
}

// AssertEqual allows Polynomials to be compared. This fulfills
// geomtest.AssertEqualizer.
func (p Poly) AssertEqual(to interface{}, t cmpr.Tolerance) error {
	p2, ok := to.(Poly)
	if !ok {
		return geomerr.TypeMismatch(p, to)
	}

	ln := p.Len()
	if ln2 := p2.Len(); ln2 > ln {
		ln = ln2
	}
	for i := 0; i < ln; i++ {
		if p.Coefficient(i) != p2.Coefficient(i) {
			return geomerr.NotEqual(p, p2)
		}
	}

	return nil
}

// Divide creates a new polynomial by dividing p by (x-n). The float64 returned
// is the remainder. If (x-n) is a root of p this value will be 0.
func (p Poly) Divide(n float64, buf []float64) (Poly, float64) {
	ln := p.Len() - 1
	out := Slice(fbuf.Slice(ln, buf))
	r := p.Coefficient(ln)
	for i := ln - 1; i >= 0; i-- {
		out[i], r = r, p.Coefficient(i)+r*n
	}
	return Poly{out}, r
}

// Add p and p2 using the Sum coefficients.
func (p Poly) Add(p2 Poly) Poly {
	return Poly{Sum{p, p2}}
}

// Scale will return an instace of the Scale Coefficient wrapper.
func (p Poly) Scale(s float64) Poly {
	return Poly{Scale{
		By:           s,
		Coefficients: p,
	}}
}
