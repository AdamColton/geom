package d2list

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/list"
)

type CurveList list.List[d2.Pt1]

type CurveSlice = list.Slice[d2.Pt1]

type CurveGenerator interface {
	GenerateCurve(pts PointList) d2.Pt1
}

type CurveGeneratorList list.List[CurveGenerator]

type CurveGeneratorSlice = list.Slice[CurveGenerator]

func GenerateCurveSlice(pts PointList, l CurveGeneratorList) CurveSlice {
	out := make(CurveSlice, l.Len())
	for i := range out {
		out[i] = l.Idx(i).GenerateCurve(pts)
	}
	return out
}
