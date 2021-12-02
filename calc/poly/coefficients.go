package poly

import "github.com/adamcolton/geom/calc/fbuf"

// Coefficients wraps the concept of a list of float64. It can express the order
// of the polynomial and return any coeffcient.
type Coefficients interface {
	Coefficient(idx int) float64
	Len() int
}

// Slice fulfills Coefficients with a []float64.
type Slice []float64

// Buf creates an instance of Poly with c capacity and a value of 1. This is
// useful when taking the product of several polynomials.
func Buf(c int, buf []float64) Slice {
	return append(fbuf.Empty(c, buf), 1)
}

// Coefficient at idx. If the idx is greater than the length of the
// polynomial, then a 0 is returned.
func (s Slice) Coefficient(idx int) float64 {
	if idx >= len(s) {
		return 0
	}
	return s[idx]
}

// Len of the polynomial is equal to the length of the slice.
func (s Slice) Len() int {
	return len(s)
}

// Empty constructs an empty polynomial.
type Empty struct{}

// Coefficient always returns 0
func (Empty) Coefficient(idx int) float64 {
	return 0
}

// Len always returns 0
func (Empty) Len() int {
	return 0
}

// D0 is a degree 0 polynomial - a constant.
type D0 float64

// Coefficient returns underlying float64 if the idx is 0, otherwise it returns
// 0.
func (d D0) Coefficient(idx int) float64 {
	if idx == 0 {
		return float64(d)
	}
	return 0
}

// Len is always 1
func (D0) Len() int {
	return 1
}

// D1 is a degree 1 polynomial with the first coefficient equal to 1.
type D1 float64

// Coefficient returns the underlying float64 if idx is 0 and returns 1 if the
// idx is 1.
func (d D1) Coefficient(idx int) float64 {
	if idx == 0 {
		return float64(d)
	}
	if idx == 1 {
		return 1
	}
	return 0
}

// Len is always equal to 2
func (D1) Len() int {
	return 2
}

// Sum of 2 Coefficients
type Sum [2]Coefficients

// Coefficient at idx is the sum of the underlying Coefficients at idx.
func (s Sum) Coefficient(idx int) float64 {
	return s[0].Coefficient(idx) + s[1].Coefficient(idx)
}

// Len is the greater len of the 2 Coefficients.
func (s Sum) Len() int {
	ln := s[0].Len()
	if ln2 := s[1].Len(); ln2 > ln {
		return ln2
	}
	return ln
}

// Scale Coefficients by a constant value
type Scale struct {
	By float64
	Coefficients
}

// Coefficient is product of scale factor and the underlying Coefficient at idx.
func (s Scale) Coefficient(idx int) float64 {
	return s.Coefficients.Coefficient(idx) * s.By
}
