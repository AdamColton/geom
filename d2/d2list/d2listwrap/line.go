package d2listwrap

import (
	"github.com/adamcolton/geom/d2/curve/line"
)

type PtrLine struct{ line.Line }

type Line struct {
	PointList
	*PtrLine
}

func (l *Line) Update() {
	l.PtrLine = &PtrLine{
		Line: line.New(l.Idx(0), l.Idx(1)),
	}
}
