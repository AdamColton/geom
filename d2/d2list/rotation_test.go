package d2list

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomtest"
)

func TestRotationArray(t *testing.T) {
	center := d2.Pt{1, 1}
	offset := d2.V{2, 0}
	r := NewRotation(center, angle.Deg(180), 5, offset)

	ps := PointSlice{{1, 1}}
	sr2 := math.Sqrt2
	expect := PointSlice{
		{3, 1},
		{sr2 + 1, sr2 + 1},
		{1, 3},
		{1 - sr2, sr2 + 1},
		{-1, 1},
	}

	geomtest.Equal(t, expect, NewTransformPointList(r, ps))

}
