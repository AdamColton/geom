package d3

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/geomtest"

	"github.com/stretchr/testify/assert"
)

func TestQ(t *testing.T) {
	s2inv := 1.0 / math.Sqrt2
	tt := []struct {
		in       Pt
		q        Q
		expected Pt
	}{
		{
			q:        Q{1, 0, 0, 0},
			in:       Pt{1, 0, 0},
			expected: Pt{1, 0, 0},
		}, {
			q:        Q{1, 0, 0, 0},
			in:       Pt{0, 1, 0},
			expected: Pt{0, 1, 0},
		}, {
			q:        Q{1, 0, 0, 0},
			in:       Pt{0, 0, 1},
			expected: Pt{0, 0, 1},
		}, {
			q:        Q{0, 0, 0, 1},
			in:       Pt{1, 0, 0},
			expected: Pt{-1, 0, 0},
		}, {
			q:        Q{0, 0, 0, 1},
			in:       Pt{0, 1, 0},
			expected: Pt{0, -1, 0},
		}, {
			q:        Q{0, 0, 1, 0},
			in:       Pt{1, 0, 0},
			expected: Pt{-1, 0, 0},
		}, {
			q:        Q{0, 0, 1, 0},
			in:       Pt{0, 0, 1},
			expected: Pt{0, 0, -1},
		}, {
			q:        Q{s2inv, 0, 0, s2inv},
			in:       Pt{1, 0, 0},
			expected: Pt{0, -1, 0},
		}, {
			q:        Q{s2inv, 0, 0, s2inv},
			in:       Pt{0, 1, 0},
			expected: Pt{1, 0, 0},
		}, {
			q:        Q{s2inv, 0, 0, -s2inv},
			in:       Pt{1, 0, 0},
			expected: Pt{0, 1, 0},
		}, {
			q:        Q{1, 0, 0, -1}.Normalize(),
			in:       Pt{1, 0, 0},
			expected: Pt{0, 1, 0},
		},
		{
			q:        Q{1, 0, 0, 0}.Product(Q{sqrtHalf, 0, 0, -sqrtHalf}).Normalize(),
			in:       Pt{1, 0, 0},
			expected: Pt{0, 1, 0},
		},
	}

	for _, tc := range tt {
		t.Run(tc.expected.String(), func(t *testing.T) {
			t.Log(tc.q.String())
			geomtest.Equal(t, tc.expected, tc.q.T().Pt(tc.in))
		})
	}
}

func TestNormalizeExceptions(t *testing.T) {
	assert.Equal(t, Q{}, Q{}.Normalize())
	assert.Equal(t, Q{1, 0, 0, 0}, Q{1, 0, 0, 0}.Normalize())
}

func TestQInv(t *testing.T) {
	s, c := angle.Deg(30).Sincos()
	tt := map[string]Q{
		"ident": Q{1, 0, 0, 0},
		"X":     Q{0, 1, 0, 0},
		"Y":     Q{0, 0, 1, 0},
		"Z":     Q{0, 0, 0, 1},
		"X30":   Q{s, c, 0, 0},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			inv, _ := tc.T().Inversion()
			geomtest.Equal(t, inv, tc.TInv())
		})
	}
}

func TestQRotate(t *testing.T) {
	tt := map[string]struct {
		Q
		In       Pt
		expected Pt
	}{
		"Z90": {
			Q:        QZ(angle.Deg(90)),
			In:       Pt{1, 0, 0},
			expected: Pt{0, 1, 0},
		},
		"Z45": {
			Q:        QZ(angle.Deg(45)),
			In:       Pt{1, 0, 0},
			expected: Pt{sqrtHalf, sqrtHalf, 0},
		},
		"X90": {
			Q:        QX(angle.Deg(90)),
			In:       Pt{0, 1, 0},
			expected: Pt{0, 0, 1},
		},
		"Y90": {
			Q:        QY(angle.Deg(90)),
			In:       Pt{1, 0, 0},
			expected: Pt{0, 0, -1},
		},
		"Z": {
			Q:        QY(angle.Deg(90)),
			In:       Pt{0, 0, 1},
			expected: Pt{1, 0, 0},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			geomtest.Equal(t, tc.expected, tc.Q.T().Pt(tc.In))
		})
	}
}

func TestQV(t *testing.T) {
	tt := map[string]V{
		"Basic2D":            V{2, 1, 0}.Normal(),
		"Basic2D-neg":        V{-2, 1, 0}.Normal(),
		"Y":                  V{0, 1, 0},
		"3D":                 V{0, 0, 1}.Normal(),
		"3Dcomplex":          V{1, 2, 3}.Normal(),
		"3DcomplexNonnormal": V{1, 2, 3},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			v := QV(tc).T().V(V{1, 0, 0})
			geomtest.Equal(t, tc.Normal(), v)
		})
	}
}
