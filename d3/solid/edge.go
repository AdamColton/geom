package solid

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
)

type Edge [2]d3.Pt

func NewEdge(a, b d3.Point) Edge {
	pa, pb := a.Pt(), b.Pt()
	if pa.X > pb.X {
		pa, pb = pb, pa
	} else if pa.X == pb.X {
		if pa.Y > pb.Y {
			pa, pb = pb, pa
		} else if pa.Y == pb.Y {
			if pa.Z > pb.Z {
				pa, pb = pb, pa
			}
		}
	}
	return Edge{pa, pb}
}

func (e Edge) Pt(t float64) d3.Pt {
	return line.New(e[0], e[1]).Pt(t)
}

type EdgeMesh struct {
	edges   map[Edge]byte
	singles uint
}

func NewEdgeMesh() *EdgeMesh {
	return &EdgeMesh{
		edges: make(map[Edge]byte),
	}
}

type ErrEdgeOverUsed struct{}

func (ErrEdgeOverUsed) Error() string {
	return "Within a mesh, an edge should appear no more than twice"
}

type ErrTwoPoints struct{}

func (ErrTwoPoints) Error() string {
	return "At least two points are required"
}

func (em *EdgeMesh) add(e Edge) error {
	switch em.edges[e] {
	case 0:
		em.edges[e] = 1
		em.singles++
	case 1:
		em.edges[e] = 2
		em.singles--
	case 2:
		return ErrEdgeOverUsed{}
	}
	return nil
}

func (em *EdgeMesh) Add(pts ...d3.Point) error {
	if len(pts) < 2 {
		return ErrTwoPoints{}
	}
	if len(pts) == 2 {
		em.add(NewEdge(pts[0], pts[1]))
		return nil
	}
	ln := len(pts)
	for i, a := range pts {
		b := pts[(i+1)%ln]
		em.add(NewEdge(a, b))
	}
	return nil
}

func (em *EdgeMesh) Solid() bool {
	return len(em.edges) > 0 && em.singles == 0
}
