package d2list

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/list"
)

type PointList list.List[d2.Pt]

type PointSlice = list.Slice[d2.Pt]
