package d2list

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
)

// NewRotation creates a TransformList to create a "rotation array" (in CAD
// terms).
func NewRotation(center d2.Pt, arc angle.Rad, n int, v d2.V) TransformList {
	r := Rotation{
		Rad: arc,
		N:   n,
	}
	return Centered(center, Translate(v, r.TransformList()))
}

type Rotation struct {
	angle.Rad
	N int
}

func (r Rotation) TransformList() TransformList {
	perStep := r.Rad / angle.Rad(r.N-1)
	return TransformPowArray{
		T: d2.Rotate{perStep}.T(),
		N: r.N,
	}
}

func Centered(center d2.Pt, l TransformList) TransformProduct {
	p := d2.Translate(d2.Pt{}.Subtract(center)).Pair()
	return TransformPair(p, l)
}

func Translate(v d2.V, l TransformList) TransformList {
	t := TransformSlice{d2.Translate(v).T()}
	return NewTrasformProduct(t, l)
}
