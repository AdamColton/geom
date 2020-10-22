package scene

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d3"
)

type Camera struct {
	d3.Pt
	d3.Q
	Angle         angle.Rad
	Width, Height int
}
