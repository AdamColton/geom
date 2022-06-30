package d2list

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/list"
)

type CurveList list.List[d2.Pt1]

type CurveSlice = list.Slice[d2.Pt1]
