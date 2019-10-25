package cc

import (
	"testing"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

func TestBuilder(t *testing.T) {
	m := mesh.Extrude([]d3.Pt{
		{0, 0, 0},
		{1, 0, 0},
		{1.25, 0.5, 0},
		{1, 1, 0},
		{0, 1, 0},
		{-0.25, 0.5, 0},
	}, d3.V{0, 0, 3})

	b := NewBuilder()
	for idx := range m.Polygons {
		offset := 0.1
		if idx == 0 || idx == 1 {
			offset = 0.33
		}
		f := m.Face(idx)
		fct := make(Facet, len(f))
		for i, p := range f {
			fct[i] = Vertex{
				Pt:   p,
				Next: offset,
				Prev: offset,
			}
		}
		b.Add(fct)
	}

	b.setInnerFacets()
	b.setEdgePts()
}
