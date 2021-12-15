package box

import (
	"testing"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/adamcolton/geom/geomerr"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestAssert(t *testing.T) {
	a := &Box{
		{0, 0, 0},
		{1, 1, 1},
	}
	b := &Box{
		{0, 0, 0},
		{1, 1, 1},
	}

	err := geomtest.AssertEqual(a, b, 1e-10)
	assert.NoError(t, err)

	err = geomtest.AssertEqual(a, 1.0, 1e-10)
	assert.IsType(t, geomerr.ErrTypeMismatch{}, err)
}

func TestBox(t *testing.T) {
	bx := New(d3.Pt{0, 1, 0.5}, d3.Pt{1, 0.5, 0}, d3.Pt{0.5, 0, 1})
	expected := &Box{
		{0, 0, 0},
		{1, 1, 1},
	}

	geomtest.Equal(t, expected, bx)

	bx.Add(d3.Pt{-1, 0, 2})
	expected = &Box{
		{-1, 0, 0},
		{1, 1, 2},
	}
	geomtest.Equal(t, expected, bx)
}

func TestLineIntersection(t *testing.T) {
	bx := New(d3.Pt{0, 0, 0}, d3.Pt{1, 1, 1})
	tt := map[string]struct {
		line.Line
		expected float64
	}{
		"Z0": {
			Line:     line.New(d3.Pt{0.5, 0.5, -1}, d3.Pt{0.5, 0.5, 2}),
			expected: 1.0 / 3.0,
		},
		"Z1": {
			Line:     line.New(d3.Pt{0.5, 0.5, 2}, d3.Pt{0.5, 0.5, -1}),
			expected: 1.0 / 3.0,
		},
		"Z-angle": {
			Line:     line.New(d3.Pt{0.5, 0.5, -1}, d3.Pt{0.75, 0.75, 2}),
			expected: 1.0 / 3.0,
		},
		"X0": {
			Line:     line.New(d3.Pt{-1, 0.5, 0.5}, d3.Pt{2, 0.5, 0.5}),
			expected: 1.0 / 3.0,
		},
		"X1": {
			Line:     line.New(d3.Pt{2, 0.5, 0.5}, d3.Pt{-1, 0.5, 0.5}),
			expected: 1.0 / 3.0,
		},
		"X-angle": {
			Line:     line.New(d3.Pt{2, 0.5, 0.5}, d3.Pt{-1, 0.75, 0.75}),
			expected: 1.0 / 3.0,
		},
		"Y0": {
			Line:     line.New(d3.Pt{0.5, -1, 0.5}, d3.Pt{0.5, 2, 0.5}),
			expected: 1.0 / 3.0,
		},
		"Y1": {
			Line:     line.New(d3.Pt{0.5, 2, 0.5}, d3.Pt{0.5, -1, 0.5}),
			expected: 1.0 / 3.0,
		},
		"Y-angle": {
			Line:     line.New(d3.Pt{0.5, 2, 0.5}, d3.Pt{0.75, -1, 0.75}),
			expected: 1.0 / 3.0,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			i, _ := bx.LineIntersection(tc.Line)
			assert.Equal(t, tc.expected, i)
		})
	}
}
