package bezier

import (
	"sort"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/curve/poly"
	"github.com/adamcolton/geom/d2/shape/box"
)

// Blossom point for the control points of a bezier curve
func (b Bezier) Blossom(fs ...float64) d2.Pt {
	// https://en.wikipedia.org/wiki/Blossom_(functional)
	return b.newBuf(nil, fs).blossom()
}

// BlossomBuf computes the Blossom point for the control points of a bezier
// curve using the provided buffer. Reusing a buffer can increase performance.
func (b Bezier) BlossomBuf(ptBuf []d2.Pt, fs ...float64) d2.Pt {
	return b.newBuf(ptBuf, fs).blossom()
}

// Segment returns a bezier curve that whose start and end is relative to the
// base curve. So calling b.(0.2, 0.7) will return a curve that exactly matches
// b from 0.2 to 0.7.
func (b Bezier) Segment(start, end float64) Bezier {
	return b.newBuf(nil, nil).segment(start, end).Bezier
}

// SegmentBuf returns a bezier curve that whose start and end is relative to the
// base curve. Providing a buf reduces the overhead.
func (b Bezier) SegmentBuf(start, end float64, ptBuf []d2.Pt, floatBuf []float64) Bezier {
	return b.newBuf(ptBuf, floatBuf).segment(start, end).Bezier
}

// LineIntersections fulfills line.LineIntersector returning the intersection
// points relative to the line.
func (b Bezier) LineIntersections(l line.Line, buf []float64) []float64 {
	return poly.NewBezier(b).LineIntersections(l, buf)
}

// BezierIntersections returns the intersection points relative to the Bezier
// curve.
func (b Bezier) BezierIntersections(l line.Line) []float64 {
	return poly.NewBezier(b).PolyLineIntersections(l, nil)
}

const maxSize = 1e-20

type buf struct {
	fs  []float64
	pts []d2.Pt
	box []float64
	Bezier
}

func (b Bezier) newBuf(pts []d2.Pt, fs []float64) buf {
	ln := len(b)
	if ptsLn := len(pts); ptsLn < ln {
		pts = make([]d2.Pt, ln)
	} else if ptsLn > ln {
		pts = pts[:ln]
	}
	ln--
	if fsLn := len(fs); fsLn > ln {
		fs = fs[:ln]
	} else if fsLn < ln {
		fs = make([]float64, ln)
	}
	return buf{
		pts:    pts,
		fs:     fs,
		Bezier: b,
	}
}

func (b buf) blossom() d2.Pt {
	ln := len(b.pts)
	copy(b.pts, b.Bezier)
	for _, f := range b.fs {
		ln--
		for i, pt := range b.pts[:ln] {
			b.pts[i] = pt.Add(b.pts[i+1].Subtract(pt).Multiply(f))
		}
	}
	return b.pts[0]
}

func (b buf) segment(start, end float64) buf {
	ln := len(b.Bezier)
	out := make(Bezier, ln)

	for j := range b.fs {
		b.fs[j] = start
	}
	out[0] = b.blossom()
	for i := range b.fs {
		b.fs[i] = end
		out[i+1] = b.blossom()
	}
	return buf{
		fs:     b.fs,
		pts:    b.pts,
		Bezier: out,
	}
}

func (b buf) bezier(l line.Line, t0, t1 float64) []float64 {
	bx := box.New(b.Bezier...)
	b.box = bx.LineIntersections(l, b.box[:0])
	if len(b.box) == 0 {
		return nil
	}
	tc := (t0 + t1) / 2.0
	if bx.V().Mag2() < maxSize {
		return []float64{tc}
	}
	return append(b.segment(0, 0.5).bezier(l, t0, tc), b.segment(0.5, 1).bezier(l, tc, t1)...)
}

func (b buf) line(l line.Line, max int, tBuf []float64) []float64 {
	bx := box.New(b.Bezier...)
	b.box = bx.LineIntersections(l, b.box[:0])
	if len(b.box) == 0 {
		return tBuf
	}
	if bx.V().Mag2() < maxSize {
		if len(b.box) > 1 {
			tBuf = append(tBuf, (b.box[0]+b.box[1])/2)
		} else {
			tBuf = append(tBuf, b.box[0])
		}
		return tBuf
	}

	tBuf = b.segment(0, 0.5).line(l, max, tBuf)
	if max == 0 || len(tBuf) < max {
		tBuf = b.segment(0.5, 1).line(l, max, tBuf)
	}

	return tBuf
}

const small cmpr.Tolerance = 1e-10
const smallish cmpr.Tolerance = 1e-5

func removeDups(s []float64) []float64 {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })

	prev := 0
	for i := 1; i < len(s); i++ {
		if !small.Equal(s[prev], s[i]) {
			prev++
			s[prev] = s[i]
		}
	}
	return s[:prev+1]
}

// BezierIntersections returns the intersection points relative to the Bezier
// curve.
func (b Bezier) Intersections(b2 Bezier) [][2]float64 {
	out := bez(b.newBuf(nil, nil), b2.newBuf(nil, nil), 0, 1, 0, 1)
	if len(out) == 0 {
		return out
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i][0] < out[j][0] ||
			(out[i][0] == out[j][0] && out[i][1] < out[j][1])
	})
	prev := 0
	for i := 1; i < len(out); i++ {
		if !smallish.Equal(out[prev][0], out[i][0]) || !smallish.Equal(out[prev][1], out[i][1]) {
			prev++
			out[prev] = out[i]
		}
	}

	return out[:prev+1]
}

func bez(b0, b1 buf, t0_0, t0_1, t1_0, t1_1 float64) [][2]float64 {
	bx0 := box.New(b0.Bezier...)
	bx1 := box.New(b1.Bezier...)
	if !bx0.Overlaps(bx1) {
		return nil
	}

	t0_c := (t0_0 + t0_1) / 2
	t1_c := (t1_0 + t1_1) / 2

	if bx0.V().Mag2() < 1e-10 || bx1.V().Mag2() < 1e-10 {
		return [][2]float64{{t0_c, t1_c}}
	}

	b0_0 := b0.segment(0, 0.5)
	b0_1 := b0.segment(0.5, 1)
	b1_0 := b1.segment(0, 0.5)
	b1_1 := b1.segment(0.5, 1)

	out := bez(b0_0, b1_0, t0_0, t0_c, t1_0, t1_c)
	out = append(out, bez(b0_0, b1_1, t0_0, t0_c, t1_c, t1_1)...)
	out = append(out, bez(b0_1, b1_0, t0_c, t0_1, t1_0, t1_c)...)
	out = append(out, bez(b0_1, b1_1, t0_c, t0_1, t1_c, t1_1)...)

	return out
}
