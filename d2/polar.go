package d2

import (
	"github.com/adamcolton/geom/angle"
)

// Polar is an abstract D2 using polar coordinates. It can be useful for
// defining Pt or V values.
type Polar struct {
	M float64
	A angle.Rad
}

// D2 representation of the Polar
func (p Polar) D2() D2 {
	s, c := p.A.Sincos()
	return D2{c * p.M, s * p.M}
}

// Pt representation of the Polar
func (p Polar) Pt() Pt {
	return Pt(p.D2())
}

// V representation of the Polar
func (p Polar) V() V {
	return V(p.D2())
}
