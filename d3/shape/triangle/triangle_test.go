package triangle

import (
	"testing"

	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/stretchr/testify/assert"
)

func TestTriangleIntersection(t *testing.T) {
	tt := map[string]struct {
		t        *Triangle
		l        line.Line
		expected Intersection
	}{
		"basic": {
			t: &Triangle{
				{0, 0, 0},
				{1, 0, 0},
				{0, 1, 0},
			},
			l: line.New(d3.Pt{0.25, 0.25, 2}, d3.Pt{0.25, 0.25, 1}),
			expected: Intersection{
				U:    0.25,
				V:    0.25,
				T:    2,
				Does: true,
			},
		},
		"parallel": {
			t: &Triangle{
				{0, 0, 0},
				{1, 0, 0},
				{0, 1, 0},
			},
			l: line.New(d3.Pt{0, 0, 1}, d3.Pt{1, 1, 1}),
			expected: Intersection{
				U:    0,
				V:    0,
				T:    0,
				Does: false,
			},
		},
		"u-outside": {
			t: &Triangle{
				{0, 0, 0},
				{1, 0, 0},
				{0, 1, 0},
			},
			l: line.New(d3.Pt{2, 0.5, 0}, d3.Pt{2, 0.5, 1}),
			expected: Intersection{
				U:    2,
				V:    0,
				T:    0,
				Does: false,
			},
		},
		"v-outside": {
			t: &Triangle{
				{0, 0, 0},
				{1, 0, 0},
				{0, 1, 0},
			},
			l: line.New(d3.Pt{0.5, 2, 0}, d3.Pt{0.5, 2, 1}),
			expected: Intersection{
				U:    0.5,
				V:    2,
				T:    0,
				Does: false,
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			ti := tc.t.Intersector()
			assert.Equal(t, tc.expected, tc.t.Intersection(tc.l))
			assert.Equal(t, tc.expected, ti.Intersection(tc.l))
			assert.Equal(t, tc.expected, ti.RawIntersection(tc.l))
			tHit, hit := tc.t.LineIntersection(tc.l)
			assert.Equal(t, tc.expected.T, tHit)
			assert.Equal(t, tc.expected.Does, hit)
			if tc.expected.Does {
				assert.Equal(t, []float64{tc.expected.T}, tc.t.Intersections(tc.l))
			} else {
				assert.Nil(t, tc.t.Intersections(tc.l))
			}
		})
	}
}

func TestBT(t *testing.T) {
	tt := map[string]struct {
		t        *Triangle
		o, u     int
		b        barycentric.B
		expected d3.Pt
	}{
		"Basic": {
			t:        &Triangle{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}},
			o:        0,
			u:        1,
			b:        barycentric.B{0.5, 0.5},
			expected: d3.Pt{0.5, 0.5, 0},
		},
		"Complex": {
			t:        &Triangle{{1, 2, 3}, {2, 4, 6}, {3, 6, 9}},
			o:        0,
			u:        1,
			b:        barycentric.B{0.1, 0.2},
			expected: d3.Pt{X: 1.5, Y: 3, Z: 4.5},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			bt := tc.t.BT(tc.o, tc.u)
			assert.Equal(t, tc.expected, bt.PtB(tc.b))
			assert.Equal(t, tc.t, bt.Triangle())
		})
	}
}

func TestValidate(t *testing.T) {
	tri := &Triangle{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
	}
	assert.NoError(t, tri.Validate())

	tri = &Triangle{
		{0, 0, 0},
		{1, 0, 0},
		{0, 0, 0},
	}
	assert.Equal(t, "At least 2 verticies have the same value", tri.Validate().Error())
}

func TestInvalidBT(t *testing.T) {
	tt := map[string]struct {
		o, u int
	}{
		"o<0": {
			o: -1,
			u: 0,
		},
		"o>2": {
			o: 3,
			u: 0,
		},
		"u<0": {
			o: 1,
			u: -1,
		},
		"u>2": {
			o: 1,
			u: 5,
		},
		"o==u": {
			o: 1,
			u: 1,
		},
	}
	tri := &Triangle{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
	}
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Nil(t, tri.BT(tc.o, tc.u))
		})
	}
}

func TestNormal(t *testing.T) {
	tri := &Triangle{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
	}
	assert.Equal(t, d3.V{0, 0, 1}, tri.Normal())
}
