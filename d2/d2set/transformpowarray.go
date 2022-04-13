package d2set

import (
	"github.com/adamcolton/geom/d2"
)

type TransformPowArray struct {
	*d2.T
	Offset, N int
}

func (ta TransformPowArray) Len() int {
	return ta.N
}

func (ta TransformPowArray) Get(n int) *d2.T {
	return ta.Pow(uint(n + ta.Offset))
}

func (ta TransformPowArray) ToTransformSlice() TransformSlice {
	if ta.N <= 0 {
		return nil
	}
	out := make(TransformSlice, ta.N)
	var t *d2.T
	if ta.Offset > 0 {
		t = ta.Pow(uint(ta.Offset))
	} else {
		t = d2.IndentityTransform()
	}

	for i := 0; i < ta.N-1; i++ {
		out[i] = t
		t = t.T(ta.T)
	}
	out[ta.N-1] = t

	return out
}
