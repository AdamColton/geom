package d2set

import (
	"testing"

	"github.com/adamcolton/geom/d2"
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
