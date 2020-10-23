package scene

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d3"
)

type Camera struct {
	d3.Pt
	Rot *d3.T
	// Angle corresponds to Width.
	Angle angle.Rad
	Size  grid.Pt
}

func NewCamera(pt d3.Pt, a angle.Rad) *Camera {
	return &Camera{
		Pt:    pt,
		Angle: a,
	}
}

func (c *Camera) SetSize(w, h int) *Camera {
	c.Size = grid.Pt{X: w, Y: h}
	return c
}

func (c *Camera) ByAspect(w int, a grid.Aspect) *Camera {
	c.Size = a.Pt(w)
	return c
}

func (c *Camera) Widescreen(w int) *Camera {
	return c.ByAspect(w, grid.Widescreen)
}

func (c *Camera) Square(w int) *Camera {
	return c.ByAspect(w, grid.Square)
}

func (c *Camera) Fullscreen(w int) *Camera {
	return c.ByAspect(w, grid.Fullscreen)
}

type noramlizer interface {
	Normalize()
}

func (c *Camera) SetRot(r d3.TGen) *Camera {
	if n, ok := r.(noramlizer); ok {
		n.Normalize()
	}
	c.Rot = r.T()
	return c
}

type CameraFactory interface {
	Camera(frameIdx int) *Camera
}
