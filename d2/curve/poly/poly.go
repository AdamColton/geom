package poly

import (
	"github.com/adamcolton/geom/calc/comb"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

type Poly []d2.V

func (p Poly) Pt1(t0 float64) d2.Pt {
	tn := 1.0
	var sum d2.Pt
	for _, v := range p {
		sum = sum.Add(v.Multiply(tn))
		tn *= t0
	}
	return sum
}

func (p Poly) V1(t0 float64) d2.V {
	var sum d2.V
	tn := 1.0
	for i, c := range p[1:] {
		sum = sum.Add(c.Multiply(float64(i+1) * tn))
		tn *= t0
	}
	return sum
}

type FirstOrder struct {
	P Poly
}

func (p Poly) V1c0() d2.V1 {
	out := make(Poly, len(p)-1)
	for i, c := range p[1:] {
		out[i] = c.Multiply(float64(i + 1))
	}
	return FirstOrder{out}
}

func (fo FirstOrder) V1(t0 float64) d2.V {
	return fo.P.Pt1(t0).V()
}

func (p Poly) Multiply(p2 Poly) Poly {
	out := make(Poly, len(p)+len(p2)-1)

	for i, v := range p {
		for j, v2 := range p2 {
			t := i + j
			out[t] = out[t].Add(v.Product(v2))
		}
	}
	return out
}

func (p Poly) Add(p2 Poly) Poly {
	ln := len(p)
	if ln2 := len(p2); ln2 > ln {
		ln = ln2
	}
	out := make(Poly, ln)
	for i := range out {
		if i < len(p) {
			out[i] = p[i]
		}
		if i < len(p2) {
			out[i] = out[i].Add(p[i])
		}
	}
	return out
}

func Bezier(pts []d2.Pt) Poly {
	l := len(pts) - 1

	// B(t) = ∑ binomialCo(l,i) * (1-t)^(l-i) * t^(i) * points[i]

	sum := make(Poly, l+1)
	for i, pt := range pts {
		v := pt.V().Multiply(float64(comb.Binomial(l, i)))
		sign := 1.0
		for term := 0; term < l-i+1; term++ {
			s := sign * float64(comb.Binomial(l-i, term))
			sign = -sign
			sum[term+i] = sum[term+i].Add(v.Multiply(s))
		}
	}
	return sum
}

func (p Poly) LineIntersections(l line.Line, buf []float64) []float64 {
	return nil
}
