package d2

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/angle"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	p := Pt{1, 2}
	v := V{3, 4}
	sum := p.Add(v)
	assert.Equal(t, Pt{4, 6}, sum)
}

func TestPtBaseFuncs(t *testing.T) {
	p := Pt{3, 4}
	assert.Equal(t, p, p.Pt())
	assert.Equal(t, V{3, 4}, p.V())
	assert.Equal(t, Polar{5, p.Angle()}, p.Polar())
	assert.Equal(t, 25.0, p.Mag2())
	assert.Equal(t, 5.0, p.Mag())
	assert.Equal(t, Pt{6, 8}, p.Multiply(2))
}

func TestVBaseFuncs(t *testing.T) {
	v := V{3, 4}
	assert.Equal(t, v, v.V())
	assert.Equal(t, Pt{3, 4}, v.Pt())
	assert.Equal(t, Polar{5, v.Angle()}, v.Polar())
	assert.Equal(t, 25.0, v.Mag2())
	assert.Equal(t, 5.0, v.Mag())
	assert.Equal(t, V{6, 8}, v.Multiply(2))
	assert.Equal(t, V{6, 12}, v.Product(V{2, 3}))
	assert.Equal(t, V{5, 7}, v.Add(V{2, 3}))
	assert.Equal(t, V{1, 3}, v.Subtract(V{2, 1}))
	assert.Equal(t, V{2, 1}, V{-2, -1}.Abs())
}

func TestPolar(t *testing.T) {
	p := Polar{math.Sqrt2 * 2, angle.Deg(45)}
	v := p.V()
	assert.InDelta(t, 2.0, v.X, 1e-10)
	assert.InDelta(t, 2.0, v.Y, 1e-10)
	pt := p.Pt()
	assert.InDelta(t, 2.0, pt.X, 1e-10)
	assert.InDelta(t, 2.0, pt.Y, 1e-10)
}

func TestSubtract(t *testing.T) {
	p1, p2 := Pt{1, 2}, Pt{4, 3}
	v := p1.Subtract(p2)
	assert.Equal(t, V{-3, -1}, v)
}

func TestAngle(t *testing.T) {
	tt := []struct {
		V
		expected angle.Rad
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

func TestMinMax(t *testing.T) {
	m, M := MinMax(Pt{1, 5}, Pt{3, 2}, Pt{2, 3})
	assert.Equal(t, Pt{1, 2}, m)
	assert.Equal(t, Pt{3, 5}, M)

	m, M = MinMax(Pt{3, 2}, Pt{1, 5}, Pt{2, 3})
	assert.Equal(t, Pt{1, 2}, m)
	assert.Equal(t, Pt{3, 5}, M)

	m, M = MinMax()
	assert.Equal(t, Pt{0, 0}, m)
	assert.Equal(t, Pt{0, 0}, M)
}

func TestCross(t *testing.T) {
	tt := map[string]struct {
		a, b     V
		expected float64
	}{
		"x,y-->1": {
			a:        V{1, 0},
			b:        V{0, 1},
			expected: 1,
		},
		"equal-->0": {
			a:        V{1, 2},
			b:        V{1, 2},
			expected: 0,
		},
		"colinear-->0": {
			a:        V{1, 2},
			b:        V{3, 6},
			expected: 0,
		},
		"opposite-->0": {
			a:        V{1, 2},
			b:        V{-1, -2},
			expected: 0,
		},
		"90deg-->1": {
			a:        Polar{1, angle.Deg(20)}.V(),
			b:        Polar{1, angle.Deg(110)}.V(),
			expected: 1,
		},
		"-90deg-->-1": {
			a:        Polar{1, angle.Deg(110)}.V(),
			b:        Polar{1, angle.Deg(20)}.V(),
			expected: -1,
		},
		"2x90deg-->4": {
			a:        Polar{2, angle.Deg(20)}.V(),
			b:        Polar{2, angle.Deg(110)}.V(),
			expected: 4,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.a.Cross(tc.b))
		})
	}
}

func TestDot(t *testing.T) {
	tt := map[string]struct {
		a, b     V
		expected float64
	}{
		"x,y-->0": {
			a:        V{1, 0},
			b:        V{0, 1},
			expected: 0,
		},
		"x,x-->1": {
			a:        V{1, 0},
			b:        V{1, 0},
			expected: 1,
		},
		"aX,bX-->abX": {
			a:        V{2, 0},
			b:        V{3, 0},
			expected: 6,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.a.Dot(tc.b))
		})
	}
}
