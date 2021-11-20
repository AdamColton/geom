package descent

import (
	"fmt"
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/d2/curve/distance"
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
	fn, exact := distance.New(bez, l)

	s := Solver{
		Ln: 2,
		Fn: func(x []float64) float64 {
			return fn(x[0], x[1])
		},
	}
	s.SetDFn(nil)

	scale := (grid.Scale{
		X: 1.0 / 10.0,
		Y: 1.0 / 10.0,
	}).T
	buf := make([]float64, 2)
	x := make([]float64, 2)
	expected := make([]float64, 2)
	for pt := range (grid.Pt{10, 10}).Iter().Chan() {
		x[0], x[1] = scale(pt)
		expected[0], expected[1] = exact(x[0], x[1])
		got := s.DFn(x, buf)
		for i, g := range got {
			assert.InDelta(t, expected[i], g, 4e-4)
		}
	}
}

func TestPartialDerivative2(t *testing.T) {
	bez := bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	l := line.New(d2.Pt{0, 200}, d2.Pt{500, 300})
	fn, exact := distance.New2(bez, l)

	s := Solver{
		Ln: 2,
		Fn: func(x []float64) float64 {
			return fn(x[0], x[1])
		},
	}
	s.SetDFn(nil)

	scale := (grid.Scale{
		X: 1.0 / 10.0,
		Y: 1.0 / 10.0,
	}).T
	buf := make([]float64, 2)
	x := make([]float64, 2)
	expected := make([]float64, 2)
	for pt := range (grid.Pt{10, 10}).Iter().Chan() {
		x[0], x[1] = scale(pt)
		expected[0], expected[1] = exact(x[0], x[1])
		got := s.DFn(x, buf)
		for i, g := range got {
			assert.InDelta(t, expected[i], g, 0.1)
		}
	}
}

func TestDescent(t *testing.T) {
	bez := bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	l := line.New(d2.Pt{0, 200}, d2.Pt{500, 300})
	fn, dfn := distance.New(bez, l)

	s := &Solver{
		Ln: 2,
		Fn: func(x []float64) float64 {
			return fn(x[0], x[1])
		},
		DFn: func(x, buf []float64) []float64 {
			buf[0], buf[1] = dfn(x[0], x[1])
			return buf
		},
	}
	s.SetDFn(nil)

	x := []float64{0, 0}
	buf1 := []float64{0, 0}
	buf2 := []float64{0, 0}
	for i := 0; i < 30; i++ {
		s.Step(x, buf1, buf2)
	}

	assert.InDelta(t, 0.0, s.Fn(x), 3e-3)
}

func TestSecant(t *testing.T) {
	bez := bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	l := line.New(d2.Pt{0, 200}, d2.Pt{500, 300})
	var fn Fn = func(x []float64) float64 {
		return bez.Pt1(x[0]).Distance(l.Pt1(x[1]))
	}

	x0, x1 := []float64{0, 0}, []float64{0.01, 0.01}
	for i := 0; i < 100; i++ {
		x0, x1 = fn.SecantStep(x0, x1), x0
		fmt.Println(x0, x1, fn(x0))
	}

	assert.InDelta(t, 0.0, fn(x0), 3e-3)

	assert.Equal(t, x0, fn.D(x0, x1))
	assert.Equal(t, x0, fn.D(x1, x0))
}
