package d2set

import (
	"github.com/adamcolton/geom/calc"
	"github.com/adamcolton/geom/d2"
)

type Combinator interface {
	Len(a, b int) int
	Idx(n, a, b int) (int, int)
}

type CrossComb struct{}

func (CrossComb) Len(a, b int) int {
	return a * b
}

func (CrossComb) Idx(n, a, b int) (int, int) {
	return n % a, n / a
}

type ModComb struct{}

func (ModComb) Len(a, b int) int {
	return calc.LCM(a, b)
}

func (ModComb) Idx(n, a, b int) (int, int) {
	return n % a, n % b
}

type TransformPointList struct {
	PointList
	TransformList
	Combinator
}

func (tpl TransformPointList) Len() int {
	return tpl.GetComb().Len(tpl.PointList.Len(), tpl.TransformList.Len())
}

func (tpl TransformPointList) Get(n int) d2.Pt {
	pln := tpl.PointList.Len()
	tln := tpl.TransformList.Len()
	pIdx, tIdx := tpl.GetComb().Idx(n, pln, tln)
	pt := tpl.PointList.Get(pIdx)
	t := tpl.TransformList.Get(tIdx)
	return t.Pt(pt)
}

func (tpl TransformPointList) GetComb() Combinator {
	if tpl.Combinator == nil {
		return CrossComb{}
	}
	return tpl.Combinator
}

func (tpl TransformPointList) ToPointSlice() PointSlice {
	pln := tpl.PointList.Len()
	tln := tpl.TransformList.Len()
	ln := pln * tln
	if ln <= 0 {
		return nil
	}
	out := make(PointSlice, 0, ln)
	ts := NewTransformSlice(tpl.TransformList)
	ps := NewPointSlice(tpl.PointList)

	for _, t := range ts {
		for _, p := range ps {
			out = append(out, t.Pt(p))
		}
	}

	return out
}
