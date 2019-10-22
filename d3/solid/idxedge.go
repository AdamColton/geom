package solid

import (
	"strconv"
	"strings"
)

type IdxEdge [2]uint32

func (e IdxEdge) String() string {
	return strings.Join([]string{
		"[", strconv.Itoa(int(e[0])), ", ", strconv.Itoa(int(e[1])), "]",
	}, "")
}

func NewIdxEdge(a, b uint32) IdxEdge {
	if b < a {
		return IdxEdge{b, a}
	}
	return IdxEdge{a, b}
}

type IdxEdgeMesh struct {
	edges   map[IdxEdge]byte
	singles uint
}

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

func (em *IdxEdgeMesh) Solid() bool {
	return len(em.edges) > 0 && em.singles == 0
}

func (em *IdxEdgeMesh) Edges() []IdxEdge {
	es := make([]IdxEdge, 0, len(em.edges))
	for e := range em.edges {
		es = append(es, e)
	}
	return es
}

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
