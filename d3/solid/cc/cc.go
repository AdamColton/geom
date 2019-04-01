package cc

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

type ccMesh struct {
	mesh.Mesh
	edgeCtr    uint32
	edge2Idx   map[solid.IdxEdge]uint32
	pt2edge    map[uint32]map[uint32]uint32 // maps ptIdx to other ptIndex it shares an edge with
	edge2face  map[solid.IdxEdge][]uint32   //maps edges to face indexes
	pt2face    map[uint32][]uint32
	facePoints []d3.Pt
	edgePoints []d3.Pt
	baryPoints []d3.Pt
}

func (b *ccMesh) populateEdges() {
	b.edge2face = make(map[solid.IdxEdge][]uint32)
	b.pt2edge = make(map[uint32]map[uint32]uint32)
	b.pt2face = make(map[uint32][]uint32)
	b.edge2Idx = make(map[solid.IdxEdge]uint32)
	for i, f := range b.Polygons {
		pIdx := f[len(f)-1]
		for _, cIdx := range f {
			b.addFaceEdge(pIdx, cIdx, solid.NewIdxEdge(pIdx, cIdx), uint32(i))
			pIdx = cIdx
		}
	}
}

func (b *ccMesh) addFaceEdge(ptIdx1, ptIdx2 uint32, e solid.IdxEdge, fIdx uint32) {
	m1 := b.pt2edge[ptIdx1]
	if m1 == nil {
		m1 = make(map[uint32]uint32)
		b.pt2edge[ptIdx1] = m1
	}
	m2 := b.pt2edge[ptIdx2]
	if m2 == nil {
		m2 = make(map[uint32]uint32)
		b.pt2edge[ptIdx2] = m2
	}
	if _, indexed := m2[ptIdx1]; !indexed {
		idx := b.edgeCtr
		b.edgeCtr++
		m2[ptIdx1] = idx
		m1[ptIdx2] = idx
		b.edge2Idx[e] = idx
	}
	b.edge2face[e] = append(b.edge2face[e], fIdx)
	b.pt2face[ptIdx1] = append(b.pt2face[ptIdx1], fIdx)
}

func (b *ccMesh) setFacePoints() {
	b.facePoints = make([]d3.Pt, len(b.Polygons))
	for i, f := range b.Polygons {
		p := &affinePoint{}
		for _, idx := range f {
			p.add(b.Pts[idx])
		}
		b.facePoints[i] = p.Get()
	}
}

func (b *ccMesh) setEdgePoints() {
	b.edgePoints = make([]d3.Pt, len(b.edge2face))
	for e, fs := range b.edge2face {
		p := &affinePoint{}
		p.add(b.Pts[e[0]])
		p.add(b.Pts[e[1]])
		for _, fIdx := range fs {
			p.add(b.facePoints[fIdx])
		}
		eIdx := b.edge2Idx[e]
		b.edgePoints[eIdx] = p.Get()
	}
}

func (cc *ccMesh) setBaryPoints() {
	cc.baryPoints = make([]d3.Pt, len(cc.Pts))
	for i, p := range cc.Pts {
		cc.setBaryPoint(uint32(i), p)
	}
}

func (cc *ccMesh) setBaryPoint(i uint32, p d3.Pt) {
	r := &affinePoint{}
	edges := cc.pt2edge[i]
	r.weight(p, float64(len(edges)))
	for p2Idx := range edges {
		r.add(cc.Pts[p2Idx])
	}

	f := &affinePoint{}
	for _, fIdx := range cc.pt2face[i] {
		fPt := cc.facePoints[int(fIdx)]
		f.add(fPt)
	}

	b := &affinePoint{}
	b.weight(f.Get(), 1/f.sum)
	b.weight(r.Get(), 2/f.sum)
	b.weight(p, (f.sum-3)/f.sum)
	cc.baryPoints[i] = b.Get()
}

func (cc *ccMesh) subdivide() mesh.Mesh {
	var m mesh.Mesh
	m.Pts = append(m.Pts, cc.baryPoints...)
	m.Pts = append(m.Pts, cc.edgePoints...)
	m.Pts = append(m.Pts, cc.facePoints...)

	for i, f := range cc.Polygons {
		cc.subdivideFace(uint32(i), f, &m)
	}
	return m
}

// each new facet is defined by the facepoint, one of the original points and
// two edge points
func (cc *ccMesh) subdivideFace(i uint32, f []uint32, m *mesh.Mesh) {
	bpLn := uint32(len(cc.baryPoints))
	epLn := uint32(len(cc.edgePoints))

	fpIdx := bpLn + epLn + i

	ln := len(f)

	prevEIdx := cc.edge2Idx[solid.NewIdxEdge(f[0], f[ln-1])] + bpLn
	for i, cIdx := range f {
		nIdx := f[(i+1)%ln]
		e := solid.NewIdxEdge(cIdx, nIdx)
		curEIdx := cc.edge2Idx[e] + bpLn
		m.Polygons = append(m.Polygons, []uint32{
			fpIdx,
			prevEIdx,
			cIdx,
			curEIdx,
		})
		prevEIdx = curEIdx
	}

}

func Subdivide(m mesh.Mesh, n int) mesh.Mesh {
	if n <= 0 {
		return m
	}
	cc := ccMesh{
		Mesh: m,
	}

	cc.populateEdges()
	cc.setFacePoints()
	cc.setEdgePoints()
	cc.setBaryPoints()
	return Subdivide(cc.subdivide(), n-1)
}
