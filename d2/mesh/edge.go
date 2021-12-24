package mesh

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

type Edge [2]d2.Pt

func NewEdge(pt0, pt1 d2.Pt) Edge {
	if pt0.X < pt1.X {
		return Edge{pt0, pt1}
	}
	if pt0.X > pt1.X {
		return Edge{pt1, pt0}
	}
	if pt0.Y < pt1.Y {
		return Edge{pt0, pt1}
	}
	return Edge{pt1, pt0}
}

func Edges(pts ...d2.Pt) []Edge {
	if len(pts) == 2 {
		return []Edge{NewEdge(pts[0], pts[1])}
	}
	ln := len(pts)
	edges := make([]Edge, ln)
	for i, a := range pts {
		b := pts[(i+1)%ln]
		edges[i] = NewEdge(a, b)
	}
	return edges
}

func (e Edge) Line() line.Line {
	return line.New(e[0], e[1])
}

func (e Edge) Length() float64 {
	return e[0].Distance(e[1])
}

type IdxEdge [2]uint32

func NewIdxEdge(a, b uint32) IdxEdge {
	if b < a {
		return IdxEdge{b, a}
	}
	return IdxEdge{a, b}
}
