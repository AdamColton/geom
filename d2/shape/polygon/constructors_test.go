package polygon

import (
	"testing"

	"github.com/adamcolton/geom/geomtest"

	"github.com/adamcolton/geom/d2"
)

func TestRectangleTwoPoints(t *testing.T) {
	p := RectangleTwoPoints(d2.Pt{2.5, 8.1}, d2.Pt{5.3, 1.1})
	expected := Polygon{
		{2.5, 8.1},
		{5.3, 8.1},
		{5.3, 1.1},
		{2.5, 1.1},
	}
	geomtest.Equal(t, []d2.Pt(expected), []d2.Pt(p))
}
