package d2listwrap

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/list"
)

type Bezier struct {
	PointList
	bezier.Bezier
}

func (b *Bezier) Update() {
	b.Bezier = bezier.Bezier(list.NewSlice[d2.Pt](b))
}
