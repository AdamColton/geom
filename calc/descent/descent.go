package descent

import (
	"github.com/adamcolton/geom/calc/cmpr"
)

type Fn func([]float64) float64

type Solver struct {
	Ln  int
	Fn  Fn
	DFn []Fn
}

func (s *Solver) SetDFn() {
	if len(s.DFn) < s.Ln {
		cp := make([]Fn, s.Ln)
		copy(cp, s.DFn)
		s.DFn = cp
	}
	for i := 0; i < s.Ln; i++ {
		if s.DFn[i] == nil {
			s.DFn[i] = s.Fn.CurryPartialDerivative(i)
		}
	}
}

func (s *Solver) Derivative(x, buf []float64) []float64 {
	var out []float64
	if len(buf) < s.Ln {
		out = make([]float64, s.Ln)
	} else {
		out = buf[:s.Ln]
	}

	for i, df := range s.DFn {
		out[i] = df(x)
	}
	return out
}

// m is momentum
func (s *Solver) Step(x, buf1, buf2 []float64) []float64 {
	d := s.Derivative(x, buf1)
	for i := range d {
		d[i] *= -1
	}
	g := s.G(x, d, buf2)
	for i := range d {
		d[i] = d[i] * g
		x[i] += d[i]
	}
	return d
}

func (s *Solver) G(x, d, buf []float64) float64 {
	var dx []float64
	if len(buf) < s.Ln {
		dx = make([]float64, s.Ln)
	} else {
		dx = buf[:s.Ln]
	}
	var fn SFn = func(g float64) float64 {
		for i, di := range d {
			dx[i] = x[i] + di*g
		}
		return s.Fn(dx)
	}
	return fn.G()
}

func (fn Fn) PartialDerivative(x []float64, idx int) float64 {
	// d is the delta between x+step and x-step
	// first we're looking for a step size that gets d close to zero
	const small cmpr.Tolerance = 1e-3
	step := 1e-2
	d := 1.0
	xi := x[idx]
	for !small.Zero(d) {
		step /= 2
		x[idx] = xi + step
		d = fn(x)
		x[idx] = xi - step
		d -= fn(x)
		x[idx] = xi
	}
	return d / (2 * step)
}

func (fn Fn) CurryPartialDerivative(idx int) Fn {
	return func(x []float64) float64 {
		return fn.PartialDerivative(x, idx)
	}
}

type SFn func(float64) float64

func (fn SFn) D(x float64) float64 {
	const small cmpr.Tolerance = 1e-3
	step := 1e-2
	d := 1.0
	for !small.Zero(step) {
		step /= 2
		d = fn(x+step) - fn(x-step)
	}
	return d / (2 * float64(step))
}

func (fn SFn) G() float64 {
	return fn.ReduceG(fn.ExpandG())
}
func (fn SFn) ExpandG() (g0, g1 float64) {
	step := 1.0
	g1 = 1.0
	f0 := fn(g0)
	for i := 0; i < 30; i++ {
		g1 = g0 + step
		f1 := fn(g1)
		if f1 > f0 {
			return
		}
		g0, f0 = g1, f1
		step *= 2
	}
	return
}

// g0 should be the best estimate of g
// I think I could find a better escape condition by looking at incremental
// improvements.
func (fn SFn) ReduceG(g0, g1 float64) float64 {
	f0, f1 := fn(g0), fn(g1)
	f := f0
	for i := 0; i < 30; i++ {
		if f1 > f0 {
			f0, f1, g0, g1 = f1, f0, g1, g0
		}
		if f/f0 > 10 {
			return g1
		}
		g0 = (g0 + g1) / 2
		f0 = fn(g0)
	}
	if g1 > 0.0 {
		return g1
	}
	return g0
}

func (fn SFn) Step(g float64) (float64, float64, bool) {
	const small cmpr.Tolerance = 1e-3
	f := fn(g)
	df_dg := fn.D(g)
	if small.Zero(df_dg) {
		return 0, f, true
	}
	step := f / df_dg
	return step, f, false
}
