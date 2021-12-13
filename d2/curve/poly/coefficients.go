package poly

import (
	"github.com/adamcolton/geom/calc/comb"
	poly1d "github.com/adamcolton/geom/calc/poly"
	"github.com/adamcolton/geom/d2"
)

// Coefficients wraps the concept of a list of d2.V. It must return the length
// and be able to return the coefficient at any index.
type Coefficients interface {
	Coefficient(idx int) d2.V
	Len() int
}

// X converts and instance of Coefficients into a set of 1D Coefficients using
// their X components.
type X struct{ Coefficients }

// Coefficient returns the X value of the underlying Coefficients.
func (x X) Coefficient(idx int) float64 {
	if idx >= x.Len() || idx < 0 {
		return 0
	}
	return x.Coefficients.Coefficient(idx).X
}

// Len returns the Len of the underlying Coefficients.
func (x X) Len() int {
	return x.Coefficients.Len()
}

// Y converts and instance of Coefficients into a set of 1D Coefficients using
// their Y components.
type Y struct{ Coefficients }

// Coefficient returns the Y value of the underlying Coefficients.
func (y Y) Coefficient(idx int) float64 {
	if idx >= y.Len() || idx < 0 {
		return 0
	}
	return y.Coefficients.Coefficient(idx).Y
}

// Len returns the Len of the underlying Coefficients.
func (y Y) Len() int {
	return y.Coefficients.Len()
}

// Slice fulfills Coefficients using a Slice.
type Slice []d2.V

// Buf returns and instance of Slice. It will use the provided buffer if
// possible.
func Buf(ln int, buf []d2.V) Slice {
	if cap(buf) >= ln {
		buf = buf[:ln]
		for i := range buf {
			buf[i].X, buf[i].Y = 0, 0
		}
		return buf
	}
	return make(Slice, ln)
}

// Coefficient returns the d2.V at the given index if it is in range, otherwise
// it returns d2.V{0,0}.
func (s Slice) Coefficient(idx int) d2.V {
	if idx >= s.Len() || idx < 0 {
		return d2.V{}
	}
	return s[idx]
}

// Len returns the length of the underlying slice.
func (s Slice) Len() int {
	return len(s)
}

// Sum fulfills Coefficients by adding the two underlying Coefficients together.
type Sum [2]Coefficients

// Coefficient returns the sum of both Coefficients at the given index.
func (s Sum) Coefficient(idx int) d2.V {
	return s[0].Coefficient(idx).Add(s[1].Coefficient(idx))
}

// Len returns the longer of the two underlying Coefficients.
func (s Sum) Len() int {
	ln0 := s[0].Len()
	if ln1 := s[1].Len(); ln1 > ln0 {
		return ln1
	}
	return ln0
}

// Produce of two Coefficients.
type Product [2]Coefficients

// Coefficient of the product at the given index.
func (p Product) Coefficient(idx int) d2.V {
	return d2.V{
		X: poly1d.Product{X{p[0]}, X{p[1]}}.Coefficient(idx),
		Y: poly1d.Product{Y{p[0]}, Y{p[1]}}.Coefficient(idx),
	}
}

// Len of the product of the Coefficients.
func (p Product) Len() int {
	return p[0].Len() + p[1].Len() - 1
}

// Derivative of the underlying Coefficients.
type Derivative struct {
	Coefficients
}

// Coefficient at idx is (idx+1)*Coefficient(idx+1).
func (d Derivative) Coefficient(idx int) d2.V {
	idx++
	return d.Coefficients.Coefficient(idx).Multiply(float64(idx))
}

// Len is one less than the underlying Coefficients.
func (d Derivative) Len() int {
	return d.Coefficients.Len() - 1
}

// Bezier fulfills Coefficients representing a Bezier as a polynomial.
type Bezier []d2.Pt

// NewBezier creates a Polynomial from the points used to define a Bezier curve.
func NewBezier(pts []d2.Pt) Poly {
	return Poly{Bezier(pts)}
}

// Len of the underlying slice
func (b Bezier) Len() int {
	return len(b)
}

var signTab = [2]float64{1, -1}

// Coefficient at the given index of the polynomial representation of the Bezier
// curve.
func (b Bezier) Coefficient(idx int) d2.V {
	// B(t) = ∑ binomialCo(l,i) * (1-t)^(l-i) * t^(i) * points[i]
	l := len(b) - 1
	var sum d2.V

	for i, pt := range b {
		term := idx - i
		if term < 0 {
			break
		}
		v := pt.V().Multiply(float64(comb.Binomial(l, i)))
		s := signTab[term&1] * float64(comb.Binomial(l-i, term))
		sum = sum.Add(v.Multiply(s))
	}
	return sum
}
