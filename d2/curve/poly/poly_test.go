package poly_test

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/curve/poly"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestLine(t *testing.T) {
	l := line.New(d2.Pt{1, 2}, d2.Pt{7, 4})
	p := poly.New(l.T0.V(), l.D)
	x := p.X()
	y := p.Y()

	geomtest.Equal(t, d2.V{}, p.Coefficient(2))
	geomtest.Equal(t, 0.0, x.Coefficient(2))
	geomtest.Equal(t, 0.0, y.Coefficient(2))

	for i := 0.0; i < 1.0; i += 0.05 {
		pt := l.Pt1(i)
		geomtest.Equal(t, pt, p.Pt1(i))
		geomtest.Equal(t, pt.X, x.F(i))
		geomtest.Equal(t, pt.Y, y.F(i))
	}
}

func TestAdd(t *testing.T) {
	p1 := poly.New(d2.V{1, 2}, d2.V{3, 4})
	p2 := poly.New(d2.V{5, 6}, d2.V{7, 8})
	s := p1.Add(p2)

	for i := 0.0; i < 1.0; i += 0.05 {
		geomtest.Equal(t, s.Pt1(i), p1.Pt1(i).Add(d2.V(p2.Pt1(i))))
	}

	p1 = poly.New(d2.V{1, 2}, d2.V{3, 4})
	p2 = poly.New(d2.V{5, 6})

	assert.Equal(t, 2, p1.Add(p2).Len())
	assert.Equal(t, 2, p2.Add(p1).Len())
}

func TestMultiply(t *testing.T) {
	l1 := line.New(d2.Pt{1, 2}, d2.Pt{7, 4})
	l2 := line.New(d2.Pt{3, 1}, d2.Pt{4, 5})

	p1 := poly.New(l1.T0.V(), l1.D)
	p2 := poly.New(l2.T0.V(), l2.D)

	m := p1.Multiply(p2)

	for i := 0.0; i < 1.0; i += 0.05 {
		pt1, pt2 := l1.Pt1(i), l2.Pt1(i)
		expected := d2.Pt{pt1.X * pt2.X, pt1.Y * pt2.Y}
		geomtest.Equal(t, expected, m.Pt1(i))
	}
}
