package mesh

import (
	"bytes"
	"math"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/shape/polygon"
	"github.com/adamcolton/geom/d3/solid"
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

func (m Mesh) String() string {
	buf := bytes.NewBuffer(nil)
	m.WriteObj(buf)
	return buf.String()
}

func (m Mesh) Edges() []solid.IdxEdge {
	es := solid.NewIdxEdgeMesh()
	for _, p := range m.Polygons {
		es.Add(p...)
	}
	return es.Edges()
}

type TriangleMesh struct {
	Pts      []d3.Pt
	Polygons [][][3]uint32
}

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

func (m TriangleMesh) Edges() ([]solid.IdxEdge, []solid.IdxEdge) {
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

func (m TriangleMesh) T(t *d3.T) TriangleMesh {
	m2 := TriangleMesh{
		Pts:      make([]d3.Pt, len(m.Pts)),
		Polygons: make([][][3]uint32, len(m.Polygons)),
	}

	for i, p := range m.Pts {
		m2.Pts[i] = t.PtScl(p)
	}
	m2.Polygons = m.Polygons

	return m2
}

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

func (m TriangleMesh) RoundXY() {
	for i, p := range m.Pts {
		m.Pts[i].X = math.Round(p.X)
		m.Pts[i].Y = math.Round(p.Y)
	}
}
