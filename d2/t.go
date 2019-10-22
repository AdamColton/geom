package d2

import (
	"strconv"
	"strings"

	"github.com/adamcolton/geom/angle"
)

/*

| a b c |   | x |
| d e f | * | y | = | ax+by+c dx+ey+f gx+hy+i|
| g h i |   | 1 |

Because of the array syntax, (x,y) get flipped either in the layout or in the
index. I've chosen the index. Therefor

| (0,0) (1,0) (2,0) |   | [0][0] [0][1] [0][2] |
| (0,1) (1,1) (2,1) | = | [1][0] [1][1] [1][2] |
| (0,2) (1,2) (2,2) |   | [2][0] [2][1] [2][2] |

| a b c |   | j k l |   | aj+bm+cp ak+bn+cq al+bo+cr |
| d e f | * | m n o | = | dj+em+fp dk+en+fq dl+eo+fr |
| g h i |   | p q r |   | gj+hm+ip gk+hn+iq gl+ho+ir |

*/

// T represets a transform matrix
type T [3][3]float64

// IndentityTransform returns an Identity matrix
func IndentityTransform() T {
	return T{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
}

// PtF applies the transform to a Pt, returning the resulting Pt and scale.
func (t T) PtF(pt Pt) (Pt, float64) {
	return Pt{
		pt.X*t[0][0] + pt.Y*t[0][1] + t[0][2],
		pt.X*t[1][0] + pt.Y*t[1][1] + t[1][2],
	}, pt.X*t[2][0] + pt.Y*t[2][1] + t[2][2]
}

// Pt applies the transform to a Pt.
func (t T) Pt(pt Pt) Pt {
	return Pt{
		pt.X*t[0][0] + pt.Y*t[0][1] + t[0][2],
		pt.X*t[1][0] + pt.Y*t[1][1] + t[1][2],
	}
}

// Slice applies the transform to a slice of Pts
func (t T) Slice(pts []Pt) []Pt {
	out := make([]Pt, len(pts))
	for i, pt := range pts {
		out[i] = t.Pt(pt)
	}
	return out
}

// V applies the transform to a V
func (t T) V(v V) V {
	return V{
		v.X*t[0][0] + v.Y*t[0][1] + t[0][2],
		v.X*t[1][0] + v.Y*t[1][1] + t[1][2],
	}
}

// VF applies the transform and returns a V and the scale
func (t T) VF(v V) (V, float64) {
	return V{
		v.X*t[0][0] + v.Y*t[0][1] + t[0][2],
		v.X*t[1][0] + v.Y*t[1][1] + t[1][2],
	}, v.X*t[2][0] + v.Y*t[2][1] + t[2][2]
}

// T returns the product of t with t2
func (t T) T(t2 T) T {
	return T{
		{
			t[0][0]*t2[0][0] + t[1][0]*t2[0][1] + t[2][0]*t2[0][2],
			t[0][1]*t2[0][0] + t[1][1]*t2[0][1] + t[2][1]*t2[0][2],
			t[0][2]*t2[0][0] + t[1][2]*t2[0][1] + t[2][2]*t2[0][2],
		}, {
			t[0][0]*t2[1][0] + t[1][0]*t2[1][1] + t[2][0]*t2[1][2],
			t[0][1]*t2[1][0] + t[1][1]*t2[1][1] + t[2][1]*t2[1][2],
			t[0][2]*t2[1][0] + t[1][2]*t2[1][1] + t[2][2]*t2[1][2],
		}, {
			t[0][0]*t2[2][0] + t[1][0]*t2[2][1] + t[2][0]*t2[2][2],
			t[0][1]*t2[2][0] + t[1][1]*t2[2][1] + t[2][1]*t2[2][2],
			t[0][2]*t2[2][0] + t[1][2]*t2[2][1] + t[2][2]*t2[2][2],
		},
	}
}

// Scale generates a scale transform
type Scale V

// T returns the scale transform
func (s Scale) T() T {
	return T{
		{s.X, 0, 0},
		{0, s.Y, 0},
		{0, 0, 1},
	}
}

// TInv returns the inverse of the scale transform
func (s Scale) TInv() T {
	return T{
		{1.0 / s.X, 0, 0},
		{0, 1.0 / s.Y, 0},
		{0, 0, 1},
	}
}

// Pair returns the Scale transform and it's inverse
func (s Scale) Pair() [2]T {
	return [2]T{
		{
			{s.X, 0, 0},
			{0, s.Y, 0},
			{0, 0, 1},
		}, {
			{1.0 / s.X, 0, 0},
			{0, 1.0 / s.Y, 0},
			{0, 0, 1},
		},
	}
}

// Rotate generates a rotation transform
type Rotate angle.Rad

// T returns the rotation transform
func (r Rotate) T() T {
	s, c := angle.Rad(r).Sincos()
	return T{
		{c, -s, 0},
		{s, c, 0},
		{0, 0, 1},
	}
}

// TInv returns the inverse of the rotation transform
func (r Rotate) TInv() T {
	s, c := angle.Rad(r).Sincos()
	return T{
		{c, s, 0},
		{-s, c, 0},
		{0, 0, 1},
	}
}

// Pair returns the rotation transform and it's inverse
func (r Rotate) Pair() [2]T {
	s, c := angle.Rad(r).Sincos()
	return [2]T{
		{
			{c, -s, 0},
			{s, c, 0},
			{0, 0, 1},
		}, {
			{c, s, 0},
			{-s, c, 0},
			{0, 0, 1},
		},
	}
}

// Translate generates a translation transform
type Translate V

// T returns the translation transform
func (t Translate) T() T {
	return T{
		{1, 0, t.X},
		{0, 1, t.Y},
		{0, 0, 1},
	}
}

// TInv returns the inverse of the translation transform.
func (t Translate) TInv() T {
	return T{
		{1, 0, -t.X},
		{0, 1, -t.Y},
		{0, 0, 1},
	}
}

// Pair returns the translation transform and it's inverse.
func (t Translate) Pair() [2]T {
	return [2]T{
		{
			{1, 0, t.X},
			{0, 1, t.Y},
			{0, 0, 1},
		}, {
			{1, 0, -t.X},
			{0, 1, -t.Y},
			{0, 0, 1},
		},
	}
}

// Chain combines multiple TGens into one
type Chain []TGen

// T does a forward multiplication through the chain returning the transform
func (c Chain) T() T {
	if len(c) == 0 {
		return IndentityTransform()
	}
	if len(c) == 1 {
		return c[0].T()
	}
	t := c[0].T().T(c[1].T())
	for _, nxt := range c[2:] {
		t = t.T(nxt.T())
	}
	return t
}

// TInv does a reverse multiplication through the chain returning the inverse of
// the transform.
func (c Chain) TInv() T {
	ln := len(c)
	if ln == 0 {
		return IndentityTransform()
	}
	if ln == 1 {
		return c[0].TInv()
	}
	t := c[ln-1].TInv().T(c[ln-2].TInv())
	for i := ln - 3; i >= 0; i-- {
		t = t.T(c[i].TInv())
	}
	return t
}

// Pair calls pair on all the TGen in the chain and computes both the transform
// and it's inverse.
func (c Chain) Pair() [2]T {
	ln := len(c)
	if ln == 0 {
		return [2]T{IndentityTransform(), IndentityTransform()}
	}
	if ln == 1 {
		return c[0].Pair()
	}
	ps := make([][2]T, ln)
	for i, t := range c {
		ps[i] = t.Pair()
	}
	out := [2]T{
		ps[0][0].T(ps[1][0]),
		ps[ln-1][1].T(ps[ln-2][1]),
	}
	ln--
	for i := 2; i <= ln; i++ {
		out[0] = out[0].T(ps[i][0])
		out[1] = out[1].T(ps[ln-i][1])
	}
	return out
}

func (t T) String() string {
	return strings.Join([]string{
		"T[ (",
		strconv.FormatFloat(t[0][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[0][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[0][2], 'f', Prec, 64),
		"), (",
		strconv.FormatFloat(t[1][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][2], 'f', Prec, 64),
		"), (",
		strconv.FormatFloat(t[2][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][2], 'f', Prec, 64),
		") ]",
	}, "")
}
