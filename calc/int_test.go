package calc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGCD(t *testing.T) {
	tt := map[string]struct {
		a, b, expected int
	}{
		"15,20,5": {
			a:        15,
			b:        20,
			expected: 5,
		},
		"154,165,11": {
			a:        154,
			b:        165,
			expected: 11,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, GCD(tc.a, tc.b))
		})
	}
}

func TestLCM(t *testing.T) {
	// 2,3,5,7,11

	tt := map[string]struct {
		a, b, expected int
	}{
		"2,4,4": {
			a:        2,
			b:        4,
			expected: 4,
		},
		"110,105,2310": {
			a:        110,
			b:        105,
			expected: 2310,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, LCM(tc.a, tc.b))
		})
	}
}
