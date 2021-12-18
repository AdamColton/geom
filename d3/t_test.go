package d3

import (
	"math/rand"
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/geomerr"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestTAssertEqual(t *testing.T) {
	t1 := &T{}
	t2 := &T{}

	err := t1.AssertEqual(t2, 1e-10)
	assert.NoError(t, err)

	t2[2][2] = 1
	err = t1.AssertEqual(t2, 1e-10)
	assert.Equal(t, "\t2: \t2: Expected 0 got 1", err.Error())

	err = t1.AssertEqual(1.0, 1e-10)
	assert.IsType(t, geomerr.ErrTypeMismatch{}, err)
}

func TestInverse(t *testing.T) {
	tr := Scale(V{1, 2, 3}).T()
	trI, ok := tr.Inversion()
	assert.True(t, ok)

	assert.Equal(t, Identity(), tr.T(trI))
	ident := Identity()

	for i := 0; i < 100; i++ {
		for x := 0; x < 4; x++ {
			for y := 0; y < 4; y++ {
				tr[x][y] = rand.Float64()
			}
		}

		trI, ok = tr.Inversion()
		if ok {
			shouldBeIdent := tr.T(trI)
			for x := 0; x < 4; x++ {
				for y := 0; y < 4; y++ {
					tr[x][y] = rand.Float64()
					assert.InDelta(t, ident[x][y], shouldBeIdent[x][y], 1e-5)
				}
			}
		}
	}
}

func TestRotate(t *testing.T) {

	EqualPt(t, Pt{0, 1, 0}, Rotation{angle.Rot(0.25), XY}.T().Pt(Pt{1, 0, 0}))
	EqualPt(t, Pt{-1, 0, 0}, Rotation{angle.Rot(0.50), XY}.T().Pt(Pt{1, 0, 0}))
	EqualPt(t, Pt{0, -1, 0}, Rotation{angle.Rot(0.75), XY}.T().Pt(Pt{1, 0, 0}))

	EqualPt(t, Pt{0, 0, 1}, Rotation{angle.Rot(0.25), XZ}.T().Pt(Pt{1, 0, 0}))
	EqualPt(t, Pt{-1, 0, 0}, Rotation{angle.Rot(0.50), XZ}.T().Pt(Pt{1, 0, 0}))
	EqualPt(t, Pt{0, 0, -1}, Rotation{angle.Rot(0.75), XZ}.T().Pt(Pt{1, 0, 0}))

	EqualPt(t, Pt{0, 0, 1}, Rotation{angle.Rot(0.25), YZ}.T().Pt(Pt{0, 1, 0}))
	EqualPt(t, Pt{0, -1, 0}, Rotation{angle.Rot(0.50), YZ}.T().Pt(Pt{0, 1, 0}))
	EqualPt(t, Pt{0, 0, -1}, Rotation{angle.Rot(0.75), YZ}.T().Pt(Pt{0, 1, 0}))

	for r := (Rotation{angle.Rot(0), XY}); r.Angle.Rot() < 1.0; r.Angle += angle.Rot(0.01) {
		p := r.Pair()
		geomtest.Equal(t, Identity(), p[0].T(p[1]))
	}

	for r := (Rotation{angle.Rot(0), XZ}); r.Angle.Rot() < 1.0; r.Angle += angle.Rot(0.01) {
		p := r.Pair()
		geomtest.Equal(t, Identity(), p[0].T(p[1]))
	}

	for r := (Rotation{angle.Rot(0), YZ}); r.Angle.Rot() < 1.0; r.Angle += angle.Rot(0.01) {
		p := r.Pair()
		geomtest.Equal(t, Identity(), p[0].T(p[1]))
	}
}

func TestT(t *testing.T) {
	tt := []struct {
		p        Pt
		t        *T
		expected Pt
	}{
		{
			p:        Pt{1, 0, 0},
			t:        Identity(),
			expected: Pt{1, 0, 0},
		}, {
			p:        Pt{1, 0, 0},
			t:        Translate(V{1, 2, 3}).T(),
			expected: Pt{2, 2, 3},
		}, {
			p:        Pt{1, 0, 0},
			t:        Rotation{angle.Rot(0.25), XY}.T(),
			expected: Pt{0, 1, 0},
		}, {
			p:        Pt{1, 0, 0},
			t:        Rotation{angle.Rot(0.25), XZ}.T(),
			expected: Pt{0, 0, 1},
		}, {
			p: Pt{2, 1, 0},
			t: NewTSet().
				AddBoth(Translate(V{-1, -1, 0}).Pair()).
				Add(Rotation{angle.Rot(0.25), XZ}.T()).
				Get(),
			expected: Pt{2, 2, 2}.Multiply(0.5),
		}, {
			p:        Pt{1, 2, 3},
			t:        ScaleF(2).T(),
			expected: Pt{2, 4, 6},
		},
	}

	for _, tc := range tt {
		t.Run(tc.t.String(), func(t *testing.T) {
			t.Log(tc.p)
			p, w := tc.t.PtF(tc.p)
			assert.Equal(t, 1.0, w) // for these, w should always be 1
			EqualPt(t, tc.expected, p)
			EqualPt(t, tc.expected, tc.t.Pt(tc.p))
			EqualPt(t, tc.expected, tc.t.PtScl(tc.p))

			v := V(tc.p)
			tv, w := tc.t.VF(v)
			assert.Equal(t, 1.0, w)
			EqualV(t, V(tc.expected), tv)
			EqualV(t, V(tc.expected), tc.t.V(v))
		})
	}
}

func TestTGen(t *testing.T) {
	tt := map[string]TGen{
		"Scale":     Scale(V{1, 2, 3}),
		"Translate": Translate(V{1, 2, 3}),
		"RotateXY":  Rotation{angle.Deg(31), XY},
		"RotateXZ":  Rotation{angle.Deg(31), XZ},
		"RotateYZ":  Rotation{angle.Deg(31), YZ},
	}

	id := Identity()
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			p := GetTPair(tc)
<<<<<<< HEAD
			geomtest.Equal(t, p[0], tc.T())
			geomtest.Equal(t, p[1], GetTInv(tc))
			geomtest.Equal(t, id, p[0].T(p[1]))
=======
			TEqual(t, p[0], tc.T())
			TEqual(t, p[1], tc.TInv())
			TEqual(t, id, p[0].T(p[1]))
>>>>>>> Pair should not be part of TGen interface
		})
	}
}

func TestDetZero(t *testing.T) {
	_, ok := (&T{}).Inversion()
	assert.False(t, ok)
}

func TestPts(t *testing.T) {
	pts := Scale(V{1, 2, 3}).T().Pts([]Pt{{1, 1, 1}, {2, 2, 2}})
	assert.Equal(t, []Pt{{1, 2, 3}, {2, 4, 6}}, pts)
}

func TestPtsScale(t *testing.T) {
	tr := Rotation{angle.Deg(90), XY}.T()
	tr[3][3] = 2
	got := tr.PtsScl([]Pt{
		{1, 0, 0},
		{0, 1, 0},
	})
	expected := []Pt{
		{0, 0.5, 0},
		{-0.5, 0, 0},
	}
	for i, exp := range expected {
		EqualPt(t, exp, got[i])
	}
}
