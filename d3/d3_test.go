package d3

import (
	"fmt"
	"math"
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func EqualPt(t *testing.T, expected, got Pt) {
	assert.InDelta(t, expected.X, got.X, 1e-5, "X")
	assert.InDelta(t, expected.Y, got.Y, 1e-5, "Y")
	assert.InDelta(t, expected.Z, got.Z, 1e-5, "Z")
}

func EqualV(t *testing.T, expected, got V) {
	assert.InDelta(t, expected.X, got.X, 1e-5, "X")
	assert.InDelta(t, expected.Y, got.Y, 1e-5, "Y")
	assert.InDelta(t, expected.Z, got.Z, 1e-5, "Z")
}

func TestBasicMath(tt *testing.T) {
	t := geomtest.New(tt)
	t.Equal(Pt{1, 2, 3}, Pt{1, 2, 3}.Pt())
	t.Equal(V{1, 2, 3}, V{1, 2, 3}.V())

	t.Equal(Pt{3, 5, 7}, Pt{2, 3, 4}.Add(V{1, 2, 3}))
	t.Equal(V{3, 5, 7}, V{2, 3, 4}.Add(V{1, 2, 3}))

	t.Equal(V{1, 2, 3}, Pt{3, 5, 7}.Subtract(Pt{2, 3, 4}))
	t.Equal(V{1, 2, 3}, V{3, 5, 7}.Subtract(V{2, 3, 4}))

	t.Equal(Pt{4, 6, 8}, Pt{2, 3, 4}.Multiply(2))
	t.Equal(V{4, 6, 8}, V{2, 3, 4}.Multiply(2))

	t.Equal(Pt{1, 2, 3}, Pt{1.1, 2.2, 2.9}.Round())

	t.Equal(V{1, 2, 3}, V{-1, -2, -3}.Abs())
}

func TestMag(t *testing.T) {
	tt := []struct {
		M interface {
			Mag() float64
			Mag2() float64
			String() string
		}
		expected float64
	}{
		{
			M:        V{3, 4, 0},
			expected: 5,
		},
		{
			M:        V{3, 4, 12},
			expected: 13,
		},
		{
			M:        Pt{3, 4, 0},
			expected: 5,
		},
		{
			M:        Pt{3, 4, 12},
			expected: 13,
		},
		{
			M:        D3{3, 4, 0},
			expected: 5,
		},
		{
			M:        D3{3, 4, 12},
			expected: 13,
		},
		{
			M:        D3{0, 0, 0},
			expected: 0,
		},
		{
			M:        V{0, 0, 0},
			expected: 0,
		},
		{
			M:        Pt{0, 0, 0},
			expected: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.M.String(), func(t *testing.T) {
			geomtest.Equal(t, tc.expected, tc.M.Mag())
			geomtest.Equal(t, tc.expected*tc.expected, tc.M.Mag2())
		})
	}
}

func TestNormal(t *testing.T) {
	tt := []V{
		{3, 4, 0},
		{3, 4, 12},
	}

	for _, tc := range tt {
		t.Run(tc.String(), func(t *testing.T) {
			assert.Equal(t, 1.0, tc.Normal().Mag())
		})
	}

	assert.Equal(t, V{}, V{}.Normal())
}

func TestCross(t *testing.T) {
	tt := []struct {
		a, b     V
		expected V
	}{
		{
			a:        V{1, 1, 1},
			b:        V{2, 2, 2},
			expected: V{0, 0, 0},
		},
		{
			a:        V{1, 0, 0},
			b:        V{0, 1, 0},
			expected: V{0, 0, 1},
		},
		{
			a:        Rotation{angle.Deg(22.5), XY}.T().V(V{1, 0, 0}),
			b:        Rotation{angle.Deg(-22.5), XY}.T().V(V{1, 0, 0}),
			expected: V{0, 0, -0.70710678118},
		},
	}

	for _, tc := range tt {
		t.Run(tc.a.String()+tc.b.String(), func(t *testing.T) {
			geomtest.Equal(t, tc.expected, tc.a.Cross(tc.b))
		})
	}
}

func TestDot(t *testing.T) {
	tt := []struct {
		a, b     V
		expected float64
	}{
		{
			a:        V{1, 1, 1},
			b:        V{2, 2, 2},
			expected: 6,
		},
		{
			a:        V{1, 1, -2},
			b:        V{2, 2, 2},
			expected: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.a.String()+tc.b.String(), func(t *testing.T) {
			geomtest.Equal(t, tc.expected, tc.a.Dot(tc.b))
		})
	}
}

func TestAddV(t *testing.T) {
	v1 := V{1, 2, 3}
	v2 := V{2, 3, 4}
	geomtest.Equal(t, V{3, 5, 7}, v1.Add(v2))
}

func TestAddPt(t *testing.T) {
	p := Pt{1, 2, 3}
	v := V{2, 3, 4}
	geomtest.Equal(t, Pt{3, 5, 7}, p.Add(v))
}

func TestDistance(t *testing.T) {
	tt := []struct {
		a, b     Pt
		expected float64
	}{
		{
			a:        Pt{0, 0, 0},
			b:        Pt{1, 0, 0},
			expected: 1,
		},
		{
			a:        Pt{1, 3, -3},
			b:        Pt{-2, 7, 9},
			expected: 13,
		},
	}

	for _, tc := range tt {
		t.Run(tc.a.String()+tc.b.String(), func(t *testing.T) {
			geomtest.Equal(t, tc.expected, tc.a.Distance(tc.b))
		})
	}
}

func TestProject(t *testing.T) {
	tt := []struct {
		a, b       V
		bOnA, aOnB V
	}{
		{
			a:    V{1, 0, 0},
			b:    V{1, 1, 0},
			bOnA: V{1, 0, 0},
			aOnB: V{1, 1, 0},
		},
		{
			a:    V{.5, 0, 0},
			b:    V{1, 1, 0},
			bOnA: V{.25, 0, 0},
			aOnB: V{.5, .5, 0},
		},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprint(tc.a, tc.b), func(t *testing.T) {
			geomtest.Equal(t, tc.bOnA, tc.a.Project(tc.b))
			geomtest.Equal(t, tc.aOnB, tc.b.Project(tc.a))
		})
	}
}

func TestAng(t *testing.T) {
	a := math.Pi / 2
	v := V{1, 0, 0}
	v2 := Rotation{
		Angle: angle.Rad(a),
		Plane: XY,
	}.T().V(v)

	geomtest.Equal(t, a, float64(v.Ang(v2)))
}
