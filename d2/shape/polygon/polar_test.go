package polygon

import (
	"testing"

	"github.com/adamcolton/geom/geomtest"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
)

func TestPolar(t *testing.T) {
	pp := PolarPolygon{
		{1, angle.Rot(0.0 / 6.0)},
		{0.75, angle.Rot(2.0 / 6.0)},
		{1, angle.Rot(3.0 / 6.0)},
		{2, angle.Rot(5.0 / 6.0)},
		{2, angle.Rot(4.0 / 6.0)},
		{0.75, angle.Rot(1.0 / 6.0)},
	}
	pp.Sort()

	p := pp.Polygon(d2.Pt{0, 0})
	expected := Polygon{
		{X: 1, Y: 0},
		{X: 0.375, Y: 0.649519052838329},
		{X: -0.375, Y: 0.6495190528383291},
		{X: -1, Y: 0},
		{X: -1, Y: -1.7320508075688767},
		{X: 1, Y: -1.7320508075688772},
	}
	for i, e := range expected {
		geomtest.Equal(t, e, p[i])
	}
}
