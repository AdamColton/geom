package mesh

import (
	"github.com/adamcolton/geom/d3"
)

type Mesh struct {
	Pts      []d3.Pt
	Polygons [][]uint32
}
