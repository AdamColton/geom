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
	ptCtx := PointContext{
		PointList: ps,
	}

	cs := d2list.CurveSlice{
		ptCtx.Triangle(0, 1, 2),
	}
	ptCtx.Update()

	DrawCurveList(cs, ctx)
	ptCtx.Draw(ctx)
	ctx.SavePNG("test.png")
}
