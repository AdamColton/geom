package mesh

import (
	"bytes"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/shape/polygon"
	"github.com/adamcolton/geom/d3/solid"
)

// Mesh defines a solid with polygon facets
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

// T applies a transform to all the points in the mesh
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

// String version of the mesh in .obj format
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

// TriangleMesh is a mesh comprised only of triangles.
type TriangleMesh struct {
	Pts      []d3.Pt
	Polygons [][][3]uint32
}

// TriangleMesh generated from a normal mesh
func (m Mesh) TriangleMesh() (TriangleMesh, error) {
	var i int
	out := TriangleMesh{
		Pts:      m.Pts,
		Polygons: make([][][3]uint32, len(m.Polygons)),
	}

	for i = range m.Polygons {
		p3D := polygon.Polygon(m.Face(i))
		p2D, _, err := p3D.D2()
		if err != nil {
			return out, err
		}
		triangles := p2D.FindTriangles()
		byIdx := make([][3]uint32, len(triangles))
		for j, t := range triangles {
			byIdx[j] = [3]uint32{
				m.Polygons[i][t[0]],
				m.Polygons[i][t[1]],
				m.Polygons[i][t[2]],
			}
		}
		out.Polygons[i] = byIdx
	}
	return out, nil
}

// Edges returns the edges in the mesh. Original contains the edges that were in
// the original mesh and subFacet is the edges creates when facets were
// subdivided into triangles.
func (m TriangleMesh) Edges() (original []solid.IdxEdge, subFacet []solid.IdxEdge) {
	s := make(map[solid.IdxEdge]bool)
	d := make(map[solid.IdxEdge]bool)
	for _, p := range m.Polygons {
		es := solid.NewIdxEdgeMesh()
		for _, t := range p {
			es.Add(t[:]...)
		}
		ps, pd := es.SingleDouble()
		for _, se := range ps {
			s[se] = true
		}
		for _, de := range pd {
			d[de] = true
		}
	}
	var sl, dl []solid.IdxEdge
	for se := range s {
		sl = append(sl, se)
	}
	for de := range d {
		dl = append(dl, de)
	}
	return sl, dl
}

// T applies a transformation to all the points in the mesh.
func (m TriangleMesh) T(t *d3.T) TriangleMesh {
	return TriangleMesh{
		Pts:      t.PtsScl(m.Pts),
		Polygons: m.Polygons,
	}
}

// Face converts a face from index values to Pt values
func (m TriangleMesh) Face(idx int) [][3]d3.Pt {
	p := m.Polygons[idx]
	f := make([][3]d3.Pt, len(p))
	for i, idx := range p {
		f[i][0] = m.Pts[idx[0]]
		f[i][1] = m.Pts[idx[1]]
		f[i][2] = m.Pts[idx[2]]
	}
	return f
}
