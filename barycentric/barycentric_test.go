package barycentric_test

import (
	"testing"

	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/geomtest"
	"github.com/testify/assert"
)

func TestBIterator(t *testing.T) {
	bi := &barycentric.BIterator{
		Origin: 0,
		U:      1,
		Step: [2]barycentric.B{
			{.1, 0},
			{0, .1},
		},
	}

	expected := []barycentric.B{
		{0.00, 0.00},
		{0.10, 0.00},
		{0.20, 0.00},
		{0.30, 0.00},
		{0.40, 0.00},
		{0.50, 0.00},
		{0.60, 0.00},
		{0.70, 0.00},
		{0.80, 0.00},
		{0.90, 0.00},
		{1.00, 0.00},
		{0.00, 0.10},
		{0.10, 0.10},
		{0.20, 0.10},
		{0.30, 0.10},
		{0.40, 0.10},
		{0.50, 0.10},
		{0.60, 0.10},
		{0.70, 0.10},
		{0.80, 0.10},
		{0.90, 0.10},
		{0.00, 0.20},
		{0.10, 0.20},
		{0.20, 0.20},
		{0.30, 0.20},
		{0.40, 0.20},
		{0.50, 0.20},
		{0.60, 0.20},
		{0.70, 0.20},
		{0.80, 0.20},
		{0.00, 0.30},
		{0.10, 0.30},
		{0.20, 0.30},
		{0.30, 0.30},
		{0.40, 0.30},
		{0.50, 0.30},
		{0.60, 0.30},
		{0.70, 0.30},
		{0.00, 0.40},
		{0.10, 0.40},
		{0.20, 0.40},
		{0.30, 0.40},
		{0.40, 0.40},
		{0.50, 0.40},
		{0.60, 0.40},
		{0.00, 0.50},
		{0.10, 0.50},
		{0.20, 0.50},
		{0.30, 0.50},
		{0.40, 0.50},
		{0.50, 0.50},
		{0.00, 0.60},
		{0.10, 0.60},
		{0.20, 0.60},
		{0.30, 0.60},
		{0.40, 0.60},
		{0.00, 0.70},
		{0.10, 0.70},
		{0.20, 0.70},
		{0.30, 0.70},
		{0.00, 0.80},
		{0.10, 0.80},
		{0.20, 0.80},
		{0.00, 0.90},
		{0.10, 0.90},
		{0.00, 1.00},
	}

	for b, done := bi.Start(); !done; b, done = bi.Next() {
		geomtest.Equal(t, expected[bi.Idx], b)
	}

	bi = nil
	_, done := bi.Start()
	assert.True(t, done)
}

func TestEdge(t *testing.T) {
	tt := map[string]struct {
		barycentric.B
		d        float64
		expected bool
	}{
		"Origin": {
			B:        barycentric.B{0, 0},
			d:        1e-10,
			expected: true,
		},
		"Origin-offset-true": {
			B:        barycentric.B{-0.01, 0.01},
			d:        1e-1,
			expected: true,
		},
		"Origin-offset-false": {
			B:        barycentric.B{-0.01, 0.01},
			d:        1e-10,
			expected: false,
		},
		"middle-u": {
			B:        barycentric.B{0.5, 0},
			d:        1e-10,
			expected: true,
		},
		"middle-v": {
			B:        barycentric.B{0, 0.5},
			d:        1e-10,
			expected: true,
		},
		"middle-w": {
			B:        barycentric.B{0.5, 0.5},
			d:        1e-10,
			expected: true,
		},
		"middle": {
			B:        barycentric.B{0.25, 0.25},
			d:        1e-10,
			expected: false,
		},
		"wide-edge": {
			B:        barycentric.B{1.1, 0},
			d:        0.2,
			expected: true,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.B.Edge(tc.d))
		})
	}
}

func TestString(t *testing.T) {
	b := barycentric.B{.3, .1}
	assert.Equal(t, "B(0.30, 0.10)", b.String())
}

func TestV(t *testing.T) {
	tt := map[string]struct {
		origin, u int
		expected  int
	}{
		"o:0,u:1": {
			origin:   0,
			u:        1,
			expected: 2,
		},
		"o:0,u:2": {
			origin:   0,
			u:        2,
			expected: 1,
		},
		"o:1,u:0": {
			origin:   1,
			u:        0,
			expected: 2,
		},
		"o:1,u:2": {
			origin:   1,
			u:        2,
			expected: 0,
		},
		"o:2,u:0": {
			origin:   2,
			u:        0,
			expected: 1,
		},
		"o:2,u:1": {
			origin:   2,
			u:        1,
			expected: 0,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, (&barycentric.BIterator{Origin: tc.origin, U: tc.u}).V())
		})
	}
}
