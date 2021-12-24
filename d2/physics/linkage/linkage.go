package linkage

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
)

type Crank struct {
	Joints   []d2.Pt
	Segments [][2]uint
}

func (c *Crank) Joint(idx int) d2.Pt {
	if idx == 0 {
		return d2.Pt{}
	}
	return c.Joints[idx-1]
}

type Link struct {
	Cranks [2]uint
	Joints [2]uint
}

type Linkage struct {
	Cranks []*Crank
	Links  []*Link
}

type Constraint struct {
	Crank int
	Joint int
	Pt    d2.Pt
	Ang   *angle.Rad
}

type Assemble struct {
	*Linkage
	Constraints []*Constraint
}
