package d2set

import (
	"fmt"
	"math"
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestTransformArray(t *testing.T) {
	ta := TransformPowArray{
		T: d2.Translate(d2.V{-1, 1}).T(),
		N: 3,
	}
	assert.Equal(t, 3, ta.Len())

	ps := PointSlice{{1, 1}, {2, 2}}

	tpl := TransformPointList{
		PointList:     ps,
		TransformList: ta,
	}

	expect := PointSlice{
		{1, 1}, {2, 2},
		{0, 2}, {1, 3},
		{-1, 3}, {0, 4},
	}

	geomtest.Equal(t, expect, tpl)
	geomtest.Equal(t, expect, NewPointSlice(tpl))

	ta.Offset = 1
	ta.N = 2
	tpl.TransformList = ta
	geomtest.Equal(t, expect[2:], tpl)
	geomtest.Equal(t, expect[2:], NewPointSlice(tpl))
}

func TestCopyPointSlice(t *testing.T) {
	expect := PointSlice{
		{1, 1}, {2, 2},
		{0, 2}, {1, 3},
		{-1, 3}, {0, 4},
	}
	geomtest.Equal(t, expect, NewPointSlice(expect))
}

func TestRotationArray(t *testing.T) {
	r := RotationArray{
		Arc:    angle.Deg(180),
		V:      d2.V{2, 0},
		N:      5,
		Source: PointSlice{{1, 1}},
	}
	sr2 := math.Sqrt2
	expect := PointSlice{
		{3, 1},
		{sr2 + 1, sr2 + 1},
		{1, 3},
		{1 - sr2, sr2 + 1},
		{-1, 1},
	}
	geomtest.Equal(t, expect, r)
}

func TestReflect(t *testing.T) {
	r := Reflect{
		Source: PointSlice{{1, 1}, {2, 2}, {3, 3}},
		Line:   line.New(d2.Pt{4, 0}, d2.Pt{4, 1}),
	}
	expect := PointSlice{
		{1, 1}, {2, 2}, {3, 3},
		{7, 1}, {6, 2}, {5, 3},
	}
	geomtest.Equal(t, expect, r)
}

func TestPt1Source(t *testing.T) {
	s := Pt1Source{
		Pt1: line.New(d2.Pt{0, 0}, d2.Pt{5, 5}),
	}
	s.Set(0, 1, 6)
	expect := PointSlice{
		{0, 0}, {1, 1}, {2, 2},
		{3, 3}, {4, 4}, {5, 5},
	}
	geomtest.Equal(t, expect, s)
}

func TestCrossComb(t *testing.T) {
	cc := CrossComb{}
	a := 3
	b := 4
	assert.Equal(t, 12, cc.Len(a, b))
	got := make(map[string]bool)
	gen := func(a, b int) string { return fmt.Sprintf("%d_%d", a, b) }
	for i := 0; i < 12; i++ {
		ai, bi := cc.Idx(i, a, b)
		got[gen(ai, bi)] = true
	}

	for da := 0; da < a; da++ {
		for db := 0; db < b; db++ {
			s := gen(da, db)
			assert.True(t, got[s], s)
		}
	}
}

func TestModComb(t *testing.T) {
	tt := map[string]struct {
		a, b   int
		expect [][2]int
	}{
		"3_6": {
			a: 3,
			b: 6,
			expect: [][2]int{
				{0, 0},
				{1, 1},
				{2, 2},
				{0, 3},
				{1, 4},
				{2, 5},
			},
		},
	}

	mc := ModComb{}
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, len(tc.expect), mc.Len(tc.a, tc.b))
			for i, e := range tc.expect {
				ga, gb := mc.Idx(i, tc.a, tc.b)
				assert.Equal(t, e[0], ga)
				assert.Equal(t, e[1], gb)
			}
		})
	}
}
