package mesh

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape/triangle"
)

type EdgeMesh map[Edge]bool

func NewEdgeMesh() EdgeMesh {
	return make(EdgeMesh)
}

func EdgeMeshFromTrianges(ts []triangle.Triangle) EdgeMesh {
	em := make(EdgeMesh)
	for _, t := range ts {
		em.Add(t[:]...)
	}
	return em
}

type ErrTwoPoints struct{}

func (ErrTwoPoints) Error() string {
	return "At least two points are required"
}

func (em EdgeMesh) Add(pts ...d2.Pt) error {
	if len(pts) < 2 {
		return ErrTwoPoints{}
	}
	for _, e := range Edges(pts...) {
		em[e] = true
	}
	return nil
}

func (em EdgeMesh) Flip(pts ...d2.Pt) error {
	if len(pts) < 2 {
		return ErrTwoPoints{}
	}
	for _, e := range Edges(pts...) {
		em[e] = !em[e]
	}
	return nil
}
