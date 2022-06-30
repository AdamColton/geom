package list

type Product[TA, TB, TOut any] struct {
	A List[TA]
	B List[TB]
	Combinator
	Fn func(TA, TB) TOut
}

func (p Product[TA, TB, TOut]) Len() int {
	var a, b int
	if p.A != nil {
		a = p.A.Len()
	}
	if p.B != nil {
		b = p.B.Len()
	}
	return p.GetComb().Len(a, b)
}

func (p Product[TA, TB, TOut]) Idx(i int) TOut {
	a, b := p.GetComb().Idx(i, p.A.Len(), p.B.Len())
	return p.Fn(p.A.Idx(a), p.B.Idx(b))
}

func (p Product[TA, TB, TOut]) GetComb() Combinator {
	if p.Combinator == nil {
		return CrossComb{}
	}
	return p.Combinator
}

func (p Product[TA, TB, TOut]) ToSlice() Slice[TOut] {
	lnA := p.A.Len()
	lnB := p.B.Len()
	ln := p.GetComb().Len(lnA, lnB)
	if ln <= 0 {
		return nil
	}
	out := make(Slice[TOut], ln)
	sa := NewSlice(p.A)
	sb := NewSlice(p.B)

	for i := range out {
		a, b := p.GetComb().Idx(i, lnA, lnB)
		out[i] = p.Fn(sa[a], sb[b])
	}

	return out
}
