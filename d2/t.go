package d2

import (
	"strconv"
	"strings"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
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

var indentityTransform = T{
	{1, 0, 0},
	{0, 1, 0},
	{0, 0, 1},
}

type IndentityTransform struct{}

func (IndentityTransform) GetT() *T {
	return &indentityTransform
}
func (IndentityTransform) TInv() *T {
	return &indentityTransform
}
func (IndentityTransform) Pair() Pair {
	return Pair{&indentityTransform, &indentityTransform}
}

// PtF applies the transform to a Pt, returning the resulting Pt and scale.
func (t *T) PtF(pt Pt) (Pt, float64) {
	return Pt{
		pt.X*t[0][0] + pt.Y*t[0][1] + t[0][2],
		pt.X*t[1][0] + pt.Y*t[1][1] + t[1][2],
	}, pt.X*t[2][0] + pt.Y*t[2][1] + t[2][2]
}

// Pt applies the transform to a Pt.
func (t *T) Pt(pt Pt) Pt {
	return Pt{
		pt.X*t[0][0] + pt.Y*t[0][1] + t[0][2],
		pt.X*t[1][0] + pt.Y*t[1][1] + t[1][2],
	}
}

// Slice applies the transform to a slice of Pts
func (t *T) Slice(pts []Pt) []Pt {
	out := make([]Pt, len(pts))
	for i, pt := range pts {
		out[i] = t.Pt(pt)
	}
	return out
}

// V applies the transform to a V
func (t *T) V(v V) V {
	return V{
		v.X*t[0][0] + v.Y*t[0][1] + t[0][2],
		v.X*t[1][0] + v.Y*t[1][1] + t[1][2],
	}
}

// VF applies the transform and returns a V and the scale
func (t *T) VF(v V) (V, float64) {
	return V{
		v.X*t[0][0] + v.Y*t[0][1] + t[0][2],
		v.X*t[1][0] + v.Y*t[1][1] + t[1][2],
	}, v.X*t[2][0] + v.Y*t[2][1] + t[2][2]
}

// TProd returns the product of multiple transforms.
func TProd(ts ...*T) *T {
	if len(ts) == 0 {
		return &indentityTransform
	}
	t := ts[0]
	for _, t2 := range ts[1:] {
		t = t.T(t2)
	}
	return t
}

// T returns the product of t with t2
func (t *T) T(t2 *T) *T {
	return &T{
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

// AssertEqual fulfils geomtest.AssertEqualizer
func (t *T) AssertEqual(actual interface{}, tol cmpr.Tolerance) error {
	t2, ok := actual.(*T)
	if !ok {
		return geomerr.TypeMismatch(t, actual)
	}

	return geomerr.NewSliceErrs(3, -1, func(x int) error {
		return geomerr.NewSliceErrs(3, -1, func(y int) error {
			if !tol.Equal(t[x][y], t2[x][y]) {
				return geomerr.NotEqual(t[x][y], t2[x][y])
			}
			return nil
		})
	})
}

// GetT fulfills TGen.
func (t *T) GetT() *T {
	return t
}

// Pair fulfills TGen. It returns a Pair with t and it's inverse.
func (t *T) Pair() Pair {
	return Pair{
		t,
		t.TInv(),
	}
}

// TInv finds the inverse of T assuming the inverse is computable.
func (t *T) TInv() *T {
	inv, _ := t.Inversion()
	return inv
}

// Inversion computes the inverse of t and a bool indicating of the inverse is
// computable. The inverse is not computable if the determinate is 0. If the
// determinate is zero, the returned *T will be the inversion before scaling by
// 1/determinate.
func (t *T) Inversion() (*T, bool) {
	//https://stackoverflow.com/questions/983999/simple-3x3-matrix-inverse-code-c
	out := &T{}
	a := t[1][0] * t[2][2]
	b := t[1][2] * t[2][0]
	out[0][0] = (t[1][1]*t[2][2] - t[2][1]*t[1][2])
	out[0][1] = (t[0][2]*t[2][1] - t[0][1]*t[2][2])
	out[0][2] = (t[0][1]*t[1][2] - t[0][2]*t[1][1])
	out[1][0] = (b - a)
	out[1][1] = (t[0][0]*t[2][2] - t[0][2]*t[2][0])
	out[1][2] = (t[1][0]*t[0][2] - t[0][0]*t[1][2])
	out[2][0] = (t[1][0]*t[2][1] - t[1][1]*t[2][0])
	out[2][1] = (t[2][0]*t[0][1] - t[0][0]*t[2][1])
	out[2][2] = (t[0][0]*t[1][1] - t[1][0]*t[0][1])

	det := t[0][0]*(out[0][0]) - t[0][1]*(a-b) + t[0][2]*(out[2][0])
	if det == 0 {
		return out, false
	}

	det = 1.0 / det

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			out[y][x] *= det
		}
	}
	return out, true
}

// Scale generates a scale transform
type Scale V

// T returns the scale transform
func (s Scale) GetT() *T {
	return &T{
		{s.X, 0, 0},
		{0, s.Y, 0},
		{0, 0, 1},
	}
}

// TInv returns the inverse of the scale transform
func (s Scale) TInv() *T {
	return &T{
		{1.0 / s.X, 0, 0},
		{0, 1.0 / s.Y, 0},
		{0, 0, 1},
	}
}

// Pair returns the Scale transform and it's inverse
func (s Scale) Pair() Pair {
	return Pair{
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
type Rotate struct{ angle.Sincoser }

// T returns the rotation transform
func (r Rotate) GetT() *T {
	s, c := r.Sincos()
	return &T{
		{c, -s, 0},
		{s, c, 0},
		{0, 0, 1},
	}
}

// TInv returns the inverse of the rotation transform
func (r Rotate) TInv() *T {
	s, c := r.Sincos()
	return &T{
		{c, s, 0},
		{-s, c, 0},
		{0, 0, 1},
	}
}

// Pair returns the rotation transform and it's inverse
func (r Rotate) Pair() Pair {
	s, c := r.Sincos()
	return Pair{
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
func (t Translate) GetT() *T {
	return &T{
		{1, 0, t.X},
		{0, 1, t.Y},
		{0, 0, 1},
	}
}

// TInv returns the inverse of the translation transform.
func (t Translate) TInv() *T {
	return &T{
		{1, 0, -t.X},
		{0, 1, -t.Y},
		{0, 0, 1},
	}
}

// Pair returns the translation transform and it's inverse.
func (t Translate) Pair() Pair {
	return Pair{
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
func (c Chain) GetT() *T {
	if len(c) == 0 {
		return &indentityTransform
	}
	if len(c) == 1 {
		return c[0].GetT()
	}
	t := c[0].GetT().T(c[1].GetT())
	for _, nxt := range c[2:] {
		t = t.T(nxt.GetT())
	}
	return t
}

// TInv does a reverse multiplication through the chain returning the inverse of
// the transform.
func (c Chain) TInv() *T {
	ln := len(c)
	if ln == 0 {
		return &indentityTransform
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
func (c Chain) Pair() Pair {
	ln := len(c)
	if ln == 0 {
		return IndentityTransform{}.Pair()
	}
	if ln == 1 {
		return c[0].Pair()
	}
	ps := make([][2]*T, ln)
	for i, t := range c {
		ps[i] = t.Pair()
	}
	out := [2]*T{
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

func (t *T) String() string {
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

// TransformSet builds up a chain of transformaitions.
type TransformSet struct {
	Head, Middle, Tail []*T
}

// NewTSet creates a TransformSet.
func NewTSet() *TransformSet {
	return &TransformSet{}
}

// AddBoth appends the transform and it's inverse to the head and tail.
func (ts *TransformSet) AddBoth(t TGen) *TransformSet {
	p := t.Pair()
	ts.Head = append(ts.Head, p[0])
	ts.Tail = append(ts.Tail, p[1])
	return ts
}

// Add t to the middle
func (ts *TransformSet) Add(t *T) *TransformSet {
	ts.Middle = append(ts.Middle, t)
	return ts
}

// Get produces a transform produces a transform by applying the transforms in
// head, then middle then applying tail in reverse.
func (ts *TransformSet) GetT() *T {
	h := TProd(ts.Head...)
	m := TProd(ts.Middle...)
	var t *T
	if ln := len(ts.Tail); ln > 0 {
		t = ts.Tail[ln-1]
		for i := ln - 2; i >= 0; i-- {
			t = t.T(ts.Tail[i])
		}
	}
	return TProd(h, m, t)
}

// Pair represents a Transform and it's Inverse. It fullfills TGen.
type Pair [2]*T

// GetT fulfills TGen.
func (p Pair) GetT() *T {
	return p[0]
}

// TInv fulfills TGen.
func (p Pair) TInv() *T {
	return p[1]
}

// Pair fulfills TGen.
func (p Pair) Pair() Pair {
	return p
}
