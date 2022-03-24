package d2set

import (
	"github.com/adamcolton/geom/d2"
)

type TransformPowArray struct {
	Source PointSet
	*d2.T
	Offset, N int
}

func (ta TransformPowArray) Len() int {
	return ta.N * ta.Source.Len()
}

func (ta TransformPowArray) Get(n int) d2.Pt {
	ln := ta.Source.Len()
	pt := ta.Source.Get(n % ln)
	t := ta.Pow(uint((n / ln) + ta.Offset))
	return t.Pt(pt)
}

func (ta TransformPowArray) ToPointSlice() PointSlice {
	if ta.N <= 0 {
		return nil
	}
	ln := ta.Source.Len()
	n := ln * ta.N
	out := make([]d2.Pt, n)
	for i := 0; i < ln; i++ {
		out[i] = ta.Source.Get(i)
	}
	if ta.Offset > 0 {
		t := ta.Pow(uint(ta.Offset))
		for i := 0; i < ln; i++ {
			out[i] = t.Pt(out[i])
		}
	}

	for i := ln; i < n; i += ln {
		for j := 0; j < ln; j++ {
			out[i+j] = ta.T.Pt(out[i+j-ln])
		}
	}
	return out
}
