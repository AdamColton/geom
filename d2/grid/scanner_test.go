package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	tt := map[string]struct {
		Range
		Opts     []ScanOption
		expected []Pt
	}{
		"(0,0):(3,3)+x+y": {
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
		"(0,0):(3,3)+x-y": {
			Opts: []ScanOption{ScanSecondaryReversed},
			Range: Range{
				{0, 0},
				{3, 3},
			},
			expected: []Pt{
				{0, 2}, {1, 2}, {2, 2},
				{0, 1}, {1, 1}, {2, 1},
				{0, 0}, {1, 0}, {2, 0},
			},
		},
		"(0,0):(3,3)-x+y": {
			Opts: []ScanOption{ScanPrimaryReversed},
			Range: Range{
				{0, 0},
				{3, 3},
			},
			expected: []Pt{
				{2, 0}, {1, 0}, {0, 0},
				{2, 1}, {1, 1}, {0, 1},
				{2, 2}, {1, 2}, {0, 2},
			},
		},
		"(0,0):(3,3)-x-y": {
			Opts: []ScanOption{ScanPrimaryReversed, ScanSecondaryReversed},
			Range: Range{
				{0, 0},
				{3, 3},
			},
			expected: []Pt{
				{2, 2}, {1, 2}, {0, 2},
				{2, 1}, {1, 1}, {0, 1},
				{2, 0}, {1, 0}, {0, 0},
			},
		},
		"(0,0):(3,3)+y+x": {
			Opts: []ScanOption{ScanVertical},
			Range: Range{
				{0, 0},
				{3, 3},
			},
			expected: []Pt{
				{0, 0}, {0, 1}, {0, 2},
				{1, 0}, {1, 1}, {1, 2},
				{2, 0}, {2, 1}, {2, 2},
			},
		},
		"(0,0):(3,3)+y-x": {
			Opts: []ScanOption{ScanVertical, ScanSecondaryReversed},
			Range: Range{
				{0, 0},
				{3, 3},
			},
			expected: []Pt{
				{2, 0}, {2, 1}, {2, 2},
				{1, 0}, {1, 1}, {1, 2},
				{0, 0}, {0, 1}, {0, 2},
			},
		},
		"(0,0):(3,3)-y+x": {
			Opts: []ScanOption{ScanVertical, ScanPrimaryReversed},
			Range: Range{
				{0, 0},
				{3, 3},
			},
			expected: []Pt{
				{0, 2}, {0, 1}, {0, 0},
				{1, 2}, {1, 1}, {1, 0},
				{2, 2}, {2, 1}, {2, 0},
			},
		},
		"(3,3):(0,0)+x+y": {
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
			s := NewScanner(tc.Range, tc.Opts...)
			if !assert.Equal(t, tc.expected[0], s.s) {
				return
			}
			last := -1
			for done := s.Reset(); !done; done = s.Next() {
				last++
				if !assert.Equal(t, last, s.Idx()) {
					return
				}
				if !assert.Equal(t, tc.expected[s.Idx()], s.Pt()) {
					t.Log(last)
					return
				}
				if !assert.Equal(t, s.Idx(), s.PtIdx(s.Pt())) {
					t.Log(last)
					return
				}
			}
			assert.True(t, s.Done())
			assert.EqualValues(t, last, len(tc.expected)-1)
			assert.True(t, s.Contains(tc.Range[0]))
			assert.False(t, s.Contains(tc.Range[1]))
		})
	}
}
