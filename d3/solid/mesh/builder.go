package mesh

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid"
)

// Builder is inteded to help build meshes. Extrusions are more useful and this
// will probably be deprecated.
type Builder struct {
	ctr   uint32
	pts   map[d3.Pt]uint32
	edges *solid.IdxEdgeMesh
	Mesh
}

// NewBuilder creats a builder
func NewBuilder() *Builder {
	return &Builder{
		pts:   make(map[d3.Pt]uint32),
		edges: solid.NewIdxEdgeMesh(),
	}
}

// Add a facet
func (b *Builder) Add(p []d3.Pt) error {
	idxs := make([]uint32, len(p))
	for i, pt := range p {
		idxs[i] = b.addPt(pt)
	}

	b.edges.Add(idxs...)

	b.Polygons = append(b.Polygons, idxs)

	return nil
}

func (b *Builder) addPt(pt d3.Pt) uint32 {
	if i := b.pts[pt]; i > 0 {
		return i - 1
	}

	i := b.ctr
	b.ctr++
	b.pts[pt] = b.ctr
	b.Pts = append(b.Pts, pt)
	return i
}

// Solid is true if the mesh is solid
func (b *Builder) Solid() bool {
	return b.edges.Solid()
}

// Extrude creates a mesh by extruding a face.
func Extrude(face []d3.Pt, v d3.V) Mesh {
	ln := len(face)
	m := Mesh{
		Pts:      make([]d3.Pt, ln*2),
		Polygons: make([][]uint32, ln+2),
	}
	copy(m.Pts, face)
	m.Polygons[0] = make([]uint32, ln)
	m.Polygons[1] = make([]uint32, ln)

	for i, p := range face {
		m.Pts[i+ln] = p.Add(v)
		m.Polygons[0][i] = uint32(i)
		m.Polygons[1][i] = uint32(i + ln)
	}

	prev := ln - 1
	for i := range face {
		m.Polygons[i+2] = []uint32{
			uint32(prev),
			uint32(i),
			uint32(i + ln),
			uint32(prev + ln),
		}
		prev = i
	}
	return m
}
