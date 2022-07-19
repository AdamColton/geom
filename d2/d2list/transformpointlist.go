package d2list

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/list"
)

type TransformPointList = list.Product[*d2.T, d2.Pt, d2.Pt]

func TransformPointListFn(a *d2.T, b d2.Pt) d2.Pt {
	return a.Pt(b)
}

func NewTransformPointList(t TransformList, p PointList) TransformPointList {
	return TransformPointList{
		A:  t,
		B:  p,
		Fn: TransformPointListFn,
	}
}

func PointListTransform(t *d2.T, p PointList) TransformPointList {
	return NewTransformPointList(TransformSlice{t}, p)
}
