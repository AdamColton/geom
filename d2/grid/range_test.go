package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	r := Range{
		{0, 0},
		{3, 3},
	}
	assert.True(t, r.Contains(Pt{1, 1}))
	assert.True(t, r.Contains(Pt{0, 0}))
	assert.False(t, r.Contains(Pt{3, 3}))
	assert.Equal(t, Pt{0, 0}, r.Min())
	assert.Equal(t, Pt{2, 2}, r.Max())

	r = Range{
		{0, 0},
		{-3, -3},
	}
	assert.True(t, r.Contains(Pt{-1, -1}))
	assert.True(t, r.Contains(Pt{0, 0}))
	assert.False(t, r.Contains(Pt{-3, -3}))
	assert.Equal(t, Pt{-2, -2}, r.Min())
	assert.Equal(t, Pt{0, 0}, r.Max())

	r = Range{
		{-3, -3},
		{0, 0},
	}
	assert.Equal(t, Pt{-3, -3}, r.Min())
	assert.Equal(t, Pt{-1, -1}, r.Max())

	r = Range{
		{3, 3},
		{0, 0},
	}
	assert.Equal(t, Pt{1, 1}, r.Min())
	assert.Equal(t, Pt{3, 3}, r.Max())

	r = Range{
		{0, 0},
		{0, 0},
	}
	assert.False(t, r.Contains(Pt{0, 0}))

	r = Range{
		{0, 0},
		{1, 0},
	}
	assert.False(t, r.Contains(Pt{0, 0}))

}

func TestRangeScale(t *testing.T) {
	tt := map[string]struct {
		Range
		expected [][2]float64
	}{
		"0,0:3,3": {
			Range: Range{Pt{0, 0}, Pt{3, 3}},
			expected: [][2]float64{
				{0, 0}, {0.5, 0}, {1, 0},
				{0, 0.5}, {0.5, 0.5}, {1, 0.5},
				{0, 1}, {0.5, 1}, {1, 1},
			},
		},
		"3,3:0,0": {
			Range: Range{Pt{3, 3}, Pt{0, 0}},
			expected: [][2]float64{
				{0, 0}, {0.5, 0}, {1, 0},
				{0, 0.5}, {0.5, 0.5}, {1, 0.5},
				{0, 1}, {0.5, 1}, {1, 1},
			},
		},
		"0,0:1,1": {
			Range:    Range{Pt{0, 0}, Pt{1, 1}},
			expected: [][2]float64{{1, 1}},
		},
		"0,0:0,0": {
			Range:    Range{Pt{0, 0}, Pt{0, 0}},
			expected: [][2]float64{},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			s := tc.Scale()
			var got [2]float64
			for i, pt := range tc.Iter().Slice() {
				got[0], got[1] = s.T(pt)
				assert.Equal(t, tc.expected[i], got)
			}
		})
	}
}
