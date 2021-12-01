package poly_test

import (
	"testing"

	"github.com/adamcolton/geom/calc/poly"
	"github.com/adamcolton/geom/geomerr"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

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
