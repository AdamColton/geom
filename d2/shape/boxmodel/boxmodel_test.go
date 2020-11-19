package boxmodel

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/box"
	"github.com/adamcolton/geom/d2/shape/ellipse"
	"github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/stretchr/testify/assert"
)

type testShape interface {
	shape.Shape
	shape.Centroid
	shape.Area
}

func countIterator(fn func() (Iterator, box.Box, bool)) int {
	sum := 0
	for c, _, done := fn(); !done; _, done = c.Next() {
		sum++
	}
	return sum
}

func TestBasicShapes(t *testing.T) {
	tt := map[string]testShape{
		"ellipse":  ellipse.New(d2.Pt{100, 350}, d2.Pt{400, 110}, 170),
		"triangle": &triangle.Triangle{{100, 100}, {200, 400}, {400, 50}},
		"polygon":  polygon.RegularPolygonRadius(d2.Pt{250, 250}, 200, 0, 7),
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			b := New(tc, 12)
			d := tc.Centroid().Distance(b.Centroid())
			assert.InDelta(t, 0, d, 0.02)
			// the difference in the relative area should be less than 0.001
			a1, a2 := tc.Area(), b.Area()
			p := (a1 - a2) / a1
			assert.InDelta(t, 0, p, 1e-3)
			assert.Equal(t, b.Area(), b.SignedArea())

			assert.Equal(t, countIterator(b.OutsideCursor), b.Outside())
			assert.Equal(t, countIterator(b.InsideCursor), b.Inside())
			assert.Equal(t, countIterator(b.PerimeterCursor), b.Perimeter())
		})
	}
}
