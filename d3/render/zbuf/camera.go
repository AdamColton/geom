package zbuf

import (
	"math"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d3"
)

type Camera struct {
	d3.Pt
	d3.Q
	Near, Far float64
	Angle     angle.Rad
	W, H      int
}

func (c Camera) T() *d3.T {
	// https://www.youtube.com/watch?v=mpTl003EXCY&list=PL_w_qWAQZtAZhtzPI5pkAtcUVgmzdAP8g&index=5
	v := d3.Pt{}.Subtract(c.Pt)
	translate := d3.Translate(v)
	rot := c.Q.Normalize().T()

	perspective := c.Perspective()
	post := d3.Translate(d3.V{1, 1, 0}).T().T(d3.Scale(d3.V{float64(c.W) / 2, float64(c.H) / 2, -1}).T())

	return translate.T().T(rot).T(perspective).T(post)
}

func (c Camera) Perspective() *d3.T {
	a, b := c.ab()
	x := 1.0 / math.Tan(float64(c.Angle)/2.0)
	y := x * float64(c.W) / float64(c.H)
	return &d3.T{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, a, b},
		{0, 0, -1, 0},
	}
}

func (c Camera) ab() (float64, float64) {
	d := (c.Far - c.Near)
	a := (c.Far + c.Near) / d
	b := (2 * c.Near * c.Far) / d
	return a, b
}
