package d2set

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestTransformArray(t *testing.T) {
	ta := TransformArray{
		Source: PointSlice{{1, 1}, {2, 2}},
		T:      d2.Translate(d2.V{-1, 1}).T(),
		N:      3,
	}
	assert.Equal(t, 6, ta.Len())

	expect := PointSlice{
		{1, 1}, {2, 2},
		{0, 2}, {1, 3},
		{-1, 3}, {0, 4},
	}

	geomtest.Equal(t, expect, ta)
	geomtest.Equal(t, expect, NewPointSlice(ta))
}

func TestCopyPointSlice(t *testing.T) {
	expect := PointSlice{
		{1, 1}, {2, 2},
		{0, 2}, {1, 3},
		{-1, 3}, {0, 4},
	}
	geomtest.Equal(t, expect, NewPointSlice(expect))
}

func TestRotationArray(t *testing.T) {
	r := RotationArray{
		Arc:    angle.Deg(180),
		V:      d2.V{2, 0},
		N:      5,
		Source: PointSlice{{1, 1}},
	}
	sr2 := math.Sqrt2
	expect := PointSlice{
		{3, 1},
		{sr2 + 1, sr2 + 1},
		{1, 3},
		{1 - sr2, sr2 + 1},
		{-1, 1},
	}
	geomtest.Equal(t, expect, r)
}

func TestReflect(t *testing.T) {
	r := Reflect{
		Source: PointSlice{{1, 1}, {2, 2}, {3, 3}},
		Line:   line.New(d2.Pt{4, 0}, d2.Pt{4, 1}),
	}
	expect := PointSlice{
		{1, 1}, {2, 2}, {3, 3},
		{7, 1}, {6, 2}, {5, 3},
	}
	geomtest.Equal(t, expect, r)
}

func TestPt1Source(t *testing.T) {
	s := Pt1Source{
		Pt1: line.New(d2.Pt{0, 0}, d2.Pt{5, 5}),
	}
	s.Set(0, 1, 6)
	expect := PointSlice{
		{0, 0}, {1, 1}, {2, 2},
		{3, 3}, {4, 4}, {5, 5},
	}
	geomtest.Equal(t, expect, s)
}
