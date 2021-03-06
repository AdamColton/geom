package ellipsearc

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestEllipseArc(t *testing.T) {
	f1 := d2.Pt{0, 0}
	f2 := d2.Pt{2, 0}
	r := 1.0

	e := New(f1, f2, r)
	assert.Equal(t, r, e.Pt1(0.25).Distance(e.c))
	assert.InDelta(t, 0.0, e.Pt1(0).Y, 1e-10)

	// Test by definition of ellipse
	p0 := e.Pt1(0)
	d0 := f1.Distance(p0) + f2.Distance(p0)
	for i := 0.0; i <= 1.0; i += 0.2 {
		p := e.Pt1(i)
		d := f1.Distance(p) + f2.Distance(p)
		assert.InDelta(t, d0, d, 1e-10)
	}

	// Get correct foci
	tf1, tf2 := e.Foci()
	assert.InDelta(t, 0, f1.Distance(tf1), 1e-10)
	assert.InDelta(t, 0, f2.Distance(tf2), 1e-10)

	geomtest.Equal(t, d2.Pt{1, 0}, e.Centroid())

	M, m := e.Axis()
	assert.InDelta(t, math.Sqrt(2), M, 1e-5)
	assert.InDelta(t, 1, m, 1e-5)

	a, s, c := e.Angle()
	assert.InDelta(t, 0, a.Rad(), 1e-5)
	assert.InDelta(t, 0, s, 1e-5)
	assert.InDelta(t, 1, c, 1e-5)
}

func TestEllipseStandard(t *testing.T) {
	// Make sure the ellipse follows standards of Polar plane
	f1 := d2.Pt{-1, 0}
	f2 := d2.Pt{1, 0}
	r := 1.0
	e := New(f1, f2, r)

	// An angle of 0 should be in the +X direction
	assert.Equal(t, d2.Pt{math.Sqrt2, 0}, e.Pt1(0))

	// 1/4 rotation should be +Y
	assert.InDelta(t, 0.0, e.Pt1(0.25).Distance(d2.Pt{0, 1}), 1e-10)
}

func TestBoundingBox(t *testing.T) {
	es := []*EllipseArc{
		New(d2.Pt{0, 0}, d2.Pt{1, 0}, 1),
		New(d2.Pt{0, 0}, d2.Pt{1, 1}, 1.23),
		New(d2.Pt{0, 0}, d2.Pt{0, 0}, 5),
		New(d2.Pt{3, 4}, d2.Pt{8, 2}, 12),
	}

	for _, e := range es {
		bb := polygon.RectangleTwoPoints(e.BoundingBox())
		for i := 0.0; i <= 1.0; i += 0.05 {
			assert.True(t, bb.Contains(e.Pt1(i)))
		}
	}
}

func TestCartesianArc(t *testing.T) {
	f1 := d2.Pt{3, 2}
	f2 := d2.Pt{2, 1}
	r := 1.0

	e := New(f1, f2, r)

	for t0 := 0.0; t0 < 1; t0 += 0.1 {
		pt := e.Pt1(t0)
		a, b, c := e.ABC(pt.Y)
		x1 := e.c.X + (-b+math.Sqrt(b*b-4*a*c))/(2*a)
		x2 := e.c.X + (-b-math.Sqrt(b*b-4*a*c))/(2*a)
		d1 := math.Abs(pt.X - x1)
		d2 := math.Abs(pt.X - x2)
		assert.InDelta(t, 0, math.Min(d1, d2), 1e-10)
	}
}

func TestV1(t *testing.T) {
	f1 := d2.Pt{0, 0}
	f2 := d2.Pt{2, 0}
	r := 1.0

	e := New(f1, f2, r)
	geomtest.V1(t, e)
}

func TestLineIntersections(t *testing.T) {
	tt := map[string]struct {
		line.Line
		*EllipseArc
		expected []float64
	}{
		"basic": {
			Line:       line.New(d2.Pt{0, 0}, d2.Pt{5, 0}),
			EllipseArc: New(d2.Pt{2, 0}, d2.Pt{3, 0}, 1),
			expected:   []float64{0.7236067, 0.276393},
		},
		"vertical": {
			Line:       line.New(d2.Pt{2, -3}, d2.Pt{2, 3}),
			EllipseArc: New(d2.Pt{2, 0}, d2.Pt{3, 0}, 1),
			expected:   []float64{0.649071, 0.350928},
		},
		"no-intersect": {
			Line:       line.New(d2.Pt{-10, 0}, d2.Pt{-11, 100}),
			EllipseArc: New(d2.Pt{2, 0}, d2.Pt{3, 0}, 1),
			expected:   nil,
		},
		"vertical-no-intersect": {
			Line:       line.New(d2.Pt{0, 0}, d2.Pt{0, 1}),
			EllipseArc: New(d2.Pt{2, 0}, d2.Pt{3, 0}, 1),
			expected:   nil,
		},
		"pointline": {
			Line:       line.New(d2.Pt{0, 0}, d2.Pt{0, 0}),
			EllipseArc: New(d2.Pt{2, 0}, d2.Pt{3, 0}, 1),
			expected:   nil,
		},
		"slopeline": {
			Line:       line.New(d2.Pt{0, 0}, d2.Pt{5, 5}),
			EllipseArc: New(d2.Pt{2, 2}, d2.Pt{3, 2}, 1),
			expected:   []float64{0.584990, 0.303898},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			_ = n
			is := tc.EllipseArc.LineIntersections(tc.Line)
			assert.Equal(t, len(tc.expected), len(is))
			for idx, ti := range tc.expected {
				assert.InDelta(t, ti, is[idx], 1e-5)
			}
		})
	}
}
