package triangle

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/angle"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	tt := []struct {
		T1, T2 *Triangle
	}{
		{
			T1: &Triangle{
				{1, 1},
				{2, 1},
				{2, 2},
			},
			T2: &Triangle{
				{2, 1},
				{3, 1},
				{3, 2},
			},
		},
		{
			T1: &Triangle{
				{0, 0},
				{1, 0},
				{0, 1},
			},
			T2: &Triangle{
				{0, 1},
				{0, 0},
				{1, 0},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.T1.String()+tc.T2.String(), func(t *testing.T) {
			tfrm, err := Transform(tc.T1, tc.T2)
			assert.NoError(t, err)
			for i, pt := range tc.T1 {
				pt = tfrm.Pt(pt)
				assert.Equal(t, tc.T2[i], pt)
			}
		})
	}
}

func TestPt2(t *testing.T) {
	tri := &Triangle{{0, 0}, {1, 0}, {0, 1}}
	assert.Equal(t, tri[0], tri.Pt2(0, 0))

	i := grid.Pt{1, 1}.To(grid.Pt{10, 10})
	for i, ok := i.Start(); ok; ok = i.Next() {
		p := i.Pt().D2().Pt().Multiply(1.0 / 10.0)
		assert.True(t, tri.Contains(tri.Pt2(p.X, p.Y)))
	}

	var pt2 d2.Pt2
	pt2 = tri
	assert.NotNil(t, pt2)
}

func TestPt1(t *testing.T) {
	tri := &Triangle{{0, 0}, {1, 0}, {0, 1}}
	assert.Equal(t, tri[0], tri.Pt1(0))
	assert.Equal(t, tri[1], tri.Pt1(1.0/3.0))
	assert.Equal(t, tri[2], tri.Pt1(2.0/3.0))
	assert.Equal(t, tri[0], tri.Pt1(1))
}

func TestCircumCenter(t *testing.T) {
	testCases := []struct {
		center d2.Pt
		m      float64
		angles []angle.Rad
	}{
		{
			center: d2.Pt{2, 3},
			m:      2.0,
			angles: []angle.Rad{1, 2, 3},
		},
		{
			center: d2.Pt{1, 4},
			m:      1.23,
			angles: []angle.Rad{0, math.Pi, 3},
		},
	}

	for _, tc := range testCases {
		tr := &Triangle{
			tc.center.Add(d2.Polar{tc.m, tc.angles[0]}.V()),
			tc.center.Add(d2.Polar{tc.m, tc.angles[1]}.V()),
			tc.center.Add(d2.Polar{tc.m, tc.angles[2]}.V()),
		}

		assert.InDelta(t, 0, tc.center.Distance(tr.CircumCenter()), 1e-10)
	}

	tr := Triangle{}

	assert.InDelta(t, 0, d2.Pt{0, 0}.Distance(tr.CircumCenter()), 1e-10)

	tr[2] = d2.Pt{2, 2}
	assert.InDelta(t, 0, d2.Pt{1, 1}.Distance(tr.CircumCenter()), 1e-10)
}

func TestContains(t *testing.T) {
	tt := map[string]struct {
		*Triangle
		inside, outside []d2.Pt
	}{
		"Basic": {
			Triangle: &Triangle{{0, 0}, {1, 0}, {0, 1}},
			inside:   []d2.Pt{{.1, .1}, {.4, .4}, {.8, .1}, {.1, .8}},
			outside:  []d2.Pt{{-1, -1}, {.6, .6}},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			for _, i := range tc.inside {
				assert.True(t, tc.Contains(i))
			}
			for _, o := range tc.outside {
				assert.False(t, tc.Contains(o))
			}
		})
	}
}

func TestArea(t *testing.T) {
	tri := &Triangle{{0, 0}, {0, 1}, {1, 0}}
	assert.Equal(t, 0.5, tri.Area())
	assert.Equal(t, -0.5, tri.SignedArea())
}

func TestPerimeter(t *testing.T) {
	tri := &Triangle{{0, 0}, {0, 1}, {1, 0}}
	assert.Equal(t, 2.0+math.Sqrt2, tri.Perimeter())
}

func TestCentroid(t *testing.T) {
	tri := &Triangle{{0, 0}, {0, 1}, {1, 0}}
	assert.Equal(t, d2.Pt{1.0 / 3.0, 1.0 / 3.0}, tri.Centroid())
}

func TestLimit(t *testing.T) {
	tri := &Triangle{}
	assert.Equal(t, d2.LimitBounded, tri.L(1, 0))
	assert.Equal(t, d2.LimitUndefined, tri.L(2, 0))
}

func TestIntersections(t *testing.T) {
	tri := &Triangle{{0, 0}, {0, 1}, {1, 0}}
	l := line.New(d2.Pt{0, .1}, d2.Pt{1, .1})
	expected := []float64{0, .9}
	_ = line.Intersector(tri)
	assert.Equal(t, expected, tri.LineIntersections(l, nil))
	assert.Equal(t, expected[:1], tri.LineIntersections(l, []float64{0}))
}

func TestBoundingBox(t *testing.T) {
	tri := &Triangle{{0, 0}, {0, 1}, {1, 0}}
	m, M := tri.BoundingBox()
	assert.Equal(t, d2.Pt{0, 0}, m)
	assert.Equal(t, d2.Pt{1, 1}, M)
}
