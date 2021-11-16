package shape

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

type Union [2]Shape

func (a Union) Contains(pt d2.Pt) bool {
	return a[0].Contains(pt) || a[1].Contains(pt)
}

func (a Union) LineIntersections(l line.Line, buf []float64) []float64 {
	max := len(buf)
	buf = a[0].LineIntersections(l, buf[:0])
	ln := len(buf)
	for i := 0; i < ln; i++ {
		t := buf[i]
		pt := l.Pt1(t)
		if a[1].Contains(pt) {
			ln--
			buf[i] = buf[ln]
			i--
		}
	}
	buf = buf[:ln]
	if max > 0 && len(buf) >= max {
		return buf[:max]
	}

	i1 := a[1].LineIntersections(l, buf[ln:])
	for _, t := range i1 {
		pt := l.Pt1(t)
		if !a[0].Contains(pt) {
			buf = append(buf, t)
		}
	}
	if max > 0 && len(buf) >= max {
		return buf[:max]
	}
	return buf
}

// BoundingBox fulfills shape.Shape, it returns a box that contains the shape.
func (a Union) BoundingBox() (d2.Pt, d2.Pt) {
	m, M := a[0].BoundingBox()
	m1, M1 := a[1].BoundingBox()
	return d2.Min(m, m1), d2.Max(M, M1)
}
