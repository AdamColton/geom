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

// Product of two Coefficients
type Product [2]Coefficients

// Coefficient at idx is the sum of all p[i]*p2[j] where i+j == idx
func (p Product) Coefficient(idx int) float64 {
	l0 := p[0].Len()
	l1 := p[1].Len()

	var sum float64
	i := idx - l1
	if i < 0 {
		i = 0
	}
	for j := 0; i < l0 && i <= idx; i++ {
		j = idx - i
		sum += p[0].Coefficient(i) * p[1].Coefficient(j)
	}
	return sum
}

// Len is one less than the sum of the lengths.
func (p Product) Len() int {
	return p[0].Len() + p[1].Len() - 1
}

// Derivative of the Coefficients
type Derivative struct {
	Coefficients
}

// Coefficient at idx is (idx+1)*Coefficient(idx+1).
func (d Derivative) Coefficient(idx int) float64 {
	idx++
	return d.Coefficients.Coefficient(idx) * float64(idx)
}

// Len is always one less than the underlying Coefficients.
func (d Derivative) Len() int {
	return d.Coefficients.Len() - 1
}

// Integral of the underlying  Coefficients.
type Integral struct {
	Coefficients
	C float64
}

// Coefficient at idx is Coefficient(idx-1)/idx. Except at 0 where it is C.
func (i Integral) Coefficient(idx int) float64 {
	if idx == 0 {
		return i.C
	}
	return i.Coefficients.Coefficient(idx-1) / float64(idx)
}

// Len is always one more than the underlying Coefficients.
func (i Integral) Len() int {
	return i.Coefficients.Len() + 1
}

// RemoveLeadingZero simplifies a Polynomial where the leading Coefficient is
// zero. Note that this does no verification, it is only intended as a wrapper.
type RemoveLeadingZero struct{ Coefficients }

// Len is always one less than the underlying Coefficients.
func (r RemoveLeadingZero) Len() int { return r.Coefficients.Len() - 1 }
