package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLenUnit(t *testing.T) {
	tt := map[string]struct {
		expected Length
	}{
		"mm": {
			expected: 0.001,
		},
		"cm": {
			expected: 0.01,
		},
		"10'": {
			expected: 10 * Foot,
		},
		`6"`: {
			expected: 6 * Inch,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			got, err := NewLenUnit(n)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
