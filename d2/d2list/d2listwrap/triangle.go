package d2listwrap

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/d2list"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/triangle"
)

type Triangle struct {
	d2list.PointList
	*triangle.Triangle
}

func (t *Triangle) Update() {
	t.Triangle = &triangle.Triangle{
		t.Idx(0),
		t.Idx(1),
		t.Idx(2),
	}
}

type TriangleGenerator [3]int

func (tg TriangleGenerator) GenerateTriangle(pts d2list.PointList) *Triangle {
	return &Triangle{
		PointList: PointSubList{
			Idxs:      tg[:],
			PointList: pts,
		},
	}
}

func (tg TriangleGenerator) GenerateCurve(pts d2list.PointList) d2.Pt1 {
	return tg.GenerateTriangle(pts)
}

func (tg TriangleGenerator) GenerateShape(pts d2list.PointList) shape.Shape {
	return tg.GenerateTriangle(pts)
}
