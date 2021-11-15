package descent

import "github.com/adamcolton/geom/calc/cmpr"

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
			s.DFn[i] = CurryPartialDerivative(s.Fn, i)
		}
	}
}

func (s *Solver) Derivative(x []float64) []float64 {
	out := make([]float64, len(s.DFn))
	for i, df := range s.DFn {
		out[i] = df(x)
	}
	return out
}

func (s *Solver) Step(x []float64) {
	d := s.Derivative(x)
	for i, di := range d {
		x[i] -= di
	}
}

const small cmpr.Tolerance = 1e-3

func PartialDerivative(f Fn, x []float64, idx int) float64 {
	size := 1e-2
	d := 1.0
	save := x[idx]
	for !small.Zero(d) {
		size /= 10
		x[idx] += size
		d = f(x)
		x[idx] -= size
		d -= f(x)
		x[idx] = save
	}
	return d / (2 * size)
}

func CurryPartialDerivative(f Fn, idx int) Fn {
	return func(x []float64) float64 {
		return PartialDerivative(f, x, idx)
	}
}
