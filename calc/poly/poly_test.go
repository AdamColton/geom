package poly_test

import (
	"math"
	"sort"
	"testing"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/calc/poly"
	"github.com/adamcolton/geom/geomerr"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func sortFloats(fs []float64) {
	sort.Slice(fs, func(i, j int) bool { return fs[i] < fs[j] })
}

func TestCoefficient(t *testing.T) {
	p := poly.New(1, 2, 3)
	geomtest.Equal(t, 1.0, p.Coefficient(0))
	geomtest.Equal(t, 2.0, p.Coefficient(1))
	geomtest.Equal(t, 3.0, p.Coefficient(2))
	geomtest.Equal(t, 0.0, p.Coefficient(3))
	geomtest.Equal(t, 0.0, p.Coefficient(100))

	geomtest.Equal(t, []float64{1, 2, 3}, p.Buf())

	e := poly.Poly{poly.Empty{}}
	p = poly.New()
	assert.Equal(t, p, e)
	p = poly.Poly{poly.Slice(nil)}
	geomtest.Equal(t, p, e)
	geomtest.Equal(t, 0.0, e.Coefficient(0))

	d0 := poly.Poly{poly.D0(5)}
	p = poly.New(5)
	assert.Equal(t, p, d0)
	p = poly.Poly{poly.Slice{5}}
	geomtest.Equal(t, p, d0)
	geomtest.Equal(t, 0.0, d0.Coefficient(1))

	d1 := poly.Poly{poly.D1(5)}
	p = poly.New(5, 1)
	assert.Equal(t, p, d1)
	p = poly.Poly{poly.Slice{5, 1}}
	geomtest.Equal(t, p, d1)

	p = poly.New(1, 2, 3, 0, 0, 0)
	p2 := poly.New(1, 2, 3)
	geomtest.Equal(t, p, p2)

	buf := make([]float64, 3)
	b := poly.Buf(3, buf)
	p = poly.New(1)
	geomtest.Equal(t, p, poly.Poly{b})
	geomtest.Equal(t, 1.0, buf[0])
}

func TestCopy(t *testing.T) {
	buf := make([]float64, 20)
	p := poly.New(1, 2, 3)
	cp := p.Copy(buf)
	geomtest.Equal(t, p, cp)
	geomtest.Equal(t, p.Buf(), buf[:3])
	geomtest.Equal(t, 0.0, buf[4])

	cp = p.Copy(nil)
	geomtest.Equal(t, p, cp)
}

func TestF(t *testing.T) {
	p := poly.New(5)
	geomtest.Equal(t, 5.0, p.F(2.0))

	p = poly.New(5, 2)
	geomtest.Equal(t, 6.0, p.F(0.5))

	p = poly.New(5, 2, 4)
	geomtest.Equal(t, 7.0, p.F(0.5))
}

func TestAssertEqual(t *testing.T) {
	d1 := poly.Poly{poly.D1(5)}
	p := poly.Poly{poly.Slice{5, 1, 0}}

	err := d1.AssertEqual(p, 1e-10)
	assert.Nil(t, err)

	p = poly.New(1, 5)
	err = p.AssertEqual(d1, 1e-10)
	assert.IsType(t, geomerr.ErrNotEqual{}, err)

	err = p.AssertEqual(1.0, 1e-10)
	assert.IsType(t, geomerr.ErrTypeMismatch{}, err)

}

func TestDivide(t *testing.T) {
	p := poly.New(120, 154, 71, 14, 1) // (x+2)(x+3)(x+4)(x+5)
	f := 0.0

	expected := poly.New(60, 47, 12, 1)
	p, f = p.Divide(-2, p.Buf())
	geomtest.Equal(t, expected, p)
	geomtest.Equal(t, 0.0, f)
	assert.Equal(t, 4, p.Len())
	geomtest.Equal(t, 0.0, p.F(-3))
	geomtest.Equal(t, 6.0, p.F(-2))

	expected = poly.New(12, 7, 1)
	p, f = p.Divide(-5, p.Buf())
	geomtest.Equal(t, expected, p)
	geomtest.Equal(t, 0.0, f)
	assert.Equal(t, 3, p.Len())
	geomtest.Equal(t, 0.0, p.F(-3))
	geomtest.Equal(t, 2.0, p.F(-5))
}

func TestSum(t *testing.T) {
	p1 := poly.New(1, 2)
	p2 := poly.New(3, 4, 5)

	expected := poly.New(4, 6, 5)
	geomtest.Equal(t, expected, p1.Add(p2))

	assert.Equal(t, 3, p2.Add(p1).Len())
}

func TestScale(t *testing.T) {
	got := poly.New(1, 2, 3).Scale(2)
	expected := poly.New(2, 4, 6)
	geomtest.Equal(t, expected, got)

	got = poly.New(1, 2, 3).Scale(2)
	geomtest.Equal(t, expected, got)
}

func TestMultiply(t *testing.T) {
	p1 := poly.New(-1, 1)
	p2 := poly.New(1, 1)

	geomtest.Equal(t, poly.New(-1, 0, 1), p1.Multiply(p2))

	p := poly.New(1)
	p2 = poly.New(1)

	for i := 2.0; i < 6; i++ {
		x := poly.New(-i, 1)
		p = p.Multiply(x).Copy(nil)
		p2 = p2.Multiply(x)
	}
	expected := poly.New(120, -154, 71, -14, 1)
	geomtest.Equal(t, expected, p)
	geomtest.Equal(t, expected, p2)
}

func TestMultSwap(t *testing.T) {
	buf, bufa, bufb := make([]float64, 10), make([]float64, 10), make([]float64, 10)
	expa := poly.New(1, 1)
	expb := poly.New(-1, 1)
	a := expa.Copy(bufa)
	b := expb.Copy(bufb)
	swap := buf

	// buf --> a
	// bufa --> swap
	swap = a.MultSwap(b, swap)
	expa = expa.Multiply(expb)
	geomtest.Equal(t, expa, a)
	geomtest.Equal(t, a.Buf(), buf[:3]) // a should now be in buf
	assert.Equal(t, swap, bufa[:2])     // swap will have the old value of a

	// bufa --> b
	// bufb --> swap
	swap = b.MultSwap(a, swap)
	expb = expb.Multiply(expa)
	geomtest.Equal(t, expb, b)
	geomtest.Equal(t, b.Buf(), bufa[:4]) // a should now be in buf
	assert.Equal(t, swap, bufb[:2])      // swap will have the old value of a

	// bufb --> a
	// buf --> swap
	swap = a.MultSwap(b, swap)
	expa = expa.Multiply(expb)
	geomtest.Equal(t, expa, a)
	geomtest.Equal(t, a.Buf(), bufb[:6]) // a should now be in buf
	assert.Equal(t, swap, buf[:3])       // swap will have the old value of a
}

func TestExp(t *testing.T) {
	tt := map[string]struct {
		p   poly.Poly
		pow int
	}{
		"(x2+c)^5": {
			p:   poly.New(2, -3),
			pow: 5,
		},
		"(x3+x2+c)^4": {
			p:   poly.New(1, 1, 1),
			pow: 4,
		},
		"(x4+x3+x2+c)^3": {
			p:   poly.New(4, 2, -3, 1),
			pow: 3,
		},
		"(x4+x3+x2+c)^1": {
			p:   poly.New(4, 2, -3, 1),
			pow: 1,
		},
		"(x4+x3+x2+c)^2": {
			p:   poly.New(4, 2, -3, 1),
			pow: 2,
		},
		"(x4+x3+x2+c)^0": {
			p:   poly.New(4, 2, -3, 1),
			pow: 0,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			ln := tc.p.Len()*tc.pow - tc.pow + 1
			prod := poly.Poly{poly.Buf(ln, nil)}
			buf := make([]float64, ln)
			for i := 0; i < tc.pow; i++ {
				buf = prod.MultSwap(tc.p, buf)
			}
			buf = make([]float64, ln*3)
			geomtest.Equal(t, prod, tc.p.Exp(tc.pow, buf))
			buf = make([]float64, ln*2+1)
			geomtest.Equal(t, prod, tc.p.Exp(tc.pow, buf))
			buf = make([]float64, ln+1)
			geomtest.Equal(t, prod, tc.p.Exp(tc.pow, buf))
			geomtest.Equal(t, prod, tc.p.Exp(tc.pow, nil))

		})
	}

	// when no buffer is provided the returned value is equal to Poly{Empty{}}
	assert.Equal(t, poly.Poly{poly.Empty{}}, poly.New(4, 2, -3, 1).Exp(-1, nil))
	assert.Equal(t, poly.Poly{poly.D0(1)}, poly.New(4, 2, -3, 1).Exp(0, nil))

	// when a buffer is provided, it is used
	assert.Equal(t, poly.Poly{poly.Slice{}}, poly.New(4, 2, -3, 1).Exp(-1, []float64{1, 2, 3}))
	assert.Equal(t, poly.Poly{poly.Slice{1}}, poly.New(4, 2, -3, 1).Exp(0, []float64{5, 2, 3}))
}

func TestD(t *testing.T) {
	geomtest.Equal(t, poly.New(1, 8), poly.New(3, 1, 4).D())
	geomtest.Equal(t, poly.New(1, 8, 3), poly.New(3, 1, 4, 1).D())

	p := poly.New(3, 1, 4, 1)
	d := p.D()
	geomtest.Equal(t, poly.New(1, 8, 3), d)

	dc := poly.Poly{poly.Derivative{p}}

	for x := -10.0; x < 10.0; x += 0.1 {
		df := d.F(x)
		assert.Equal(t, df, p.Df(x))
		assert.Equal(t, df, dc.F(x))
	}
}

func TestIntegral(t *testing.T) {
	p := poly.New(1, 2)
	i := p.Integral(-1)
	d := i.D()
	geomtest.Equal(t, d, p)
	geomtest.Equal(t, -1.0, i.F(0))

	i = p.IntegralAt(1, 1)
	geomtest.Equal(t, 1.0, i.F(1.0))
	i = p.IntegralAt(1, 2)
	geomtest.Equal(t, 2.0, i.F(1.0))
}

func TestNewton(t *testing.T) {
	want := cmpr.Tolerance(1e-10)
	req := cmpr.Tolerance(1e-5)
	buf := make([]float64, 11)
	p := poly.Poly{poly.Buf(12, nil)}
	dbuf := poly.Buf(11, nil)

	for i := 2.0; i < 12; i++ {
		buf = p.MultSwap(poly.New(-i, 1), buf)
		d := p.D().Copy(dbuf).Coefficients
		for j := 1.9; j < i; j++ {
			for variants := 0; variants < 2; variants++ {
				if variants == 1 {
					d = nil
				}
				r, y := p.Newton(j, want, 100, d)
				geomtest.EqualInDelta(t, 0.0, y, req)
				geomtest.EqualInDelta(t, 0.0, p.F(r), req)
				geomtest.EqualInDelta(t, j+0.1, r, req)
			}
		}
	}
}

func TestHalley(t *testing.T) {
	want := cmpr.Tolerance(1e-10)
	req := cmpr.Tolerance(1e-5)
	buf := make([]float64, 11)
	p := poly.Poly{poly.Buf(12, nil)}
	dBuf, ddBuf := poly.Buf(11, nil), poly.Buf(10, nil)

	for i := 2.0; i < 12; i++ {
		buf = p.MultSwap(poly.New(-i, 1), buf)
		for j := 1.9; j < i; j++ {
			d := p.D().Copy(dBuf).Coefficients
			dd := poly.Poly{d}.D().Copy(ddBuf).Coefficients
			for variants := 0; variants < 2; variants++ {
				if variants == 1 {
					d, dd = nil, nil
				}
				r, y := p.Halley(j, want, 50, d, dd)
				geomtest.EqualInDelta(t, 0.0, y, req)
				geomtest.EqualInDelta(t, 0.0, p.F(r), req)
				geomtest.EqualInDelta(t, j+0.1, r, req)
			}
		}
	}

	// Start at a point that will cycle
	p = poly.New(0, 0, -8, 0, 1)
	r, y := p.Halley(2, want, 50, nil, nil)
	geomtest.EqualInDelta(t, 0.0, y, req)
	geomtest.EqualInDelta(t, 0.0, p.F(r), req)
	geomtest.EqualInDelta(t, math.Sqrt(8), r, req)

	// Start at a point that will have a denominator of 0
	// setup a case where 2*dp*dp - p*d2p = 0 but
	// p, dp and d2p are not 0
	x := 1.1
	d2p := poly.New(-6, 6)
	d2px := d2p.F(x)

	dp := d2p.IntegralAt(x, math.Sqrt(d2px))
	dpx := dp.F(x)
	geomtest.Equal(t, math.Pow(dpx, 2), d2px)

	p = dp.IntegralAt(x, 2)
	px := p.F(x)
	geomtest.Equal(t, 2.0, px)
	geomtest.Equal(t, 0.0, 2*dpx*dpx-px*d2px)
	r, y = p.Halley(x, want, 50, dp, d2p)
	geomtest.EqualInDelta(t, 0.0, y, req)
	geomtest.EqualInDelta(t, 0.0, p.F(r), req)

}

func TestRoots(t *testing.T) {
	ln := 17.0
	bufLn := 5*int(ln) - 6
	p := poly.Poly{poly.Buf(bufLn, nil)}
	buf := poly.Slice(make([]float64, bufLn))
	for i := 2.0; i < ln; i++ {
		buf = p.MultSwap(poly.New(-i, 1), buf)

		for variants := 0; variants < 4; variants++ {
			var roots []float64
			if variants == 0 {
				roots = p.Roots(buf)
			} else if variants == 1 {
				p0 := poly.Poly{append(p.Coefficients.(poly.Slice), 0)}
				roots = p0.Roots(buf)
			} else if variants == 2 {
				roots = p.Roots(nil)
			} else {
				for j := 1; j < int(i)-1; j++ {
					roots = p.Roots(buf[:j])
					assert.Len(t, roots, j)
				}
				continue
			}
			assert.Len(t, roots, int(i)-1)
			sortFloats(roots)
			expected := 2.0
			for _, r := range roots {
				assert.InDelta(t, expected, r, 6e-3)
				expected++
			}
		}
	}

	assert.Nil(t, poly.New(1).Roots(nil))

	p = poly.New(-2, 1, 0, 0, 0, 0, 0, 1)
	assert.Equal(t, []float64{1}, p.Roots(nil))
}

func TestQuad(t *testing.T) {
	tt := map[string]struct {
		expected []float64
		a, b, c  float64
	}{
		"2-intercepts": {
			a: 1, b: -8, c: 12,
			expected: []float64{2, 6},
		},
		"1-intercept": {
			a: 1, b: -8, c: 16,
			expected: []float64{4},
		},
		"0-intercepts": {
			a: 1, b: -8, c: 17,
			expected: nil,
		},
		"a=0": {
			b: -8, c: 16,
			expected: []float64{2},
		},
		"a=0&b=0": {
			c:        16,
			expected: nil,
		},
	}

	buf := make([]float64, 2)
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			got := poly.Quad(tc.c, tc.b, tc.a, buf)
			sortFloats(got)
			geomtest.Equal(t, tc.expected, got)
			p := poly.New(tc.c, tc.b, tc.a)
			for _, r := range got {
				geomtest.Equal(t, 0.0, p.F(r))
			}
			if len(tc.expected) > 1 {
				got = poly.Quad(tc.c, tc.b, tc.a, buf[:1])
				assert.Len(t, got, 1)
			}
		})
	}
}
