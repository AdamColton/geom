package scene

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d3"
)

type Camera struct {
	d3.Pt
	d3.Q
	Rot           *d3.T
	Angle         angle.Rad
	Width, Height int
}

func (c *Camera) SetRot() {
	if c.Rot == nil {
		c.Rot = c.Q.Normalize().T()
	}
}
