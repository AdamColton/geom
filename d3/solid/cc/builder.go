package cc

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/adamcolton/geom/d3/solid"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

type Facet []Vertex

type Vertex struct {
	d3.Pt
	Prev, Next float64
}

type Builder struct {
	facets        []Facet
	edges         *solid.IdxEdgeMesh
	points        *solid.PointSet
	innerFacetPts []d3.Pt
	innerFacetIdx []int

	edgePts []d3.Pt
	edgeMap map[solid.IdxEdge]solid.IdxEdge

	handlesMap map[solid.IdxEdge][2][2]uint32
}

func NewBuilder() *Builder {
	return &Builder{
		edges:  solid.NewIdxEdgeMesh(),
		points: solid.NewPointSet(),
	}
}

func (b *Builder) Add(f Facet) error {
	// TODO: check len(f) >2
	idxs := make([]uint32, len(f))
	for i, v := range f {
		idxs[i] = b.points.Add(v.Pt)
	}
	b.edges.Add(idxs...)
	b.facets = append(b.facets, f)
	return nil
}

func (b *Builder) Solid() bool { return b.edges.Solid() }

func (b *Builder) setInnerFacets() {
	b.innerFacetIdx = make([]int, len(b.facets))
	b.handlesMap = make(map[solid.IdxEdge][2][2]uint32)
	for i, f := range b.facets {
		b.innerFace(i, f)
	}
}

func (b *Builder) innerFace(idx int, f Facet) {
	ln := len(f)
	hdls := make([]d3.Pt, 2*ln)
	es := make([]solid.IdxEdge, ln)

	prev := f[ln-1]
	pIdx, _ := b.points.Has(prev.Pt)
	for i, v := range f {
		l := line.New(prev.Pt, v.Pt)
		idx, _ := b.points.Has(v.Pt)
		es[i] = solid.NewIdxEdge(pIdx, idx)
		hdls[i*2] = l.Pt1(prev.Next)
		hdls[i*2+1] = l.Pt1(1 - v.Prev)
		prev = v
		pIdx = idx
	}

	lines := make([]line.Line, ln)
	for i := range lines {
		a := hdls[i*2+1]
		b := hdls[(i*2+4)%(2*ln)]
		lines[i] = line.New(a, b)
	}

	// map[solid.IdxEdge][2][2]uint32 -> [edge]
	prevLn := lines[ln-1]
	b.innerFacetIdx[idx] = len(b.innerFacetPts)
	for i, cur := range lines {
		t0, t1 := prevLn.Closest(cur)
		innerFacetPtIdx := uint32(len(b.innerFacetPts))
		b.innerFacetPts = append(b.innerFacetPts, line.New(prevLn.Pt1(t0), cur.Pt1(t1)).Pt1(0.5))

		for _, e := range []solid.IdxEdge{es[i], es[(i+1)%ln]} {
			hMp := b.handlesMap[e]
			var pts *[2]uint32
			if b.points.Pts[e[0]] == f[i].Pt {
				pts = &hMp[0]
			} else {
				pts = &hMp[1]
			}
			if pts[0] == 0 {
				pts[0] = innerFacetPtIdx
			} else {
				pts[1] = innerFacetPtIdx
			}
			b.handlesMap[e] = hMp
		}

		prevLn = cur
	}
}

func (b *Builder) setEdgePts() {
	b.edgeMap = make(map[solid.IdxEdge]solid.IdxEdge)
	for e, mp := range b.handlesMap {
		el := line.New(b.points.Pts[e[0]], b.points.Pts[e[1]])
		l1 := line.New(b.innerFacetPts[mp[0][0]], b.innerFacetPts[mp[0][1]])
		l2 := line.New(b.innerFacetPts[mp[1][0]], b.innerFacetPts[mp[1][1]])

		t1, _ := el.Closest(l1)
		t2, _ := el.Closest(l2)

		idx := uint32(len(b.edgePts))
		p0, p1 := el.Pt1(t1), el.Pt1(t2)
		b.edgePts = append(b.edgePts, p0, p1)
		b.edgeMap[e] = solid.NewIdxEdge(idx, idx+1)
	}
}

func (b *Builder) mesh() mesh.Mesh {
	var m mesh.Mesh
	m.Pts = append(m.Pts, b.points.Pts...)
	m.Pts = append(m.Pts, b.innerFacetPts...)
	m.Pts = append(m.Pts, b.edgePts...)

	lnOp := uint32(len(b.points.Pts))
	eStart := lnOp + uint32(len(b.innerFacetPts))

	var edgePts []uint32

	for i, f := range b.facets {
		ln := len(f)
		edgePts = edgePts[:0]

		start := uint32(b.innerFacetIdx[i]) + lnOp
		inFc := make([]uint32, len(f))
		for j, v := range f {
			idx, _ := b.points.Has(v.Pt)
			next := f[(j+1)%ln]
			nIdx, _ := b.points.Has(next.Pt)
			e := solid.NewIdxEdge(idx, nIdx)
			ie := b.edgeMap[e]
			closer := 0
			if e[0] != idx {
				closer = 1
			}
			edgePts = append(edgePts, ie[closer]+eStart, ie[1-closer]+eStart)
			inFc[j] = uint32(j) + start
		}
		m.Polygons = append(m.Polygons, inFc)

		for j, v := range f {
			idx, _ := b.points.Has(v.Pt)
			a, b := edgePts[j*2], uint32(j)+start
			m.Polygons = append(m.Polygons, []uint32{
				idx,
				a,
				b,
				edgePts[(j*2-1+2*ln)%(2*ln)],
			}, []uint32{
				a,
				edgePts[j*2+1],
				uint32((j+1)%ln) + start,
				b,
			})
		}

	}

	return m
}
