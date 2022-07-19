package d2listwrap

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/d2list"
)

type PtrLine struct{ line.Line }

type Line struct {
	d2list.PointList
	*PtrLine
}

func (l *Line) Update() {
	l.PtrLine = &PtrLine{
		Line: line.New(l.Idx(0), l.Idx(1)),
	}
}

type LineGenerator [2]int

func (lg LineGenerator) GenerateLine(pts d2list.PointList) *Line {
	return &Line{
		PointList: PointSubList{
			Idxs:      lg[:],
			PointList: pts,
		},
	}
}

func (lg LineGenerator) GenerateCurve(pts d2list.PointList) d2.Pt1 {
	return lg.GenerateLine(pts)
}

func (lg LineGenerator) GenerateShapePerimeter(pts d2list.PointList) d2list.ShapePerimeter {
	return lg.GenerateLine(pts)
}
