package mesh

type IdxEdgeMesh map[IdxEdge]bool

func NewIdxEdgeMesh() IdxEdgeMesh {
	return make(IdxEdgeMesh)
}

func (iem IdxEdgeMesh) Add(pts ...uint32) error {
	if len(pts) < 2 {
		return ErrTwoPoints{}
	}
	if len(pts) == 2 {
		iem[NewIdxEdge(pts[0], pts[1])] = true
		return nil
	}
	ln := len(pts)
	for i, a := range pts {
		b := pts[(i+1)%ln]
		iem[NewIdxEdge(a, b)] = true
	}
	return nil
}
