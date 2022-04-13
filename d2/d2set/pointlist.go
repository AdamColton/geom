package d2set

import (
	"github.com/adamcolton/geom/d2"
)

type PointList interface {
	Len() int
	Get(n int) d2.Pt
}
