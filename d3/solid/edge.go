package solid

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
)

// Edge between two points. Edge should be ordered by calling Sort. Edge is used
// as a map key to emulate a set.
type Edge [2]d3.Pt

// NewEdge from 2 points in the correct order
func NewEdge(a, b d3.Point) Edge {
	e := Edge{a.Pt(), b.Pt()}
	e.Sort()
	return e
}

// Sort guarentees the order of the points. Should only be called once when the
// Edge is created.
func (e *Edge) Sort() {
	if e[0].X > e[1].X {
		e[0], e[1] = e[1], e[0]
	} else if e[0].X == e[1].X {
		if e[0].Y > e[1].Y {
			e[0], e[1] = e[1], e[0]
		} else if e[0].Y == e[1].Y {
			if e[0].Z > e[1].Z {
				e[0], e[1] = e[1], e[0]
			}
		}
	}
}

// Pt1 treats the edge as a line and returns the corresponding point on that
// line.
func (e Edge) Pt1(t0 float64) d3.Pt {
	return line.New(e[0], e[1]).Pt1(t0)
}

// EdgeMesh represents a mesh as a set of edges. It can detect if the mesh is
// solid or if an edge is over used.
type EdgeMesh struct {
	edges   map[Edge]byte
	singles uint
}

// NewEdgeMesh creates an empty EdgeMesh
func NewEdgeMesh() *EdgeMesh {
	return &EdgeMesh{
		edges: make(map[Edge]byte),
	}
}

// ErrEdgeOverUsed is returned if an edge appears more than twice in a mesh.
type ErrEdgeOverUsed struct{}

// Error fulfils error on ErrEdgeOverUsed
func (ErrEdgeOverUsed) Error() string {
	return "Within a mesh, an edge should appear no more than twice"
}

// ErrTwoPoints is returned if a single point is added to a mesh
type ErrTwoPoints struct{}

// Error fulfils error on ErrTwoPoints
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

// Add points to a mesh. If two points are give an edge is added. If multiple
// points are given then they are added as a polygon. So if the points A,B and C
// are given the edges (A,B), (B,C) and (C,A) will all be added.
func (em *EdgeMesh) Add(pts ...d3.Point) error {
	if len(pts) < 2 {
		return ErrTwoPoints{}
	}
	if len(pts) == 2 {
		return em.add(NewEdge(pts[0], pts[1]))
	}
	ln := len(pts)
	for i, a := range pts {
		b := pts[(i+1)%ln]
		err := em.add(NewEdge(a, b))
		if err != nil {
			return err
		}
	}
	return nil
}

// Solid returns true if the mesh has edges and each edge is used exactly twice.
func (em *EdgeMesh) Solid() bool {
	return len(em.edges) > 0 && em.singles == 0
}
