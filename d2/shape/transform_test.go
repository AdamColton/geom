package shape_test

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/adamcolton/geom/geomtest"
)

func TestTransformShape(t *testing.T) {
	tri := &triangle.Triangle{
		{0, 0}, {1, 0}, {0, 1},
	}
	translate := d2.Translate(d2.V{1, 1}).GetT()
	var s shape.Shape = shape.Transform(tri, translate)

	expected := &triangle.Triangle{
		{1, 1}, {2, 1}, {1, 2},
	}
	geomtest.Equal(t, expected, s)
}
