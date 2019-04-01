package solid

type IdxEdge [2]uint32

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
		em.add(NewIdxEdge(pts[0], pts[1]))
		return nil
	}
	ln := len(pts)
	for i, a := range pts {
		b := pts[(i+1)%ln]
		em.add(NewIdxEdge(a, b))
	}
	return nil
}

func (em *IdxEdgeMesh) Solid() bool {
	return len(em.edges) > 0 && em.singles == 0
}
