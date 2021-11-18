package descent

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/stretchr/testify/assert"
)

func BenchmarkDescent(b *testing.B) {
	bez := bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	l := line.New(d2.Pt{0, 200}, d2.Pt{500, 300})
	fn := func(t0, t1 float64) float64 {
		return bez.Pt1(t0).Distance(l.Pt1(t1))
	}
	solver := &Solver{
		Ln: 2,
		Fn: func(v []float64) float64 {
			return fn(v[0], v[1])
		},
	}
	solver.SetDFn(nil)
	buf1 := make([]float64, 2)
	buf2 := make([]float64, 2)

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		p := []float64{0.04, 0.96}
		steps := 20

		for i := 0; i < steps; i++ {
			solver.Step(p, buf1, buf2)
		}
	}
}

func TestPartialDerivative(t *testing.T) {
	bez := bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	l := line.New(d2.Pt{0, 200}, d2.Pt{500, 300})
	fn, exact := Distance2(bez, l)

	s := Solver{
		Ln: 2,
		Fn: fn,
	}
	s.SetDFn(nil)

	scale := (grid.Scale{
		X: 1.0 / 10.0,
		Y: 1.0 / 10.0,
	}).T
	buf0 := make([]float64, 2)
	buf1 := make([]float64, 2)
	for pt := range (grid.Pt{10, 10}).Iter().Chan() {
		t0, t1 := scale(pt)
		x := []float64{t0, t1}
		expected := exact(x, buf0)
		got := s.DFn(x, buf1)
		for i, g := range got {
			assert.InDelta(t, expected[i], g, 4e-4)
		}
	}
}
