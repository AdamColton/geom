package d3

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate(t *testing.T) {
	pt := Pt{1, 2, 3}
	tr := Translate(V{1, 2, 3})
	pt, _ = tr.Pt(pt)

	assert.Equal(t, Pt{2, 4, 6}, pt)
}

func TestMag(t *testing.T) {
	tt := []struct {
		V
		expected float64
	}{
		{
			V:        V{3, 4, 0},
			expected: 5,
		},
		{
			V:        V{3, 4, 12},
			expected: 13,
		},
	}

	for _, tc := range tt {
		t.Run(tc.V.String(), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.V.Mag())
		})
	}
}

func TestNormal(t *testing.T) {
	tt := []V{
		V{3, 4, 0},
		V{3, 4, 12},
	}

	for _, tc := range tt {
		t.Run(tc.String(), func(t *testing.T) {
			v, err := tc.Normal()
			assert.NoError(t, err)
			assert.Equal(t, 1.0, v.Mag())
		})
	}

	_, err := V{}.Normal()
	assert.Equal(t, ErrZeroVector{}, err)
}
