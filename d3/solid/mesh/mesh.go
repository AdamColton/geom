package mesh

import (
	"bytes"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid"
)

// Mesh defines a solid with polygon facets.
type Mesh struct {
	Pts      []d3.Pt
	Polygons [][]uint32
}

// Face tranlates the index values in a face into Pts.
func (m Mesh) Face(idx int) []d3.Pt {
	p := m.Polygons[idx]
	f := make([]d3.Pt, len(p))
	for i, idx := range p {
		f[i] = m.Pts[idx]
	}
	return f
}

// T applies a transform to all the points in the mesh.
func (m Mesh) T(t *d3.T) Mesh {
	m2 := Mesh{
		Pts:      make([]d3.Pt, len(m.Pts)),
		Polygons: make([][]uint32, len(m.Polygons)),
	}

	for i, p := range m.Pts {
		m2.Pts[i] = t.PtScl(p)
	}
	for i, p := range m.Polygons {
		m2.Polygons[i] = make([]uint32, len(p))
		copy(m2.Polygons[i], p)
	}

	return m2
}

// String version of the mesh in .obj format.
func (m Mesh) String() string {
	buf := bytes.NewBuffer(nil)
	m.WriteObj(buf)
	return buf.String()
}

// Edges converts the index values into edges making sure there are no duplicates.
func (m Mesh) Edges() []solid.IdxEdge {
	es := solid.NewIdxEdgeMesh()
	for _, p := range m.Polygons {
		es.Add(p...)
	}
	return es.Edges()
}
