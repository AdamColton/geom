package cc

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/affine"
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

func (cc *ccMesh) populateEdges() {
	cc.edge2face = make(map[solid.IdxEdge][]uint32)
	cc.pt2edge = make(map[uint32]map[uint32]uint32)
	cc.pt2face = make(map[uint32][]uint32)
	cc.edge2Idx = make(map[solid.IdxEdge]uint32)
	for i, f := range cc.Polygons {
		pIdx := f[len(f)-1]
		for _, cIdx := range f {
			cc.addFaceEdge(pIdx, cIdx, solid.NewIdxEdge(pIdx, cIdx), uint32(i))
			pIdx = cIdx
		}
	}
}

func (cc *ccMesh) addFaceEdge(ptIdx1, ptIdx2 uint32, e solid.IdxEdge, fIdx uint32) {
	m1 := cc.pt2edge[ptIdx1]
	if m1 == nil {
		m1 = make(map[uint32]uint32)
		cc.pt2edge[ptIdx1] = m1
	}
	m2 := cc.pt2edge[ptIdx2]
	if m2 == nil {
		m2 = make(map[uint32]uint32)
		cc.pt2edge[ptIdx2] = m2
	}
	if _, indexed := m2[ptIdx1]; !indexed {
		idx := cc.edgeCtr
		cc.edgeCtr++
		m2[ptIdx1] = idx
		m1[ptIdx2] = idx
		cc.edge2Idx[e] = idx
	}
	cc.edge2face[e] = append(cc.edge2face[e], fIdx)
	cc.pt2face[ptIdx1] = append(cc.pt2face[ptIdx1], fIdx)
}

func (cc *ccMesh) setFacePoints() {
	cc.facePoints = make([]d3.Pt, len(cc.Polygons))
	for i, f := range cc.Polygons {
		p := &affine.Weighted{}
		for _, idx := range f {
			p.Add(cc.Pts[idx])
		}
		cc.facePoints[i] = p.Get()
	}
}

func (cc *ccMesh) setEdgePoints() {
	cc.edgePoints = make([]d3.Pt, len(cc.edge2face))
	for e, fs := range cc.edge2face {
		p := &affine.Weighted{}
		p.Add(cc.Pts[e[0]], cc.Pts[e[1]])
		for _, fIdx := range fs {
			p.Add(cc.facePoints[fIdx])
		}
		eIdx := cc.edge2Idx[e]
		cc.edgePoints[eIdx] = p.Get()
	}
}

func (cc *ccMesh) setBaryPoints() {
	cc.baryPoints = make([]d3.Pt, len(cc.Pts))
	for i, p := range cc.Pts {
		cc.setBaryPoint(uint32(i), p)
	}
}

func (cc *ccMesh) setBaryPoint(i uint32, p d3.Pt) {
	r := &affine.Weighted{}
	edges := cc.pt2edge[i]
	r.Weight(p, float64(len(edges)))
	for p2Idx := range edges {
		r.Add(cc.Pts[p2Idx])
	}

	f := &affine.Weighted{}
	for _, fIdx := range cc.pt2face[i] {
		fPt := cc.facePoints[int(fIdx)]
		f.Add(fPt)
	}

	b := &affine.Weighted{}
	b.Weight(f.Get(), 1/f.Sum)
	b.Weight(r.Get(), 2/f.Sum)
	b.Weight(p, (f.Sum-3)/f.Sum)
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
