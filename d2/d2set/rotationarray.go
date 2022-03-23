package d2set

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
)

type RotationArray struct {
	V       d2.V
	Arc     angle.Rad
	N       int
	Source  PointSet
	Reverse bool
}

func (ra RotationArray) Len() int {
	return ra.N * ra.Source.Len()
}

func (ra RotationArray) Get(n int) d2.Pt {
	ln := ra.Source.Len()
	pt := ra.Source.Get(n % ln)
	i := uint(n / ln)
	r := angle.Rot(float64(i)) * ra.Arc / angle.Rot(float64(ra.N-1))
	if ra.Reverse {
		r = -r
	}
	t := d2.Rotate(r).T()
	v := t.V(ra.V)
	t = d2.Translate(v).T()
	return t.Pt(pt)
}
