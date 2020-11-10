package line

import (
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	line := New(d2.Pt{1, 3}, d2.Pt{2, 4})
	geomtest.Equal(t, d2.Pt{1, 3}, line.Pt1(0))
	geomtest.Equal(t, d2.Pt{2, 4}, line.Pt1(1))
	assert.Equal(t, -1.0, line.AtX(0))
	assert.Equal(t, -3.0, line.AtY(0))
	assert.Equal(t, 2.0, line.B())
	assert.Equal(t, 1.0, line.M())
	assert.Equal(t, "Line( V(1.0000, 1.0000)t + Pt(1.0000, 3.0000) )", line.String())
}

func TestBisect(t *testing.T) {
	a, b := d2.Pt{-1, 0}, d2.Pt{1, 0}
	line := Bisect(a, b)
	geomtest.Equal(t, d2.Pt{0, 0}, line.Pt1(0))
	geomtest.Equal(t, d2.Pt{0, 1}, line.Pt1(1))

	tt := [][2]d2.Pt{
		{{-1, 0}, {1, 0}},
		{{1, 1}, {2, 3}},
	}

	for _, tc := range tt {
		line := Bisect(tc[0], tc[1])
		for i := -1.0; i < 2.0; i += 0.05 {
			pt := line.Pt1(i)
			assert.InDelta(t, tc[0].Distance(pt), tc[1].Distance(pt), 1e-10)
		}
	}
}

func TestLineIntersect(t *testing.T) {
	_ = Intersector(Line{})

	testCases := map[string]struct {
		points    []d2.Pt
		expectNil bool
	}{
		"normal line": {
			points: []d2.Pt{{0, 1}, {1, 2}, {1, 0}, {2, 3}},
		},
		"first line is vertical": {
			points: []d2.Pt{{0, 1}, {0, 2}, {1, 0}, {2, 3}},
		},
		"second line is vertical": {
			points: []d2.Pt{{0, 1}, {1, 2}, {1, 0}, {1, 3}},
		},
		"lines are parallel": {
			points:    []d2.Pt{{0, 1}, {1, 2}, {1, 2}, {2, 3}},
			expectNil: true,
		},
		"lines are parallel but reversed": {
			points:    []d2.Pt{{0, 1}, {1, 2}, {2, 3}, {1, 2}},
			expectNil: true,
		},
		"first line is a point": {
			points:    []d2.Pt{{0, 1}, {0, 1}, {1, 2}, {2, 3}},
			expectNil: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			l0 := New(tc.points[0], tc.points[1])
			l1 := New(tc.points[2], tc.points[3])
			i1 := l0.LineIntersections(l1, nil)
			i0 := l1.LineIntersections(l0, nil)
			t0, t1, ok := l1.Intersection(l0)
			assert.Equal(t, tc.expectNil, !ok)
			if tc.expectNil {
				assert.Nil(t, i0)
				assert.Nil(t, i1)
			} else {
				assert.Equal(t, i0[0], t0)
				assert.Equal(t, i1[0], t1)
				geomtest.Equal(t, l0.Pt1(t0), l1.Pt1(t1))
			}
		})
	}
}

func TestClosest(t *testing.T) {
	l := New(d2.Pt{1, 2}, d2.Pt{1, 1})
	p := l.Closest(d2.Pt{0, 0})
	geomtest.Equal(t, d2.Pt{1, 0}, p)

	l = New(d2.Pt{0, 0}, d2.Pt{1, 1})
	p = l.Closest(d2.Pt{0, 1})
	geomtest.Equal(t, d2.Pt{0.5, 0.5}, p)

	l = New(d2.Pt{0, 0}, d2.Pt{1, 3})
	p = l.Closest(d2.Pt{-3, 1})
	geomtest.Equal(t, d2.Pt{0, 0}, p)

	l = New(d2.Pt{0, 0}, d2.Pt{0, 0})
	p = l.Closest(d2.Pt{-3, 1})
	geomtest.Equal(t, d2.Pt{0, 0}, p)
}

func TestLineFulfillsX1(t *testing.T) {
	var dc d2.Pt1V1
	l := New(d2.Pt{1, 2}, d2.Pt{1, 1})
	dc = l
	assert.NotNil(t, dc)
}

func TestV1(t *testing.T) {
	geomtest.Equal(t, d2.AssertV1{}, New(d2.Pt{0, 0}, d2.Pt{1, 3}))
}

type mockPt1V1 struct{}

func (mockPt1V1) Pt1(t0 float64) d2.Pt {
	return d2.Pt{2 * t0, t0 * t0}
}

func (mockPt1V1) V1(t0 float64) d2.V {
	return d2.V{2, 2 * t0}
}

func TestTangentLine(t *testing.T) {
	geomtest.Equal(t, d2.AssertV1{}, mockPt1V1{}) // confirm that the mock is valid
	l := TangentLine(mockPt1V1{}, .5)
	geomtest.Equal(t, Line{d2.Pt{1, .25}, d2.V{2, 1}}, l)
}

func TestLimits(t *testing.T) {
	l := Line{}
	assert.Equal(t, d2.LimitUnbounded, l.L(1, 1))
	assert.Equal(t, d2.LimitUndefined, l.L(1, 0))
	assert.Equal(t, d2.LimitUndefined, l.L(2, 2))
	assert.Equal(t, d2.LimitUnbounded, l.VL(1, 1))
	assert.Equal(t, d2.LimitUndefined, l.VL(1, 0))
}

func TestTransform(t *testing.T) {
	l := New(d2.Pt{0, 0}, d2.Pt{1, 0}).T(d2.Rotate(angle.Deg(90)).T())
	geomtest.Equal(t, d2.Pt{0, 0}, l.Pt1(0))
	geomtest.Equal(t, d2.Pt{0, 1}, l.Pt1(1))
}

func TestCentroid(t *testing.T) {
	geomtest.Equal(t, d2.Pt{1, 1}, New(d2.Pt{0, 0}, d2.Pt{2, 2}).Centroid())
}

func TestCross(t *testing.T) {
	tt := map[string]struct {
		Line
		d2.Pt
		expected float64
	}{
		"basic": {
			Line:     New(d2.Pt{0, 0}, d2.Pt{1, 0}),
			Pt:       d2.Pt{1, 1},
			expected: 1,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.Line.Cross(tc.Pt))
		})
	}
}
