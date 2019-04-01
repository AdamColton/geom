package mesh

import (
	"github.com/adamcolton/geom/d3"
)

type Mesh struct {
	Pts      []d3.Pt
	Polygons [][]uint32
}

func (m Mesh) Face(idx int) []d3.Pt {
	p := m.Polygons[idx]
	f := make([]d3.Pt, len(p))
	for i, idx := range p {
		f[i] = m.Pts[idx]
	}
	return f
}
