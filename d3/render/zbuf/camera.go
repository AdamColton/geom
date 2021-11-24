package zbuf

import (
	"math"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/scene"
)

type Camera struct {
	Near, Far float64
	*scene.Camera
}

// Move so that the field of view is (0,1) instead of (-1,1) in x, y and z.
var scaleAndMove = d3.Scale(d3.V{0.5, 0.5, -0.5}).T().T(d3.Translate(d3.V{0.5, 0.5, 0.5}).T())

func (c *Camera) T() *d3.T {
	// https://www.youtube.com/watch?v=mpTl003EXCY&list=PL_w_qWAQZtAZhtzPI5pkAtcUVgmzdAP8g&index=5
	v := d3.Pt{}.Subtract(c.Pt)
	translate := d3.Translate(v).T()
	c.SetRot()
	perspective := c.Perspective()

	return d3.TProd(translate, c.Rot, perspective, scaleAndMove)
}

func (c *Camera) Perspective() *d3.T {
	a, b := c.ab()
	x := 1.0 / math.Tan(float64(c.Angle)/2.0)
	y := x * float64(c.Width) / float64(c.Height)
	return &d3.T{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, a, b},
		{0, 0, -1, 0},
	}
}

func (c *Camera) ab() (float64, float64) {
	d := (c.Far - c.Near)
	a := (c.Far + c.Near) / d
	b := (2 * c.Near * c.Far) / d
	return a, b
}
