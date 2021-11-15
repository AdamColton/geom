package bezier

import (
	"github.com/adamcolton/geom/d2"
)

// CompositeCurve stitches multiple curves together into a single parametric
// curve.
type CompositeBezier []Bezier

// Segment is a segment of a composite bezier curve where each segment is
// defined by 4 points. The individual segments are 3 points because the 4th
// point is taken from the previous segment. The first value is the handle to
// the previous point, the second value is handle on the end of the segment and
// the third point the end of the segment.
type Segment [3]d2.V

// NewRelativeCompositeBezier defines a curve using relative values. A segment
// is composed of 3 values, the first 2 are handles and defined relative to the
// point they are connected to. The 3rd point is the next control point defined
// relative to the previous control point.
func NewRelativeCompositeBezier(segments []Segment, transformation *d2.T) CompositeBezier {
	cb := make(CompositeBezier, len(segments))
	var prev d2.Pt
	for i, seg := range segments {
		pts := make([]d2.Pt, 4)
		pts[0] = prev
		pts[1] = prev.Add(seg[0])
		pts[3] = prev.Add(seg[2])
		pts[2] = pts[3].Add(seg[1])
		prev = pts[3]
		cb[i] = Bezier(transformation.Slice(pts))
	}
	return cb
}

func (cb CompositeBezier) Pt1(t0 float64) d2.Pt {
	ln := len(cb)
	if ln == 0 {
		return d2.Pt{}
	}
	if ln == 1 {
		return cb[0].Pt1(t0)
	}

	// 4 points = 3 segments 0:2
	ts := t0 * float64(ln)
	ti := int(ts)
	if ti >= ln {
		ti = ln - 1
	} else if ti < 0 {
		ti = 0
	}
	ts -= float64(ti)

	return cb[ti].Pt1(ts)
}

// // F returns a parametric point along the curve.
// func (cc CompositeCurve) F(t float64) F {
// 	ln := len(cc)
// 	if ln == 0 {
// 		return F{}
// 	}
// 	if ln == 1 {
// 		return cc[0](t)
// 	}

// 	// 4 points = 3 segments 0:2
// 	ts := t * float64(ln)
// 	ti := int(ts)
// 	if ti >= ln {
// 		ti = ln - 1
// 	} else if ti < 0 {
// 		ti = 0
// 	}
// 	ts -= float64(ti)

// 	return cc[ti](ts)
// }
