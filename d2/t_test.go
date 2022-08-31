package d2

import (
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/geomerr"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	tt := []struct {
		*T
		expected Pt
	}{
		{
			T:        Translate(V{1, 2}).T(),
			expected: Pt{2, 3},
		},
		{
			T:        Rotate{angle.Rot(0.25)}.T(),
			expected: Pt{-1, 1},
		}, {
			T:        Scale(V{2, 3}).T(),
			expected: Pt{2, 3},
		},
		{
			T:        Chain{Rotate{angle.Rot(0.25)}, Translate(V{2, 2}), Scale(V{2, 3})}.T(),
			expected: Pt{2, 9},
		},
	}

	pt := Pt{1, 1}
	v := V{1, 1}
	for _, tc := range tt {
		t.Run(tc.expected.String(), func(t *testing.T) {
			geomtest.Equal(t, tc.expected, tc.T.Pt(pt))
			pf, _ := tc.T.PtF(pt)
			geomtest.Equal(t, tc.expected, pf)

			geomtest.Equal(t, tc.expected.V(), tc.T.V(v))
			vf, _ := tc.T.VF(V{1, 1})
			geomtest.Equal(t, tc.expected.V(), vf)
		})
	}
}

func TestTGen(t *testing.T) {
	tt := map[string]TGen{
		"scale":     Scale(V{3, 4}),
		"rotate":    Rotate{angle.Deg(87)},
		"translate": Translate(V{3, 7}),
		"chain":     Chain{Rotate{angle.Rot(0.5)}, Translate(V{2, 2}), Scale(V{2, 3})},
		"chain0":    Chain{},
		"chain1":    Chain{Scale(V{12, 13})},
	}

	for name, gen := range tt {
		t.Run(name, func(t *testing.T) {
			tr, ti := gen.T(), gen.TInv()
			p := gen.Pair()
			geomtest.Equal(t, tr, p[0])
			geomtest.Equal(t, ti, p[1])
			geomtest.Equal(t, IndentityTransform(), tr.T(ti))
		})
	}
}

func TestTString(t *testing.T) {
	tr := T{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	assert.Equal(t, "T[ (1.0000, 2.0000, 3.0000), (4.0000, 5.0000, 6.0000), (7.0000, 8.0000, 9.0000) ]", tr.String())
}

func TestTSlice(t *testing.T) {
	tr := Translate(V{3, 4}).T()
	got := tr.Slice([]Pt{{1, 1}, {2, 2}, {3, 3}})
	expected := []Pt{{4, 5}, {5, 6}, {6, 7}}
	geomtest.Equal(t, expected, got)

}

func TestTProd(t *testing.T) {
	tt := map[string]struct {
		prod     []*T
		expected *T
	}{
		"basic": {
			prod: []*T{
				{
					{2, 3, 5},
					{7, 11, 13},
					{0, 0, 1},
				}, {
					{17, 19, 23},
					{29, 31, 37},
					{0, 0, 1},
				},
			},
			expected: &T{
				{167, 260, 355},
				{275, 428, 585},
				{0, 0, 1},
			},
		},
		"triple": {
			prod: []*T{
				{
					{2, 3, 5},
					{7, 11, 13},
					{0, 0, 1},
				}, {
					{17, 19, 23},
					{29, 31, 37},
					{0, 0, 1},
				}, {
					{41, 43, 47},
					{53, 59, 61},
					{0, 0, 1},
				},
			},
			expected: &T{
				{18672, 29064, 39757},
				{25076, 39032, 53391},
				{0, 0, 1},
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			geomtest.Equal(t, tc.expected, TProd(tc.prod...))
		})
	}
}

func TestTAssertEqual(t *testing.T) {
	tr := Translate(V{3, 4}).T()
	err := tr.AssertEqual(Pt{1, 1}, 1e-6)
	assert.Equal(t, geomerr.TypeMismatch(tr, Pt{1, 1}), err)

	cp := *tr
	cp[1][1] += 2
	err = tr.AssertEqual(&cp, 1e-6)
	serr, ok := err.(geomerr.SliceErrs)
	if assert.True(t, ok) && assert.Len(t, serr, 1) && assert.Equal(t, serr[0].Index, 1) {
		serr, ok = serr[0].Err.(geomerr.SliceErrs)
		if assert.True(t, ok) && assert.Len(t, serr, 1) && assert.Equal(t, serr[0].Index, 1) {
			assert.Equal(t, geomerr.NotEqual(tr[1][1], cp[1][1]), serr[0].Err)
		}
	}
}

func TestTransformSet(t *testing.T) {
	trans := Translate(V{1, 2})
	scale := Scale(V{3, 4})
	rot := Rotate{angle.Rot(0.25)}.T()

	tr := NewTSet().
		AddBoth(trans).
		AddBoth(scale).
		Add(rot).
		GetT()

	tp := trans.Pair()
	sp := scale.Pair()
	expected := TProd(tp[0], sp[0], rot, sp[1], tp[1])

	geomtest.Equal(t, expected, tr)
}
