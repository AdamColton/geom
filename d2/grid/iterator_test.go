package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	tt := map[string]struct {
		Range
		expected []Pt
	}{
		"(0,0):(3,3)": {
			Range: Range{
				{0, 0},
				{3, 3},
			},
			expected: []Pt{
				{0, 0}, {1, 0}, {2, 0},
				{0, 1}, {1, 1}, {2, 1},
				{0, 2}, {1, 2}, {2, 2},
			},
		},
		"(3,3):(0,0)": {
			Range: Range{
				{3, 3},
				{0, 0},
			},
			expected: []Pt{
				{3, 3}, {2, 3}, {1, 3},
				{3, 2}, {2, 2}, {1, 2},
				{3, 1}, {2, 1}, {1, 1},
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			last := -1
			for i, done := tc.Iter().Start(); !done; done = i.Next() {
				if !assert.Equal(t, last+1, i.Idx()) {
					break
				}
				last++
				if !assert.Equal(t, tc.expected[i.Idx()], i.Pt()) {
					t.Log(i.Idx())
					break
				}
			}
			assert.Equal(t, tc.Iter().Size().Area(), last+1)

			assert.Equal(t, tc.expected, tc.Iter().Slice())

			tc.Iter().Each(func(i int, pt Pt) {
				assert.Equal(t, tc.expected[i], pt)
			})

			assert.True(t, tc.Iter().Until(func(i int, pt Pt) bool {
				assert.Equal(t, tc.expected[i], pt)
				return i == len(tc.expected)-1
			}))
			assert.False(t, tc.Iter().Until(func(i int, pt Pt) bool {
				assert.Equal(t, tc.expected[i], pt)
				return false
			}))

			i := 0
			for got := range tc.Iter().Chan() {
				assert.Equal(t, tc.expected[i], got)
				i++
			}
		})
	}
}
