package poly_test

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/curve/poly"
	"github.com/adamcolton/geom/geomtest"
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
