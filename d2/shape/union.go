package shape

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

type Union [2]Shape

func (a Union) Contains(pt d2.Pt) bool {
	return a[0].Contains(pt) || a[1].Contains(pt)
}

func (a Union) LineIntersections(l line.Line) []float64 {
	i0, i1 := a[0].LineIntersections(l), a[1].LineIntersections(l)
	out := make([]float64, 0, len(i0)+len(i1))
	for _, t := range i0 {
		pt := l.Pt1(t)
		if !a[1].Contains(pt) {
			out = append(out, t)
		}
	}
	for _, t := range i1 {
		pt := l.Pt1(t)
		if !a[0].Contains(pt) {
			out = append(out, t)
		}
	}
	return out
}

func (a Union) BoundingBox() (d2.Pt, d2.Pt) {
	m, M := a[0].BoundingBox()
	m1, M1 := a[1].BoundingBox()
	return d2.Min(m, m1), d2.Max(M, M1)
}
