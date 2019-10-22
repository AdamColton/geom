package solid

import (
	"strconv"
	"strings"
)

// IdxEdge represents an edge as two index points
type IdxEdge [2]uint32

// String fulfils Stringer on IdxEdge
func (e IdxEdge) String() string {
	return strings.Join([]string{
		"[", strconv.Itoa(int(e[0])), ", ", strconv.Itoa(int(e[1])), "]",
	}, "")
}

// NewIdxEdge creates a well ordered IdxEdge
func NewIdxEdge(a, b uint32) IdxEdge {
	e := IdxEdge{a, b}
	e.Sort()
	return e
}

// Sort guarentees that the IdxEdge is well ordered
func (e *IdxEdge) Sort() {
	if e[0] > e[1] {
		e[0], e[1] = e[1], e[0]
	}
}

// IdxEdgeMesh is a mesh defined by indexed points.
type IdxEdgeMesh struct {
	edges   map[IdxEdge]byte
	singles uint
}

// NewIdxEdgeMesh returns an empty IdxEdgeMesh
func NewIdxEdgeMesh() *IdxEdgeMesh {
	return &IdxEdgeMesh{
		edges: make(map[IdxEdge]byte),
	}
}

func (em *IdxEdgeMesh) add(e IdxEdge) error {
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

// Add index points to a mesh. If two points are give an edge is added. If
// multiple points are given then they are added as a polygon. So if the indexes
// are A,B and C are given the edges (A,B), (B,C) and (C,A) will all be added.
func (em *IdxEdgeMesh) Add(pts ...uint32) error {
	if len(pts) < 2 {
		return ErrTwoPoints{}
	}
	if len(pts) == 2 {
		return em.add(NewIdxEdge(pts[0], pts[1]))
	}
	ln := len(pts)
	for i, a := range pts {
		b := pts[(i+1)%ln]
		err := em.add(NewIdxEdge(a, b))
		if err != nil {
			return err
		}
	}
	return nil
}

// Solid returns true if the mesh has edges and each edge is used exactly twice.
func (em *IdxEdgeMesh) Solid() bool {
	return len(em.edges) > 0 && em.singles == 0
}

// Edges returns all the edges in the mesh as a slice
func (em *IdxEdgeMesh) Edges() []IdxEdge {
	es := make([]IdxEdge, 0, len(em.edges))
	for e := range em.edges {
		es = append(es, e)
	}
	return es
}

// SingleDouble returns all the edges in two lists, the edges that appear once
// and the edges that appear twice.
func (em *IdxEdgeMesh) SingleDouble() ([]IdxEdge, []IdxEdge) {
	s := make([]IdxEdge, 0, em.singles)
	d := make([]IdxEdge, 0, len(em.edges)-int(em.singles))
	for e, c := range em.edges {
		if c == 1 {
			s = append(s, e)
		} else {
			d = append(d, e)
		}
	}
	return s, d
}
