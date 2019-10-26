package d3

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQ(t *testing.T) {
	s2inv := 1.0 / math.Sqrt2
	tt := []struct {
		in       Pt
		q        Q
		expected Pt
	}{
		{
			q:        Q{1, 0, 0, 0},
			in:       Pt{1, 0, 0},
			expected: Pt{1, 0, 0},
		}, {
			q:        Q{1, 0, 0, 0},
			in:       Pt{0, 1, 0},
			expected: Pt{0, 1, 0},
		}, {
			q:        Q{1, 0, 0, 0},
			in:       Pt{0, 0, 1},
			expected: Pt{0, 0, 1},
		}, {
			q:        Q{0, 0, 0, 1},
			in:       Pt{1, 0, 0},
			expected: Pt{-1, 0, 0},
		}, {
			q:        Q{0, 0, 0, 1},
			in:       Pt{0, 1, 0},
			expected: Pt{0, -1, 0},
		}, {
			q:        Q{0, 0, 1, 0},
			in:       Pt{1, 0, 0},
			expected: Pt{-1, 0, 0},
		}, {
			q:        Q{0, 0, 1, 0},
			in:       Pt{0, 0, 1},
			expected: Pt{0, 0, -1},
		}, {
			q:        Q{s2inv, 0, 0, s2inv},
			in:       Pt{1, 0, 0},
			expected: Pt{0, -1, 0},
		}, {
			q:        Q{s2inv, 0, 0, s2inv},
			in:       Pt{0, 1, 0},
			expected: Pt{1, 0, 0},
		}, {
			q:        Q{s2inv, 0, 0, -s2inv},
			in:       Pt{1, 0, 0},
			expected: Pt{0, 1, 0},
		}, {
			q:        Q{1, 0, 0, -1}.Normalize(),
			in:       Pt{1, 0, 0},
			expected: Pt{0, 1, 0},
		},
		{
			q:        Q{1, 0, 0, 0}.Product(Q{1, 0, 0, -1}).Normalize(),
			in:       Pt{1, 0, 0},
			expected: Pt{0, 1, 0},
		},
	}

	for _, tc := range tt {
		t.Run(tc.expected.String(), func(t *testing.T) {
			t.Log(tc.q.String())
			EqualPt(t, tc.expected, tc.q.T().Pt(tc.in))
		})
	}
}

func TestNormalizeZero(t *testing.T) {
	assert.Equal(t, Q{}, Q{}.Normalize())
}
