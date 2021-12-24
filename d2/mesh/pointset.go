package mesh

import (
	"github.com/adamcolton/geom/d2"
)

type PointSet struct {
	pt2idx map[d2.Pt]uint32
	Pts    []d2.Pt
}

func NewPointSet() *PointSet {
	return &PointSet{
		pt2idx: make(map[d2.Pt]uint32),
	}
}

func (ps *PointSet) Add(pt d2.Pt) uint32 {
	if idx, found := ps.pt2idx[pt]; found {
		return idx
	}
	idx := uint32(len(ps.Pts))
	ps.pt2idx[pt] = idx
	ps.Pts = append(ps.Pts, pt)
	return idx
}

func (ps *PointSet) Has(pt d2.Pt) (uint32, bool) {
	idx, found := ps.pt2idx[pt]
	return idx, found
}
