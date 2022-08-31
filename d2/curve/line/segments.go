package line

import (
	"github.com/adamcolton/geom/d2"
)

// Segments links together a series of points and fulfils Curve.
type Segments []d2.Pt

// Pt1 returns a point on the line segments. All segments are weighted equally
// regardless of actual length.
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

// LineIntersections fulfills Intersections, returning the points that intersect
// the line.
func (ls Segments) LineIntersections(l Line, buf []float64) []float64 {
	if len(ls) < 2 {
		return buf[:0]
	}
	max := len(buf)
	buf = buf[:0]
	prev := ls[0]
	for _, pt := range ls[1:] {
		seg := New(prev, pt)
		prev = pt
		t0, t1, ok := seg.Intersection(l)
		if !ok || t0 < 0 || t0 >= 1 {
			continue
		}
		buf = append(buf, t1)
		if max > 0 && len(buf) == max {
			return buf
		}
	}
	return buf
}
