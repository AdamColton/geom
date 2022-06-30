package d2listwrap

import (
	"github.com/adamcolton/geom/d2/shape/triangle"
)

type Triangle struct {
	PointList
	*triangle.Triangle
}

func (t *Triangle) Update() {
	t.Triangle = &triangle.Triangle{
		t.Idx(0),
		t.Idx(1),
		t.Idx(2),
	}
}
