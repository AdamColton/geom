package descent

import (
	"github.com/adamcolton/geom/calc/cmpr"
)

type Fn func([]float64) float64
type DFn func(x, buf []float64) []float64

type Solver struct {
	Ln  int
	Fn  Fn
	DFn DFn
}

func (s *Solver) SetDFn(partials []Fn) {
	if len(partials) < s.Ln {
		cp := make([]Fn, s.Ln)
		copy(cp, partials)
		partials = cp
	}
	for i, p := range partials {
		if p == nil {
			partials[i] = s.Fn.CurryPartialDerivative(i)
		}
	}
	s.DFn = func(x, buf []float64) []float64 {
		var out []float64
		if len(buf) < s.Ln {
			out = make([]float64, s.Ln)
		} else {
			out = buf[:s.Ln]
		}

		for i, p := range partials {
			out[i] = p(x)
		}

		return out
	}
}

func (s *Solver) Step(x, buf1, buf2 []float64) {
	d := s.DFn(x, buf1)
	for i := range d {
		d[i] *= -1
	}
	g := s.G(x, d, buf2)
	for i := range d {
		d[i] = d[i] * g
		x[i] += d[i]
	}
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

func (fn Fn) SecantStep(x0, x1 []float64) []float64 {
	dfx0 := fn.D(x0, x1)
	dfx1 := fn.D(x1, x0)
	out := make([]float64, len(x0))
	for i := range out {
		out[i] = (x0[i] - x1[i]) / (dfx1[i] - dfx0[i])
	}
	return out
}

func (fn Fn) D(x0, x1 []float64) []float64 {
	// f(x0+x1[i]) - f(x0)/d
	out := make([]float64, len(x0))
	fx0 := fn(x0)
	for i, x := range x0 {
		x0[i] = x1[i]
		out[i] = (fn(x0) - fx0) / (x1[i] - x)
		x0[i] = x
	}
	return out
}

func (fn Fn) SecantStep2(x0, x1 []float64) []float64 {
	// dx = x1-x0
	// df[i] --> x=x0; x[i]=x1[i] f(x)-f(x0) --- I think
	// dx[i]
	// x2[x] = x1

	x2 := make([]float64, len(x1))
	f0, f1 := fn(x0), fn(x1)
	df := f1 - f0
	for i := range x2 {
		m := df / (x1[i] - x0[i])
		b := f0 - m*x0[i]
		x2[i] = -b / m
	}
	return x2
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

// this may have an issue
// if we step over the min, but still to a lower point
// and then the next step is higher, we return a range that doesn't contain the min
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
	prevInc := 0.0
	belowThresh := false
	// so I think I need to track d and not break while it is decreasing
	for i := 0; i < 30; i++ {
		if f1 > f0 {
			f0, f1, g0, g1 = f1, f0, g1, g0
		}
		inc := f0 - f1
		d := inc / prevInc
		if belowThresh && d < 1 && d > 0.90 {
			return g1
		}
		belowThresh = belowThresh || d < 0.90
		prevInc = inc
		g0 = (g0 + g1) / 2
		f0 = fn(g0)
	}
	if g0 > 0.0 {
		return g0
	}
	return g1
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
