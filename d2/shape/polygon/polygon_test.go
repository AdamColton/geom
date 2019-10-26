package polygon

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestPolygonSignedArea(t *testing.T) {
	tri := triangle.Triangle{
		{0, 0},
		{1, 0},
		{0, 1},
	}
	p := Polygon(tri[:])
	assert.InDelta(t, tri.SignedArea(), p.SignedArea(), 1e-10)

	square := Polygon{
		d2.Pt{0, 0},
		d2.Pt{1, 0},
		d2.Pt{1, 1},
		d2.Pt{0, 1},
	}
	assert.InDelta(t, 1.0, square.SignedArea(), 1e-10)
}

func TestPolygonCentroid(t *testing.T) {
	square := Polygon{
		d2.Pt{0, 0},
		d2.Pt{1, 0},
		d2.Pt{1, 1},
		d2.Pt{0, 1},
	}
	assert.Equal(t, d2.Pt{0.5, 0.5}, square.Centroid())
}

func TestContains(t *testing.T) {
	tt := map[string]struct {
		Polygon
		does, doesnt []d2.Pt
	}{
		"arrowhead": {
			Polygon: Polygon{
				d2.Pt{0, 0},
				d2.Pt{2, 2},
				d2.Pt{1, 0},
				d2.Pt{2, -2},
			},
			does:   []d2.Pt{{0.5, 0}, {1, 1}},
			doesnt: []d2.Pt{{0.5, 1}, {2, 0}, {1.5, 0}},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			ll := NewLL(tc.Polygon)
			for _, p := range tc.does {
				assert.True(t, tc.Contains(p))
				assert.True(t, ll.Contains(p))
			}
			for _, p := range tc.doesnt {
				assert.False(t, tc.Contains(p))
				assert.False(t, ll.Contains(p))
			}
		})
	}
}

func TestPolygonSurface(t *testing.T) {
	i := grid.Pt{20, 20}.Iter()
	s := grid.Scale{0.05, 0.05, 0.025, 0.025}
	for sides := 3; sides < 7; sides++ {
		p := RegularPolygonSideLength(d2.Pt{}, 1, 0, sides)
		// Note that we're iterating over interior points only, not perimeter points
		// where t0 or t1 = 0.0 or 1.0.
		for _, ok := i.Start(); ok; ok = i.Next() {
			pt := p.Pt2(s.T(i.Pt()))
			assert.True(t, p.Contains(pt))
		}
	}
}

func TestUnitSquareSurface(t *testing.T) {
	// The surface function of the unit square should map to itself
	unitSquare := Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	s := grid.Scale{0.01, 0.01, 0, 0}
	grid.Pt{101, 101}.Iter().Each(func(idx int, gpt grid.Pt) {
		t0, t1 := s.T(gpt)
		geomtest.Equal(t, d2.Pt{t0, t1}, unitSquare.Pt2(t0, t1))
	})
}

func TestRegularPolygonSideLength(t *testing.T) {
	actual := RegularPolygonSideLength(d2.Pt{}, 1, 0, 4)
	expected := Polygon{d2.Pt{0.5000, -0.5000}, {0.5000, 0.5000}, {-0.5000, 0.5000}, {-0.5000, -0.5000}}
	for i, p := range expected {
		geomtest.Equal(t, p, actual[i])
	}
}

func TestRegularPolygonRadius(t *testing.T) {
	actual := RegularPolygonRadius(d2.Pt{}, 1, 0, 4)
	expected := Polygon{d2.Pt{1.0000, 0.0000}, {0.0000, 1.0000}, {-1.0000, 0.0000}, {-0.0000, -1.0000}}
	for i, p := range expected {
		geomtest.Equal(t, p, actual[i])
	}
}

func TestCountAngles(t *testing.T) {
	p := Polygon{{0, 1}, {0.5, 0.5}, {0, 0}, {1, 0.5}}
	got := make([]int, len(p))
	ccw, cw := p.CountAngles(got)
	assert.Equal(t, 3, ccw)
	assert.Equal(t, 1, cw)

	expected := []int{0, 1, 3, 2}
	assert.Equal(t, expected, got)
}

func TestPerimeter(t *testing.T) {
	p := Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	assert.Equal(t, 4.0, p.Perimeter())
}

func TestConvex(t *testing.T) {
	tt := []struct {
		Polygon Polygon
		Convex  bool
	}{
		{
			Polygon: Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			Convex:  true,
		},
		{
			Polygon: Polygon{{0, 0}, {0, 1}, {1, 1}, {1, 0}},
			Convex:  true,
		},
		{
			Polygon: Polygon{{0, 1}, {2, 2}, {1, 1}, {2, 0}},
			Convex:  false,
		},
		{
			Polygon: Polygon{{0, 0}, {1, 1}, {0, 2}, {2, 1}},
			Convex:  false,
		},
		{
			Polygon: Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0.5, 0.5}},
			Convex:  false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Polygon.String(), func(t *testing.T) {
			assert.Equal(t, tc.Convex, tc.Polygon.Convex())
		})
	}
}

func TestNonIntersecting(t *testing.T) {
	tt := []struct {
		Polygon         Polygon
		NonIntersecting bool
	}{
		{
			Polygon:         Polygon{{0, 0}, {1, 0}, {1, 1}},
			NonIntersecting: true,
		},
		{
			Polygon:         Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			NonIntersecting: true,
		},
		{
			Polygon:         Polygon{{0, 0}, {1, 1}, {1, 0}, {0, 1}},
			NonIntersecting: false,
		},
		{
			Polygon:         Polygon{{0, 1}, {0, 2}, {1, 0}, {2, 2}, {2, 1}},
			NonIntersecting: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Polygon.String(), func(t *testing.T) {
			assert.Equal(t, tc.NonIntersecting, tc.Polygon.NonIntersecting())
		})
	}
}

func TestReverse(t *testing.T) {
	tt := []struct {
		p        Polygon
		expected Polygon
	}{
		{
			p:        Polygon{{0, 0}, {1, 0}, {1, 1}},
			expected: Polygon{{1, 1}, {1, 0}, {0, 0}},
		},
		{
			p:        Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			expected: Polygon{{0, 1}, {1, 1}, {1, 0}, {0, 0}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.p.String(), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.p.Reverse())
		})
	}
}

func TestFindTriangles(t *testing.T) {
	tt := map[string]struct {
		Polygon
		expected [][3]uint32
	}{
		"simple": {
			Polygon: Polygon{{0, 0}, {2, 1}, {0, 2}, {1, 1}},
			expected: [][3]uint32{
				{1, 2, 3},
				{3, 0, 1},
			},
		},
		"heart": {
			Polygon: Polygon{
				{0, 0},
				{1, 1},
				{2, 1},
				{3, 0},
				{0, -3},
				{-3, 0},
				{-2, 1},
				{-1, 1},
			},
			expected: [][3]uint32{
				{1, 2, 3},
				{4, 5, 6},
				{6, 7, 0},
				{0, 1, 3},
				{4, 6, 0},
				{0, 3, 4},
			},
		},
		"arrow": {
			Polygon: Polygon{
				{0, 2},
				{1.5, 3.5},
				{3, 2},
				{2, 2},
				{2, 0},
				{1, 0},
				{1, 2},
			},
			expected: [][3]uint32{
				{1, 2, 3},
				{3, 4, 5},
				{6, 0, 1},
				{1, 3, 5},
				{5, 6, 1},
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.FindTriangles())
		})
	}
}

func TestPolygonPt1(t *testing.T) {
	p := Polygon{
		d2.Pt{0, 0},
		d2.Pt{2, 2},
		d2.Pt{1, 0},
		d2.Pt{2, -2},
	}

	geomtest.Equal(t, p[0], p.Pt1(0))
	geomtest.Equal(t, p[0], p.Pt1(1))
	geomtest.Equal(t, p[1], p.Pt1(0.25))
	geomtest.Equal(t, p[2], p.Pt1(0.5))
	geomtest.Equal(t, p[3], p.Pt1(0.75))

	geomtest.Equal(t, line.New(p[0], p[1]).Pt1(0.5), p.Pt1(0.125))
	geomtest.Equal(t, line.New(p[1], p[2]).Pt1(0.5), p.Pt1(0.375))
	geomtest.Equal(t, line.New(p[2], p[3]).Pt1(0.5), p.Pt1(0.625))
	geomtest.Equal(t, line.New(p[3], p[0]).Pt1(0.5), p.Pt1(0.875))

	geomtest.Equal(t, p.Pt1(0.25), p.Pt1(-0.75))
	geomtest.Equal(t, p.Pt1(0.25), p.Pt1(-1.75))
	geomtest.Equal(t, p.Pt1(0.25), p.Pt1(10.25))
}

func TestIntersections(t *testing.T) {
	tt := map[string]struct {
		Polygon
		line.Line
		expected []float64
	}{
		"basic": {
			Polygon:  Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			Line:     line.New(d2.Pt{-1, 0.5}, d2.Pt{3, 0.5}),
			expected: []float64{0.25, 0.5},
		},
		"no-intersects": {
			Polygon:  Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			Line:     line.New(d2.Pt{-1, 1.5}, d2.Pt{3, 1.5}),
			expected: nil,
		},
		"corners": {
			Polygon:  Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			Line:     line.New(d2.Pt{-1, -1}, d2.Pt{3, 3}),
			expected: []float64{0.25, 0.5},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			is := tc.Polygon.Intersections(tc.Line)
			assert.Equal(t, tc.expected, is)
			if len(is) > 0 {
				assert.True(t, NewLL(tc.Polygon).DoesIntersect(tc.Line))
			} else {
				assert.False(t, NewLL(tc.Polygon).DoesIntersect(tc.Line))
			}
		})
	}
}

func TestCollision(t *testing.T) {
	tt := map[string]struct {
		Polygon
		line.Line
		explineT float64
		expIdx   int
		expSideT float64
	}{
		"basic": {
			Polygon:  Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			Line:     line.New(d2.Pt{-1, 0.5}, d2.Pt{3, 0.5}),
			explineT: 0.25,
			expIdx:   3,
			expSideT: 0.5,
		},
		"basic-line-reverse": {
			Polygon:  Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			Line:     line.New(d2.Pt{3, 0.5}, d2.Pt{-1, 0.5}),
			explineT: 0.5,
			expIdx:   1,
			expSideT: 0.5,
		},
		"no-intersects": {
			Polygon:  Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			Line:     line.New(d2.Pt{-1, 1.5}, d2.Pt{3, 1.5}),
			explineT: 0.0,
			expIdx:   -1,
			expSideT: 0.0,
		},
		"corners": {
			Polygon:  Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			Line:     line.New(d2.Pt{-1, -1}, d2.Pt{3, 3}),
			explineT: 0.25,
			expIdx:   0,
			expSideT: 0,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			lt, idx, side := tc.Polygon.Collision(tc.Line)
			assert.Equal(t, tc.explineT, lt)
			assert.Equal(t, tc.expIdx, idx)
			assert.Equal(t, tc.expSideT, side)
		})
	}
}
