package render

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
}

func (c Camera) T() *d3.T {
	// https://www.youtube.com/watch?v=mpTl003EXCY&list=PL_w_qWAQZtAZhtzPI5pkAtcUVgmzdAP8g&index=5
	v := d3.Pt{}.Subtract(c.Pt)
	translate := d3.Translate(v)
	rot := c.Q.Normalize().T()

	perspective := c.Perspective()

	return translate.T().T(rot).T(perspective)
}

func (c Camera) Perspective() *d3.T {
	a, b := c.ab()
	m := 1.0 / math.Tan(float64(c.Angle)/2.0)
	return &d3.T{
		{m, 0, 0, 0},
		{0, m, 0, 0},
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
