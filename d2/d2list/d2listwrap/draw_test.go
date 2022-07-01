package d2listwrap

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/d2list"
	"github.com/adamcolton/geom/d2/draw"
	"github.com/adamcolton/geom/d2/grid"
)

func Test(t *testing.T) {
	gen := draw.ContextGenerator{
		Size:  grid.Pt{500, 500},
		Clear: draw.Color(1, 1, 1),
		Set:   draw.Color(0, 0, 0),
	}
	ctx := gen.Generate()

	ps := d2list.NewTransformPointList(
		d2list.TransformSlice{d2.Scale{100, 100}.T()},
		d2list.PointSlice{
			{1, 1},
			{2, 1},
			{1, 2},
		},
	)

	cs := d2list.CurveSlice{
		&Triangle{
			PointList: PointList{
				PointList: ps,
				Idxs:      []int{0, 1, 2},
			},
		},
	}
	UpdateCurveList(cs)

	DrawCurveList(cs, ctx)
	ctx.SetRGB(1, 0, 0)
	DrawPointList(ps, 3, ctx)
	ctx.SavePNG("test.png")
}
