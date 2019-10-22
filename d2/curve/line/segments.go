package line

import (
	"github.com/adamcolton/geom/d2"
)

// LineSegments links together a series of points and fulfils Curve.
type Segments []d2.Pt

func (ls Segments) Pt1(t float64) d2.Pt {
	ln := len(ls)
	if ln == 0 {
		return d2.Pt{}
	}
	if ln == 1 {
		return ls[0]
	}

	// 4 points = 3 segments 0:2
	ts := t * float64(ln-1)
	ti := int(ts)
	if ti > ln-2 {
		ti = ln - 2
	} else if ti < 0 {
		ti = 0
	}
	return New(ls[ti], ls[ti+1]).Pt1(ts - float64(ti))
}

func (ls Segments) Intersections(l2 Line) []float64 {
	if len(ls) < 2 {
		return nil
	}
	var out []float64
	prev := ls[0]
	for _, pt := range ls[1:] {
		l := New(prev, pt)
		prev = pt
		i, ok := l2.LineIntersection(l)
		if !ok || i < 0 || i >= 1 {
			continue
		}
		i, _ = l.LineIntersection(l2)
		out = append(out, i)
	}
	return out
}
