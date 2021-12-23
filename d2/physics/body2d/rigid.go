package body2d

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape"
)

type Rigid struct {
	shape.Shape // Bounds with (0,0) = Pt
}

type Position struct {
	d2.V            // Location
	Angle angle.Rad // rotation
}

func (p *Position) T() *d2.T {
	rt := d2.Rotate(p.Angle).T()
	tt := d2.Translate(p.V).T()

	return rt.T(tt)
}

func (p *Position) T2() *d2.T {
	s, c := angle.Rad(p.Angle).Sincos()
	return &d2.T{
		{c, s, c*p.X - s*p.Y},
		{-s, c, s*p.X + c*p.Y},
		{0, 0, 1},
	}
}
