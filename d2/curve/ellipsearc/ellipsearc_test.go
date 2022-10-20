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
	geomtest.Equal(t, 0.0, e.Pt1(0).Y)

	// Test by definition of ellipse
	p0 := e.Pt1(0)
	d0 := f1.Distance(p0) + f2.Distance(p0)
	for i := 0.0; i <= 1.0; i += 0.2 {
		p := e.Pt1(i)
		d := f1.Distance(p) + f2.Distance(p)
		geomtest.Equal(t, d0, d)
	}

	// Get correct foci
	tf1, tf2 := e.Foci()
	geomtest.Equal(t, 0.0, f1.Distance(tf1))
	geomtest.Equal(t, 0.0, f2.Distance(tf2))

	geomtest.Equal(t, d2.Pt{1, 0}, e.Centroid())

	M, m := e.Axis()
	geomtest.Equal(t, math.Sqrt(2), M)
	geomtest.Equal(t, 1.0, m)

	a, s, c := e.Angle()
	geomtest.Equal(t, 0.0, a.Rad())
	geomtest.Equal(t, 0.0, s)
	geomtest.Equal(t, 1.0, c)
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
	geomtest.Equal(t, 0.0, e.Pt1(0.25).Distance(d2.Pt{0, 1}))
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
		geomtest.Equal(t, 0.0, math.Min(d1, d2))
	}
}

func TestV1(t *testing.T) {
	f1 := d2.Pt{0, 0}
	f2 := d2.Pt{2, 0}
	r := 1.0

	e := New(f1, f2, r)
	geomtest.Equal(t, d2.AssertV1{}, e)
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
			expected:   []float64{0.7236067977499789, 0.276393202250021},
		},
		"vertical": {
			Line:       line.New(d2.Pt{2, -3}, d2.Pt{2, 3}),
			EllipseArc: New(d2.Pt{2, 0}, d2.Pt{3, 0}, 1),
			expected:   []float64{0.649071198499986, 0.35092880150001404},
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
			expected:   []float64{0.5849901182297057, 0.3038987706591831},
		},
		"onePoint": {
			Line:       line.New(d2.Pt{0, 1}, d2.Pt{1, 1}),
			EllipseArc: New(d2.Pt{0, 0}, d2.Pt{0, 0}, 1),
			expected:   []float64{0},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			li := line.Intersector(tc.EllipseArc)
			is := li.LineIntersections(tc.Line, nil)
			geomtest.Equal(t, tc.expected, is)
			assert.Equal(t, len(tc.expected), len(is))

			if len(tc.expected) > 0 {
				is = li.LineIntersections(tc.Line, []float64{0})
				assert.Equal(t, tc.expected[:1], is)
			}
		})
	}
}
