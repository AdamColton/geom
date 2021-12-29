package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMag(t *testing.T) {
	tt := map[string]struct {
		in       Prefix
		expected Prefix
	}{
		"Kilo": {
			in:       Kilo,
			expected: Kilo,
		},
		"2k": {
			in:       2 * Kilo,
			expected: Kilo,
		},
		"2m": {
			in:       2 * Mili,
			expected: Mili,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.in.Mag())
		})
	}
}

func TestPrefixString(t *testing.T) {
	tt := map[string]struct {
		in       Prefix
		expected string
	}{
		"k": {
			in:       Kilo,
			expected: "k",
		},
		"2k": {
			in:       2 * Kilo,
			expected: "2k",
		},
		"2.1c": {
			in:       2.1 * Centi,
			expected: "2.1c",
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.in.String())
		})
	}
}
