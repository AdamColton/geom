package box_test

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/generate"
	"github.com/adamcolton/geom/d2/shape/box"
	"github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestBox(t *testing.T) {
	b := box.New(d2.Pt{0, 0}, d2.Pt{1, 1})
	assert.Equal(t, 1.0, b.Area())
	assert.Equal(t, 1.0, b.SignedArea())
	geomtest.Equal(t, d2.Pt{0.5, 0.5}, b.Centroid())

	for i, s := range b.Sides() {
		si := b.Side(i)
		geomtest.Equal(t, s.T0, si.T0)
		geomtest.Equal(t, s.D, si.D)
	}

	assert.Equal(t, 4.0, b.Perimeter())
	m, M := b.BoundingBox()
	geomtest.Equal(t, b[0], m)
	geomtest.Equal(t, b[1], M)

	a := polygon.AssertConvexHuller([]d2.Pt{
		{0, 0}, {1, 0}, {1, 1}, {0, 1},
		{0.5, 0.5}, {0.25, 0.25}, {0.75, 0.25}, {0.75, 0.75}, {0.25, 0.75},
	})
	geomtest.Equal(t, a, b)
	for i, pt := range b.ConvexHull() {
		geomtest.Equal(t, b.Vertex(i), pt, i)
	}

	// malformed but allows testing of SignedArea
	b[0] = d2.Pt{0, 2}
	assert.Equal(t, 1.0, b.Area())
	assert.Equal(t, -1.0, b.SignedArea())
}

func TestBoxContains(t *testing.T) {
	ln := 100
	scale := 5.0
	pts := make([]d2.Pt, ln)
	for i := range pts {
		pts[i] = generate.Pt().Multiply(scale)
	}
	box := box.New(pts...)
	for _, p := range pts {
		assert.True(t, box.Contains(p))
	}
}

func TestLineIntersections(t *testing.T) {
	b := box.New(d2.Pt{0, 0}, d2.Pt{1, 1})
	tt := map[string]struct {
		expected []d2.Pt
		line.Line
	}{
		"corners": {
			expected: []d2.Pt{{0, 0}, {1, 1}},
			Line: line.Line{
				d2.Pt{0.1, 0.1},
				d2.V{0.2, 0.2},
			},
		},
		"horizontal": {
			expected: []d2.Pt{{1, 0.5}, {0, 0.5}},
			Line: line.Line{
				d2.Pt{0.5, 0.5},
				d2.V{0.2, 0},
			},
		},
		"vertical": {
			expected: []d2.Pt{{0.5, 0}, {0.5, 1}},
			Line: line.Line{
				d2.Pt{0.5, 0.5},
				d2.V{0, 1},
			},
		},
		"miss": {
			expected: nil,
			Line: line.Line{
				d2.Pt{2, 2},
				d2.V{1, -1},
			},
		},
		"single": {
			expected: []d2.Pt{{1, 1}},
			Line: line.Line{
				d2.Pt{1, 1},
				d2.V{1, -1},
			},
		},
	}

	li := line.Intersector(b)
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			is := li.LineIntersections(tc.Line, nil)
			assert.Len(t, is, len(tc.expected))
			for i, expected := range tc.expected {
				geomtest.Equal(t, expected, tc.Line.Pt1(is[i]))
			}

			if len(tc.expected) > 0 {
				one := li.LineIntersections(tc.Line, []float64{0})
				geomtest.Equal(t, tc.expected[0], tc.Line.Pt1(one[0]))
			}
		})
	}
}
