package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModCombinator(t *testing.T) {
	tt := map[string]struct {
		expected [][2]int
		A, B     int
	}{
		"one-to-one": {
			A: 5,
			B: 5,
			expected: [][2]int{
				{0, 0},
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
			},
		},
		"two-to-one": {
			A: 4,
			B: 2,
			expected: [][2]int{
				{0, 0},
				{1, 1},
				{2, 0},
				{3, 1},
			},
		},
		"one-to-two": {
			A: 2,
			B: 4,
			expected: [][2]int{
				{0, 0},
				{1, 1},
				{0, 2},
				{1, 3},
			},
		},
	}

	mc := ModComb{}
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, len(tc.expected), mc.Len(tc.A, tc.B))
			for i, e := range tc.expected {
				a, b := mc.Idx(i, tc.A, tc.B)
				assert.Equal(t, e[0], a)
				assert.Equal(t, e[1], b)
			}
		})
	}
}

func TestCrossComb(t *testing.T) {
	a, b := 4, 5
	cc := CrossComb{}
	assert.Equal(t, a*b, cc.Len(a, b))
	n := 0
	for j := 0; j < b; j++ {
		for i := 0; i < a; i++ {
			x, y := cc.Idx(n, a, b)
			assert.Equal(t, i, x, "i: %d", n)
			assert.Equal(t, j, y, "j: %d", n)
			n++
		}
	}
}
