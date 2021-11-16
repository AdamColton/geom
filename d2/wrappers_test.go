package d2

import (
	"testing"

	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

type mockPt1 struct{}

func (mockPt1) Pt1(t0 float64) Pt {
	return Pt{2 * t0, t0 * t0}
}

func TestV1Wrapper(t *testing.T) {
	var w Pt1V1
	w = V1Wrapper{mockPt1{}}

	vApproxEqual(t, V{2, 0}, w.V1(0))
	vApproxEqual(t, V{2, 0.2}, w.V1(0.1))
	vApproxEqual(t, V{2, 1}, w.V1(0.5))
	vApproxEqual(t, V{2, 2}, w.V1(1))

	assert.Equal(t, Pt{1, 0.25}, w.Pt1(0.5))

	geomtest.Equal(t, AssertV1{}, w)
}

func vApproxEqual(t *testing.T, v1, v2 V) {
	assert.InDelta(t, v1.X, v2.X, 1e-4)
	assert.InDelta(t, v1.Y, v2.Y, 1e-4)
}

type mockPt1V1 struct{}

func (mockPt1V1) Pt1(t float64) Pt {
	return Pt{}
}
func (mockPt1V1) V1(t float64) V {
	return V{}
}

type mockV1 struct{}

func (mockV1) V1(t float64) V {
	return V{}
}

type mockPt1V1c0 struct{}

func (mockPt1V1c0) Pt1(t float64) Pt {
	return Pt{}
}
func (mockPt1V1c0) V1c0() V1 {
	return mockV1{}
}

func TestGetV1(t *testing.T) {
	var v1 V1

	v1 = GetV1(mockPt1{})
	_, ok := v1.(V1Wrapper)
	assert.True(t, ok)

	v1 = GetV1(mockPt1V1{})
	_, ok = v1.(mockPt1V1)
	assert.True(t, ok)

	v1 = GetV1(mockPt1V1c0{})
	_, ok = v1.(mockV1)
	assert.True(t, ok)
}

type mockPt2 struct{}

func (mockPt2) Pt2(t0, t1 float64) Pt {
	return Pt{t0, t1}
}

type mockPt2c1 struct{}

func (mockPt2c1) Pt2(t0, t1 float64) Pt {
	return Pt{}
}

func (mockPt2c1) Pt2c1(t0 float64) Pt1 {
	return mockPt1{}
}

func TestGetPt2c1(t *testing.T) {
	var pt2c1 Pt2c1

	pt2c1 = GetPt2c1(mockPt2{})
	w, ok := pt2c1.(Pt2c1Wrapper)
	assert.True(t, ok)
	c1 := w.Pt2c1(.5)
	assert.Equal(t, Pt{.5, .9}, c1.Pt1(.9))

	pt2c1 = GetPt2c1(mockPt2c1{})
	_, ok = pt2c1.(mockPt2c1)
	assert.True(t, ok)
}
