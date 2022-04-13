package d2set

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

type Reflect struct {
	line.Line
	Source PointList
}

func (r Reflect) Len() int {
	return 2 * r.Source.Len()
}

func (r Reflect) Get(n int) d2.Pt {
	ln := r.Source.Len()
	if n < ln {
		return r.Source.Get(n)
	}
	pt := r.Source.Get(n - ln)
	return r.Line.Reflect(pt)
}
