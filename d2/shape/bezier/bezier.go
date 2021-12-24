package bezier

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/d2/curve/line"
)

// BezShape is a shape whose sides are Bezier curves.
type BezShape []bezier.Bezier

func New(pts ...[]d2.Pt) BezShape {
	ln := len(pts)
	out := make([]bezier.Bezier, ln)
	for i, side := range pts {
		sln := len(side)
		out[i] = make(bezier.Bezier, sln+1)
		copy(out[i], side)
		out[i][sln] = pts[(i+1)%ln][0]
	}
	return out
}

// Validate checks that the end of each curve is equal to the start of the
// next and that the last point is equal to the first point.
func (b BezShape) Validate() bool {
	for i, c := range b[1:] {
		if b[i][len(b[i])-1].Subtract(c[0]).Mag2() > 1e-5 {
			return false
		}
	}
	ln := len(b)
	last := b[ln-1]
	return b[0][0].Subtract(last[len(last)-1]).Mag2() < 1e-5
}

// BoundingBox fulfills shape.BoundingBoxer.
func (b BezShape) BoundingBox() (min d2.Pt, max d2.Pt) {
	m, M := d2.MinMax(b[0]...)
	for _, pts := range b[1:] {
		mi, Mi := d2.MinMax(pts...)
		m = d2.Min(m, mi)
		M = d2.Min(M, Mi)
	}
	return m, M
}

// Contains fulfills shape.Container, returns true if the point is inside the
// shape.
func (b BezShape) Contains(pt d2.Pt) bool {
	l := line.Line{
		T0: pt,
		D:  d2.V{1, 2},
	}
	var i int
	var buf []float64
	for _, c := range b {
		buf = c.LineIntersections(l, buf[:0])
		for _, p := range buf {
			if p > 0 {
				i++
			}
		}
	}
	return i%2 == 1
}

// LineIntersections fulfills line.Intersector.
func (b BezShape) LineIntersections(l line.Line, buf []float64) []float64 {
	max := len(buf)
	buf = buf[:0]
	for _, c := range b {
		buf = append(buf, c.LineIntersections(l, buf[len(buf):])...)
		if max > 0 && len(buf) > max {
			return buf[:max]
		}
	}
	return buf
}
