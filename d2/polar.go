package d2

import (
	"github.com/adamcolton/geom/angle"
)

type Polar struct {
	M float64
	A angle.Rad
}

func (p Polar) D2() D2 {
	s, c := p.A.Sincos()
	return D2{c * p.M, s * p.M}
}

func (p Polar) Pt() Pt {
	return Pt(p.D2())
}

func (p Polar) V() V {
	return V(p.D2())
}
