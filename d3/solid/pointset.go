package solid

import (
	"github.com/adamcolton/geom/d3"
)

// PointSet converts points to indexes.
type PointSet struct {
	pt2idx map[d3.Pt]uint32
	Pts    []d3.Pt
}

// NewPointSet returns an empty PointSet
func NewPointSet() *PointSet {
	return &PointSet{
		pt2idx: make(map[d3.Pt]uint32),
	}
}

// Add a point and return it's index value. If the point is already defined,
// it's index value is returned, but the point is not duplicated.
func (ps *PointSet) Add(pt d3.Pt) uint32 {
	if idx, found := ps.pt2idx[pt]; found {
		return idx
	}
	idx := uint32(len(ps.Pts))
	ps.pt2idx[pt] = idx
	ps.Pts = append(ps.Pts, pt)
	return idx
}

// Has checks if the pointset has a point and returns it's index value.
func (ps *PointSet) Has(pt d3.Pt) (uint32, bool) {
	idx, found := ps.pt2idx[pt]
	return idx, found
}
