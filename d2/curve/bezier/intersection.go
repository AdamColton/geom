package bezier

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/box"
)

// Blossom point for the control points of a bezier curve
func (b Bezier) Blossom(fs ...float64) d2.Pt {
	// https://en.wikipedia.org/wiki/Blossom_(functional)
	return b.BlossomBuf(make([]d2.Pt, len(b)), fs...)
}

// BlossomBuf computes the Blossom point for the control points of a bezier
// curve using the provided buffer. Reusing a buffer can increase performance.
func (b Bezier) BlossomBuf(buf []d2.Pt, fs ...float64) d2.Pt {
	ln := len(buf)
	copy(buf, b)
	for _, f := range fs {
		ln--
		for i, pt := range buf[:ln] {
			buf[i] = pt.Add(buf[i+1].Subtract(pt).Multiply(f))
		}
	}
	return buf[0]
}

// Segment returns a bezier curve that whose start and end is relative to the
// base curve. So calling b.(0.2, 0.7) will return a curve that exactly matches
// b from 0.2 to 0.7.
func (b Bezier) Segment(start, end float64) Bezier {
	return b.SegmentBuf(start, end, nil)
}

// SegmentBuf returns a bezier curve that whose start and end is relative to the
// base curve. Providing a buf reduces the overhead.
func (b Bezier) SegmentBuf(start, end float64, buf []d2.Pt) Bezier {
	fs := make([]float64, len(b)-1)
	for i := range fs {
		fs[i] = start
	}
	ln := len(b)
	out := make(Bezier, ln)
	if ln > len(buf) {
		buf = make([]d2.Pt, ln)
	} else if ln < len(buf) {
		buf = buf[:ln]
	}
	for i := range out {
		if i > 0 {
			fs[i-1] = end
		}
		out[i] = b.BlossomBuf(buf, fs...)
	}
	return out
}

// LineIntersections fulfills line.LineIntersector returning the intersection
// points relative to the line.
func (b Bezier) LineIntersections(l line.Line) []float64 {
	m, M := d2.MinMax(b...)
	t, ok := box.Box{m, M}.LineIntersection(l)
	if !ok {
		return nil
	}
	if m.Distance(M) < maxSize {
		return []float64{t}
	}
	buf := make([]d2.Pt, len(b))
	return append(b.SegmentBuf(0, 0.5, buf).LineIntersections(l), b.SegmentBuf(0.5, 1, buf).LineIntersections(l)...)
}

const maxSize = 1e-10

// BezierIntersections returns the intersection points relative to the Bezier
// curve.
func (b Bezier) BezierIntersections(l line.Line) []float64 {
	return b.bezierIntersections(l, 0, 1, make([]d2.Pt, len(b)))
}

func (b Bezier) bezierIntersections(l line.Line, t0, t1 float64, buf []d2.Pt) []float64 {
	m, M := d2.MinMax(b...)
	_, ok := box.Box{m, M}.LineIntersection(l)
	if !ok {
		return nil
	}
	tc := (t0 + t1) / 2.0
	if m.Distance(M) < maxSize {
		return []float64{tc}
	}
	return append(b.SegmentBuf(0, 0.5, buf).bezierIntersections(l, t0, tc, buf), b.SegmentBuf(0.5, 1, buf).bezierIntersections(l, tc, t1, buf)...)
}
