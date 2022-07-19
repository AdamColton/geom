package d2listwrap

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/d2/d2list"
	"github.com/adamcolton/geom/list"
)

type Bezier struct {
	d2list.PointList
	bezier.Bezier
}

func (b *Bezier) Update() {
	b.Bezier = bezier.Bezier(list.NewSlice[d2.Pt](b))
}

type BezierGenerator []int

func (bg BezierGenerator) GenerateBezier(pts d2list.PointList) *Bezier {
	return &Bezier{
		PointList: PointSubList{
			Idxs:      bg,
			PointList: pts,
		},
	}
}

func (bg BezierGenerator) GenerateCurve(pts d2list.PointList) d2.Pt1 {
	return bg.GenerateBezier(pts)
}

func (bg BezierGenerator) GenerateShapePerimeter(pts d2list.PointList) d2list.ShapePerimeter {
	return bg.GenerateBezier(pts)
}
