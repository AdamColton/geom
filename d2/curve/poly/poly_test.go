package poly

import (
	"testing"

	"github.com/adamcolton/geom/geomtest"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/stretchr/testify/assert"
)

func TestLine(t *testing.T) {
	l := line.New(d2.Pt{1, 2}, d2.Pt{7, 4})
	p := Poly{l.T0.V(), l.D}

	for i := 0.0; i < 1.0; i += 0.05 {
		assert.Equal(t, l.Pt1(i), p.Pt1(i))
	}
}

func TestMultiply(t *testing.T) {
	l1 := line.New(d2.Pt{1, 2}, d2.Pt{7, 4})
	l2 := line.New(d2.Pt{3, 1}, d2.Pt{4, 5})

	p1 := Poly{l1.T0.V(), l1.D}
	p2 := Poly{l2.T0.V(), l2.D}

	m := p1.Multiply(p2)

	for i := 0.0; i < 1.0; i += 0.05 {
		pt1, pt2 := l1.Pt1(i), l2.Pt1(i)
		expected := d2.Pt{pt1.X * pt2.X, pt1.Y * pt2.Y}
		geomtest.Equal(t, expected, m.Pt1(i))
	}
}

func TestDerivative(t *testing.T) {
	p := Bezier([]d2.Pt{{0, 0}, {0.5, 1}, {1, 0}})
	v1 := p.V1c0()
	dp := d2.V1Wrapper{p}

	for i := -10.0; i < 10.0; i += 0.05 {
		v := v1.V1(i)
		d := dp.V1(i)
		assert.InDelta(t, d.X, v.X, 1e-4)
		assert.InDelta(t, d.Y, v.Y, 1e-4)

		v = p.V1(i)
		assert.InDelta(t, d.X, v.X, 1e-4)
		assert.InDelta(t, d.Y, v.Y, 1e-4)
	}
}

func TestBezier(t *testing.T) {
	bp := Bezier([]d2.Pt{{0, 0}, {0.5, 1}, {1, 0}})
	bb := bezier.Bezier([]d2.Pt{{0, 0}, {0.5, 1}, {1, 0}})

	for i := 0.0; i < 1.0; i += 0.05 {
		geomtest.Equal(t, bb.Pt1(i), bp.Pt1(i))
	}
}

func TestQuad(t *testing.T) {
	pts := []d2.Pt{{0, 0}, {0.5, 1}, {1, 0}}
	bb := bezier.Bezier(pts)
	l := line.New(d2.Pt{-1, -1}, d2.Pt{2, 2})

	expected := bb.LineIntersections(l, nil)
	assert.Len(t, expected, 2)

	bp := Bezier(pts)
	lp := Poly{l.T0.V(), l.D}

	t0, t1 := Newton(bp, lp)
	geomtest.Equal(t, bp.Pt1(t0), lp.Pt1(t1))
}
