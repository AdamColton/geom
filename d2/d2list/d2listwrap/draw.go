package d2listwrap

import (
	"github.com/adamcolton/geom/d2/d2list"
	"github.com/adamcolton/geom/d2/draw"
)

func DrawCurveList(cl d2list.CurveList, ctx *draw.Context) {
	ln := cl.Len()
	for i := 0; i < ln; i++ {
		ctx.Pt1(cl.Idx(i))
	}
}

func DrawPointList(pl d2list.PointList, r float64, ctx *draw.Context) {
	ln := pl.Len()
	for i := 0; i < ln; i++ {
		ctx.Circle(pl.Idx(i), r)
	}
}
