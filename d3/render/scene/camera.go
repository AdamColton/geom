package scene

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d3"
)

// Camera can be placed in a scene when rendering.
type Camera struct {
	d3.Pt
	Rot *d3.T
	// Angle of the lens relative to the Width.
	Angle angle.Rad
	Size  grid.Pt
}

// NewCamera creates a camera, it is intended to be setup with chained calls.
func NewCamera(pt d3.Pt, a angle.Rad) *Camera {
	return &Camera{
		Pt:    pt,
		Angle: a,
	}
}

// SetSize by width and height.
func (c *Camera) SetSize(w, h int) *Camera {
	c.Size = grid.Pt{X: w, Y: h}
	return c
}

// ByAspect sets the size from a width and aspect.
func (c *Camera) ByAspect(w int, a grid.Aspect) *Camera {
	c.Size = a.Pt(w)
	return c
}

// Widescreen aspect ratio and the provided width to set the size.
func (c *Camera) Widescreen(w int) *Camera {
	return c.ByAspect(w, grid.Widescreen)
}

// Square aspect ratio and the provided width to set the size.
func (c *Camera) Square(w int) *Camera {
	return c.ByAspect(w, grid.Square)
}

// Fullscreen aspect ratio and the provided width to set the size.
func (c *Camera) Fullscreen(w int) *Camera {
	return c.ByAspect(w, grid.Fullscreen)
}

type noramlizer interface {
	Normalize()
}

// SetRot on the camera from a transform generator. If the generator has a
// Normalize method, that will be called.
func (c *Camera) SetRot(r d3.TGen) *Camera {
	if n, ok := r.(noramlizer); ok {
		n.Normalize()
	}
	c.Rot = r.T()
	return c
}

// CameraFactory produces a camera per frame.
type CameraFactory interface {
	Camera(frameIdx int) *Camera
}
