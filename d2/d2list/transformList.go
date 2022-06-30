package d2list

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/list"
)

type TransformList list.List[*d2.T]

type TransformSlice = list.Slice[*d2.T]

type TransformProduct = list.Product[*d2.T, *d2.T, *d2.T]

func TransformProductFn(a, b *d2.T) *d2.T {
	return a.T(b)
}

func NewTrasformProduct(a, b TransformList) TransformProduct {
	return TransformProduct{
		A:          a,
		B:          b,
		Combinator: list.ModComb{},
		Fn:         TransformProductFn,
	}
}

func TransformPair(pair [2]*d2.T, l TransformList) TransformProduct {
	p := NewTrasformProduct(TransformSlice{pair[0]}, l)
	return NewTrasformProduct(p, TransformSlice{pair[1]})
}
