package cc

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/adamcolton/geom/d3/solid"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

type pointSet map[d3.Point]bool

func (ps pointSet) add(pts ...d3.Point) pointSet {
	if ps == nil {
		ps = make(pointSet)
	}
	for _, pt := range pts {
		ps[pt] = true
	}
	return ps
}

type Facet []Vertex

type Vertex struct {
	d3.Pt
	Prev, Next float64
}

type Builder struct {
	facets    []Facet
	edges     *solid.EdgeMesh
	edge2face map[solid.Edge][]int
}

func NewBuilder() *Builder {
	return &Builder{
		edges:     solid.NewEdgeMesh(),
		edge2face: make(map[solid.Edge][]int),
	}
}

func (b *Builder) Add(f Facet) error {
	ln := len(f)
	pts := make([]d3.Point, ln)
	ln--
	for i, v := range f {
		pts[i] = v.Pt
	}
	err := b.edges.Add(pts...)
	if err != nil {
		return err
	}
	fIdx := len(b.facets)
	b.facets = append(b.facets, f)

	prev := f[len(f)-1].Pt
	for _, v := range f {
		cur := v.Pt
		e := solid.NewEdge(prev, cur)
		b.edge2face[e] = append(b.edge2face[e], fIdx)
	}
	return nil
}

func (b *Builder) Solid() bool { return b.edges.Solid() }

type facetMesh struct {
	pts    []d3.Pt
	facets [][]int
}

func (f Facet) polygons() facetMesh {
	ln := len(f)
	fm := facetMesh{
		pts:    make([]d3.Pt, 0, ln*4),
		facets: make([][]int, ln*2+1),
	}
	internalLines := make([]line.Line, ln)
	// populate perimeter points
	for i, v := range f {
		next := f[(i+1)%ln]
		side := line.New(v.Pt, next.Pt)
		h1, h2 := side.Pt(v.Next), side.Pt(1-next.Prev)
		fm.pts = append(fm.pts, v.Pt, h1, h2)
	}
	intStart := ln * 3
	for i := range f {
		p1 := fm.pts[i*3+2]
		p2 := fm.pts[(i*3+7)%intStart]
		internalLines[(i+1)%ln] = line.New(p1, p2)
	}
	for i := range f {
		l0, l1 := internalLines[i], internalLines[(i-1+ln)%ln]
		t0, t1 := l0.Closest(l1)
		pt := line.New(l0.Pt(t0), l1.Pt(t1)).Pt(0.5)
		fm.pts = append(fm.pts, pt)
	}
	for i := range f {
		x := i*3 - 1
		if x < 0 {
			x += intStart
		}
		fm.facets[i*2] = []int{
			i * 3,
			i*3 + 1,
			intStart + i,
			x,
		}
		fm.facets[i*2+1] = []int{
			i*3 + 1,
			i*3 + 2,
			intStart + ((i + 1) % ln),
			intStart + i,
		}
	}
	centerFace := make([]int, ln)
	for i := range centerFace {
		centerFace[i] = intStart + i
	}
	fm.facets[ln*2] = centerFace
	return fm
}

func (b *Builder) Render(iterations int) *mesh.Mesh {

	return nil
}

// func (b *Builder) toCCMesh() *ccMesh {
// 	pmb := newCCMesh()
// 	for _, f := range b.facets {
// 		fm := f.polygons()
// 		pmb.add(fm)
// 	}
// 	return pmb
// }

func (f Facet) innerFace() []d3.Pt {
	ln := len(f)
	hdls := make([]d3.Pt, 2*ln)

	prev := f[ln-1]
	for i, v := range f {
		l := line.New(prev.Pt, v.Pt)
		hdls[i*2] = l.Pt(prev.Next)
		hdls[i*2+1] = l.Pt(1 - v.Prev)
		prev = v
	}

	lines := make([]line.Line, ln)
	for i := range lines {
		a := hdls[i*2+1]
		b := hdls[(i*2+4)%(2*ln)]
		lines[i] = line.New(a, b)
	}

	prevLn := lines[ln-1]
	innerFace := make([]d3.Pt, ln)
	for i, cur := range lines {
		t0, t1 := prevLn.Closest(cur)
		innerFace[i] = line.New(prevLn.Pt(t0), cur.Pt(t1)).Pt(0.5)
		prevLn = cur
	}
	return innerFace
}

// func (b *ccMesh) add(fm facetMesh) {
// 	mp := make([]uint32, len(fm.pts))
// 	for i, p := range fm.pts {
// 		mp[i] = b.addPt(p)
// 	}
// 	for _, f := range fm.facets {
// 		cp := make([]uint32, len(f))
// 		for i, idx := range f {
// 			cp[i] = mp[idx]
// 		}
// 		b.Polygons = append(b.Polygons, cp)
// 	}
// }

// func (b *ccMesh) addPt(pt d3.Pt) uint32 {
// 	if i := b.pts[pt]; i > 0 {
// 		return i - 1
// 	}

// 	i := b.ctr
// 	b.ctr++
// 	b.pts[pt] = b.ctr
// 	b.Pts = append(b.Pts, pt)
// 	return i
// }
