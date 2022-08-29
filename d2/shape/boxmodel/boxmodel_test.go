package boxmodel

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/box"
	"github.com/adamcolton/geom/d2/shape/ellipse"
	"github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

type testShape interface {
	shape.Shape
	shape.Centroid
	shape.Area
	d2.Pt1
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

			var a polygon.AssertConvexHull
			for i := 0.0; i < 1.0; i += 0.1 {
				a = append(a, tc.Pt1(i))
			}
			for i, bx, done := b.InsideCursor(); !done; bx, done = i.Next() {
				a = append(a, bx.Centroid())
				assert.True(t, tc.Contains(bx.Centroid()))
			}

			// Some points lie just outside the hull, scaling up by just 0.1%
			// guarentees all points lie inside the hull
			scale := 1.001
			h := (&d2.TransformSet{}).
				AddBoth(d2.Translate(b.Centroid().Multiply(-1))).
				Add(d2.Scale(d2.V{scale, scale}).T()).
				GetT().
				Slice(b.ConvexHull())

			geomtest.Equal(t, a, h)
		})
	}
}
