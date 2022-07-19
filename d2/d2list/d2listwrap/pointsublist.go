package d2listwrap

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/d2list"
)

type PointSubList struct {
	d2list.PointList
	Idxs []int
}

func (p PointSubList) Len() int {
	return len(p.Idxs)
}

func (p PointSubList) Idx(idx int) d2.Pt {
	return p.PointList.Idx(p.Idxs[idx])
}
