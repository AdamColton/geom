package d2

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	p := Pt{1, 2}
	v := V{3, 4}
	sum := p.Add(v)
	assert.Equal(t, Pt{4, 6}, sum)
}

func TestSubtract(t *testing.T) {
	p1, p2 := Pt{1, 2}, Pt{4, 3}
	v := p1.Subtract(p2)
	assert.Equal(t, V{-3, -1}, v)
}

func TestAngle(t *testing.T) {
	tt := []struct {
		V
		expected float64
	}{
		{
			V:        V{1, 0},
			expected: 0,
		},
		{
			V:        V{10, 0},
			expected: 0,
		},
		{
			V:        V{0, 1},
			expected: math.Pi / 2.0,
		},
		{
			V:        V{-3, 0},
			expected: math.Pi,
		},
	}

	for _, tc := range tt {
		t.Run(tc.String(), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.V.Angle())
		})
	}
}

func TestTransform(t *testing.T) {
	tt := []struct {
		T
		expected Pt
	}{
		{
			T:        Translate(V{1, 2}),
			expected: Pt{2, 3},
		},
		{
			T:        Rotate(math.Pi / 2),
			expected: Pt{-1, 1},
		},
		{
			T:        Scale(V{2, 3}),
			expected: Pt{2, 3},
		},
		{
			T:        Rotate(math.Pi / 2).T(Translate(V{2, 2})).T(Scale(V{2, 3})),
			expected: Pt{2, 9},
		},
	}

	for _, tc := range tt {
		t.Run(tc.expected.String(), func(t *testing.T) {
			p, _ := tc.T.Pt(Pt{1, 1})
			d := p.Distance(tc.expected)
			assert.InDelta(t, 0, d, 1E-5, p.String())
		})
	}
}
