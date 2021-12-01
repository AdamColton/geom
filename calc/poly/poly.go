// Package poly performs operations on polynomials.
package poly

import (
	"math"

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

// Multiply two polynomails. Note that it is not safe to reuse either input as
// the buffer.
func (p Poly) Multiply(p2 Poly) Poly {
	return Poly{Product{p, p2}}
}

// MultSwap does a multiply and swap. It is used for effiency when doing
// consecutive multiplications. It is equivalent to:
//
// p = p.Multiply(p2)
//
// but it swaps the slice backing p with the buf after the multiplicaiton. It
// will generally be used like this:
//
// buf = p.MultSwap(p2, buf)
func (p *Poly) MultSwap(p2 Poly, buf []float64) []float64 {
	prod := p.Multiply(p2)
	out := p.Buf()
	p.Coefficients = prod.Copy(buf).Coefficients
	return out
}

// Exp raises p to the power of n. To effiently allocate the buf it should have
// capacity of 3*(len(tc.p)*tc.pow - tc.pow + 1).
func (p Poly) Exp(n int, buf []float64) Poly {
	if n < 0 {
		if cap(buf) == 0 {
			return Poly{Empty{}}
		}
		return Poly{Slice(buf[:0])}
	} else if n == 0 {
		if cap(buf) == 0 {
			return Poly{D0(1)}
		}
		return Poly{Buf(1, buf)}
	} else if n == 1 {
		return p.Copy(buf)
	} else if n == 2 {
		return p.Multiply(p).Copy(buf)
	}

	// https://en.wikipedia.org/wiki/Exponentiation_by_squaring
	//
	// Because of the repeated multiplication, to use the buffers efficiently,
	// a swap buffer is needed. So a total of 3 polynomials of length ln are
	// needed: sum, cur and swap.
	ln := p.Len()*n - n + 1

	s, buf := fbuf.Split(ln, buf)
	s = append(s, 1)
	sum := Poly{Slice(s)}

	c, buf := fbuf.Split(ln, buf)
	cur := p.Copy(c[:p.Len()])

	buf = fbuf.Slice(ln, buf)

	for {
		if n&1 == 1 {
			buf = sum.MultSwap(cur, buf)
		}
		n >>= 1
		if n == 0 {
			return sum
		}
		buf = cur.MultSwap(cur, buf)
	}
}

// D returns the derivative of p.
func (p Poly) D() Poly {
	return Poly{Derivative{p}}
}

// Df computes the value of p'(x).
func (p Poly) Df(x float64) float64 {
	return Poly{Derivative{p}}.F(x)
}

// Integral of the given polynomial with the constant set to c.
func (p Poly) Integral(c float64) Poly {
	return Poly{Integral{p, c}}
}

// Integral of the given polynomial with the constant set so that the value of
// Pt1(x) == y.
func (p Poly) IntegralAt(x, y float64) Poly {
	i := Integral{p, 0}
	i.C = y - Poly{i}.F(x)
	return Poly{i}
}

// Roots finds the real roots of the polynomial. If an algebraic solution
// exists, that will be used. Otherwise it will use Halley's method to get it
// down to an order 2 solution. Because Halley's method is an approximation,
// errors tend to compound and this seems to become unreliable above a degree 10
// polynomial. It is not safe to use p as the buffer. If the order of p>5 then
// the optimal buffer size is 5*p.Len()-6. The number of roots returned is set
// by the length of the buffer passed in. If the length is 0 then the max number
// of roots is returned.
func (p Poly) Roots(buf []float64) []float64 {
	ln := p.Len()

	if ln < 2 {
		return nil
	}
	if p.Coefficient(ln-1) == 0 {
		return Poly{RemoveLeadingZero{p.Coefficients}}.Roots(buf)
	}
	if ln == 2 {
		return append(buf[:0], -p.Coefficient(0)/p.Coefficient(1))
	}

	outLn := len(buf)
	if outLn == 0 || outLn > ln-1 {
		outLn = ln - 1
	}

	const (
		want cmpr.Tolerance = 1e-15
		need cmpr.Tolerance = 1e-2
	)

	// Note that for optimization, cp and roots are sharing the same buf of
	// length ln. This works because roots grows at the same rate cp shrinks. If
	// buf is not at least length ln, this optimization is wasted. This is also
	// why the length of roots is not set to outLn - it's already sharing space
	// with cp, so that doesn't save any buffer space.
	buf = fbuf.Empty(ln, buf)
	cur := p.Copy(buf)
	roots, buf := fbuf.Split(ln, buf)
	dbuf, buf := fbuf.Split(ln-1, buf)
	ddbuf, buf := fbuf.Split(ln-2, buf)
	d := p.D().Copy(dbuf)
	dd := d.D().Copy(ddbuf)
	dbuf, buf = fbuf.Split(ln-1, buf)
	ddbuf = fbuf.Empty(ln-2, buf)

	// cur is a polynomial that starts equal to p. Halley's method is used to
	// find roots. As roots are found they are divided out of cur. With this
	// approach, errors will accumulate. So cur is used to get close to a root
	// and then that value is passed into Halley on the original p to find the
	// actual root.
	for cur.Len() > 2 && len(roots) < outLn {
		dCur := cur.D().Copy(dbuf)
		ddCur := dCur.D().Copy(ddbuf)
		r, y := cur.Halley(0, need, 50, dCur, ddCur)
		if !need.Zero(y) {
			return roots
		}
		r, _ = p.Halley(r, want, 50, d, dd)
		cur, _ = cur.Divide(r, cur.Buf()[1:])
		roots = append(roots, r)
	}

	if ln := len(roots); ln < outLn {
		roots = append(roots, cur.Roots(roots[ln:outLn])...)
	}

	return roots
}

// Newton's method to find one root of the polynomial. The initial guess is
// passed in as x; min sets how close to 0 is acceptible and it will return if a
// value closer than that is found; steps limits the maximum number of
// iterations that will; d is the derivative. It is not required to provide d,
// but if there is a cached instance available, it reduces repeated computation.
func (p Poly) Newton(x float64, min cmpr.Tolerance, steps int, d Coefficients) (float64, float64) {
	const (
		small cmpr.Tolerance = 1e-5
	)

	if d == nil {
		d = Derivative{p}
	}
	dp := Poly{d}

	y := p.F(x)

	bestY, bestX := math.Abs(y), x
	for i := 0; i < steps && !min.Zero(y); i++ {
		if math.IsInf(y, 0) || math.IsNaN(y) {
			x += 1e-3
			continue
		}
		d := dp.F(x)
		if small.Zero(d) {
			x += 1e-3
			y = p.F(x)
			continue
		}
		d = y / d
		d *= (200 - float64(i)) / 200
		x -= d
		y = p.F(x)
		if absy := math.Abs(y); absy < bestY {
			bestX, bestY = x, absy
		}
	}
	return bestX, bestY
}

// Halley's method to find one root of the polynomial. The initial guess is
// passed in as x; min sets how close to 0 is acceptible and it will return if a
// value closer than that is found; steps limits the maximum number of
// iterations that will; d is the derivative; d2 is the second derivative. It is
// not required to provide d or d2, but if there is a cached instance available,
// it reduces repeated computation.
func (p Poly) Halley(x float64, min cmpr.Tolerance, steps int, d, d2 Coefficients) (float64, float64) {
	const (
		small cmpr.Tolerance = 1e-5
	)

	if d == nil {
		d = Derivative{p}
	}
	if d2 == nil {
		d2 = Derivative{d}
	}
	dp, ddp := Poly{d}, Poly{d2}

	y := p.F(x)
	bestY, bestX := math.Abs(y), x

	for i := 0; i < steps && !min.Zero(y); i++ {
		dy := dp.F(x)
		d2y := ddp.F(x)
		denom := 2*dy*dy - y*d2y
		if small.Zero(denom) {
			x += 1e-3
			y = p.F(x)
			continue
		}
		d := (2 * y * dy) / denom
		d *= (200 - float64(i)) / 200
		x -= d
		y = p.F(x)
		if x == bestX {
			x += 1e-3
		} else if absy := math.Abs(y); absy < bestY {
			bestX, bestY = x, absy
		}
	}
	return bestX, bestY
}
