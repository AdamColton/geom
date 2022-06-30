package d2list

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/geomtest"
)

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
