package descent

import (
	"fmt"
	"math"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
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
	for i := 0; i < 30; i++ {
		if f1 > f0 {
			f0, f1, g0, g1 = f1, f0, g1, g0
		}
		inc := f0 - f1
		if d := inc / prevInc; d < 1 && d > .90 {
			fmt.Println(d, i)
			return g1
		}
		prevInc = inc
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

func Distance2(a, b d2.Pt1) (Fn, []Fn) {
	var (
		fn Fn = func(t []float64) float64 {
			return a.Pt1(t[0]).Distance(b.Pt1(t[1]))
		}
		da      = d2.GetV1(a)
		dfn0 Fn = func(t []float64) float64 {
			// chain rule h(x) = f(g(x)) --> h'(x) = f'(g(x))*g'(x)

			// d = ( (bx-lx)^2 + (by-ly)^2 )^0.5
			// f(g) = g^0.5 :. f'(g) = ((g^-0.5)/2) * g'
			// g(h,i) = h^2 + i^2 :. g'(h,i) = 2h*h' + 2i*i'
			// h(bx) = bx-lx :. h' = bx'
			// i(by) = by-ly :. i' = by'
			// dbx/dt0 = db(t0).x
			// dby/dt0 = db(t0).y
			// d = f(g(h(bx(t0)),i(by(t0))))

			da_t0 := da.V1(t[0])
			hi := a.Pt1(t[0]).Subtract(b.Pt1(t[1]))
			dg_t0 := 2*hi.X*da_t0.X + 2*hi.Y*da_t0.Y
			g := hi.X*hi.X + hi.Y*hi.Y
			return (math.Pow(g, -0.5) / 2) * dg_t0
		}
		db      = d2.GetV1(b)
		dfn1 Fn = func(t []float64) float64 {

			db_t1 := db.V1(t[0])
			hi := a.Pt1(t[0]).Subtract(b.Pt1(t[1]))
			dg_t0 := -2*hi.X*db_t1.X - 2*hi.Y*db_t1.Y
			g := hi.X*hi.X + hi.Y*hi.Y
			return (math.Pow(g, -0.5) / 2) * dg_t0
		}
	)

	return fn, []Fn{dfn0, dfn1}
}
