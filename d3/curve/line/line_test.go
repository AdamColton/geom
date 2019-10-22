package line

import (
	"testing"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestLine(t *testing.T) {
	tt := [][2]d3.Pt{
		{d3.Pt{0, 0, 0}, d3.Pt{1, 1, 1}},
		{d3.Pt{0, 0, 0}, d3.Pt{0, 0, 0}},
	}

	for _, tc := range tt {
		t.Run(tc[0].String()+tc[1].String(), func(t *testing.T) {
			l := New(tc[0], tc[1])
			assert.Equal(t, tc[0], l.Pt1(0))
			assert.Equal(t, tc[1], l.Pt1(1))
		})
	}

	geomtest.Equal(t, d3.V{1, 1, 1}, New(d3.Pt{1, 2, 3}, d3.Pt{2, 3, 4}).V1(0))
}

func TestClosest(t *testing.T) {
	tt := []struct {
		l1, l2               Line
		expected1, expected2 float64
	}{
		{
			l1:        New(d3.Pt{0, 1, 0}, d3.Pt{0, 0, 0}),
			l2:        New(d3.Pt{1, 0, 1}, d3.Pt{0, 0, 1}),
			expected1: 1,
			expected2: 1,
		},
		{
			l1:        New(d3.Pt{0, 1, 0}, d3.Pt{0, 0, 0}),
			l2:        New(d3.Pt{1, 0, 0}, d3.Pt{0, 0, 0}),
			expected1: 1,
			expected2: 1,
		},
		{
			l1:        New(d3.Pt{1, 0.01, 0}, d3.Pt{0, 0, 0}),
			l2:        New(d3.Pt{1, 0, 0}, d3.Pt{0, 0, 0}),
			expected1: 1,
			expected2: 1,
		},
		{
			l1:        New(d3.Pt{-1, 0, 0}, d3.Pt{1, 0, 0}),
			l2:        New(d3.Pt{0, -1, 0}, d3.Pt{0, 3, 0}),
			expected1: .5,
			expected2: .25,
		},
		{
			l1:        New(d3.Pt{0, 0.1, 0}, d3.Pt{1, 0.1, 0}),
			l2:        New(d3.Pt{0.1, 0, 0}, d3.Pt{0.1, 1, 0}),
			expected1: .1,
			expected2: .1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.l1.String()+tc.l2.String(), func(t *testing.T) {
			c1, c2 := tc.l1.Closest(tc.l2)
			assert.InDelta(t, tc.expected1, c1, 1e-6)
			assert.InDelta(t, tc.expected2, c2, 1e-6)
		})
	}
}

func TestClosestEarlyReturns(t *testing.T) {
	l0 := New(d3.Pt{1, 1, 1}, d3.Pt{2, 2, 2})
	l1 := New(d3.Pt{}, d3.Pt{})

	t0, t1 := l0.Closest(l1)
	assert.Equal(t, 0.0, t0)
	assert.Equal(t, 0.0, t1)

	t0, t1 = l1.Closest(l0)
	assert.Equal(t, 0.0, t0)
	assert.Equal(t, 0.0, t1)
}
