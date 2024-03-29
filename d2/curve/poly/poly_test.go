package poly_test

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
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

func TestDerivative(t *testing.T) {
	p := poly.New(d2.V{0.0000, 0.0000}, d2.V{1.0000, 2.0000}, d2.V{0.0000, -2.0000})
	v1 := p.V1c0()
	vCp := p.V().Copy(nil)

	geomtest.EqualInDelta(t, d2.AssertV1{}, p, 1e-4)
	for i := 0.0; i <= 1.0; i += 0.05 {
		v := p.V1(i)
		geomtest.EqualInDelta(t, v, v1.V1(i), 1e-4)
		geomtest.EqualInDelta(t, v, vCp.V1(i), 1e-4)
	}
}

func TestBezier(t *testing.T) {
	bb := bezier.Bezier([]d2.Pt{{0, 0}, {0.5, 1}, {1, 0}})
	bp := poly.NewBezier(bb).Copy(nil)
	bc := poly.Poly{poly.Bezier(bb)}

	for i := 0.0; i < 1.0; i += 0.05 {
		pt := bb.Pt1(i)
		geomtest.Equal(t, pt, bp.Pt1(i))
		geomtest.Equal(t, pt, bc.Pt1(i))
	}

	buf := make([]d2.V, 3)
	bp = poly.NewBezier(bb).Copy(buf)
	for i, b := range buf {
		geomtest.Equal(t, bc.Coefficient(i), b)
	}
}

func TestLineIntersections(t *testing.T) {
	bp := poly.NewBezier([]d2.Pt{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	})
	tt := map[string]struct {
		line.Line
		Points int
	}{
		"horizontal": {
			Line:   line.New(d2.Pt{0, 200}, d2.Pt{500, 200}),
			Points: 3,
		},
		"vertical": {
			Line:   line.New(d2.Pt{500, 200}, d2.Pt{500, 300}),
			Points: 3,
		},
		"diagonal": {
			Line:   line.New(d2.Pt{0, 200}, d2.Pt{500, 300}),
			Points: 3,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			li := bp.LineIntersections(tc.Line, nil)
			assert.Len(t, li, tc.Points)
			pi := bp.PolyLineIntersections(tc.Line, nil)
			for i, t0 := range li {
				geomtest.Equal(t, tc.Pt1(t0), bp.Pt1(pi[i]))
			}
		})
	}
}
