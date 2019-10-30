package d3

import (
	"strconv"
	"strings"

	"github.com/adamcolton/geom/angle"
)

type T [4][4]float64

/*

| x |   | a b c d |
| y | * | e f g h | = | ax+by+cz+d ex+fy+gz+h ix+jy+kz+l mx+ny+pz+q |
| z |   | i j k l |
| 1 |   | m n p q |

| (0,0) (1,0) (2,0) |   | [0][0] [0][1] [0][2] |
| (0,1) (1,1) (2,1) | = | [1][0] [1][1] [1][2] |
| (0,2) (1,2) (2,2) |   | [2][0] [2][1] [2][2] |

| a b c |   | j k l |   | aj+bm+cp ak+bn+cq al+bo+cr |
| d e f | * | m n o | = | dj+em+fp dk+en+fq dl+eo+fr |
| g h i |   | p q r |   | gj+hm+ip gk+hn+iq gl+ho+ir |


| a1 b1 c1 d1 |   | a2 b2 c2 d2 |   | a1a2+b1e2+c1i2+d1m2 a1b2+b1f2+c1j2+d1n2 a1c2+b1g2+c1k2+d1p2 |
| e1 f1 g1 h1 | * | e2 f2 g2 h2 | = |
| i1 j1 k1 l1 |   | i2 j2 k2 l2 |   |
| m1 n1 p1 q1 |   | m2 n2 p2 q2 |   |
*/

func (t *T) Pt(pt Pt) Pt {
	return Pt{
		pt.X*t[0][0] + pt.Y*t[0][1] + pt.Z*t[0][2] + t[0][3],
		pt.X*t[1][0] + pt.Y*t[1][1] + pt.Z*t[1][2] + t[1][3],
		pt.X*t[2][0] + pt.Y*t[2][1] + pt.Z*t[2][2] + t[2][3],
	}
}

func (t *T) Pts(pts []Pt) []Pt {
	out := make([]Pt, len(pts))
	for i, pt := range pts {
		out[i] = t.Pt(pt)
	}
	return out
}

func (t *T) V(v V) V {
	return V{
		v.X*t[0][0] + v.Y*t[0][1] + v.Z*t[0][2] + t[0][3],
		v.X*t[1][0] + v.Y*t[1][1] + v.Z*t[1][2] + t[1][3],
		v.X*t[2][0] + v.Y*t[2][1] + v.Z*t[2][2] + t[2][3],
	}
}

func (t *T) PtF(pt Pt) (Pt, float64) {
	return Pt{
		pt.X*t[0][0] + pt.Y*t[0][1] + pt.Z*t[0][2] + t[0][3],
		pt.X*t[1][0] + pt.Y*t[1][1] + pt.Z*t[1][2] + t[1][3],
		pt.X*t[2][0] + pt.Y*t[2][1] + pt.Z*t[2][2] + t[2][3],
	}, pt.X*t[3][0] + pt.Y*t[3][1] + pt.Z*t[3][2] + t[3][3]
}

func (t *T) PtScl(pt Pt) Pt {
	w := pt.X*t[3][0] + pt.Y*t[3][1] + pt.Z*t[3][2] + t[3][3]
	return Pt{
		(pt.X*t[0][0] + pt.Y*t[0][1] + pt.Z*t[0][2] + t[0][3]) / w,
		(pt.X*t[1][0] + pt.Y*t[1][1] + pt.Z*t[1][2] + t[1][3]) / w,
		(pt.X*t[2][0] + pt.Y*t[2][1] + pt.Z*t[2][2] + t[2][3]) / w,
	}
}

func (t *T) PtsScl(pts []Pt) []Pt {
	out := make([]Pt, len(pts))
	for i, pt := range pts {
		out[i] = t.PtScl(pt)
	}
	return out
}

func (t *T) VF(v V) (V, float64) {
	return V{
		v.X*t[0][0] + v.Y*t[0][1] + v.Z*t[0][2] + t[0][3],
		v.X*t[1][0] + v.Y*t[1][1] + v.Z*t[1][2] + t[1][3],
		v.X*t[2][0] + v.Y*t[2][1] + v.Z*t[2][2] + t[2][3],
	}, v.X*t[3][0] + v.Y*t[3][1] + v.Z*t[3][2] + t[3][3]
}

func (t *T) T(t2 *T) *T {
	return &T{
		{
			t[0][0]*t2[0][0] + t[1][0]*t2[0][1] + t[2][0]*t2[0][2] + t[3][0]*t2[0][3],
			t[0][1]*t2[0][0] + t[1][1]*t2[0][1] + t[2][1]*t2[0][2] + t[3][1]*t2[0][3],
			t[0][2]*t2[0][0] + t[1][2]*t2[0][1] + t[2][2]*t2[0][2] + t[3][2]*t2[0][3],
			t[0][3]*t2[0][0] + t[1][3]*t2[0][1] + t[2][3]*t2[0][2] + t[3][3]*t2[0][3],
		}, {
			t[0][0]*t2[1][0] + t[1][0]*t2[1][1] + t[2][0]*t2[1][2] + t[3][0]*t2[1][3],
			t[0][1]*t2[1][0] + t[1][1]*t2[1][1] + t[2][1]*t2[1][2] + t[3][1]*t2[1][3],
			t[0][2]*t2[1][0] + t[1][2]*t2[1][1] + t[2][2]*t2[1][2] + t[3][2]*t2[1][3],
			t[0][3]*t2[1][0] + t[1][3]*t2[1][1] + t[2][3]*t2[1][2] + t[3][3]*t2[1][3],
		}, {
			t[0][0]*t2[2][0] + t[1][0]*t2[2][1] + t[2][0]*t2[2][2] + t[3][0]*t2[2][3],
			t[0][1]*t2[2][0] + t[1][1]*t2[2][1] + t[2][1]*t2[2][2] + t[3][1]*t2[2][3],
			t[0][2]*t2[2][0] + t[1][2]*t2[2][1] + t[2][2]*t2[2][2] + t[3][2]*t2[2][3],
			t[0][3]*t2[2][0] + t[1][3]*t2[2][1] + t[2][3]*t2[2][2] + t[3][3]*t2[2][3],
		}, {
			t[0][0]*t2[3][0] + t[1][0]*t2[3][1] + t[2][0]*t2[3][2] + t[3][0]*t2[3][3],
			t[0][1]*t2[3][0] + t[1][1]*t2[3][1] + t[2][1]*t2[3][2] + t[3][1]*t2[3][3],
			t[0][2]*t2[3][0] + t[1][2]*t2[3][1] + t[2][2]*t2[3][2] + t[3][2]*t2[3][3],
			t[0][3]*t2[3][0] + t[1][3]*t2[3][1] + t[2][3]*t2[3][2] + t[3][3]*t2[3][3],
		},
	}
}

func Identity() *T {
	return &T{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

type Scale V

func (s Scale) T() *T {
	return &T{
		{s.X, 0, 0, 0},
		{0, s.Y, 0, 0},
		{0, 0, s.Z, 0},
		{0, 0, 0, 1},
	}
}

func (s Scale) TInv() *T {
	return &T{
		{1.0 / s.X, 0, 0, 0},
		{0, 1.0 / s.Y, 0, 0},
		{0, 0, 1.0 / s.Z, 0},
		{0, 0, 0, 1},
	}
}

// Scale returns the scale transform represented by V and it's inverse
func (s Scale) Pair() [2]*T {
	return [2]*T{
		&T{
			{s.X, 0, 0, 0},
			{0, s.Y, 0, 0},
			{0, 0, s.Z, 0},
			{0, 0, 0, 1},
		},
		&T{
			{1.0 / s.X, 0, 0, 0},
			{0, 1.0 / s.Y, 0, 0},
			{0, 0, 1.0 / s.Z, 0},
			{0, 0, 0, 1},
		},
	}
}

func ScaleF(f float64) Scale {
	return Scale(V{f, f, f})
}

type Translate V

func (t Translate) T() *T {
	return &T{
		{1, 0, 0, t.X},
		{0, 1, 0, t.Y},
		{0, 0, 1, t.Z},
		{0, 0, 0, 1},
	}
}

func (t Translate) TInv() *T {
	return &T{
		{1, 0, 0, -t.X},
		{0, 1, 0, -t.Y},
		{0, 0, 1, -t.Z},
		{0, 0, 0, 1},
	}
}

// Translate returns the translate transform represented by V and it's inverse
func (t Translate) Pair() [2]*T {
	return [2]*T{
		&T{
			{1, 0, 0, t.X},
			{0, 1, 0, t.Y},
			{0, 0, 1, t.Z},
			{0, 0, 0, 1},
		}, &T{
			{1, 0, 0, -t.X},
			{0, 1, 0, -t.Y},
			{0, 0, 1, -t.Z},
			{0, 0, 0, 1},
		},
	}
}

type RotationPlane byte

const (
	XY RotationPlane = iota
	XZ
	YZ
)

type Rotation struct {
	Angle angle.Rad
	Plane RotationPlane
}

func (r Rotation) T() *T {
	s, c := r.Angle.Sincos()
	if r.Plane == XZ {
		return &T{
			{c, 0, -s, 0},
			{0, 1, 0, 0},
			{s, 0, c, 0},
			{0, 0, 0, 1},
		}
	}
	if r.Plane == YZ {
		return &T{
			{1, 0, 0, 0},
			{0, c, -s, 0},
			{0, s, c, 0},
			{0, 0, 0, 1},
		}
	}
	return &T{
		{c, -s, 0, 0},
		{s, c, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (r Rotation) TInv() *T {
	s, c := r.Angle.Sincos()
	if r.Plane == XZ {
		return &T{
			{c, 0, s, 0},
			{0, 1, 0, 0},
			{-s, 0, c, 0},
			{0, 0, 0, 1},
		}
	}
	if r.Plane == YZ {
		return &T{
			{1, 0, 0, 0},
			{0, c, s, 0},
			{0, -s, c, 0},
			{0, 0, 0, 1},
		}
	}
	return &T{
		{c, s, 0, 0},
		{-s, c, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (r Rotation) Pair() [2]*T {
	s, c := r.Angle.Sincos()
	if r.Plane == XZ {
		return [2]*T{
			&T{
				{c, 0, -s, 0},
				{0, 1, 0, 0},
				{s, 0, c, 0},
				{0, 0, 0, 1},
			}, &T{
				{c, 0, s, 0},
				{0, 1, 0, 0},
				{-s, 0, c, 0},
				{0, 0, 0, 1},
			},
		}
	}
	if r.Plane == YZ {
		return [2]*T{
			&T{
				{1, 0, 0, 0},
				{0, c, -s, 0},
				{0, s, c, 0},
				{0, 0, 0, 1},
			}, &T{
				{1, 0, 0, 0},
				{0, c, s, 0},
				{0, -s, c, 0},
				{0, 0, 0, 1},
			},
		}
	}
	return [2]*T{
		&T{
			{c, -s, 0, 0},
			{s, c, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		}, &T{
			{c, s, 0, 0},
			{-s, c, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		},
	}
}

func (t T) String() string {
	return strings.Join([]string{
		"T[ (",
		strconv.FormatFloat(t[0][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[0][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[0][2], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[0][3], 'f', Prec, 64),
		"), (",
		strconv.FormatFloat(t[1][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][2], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][3], 'f', Prec, 64),
		"), (",
		strconv.FormatFloat(t[2][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][2], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][3], 'f', Prec, 64),
		"), (",
		strconv.FormatFloat(t[3][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[3][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[3][2], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[3][3], 'f', Prec, 64),
		") ]",
	}, "")
}

func (t *T) Inversion() (*T, bool) {
	//https://stackoverflow.com/questions/1148309/inverting-a-4x4-matrix
	out := &T{}

	out[0][0] = 0 +
		t[1][1]*t[2][2]*t[3][3] -
		t[1][1]*t[2][3]*t[3][2] -
		t[2][1]*t[1][2]*t[3][3] +
		t[2][1]*t[1][3]*t[3][2] +
		t[3][1]*t[1][2]*t[2][3] -
		t[3][1]*t[1][3]*t[2][2]

	out[1][0] = 0 -
		t[1][0]*t[2][2]*t[3][3] +
		t[1][0]*t[2][3]*t[3][2] +
		t[2][0]*t[1][2]*t[3][3] -
		t[2][0]*t[1][3]*t[3][2] -
		t[3][0]*t[1][2]*t[2][3] +
		t[3][0]*t[1][3]*t[2][2]

	out[2][0] = 0 +
		t[1][0]*t[2][1]*t[3][3] -
		t[1][0]*t[2][3]*t[3][1] -
		t[2][0]*t[1][1]*t[3][3] +
		t[2][0]*t[1][3]*t[3][1] +
		t[3][0]*t[1][1]*t[2][3] -
		t[3][0]*t[1][3]*t[2][1]

	out[3][0] = 0 -
		t[1][0]*t[2][1]*t[3][2] +
		t[1][0]*t[2][2]*t[3][1] +
		t[2][0]*t[1][1]*t[3][2] -
		t[2][0]*t[1][2]*t[3][1] -
		t[3][0]*t[1][1]*t[2][2] +
		t[3][0]*t[1][2]*t[2][1]

	out[0][1] = 0 -
		t[0][1]*t[2][2]*t[3][3] +
		t[0][1]*t[2][3]*t[3][2] +
		t[2][1]*t[0][2]*t[3][3] -
		t[2][1]*t[0][3]*t[3][2] -
		t[3][1]*t[0][2]*t[2][3] +
		t[3][1]*t[0][3]*t[2][2]

	out[1][1] = 0 +
		t[0][0]*t[2][2]*t[3][3] -
		t[0][0]*t[2][3]*t[3][2] -
		t[2][0]*t[0][2]*t[3][3] +
		t[2][0]*t[0][3]*t[3][2] +
		t[3][0]*t[0][2]*t[2][3] -
		t[3][0]*t[0][3]*t[2][2]

	out[2][1] = 0 -
		t[0][0]*t[2][1]*t[3][3] +
		t[0][0]*t[2][3]*t[3][1] +
		t[2][0]*t[0][1]*t[3][3] -
		t[2][0]*t[0][3]*t[3][1] -
		t[3][0]*t[0][1]*t[2][3] +
		t[3][0]*t[0][3]*t[2][1]

	out[3][1] = 0 +
		t[0][0]*t[2][1]*t[3][2] -
		t[0][0]*t[2][2]*t[3][1] -
		t[2][0]*t[0][1]*t[3][2] +
		t[2][0]*t[0][2]*t[3][1] +
		t[3][0]*t[0][1]*t[2][2] -
		t[3][0]*t[0][2]*t[2][1]

	out[0][2] = 0 +
		t[0][1]*t[1][2]*t[3][3] -
		t[0][1]*t[1][3]*t[3][2] -
		t[1][1]*t[0][2]*t[3][3] +
		t[1][1]*t[0][3]*t[3][2] +
		t[3][1]*t[0][2]*t[1][3] -
		t[3][1]*t[0][3]*t[1][2]

	out[1][2] = 0 -
		t[0][0]*t[1][2]*t[3][3] +
		t[0][0]*t[1][3]*t[3][2] +
		t[1][0]*t[0][2]*t[3][3] -
		t[1][0]*t[0][3]*t[3][2] -
		t[3][0]*t[0][2]*t[1][3] +
		t[3][0]*t[0][3]*t[1][2]

	out[2][2] = 0 +
		t[0][0]*t[1][1]*t[3][3] -
		t[0][0]*t[1][3]*t[3][1] -
		t[1][0]*t[0][1]*t[3][3] +
		t[1][0]*t[0][3]*t[3][1] +
		t[3][0]*t[0][1]*t[1][3] -
		t[3][0]*t[0][3]*t[1][1]

	out[3][2] = 0 -
		t[0][0]*t[1][1]*t[3][2] +
		t[0][0]*t[1][2]*t[3][1] +
		t[1][0]*t[0][1]*t[3][2] -
		t[1][0]*t[0][2]*t[3][1] -
		t[3][0]*t[0][1]*t[1][2] +
		t[3][0]*t[0][2]*t[1][1]

	out[0][3] = 0 -
		t[0][1]*t[1][2]*t[2][3] +
		t[0][1]*t[1][3]*t[2][2] +
		t[1][1]*t[0][2]*t[2][3] -
		t[1][1]*t[0][3]*t[2][2] -
		t[2][1]*t[0][2]*t[1][3] +
		t[2][1]*t[0][3]*t[1][2]

	out[1][3] = 0 +
		t[0][0]*t[1][2]*t[2][3] -
		t[0][0]*t[1][3]*t[2][2] -
		t[1][0]*t[0][2]*t[2][3] +
		t[1][0]*t[0][3]*t[2][2] +
		t[2][0]*t[0][2]*t[1][3] -
		t[2][0]*t[0][3]*t[1][2]

	out[2][3] = 0 -
		t[0][0]*t[1][1]*t[2][3] +
		t[0][0]*t[1][3]*t[2][1] +
		t[1][0]*t[0][1]*t[2][3] -
		t[1][0]*t[0][3]*t[2][1] -
		t[2][0]*t[0][1]*t[1][3] +
		t[2][0]*t[0][3]*t[1][1]

	out[3][3] = 0 +
		t[0][0]*t[1][1]*t[2][2] -
		t[0][0]*t[1][2]*t[2][1] -
		t[1][0]*t[0][1]*t[2][2] +
		t[1][0]*t[0][2]*t[2][1] +
		t[2][0]*t[0][1]*t[1][2] -
		t[2][0]*t[0][2]*t[1][1]

	det := t[0][0]*out[0][0] + t[0][1]*out[1][0] + t[0][2]*out[2][0] + t[0][3]*out[3][0]

	if det == 0 {
		return out, false
	}

	det = 1.0 / det

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			out[y][x] *= det
		}
	}

	return out, true
}

type TransformSet struct {
	Head, Middle, Tail []*T
}

func NewTSet() *TransformSet {
	return &TransformSet{}
}

func (ts *TransformSet) AddBoth(t [2]*T) *TransformSet {
	ts.Head = append(ts.Head, t[0])
	ts.Tail = append(ts.Tail, t[1])
	return ts
}

func (ts *TransformSet) Add(t *T) *TransformSet {
	ts.Middle = append(ts.Middle, t)
	return ts
}

func (ts *TransformSet) Get() *T {
	t := Identity()
	for _, th := range ts.Head {
		t = t.T(th)
	}
	for _, tm := range ts.Middle {
		t = t.T(tm)
	}
	for i := len(ts.Tail) - 1; i >= 0; i-- {
		t = t.T(ts.Tail[i])
	}
	return t
}
