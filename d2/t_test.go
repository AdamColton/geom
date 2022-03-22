package d2

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/angle"
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
			T:        Rotate(math.Pi / 2).T(),
			expected: Pt{-1, 1},
		},
		{
			T:        Scale(V{2, 3}).T(),
			expected: Pt{2, 3},
		},
		{
			T:        Chain{Rotate(math.Pi / 2), Translate(V{2, 2}), Scale(V{2, 3})}.T(),
			expected: Pt{2, 9},
		},
	}

	for _, tc := range tt {
		t.Run(tc.expected.String(), func(t *testing.T) {
			p := tc.T.Pt(Pt{1, 1})
			d := p.Distance(tc.expected)
			assert.InDelta(t, 0, d, 1e-5, p.String())
			pf, _ := tc.T.PtF(Pt{1, 1})
			assert.Equal(t, p, pf)

			v := tc.T.V(V{1, 1})
			d = tc.expected.V().Subtract(v).Mag()
			assert.InDelta(t, 0, d, 1e-5, v.String())
			vf, _ := tc.T.VF(V{1, 1})
			assert.Equal(t, v, vf)
		})
	}
}

func TestTGen(t *testing.T) {
	tt := map[string]TGen{
		"scale":     Scale(V{3, 4}),
		"rotate":    Rotate(angle.Deg(87)),
		"translate": Translate(V{3, 7}),
		"chain":     Chain{Rotate(math.Pi / 2), Translate(V{2, 2}), Scale(V{2, 3})},
		"chain0":    Chain{},
		"chain1":    Chain{Scale(V{12, 13})},
	}

	for name, gen := range tt {
		t.Run(name, func(t *testing.T) {
			tr, ti := gen.T(), gen.TInv()
			p := gen.Pair()
			assert.Equal(t, tr, p[0])
			assert.Equal(t, ti, p[1])
			tApproxEqual(t, IndentityTransform(), tr.T(ti))
		})
	}
}

func tApproxEqual(t *testing.T, t1, t2 *T) {
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			assert.InDelta(t, t1[x][y], t2[x][y], 1e-10)
		}
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
	assert.Equal(t, expected, got)

}

func TestTPow(t *testing.T) {
	tt := map[string]struct {
		expected, t *T
		exp         int
	}{
		"translate": {
			t:        Translate(V{1, 2}).T(),
			exp:      5,
			expected: Translate(V{5, 10}).T(),
		},
		"rotate": {
			t:        Rotate(.01).T(),
			exp:      10,
			expected: Rotate(.1).T(),
		},
		"scale": {
			t:        Scale(V{2, 3}).T(),
			exp:      3,
			expected: Scale(V{8, 27}).T(),
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			geomtest.Equal(t, tc.expected, tc.t.Pow(uint(tc.exp)))
		})
	}
}
