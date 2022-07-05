package bezier

import (
	"math"

	"github.com/adamcolton/geom/calc/comb"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/affine"
	"github.com/adamcolton/geom/d2/shape/box"
)

// Bezier curve defined by a slice of control points
type Bezier []d2.Pt

// Pt1 fulfills d2.Pt1 and returns the parametric curve point on the bezier
// curve.
func (b Bezier) Pt1(t float64) d2.Pt {
	l := len(b) - 1
	l64 := float64(l)
	if t == 1 {
		return b[l]
	}
	if t == 0 {
		return b[0]
	}

	// B(t) = ∑ binomialCo(l,i) * (1-t)^(l-i) * t^(i) * points[i]
	// let s = (1-t)^(l-i) * t^(i)
	// then s[i] = s[i-1] * t/(1-t)
	// and s[0] = (1-t) ^ l

	ti := 1 - t
	s := math.Pow(ti, l64)
	sd := t / ti
	w := &affine.Weighted{}
	for i, p := range b {
		b := float64(comb.Binomial(l, i))
		w.Weight(p, s*b)
		s *= sd
	}
	return w.Centroid()
}

// Tangent returns a curve that is the derivative of the Bezier curve
func (b Bezier) Tangent() Tangent {
	return Tangent{
		Bezier: Bezier(diffPoints(b...)),
	}
}

// V1c0 aliases Tangent and fulfills d2.V1c0. This is more efficient if many
// tangents points are needed from one curve.
func (b Bezier) V1c0() d2.V1 {
	return b.Tangent()
}

// V1 fulfills d2.V1 and returns the derivate at t.
func (b Bezier) V1(t float64) d2.V {
	return b.Tangent().V1(t)
}

// L fulfills d2.Limiter.
func (Bezier) L(t, c int) d2.Limit {
	if t == 1 && c == 1 {
		return d2.LimitUnbounded
	}
	return d2.LimitUndefined
}

// VL fulfills d2.VLimiter.
func (Bezier) VL(t, c int) d2.Limit {
	if t == 1 && (c == 1 || c == 0) {
		return d2.LimitUnbounded
	}
	return d2.LimitUndefined
}

func (b Bezier) BoundingBox() (min, max d2.Pt) {
	return box.New(b...).BoundingBox()

}

// Tangent holds a first derivative of a Bezier curve which mathematically is
// also a bezier curve, but the returned type is V instead of Pt.
type Tangent struct {
	Bezier Bezier
}

// V1 fulfills d2.V1
func (bt Tangent) V1(t float64) d2.V {
	return bt.Bezier.Pt1(t).V()
}

func diffPoints(points ...d2.Pt) []d2.Pt {
	l := len(points) - 1
	scale := float64(l)
	dps := make([]d2.Pt, l)
	prev := points[0]
	for i, p := range points[1:] {
		dps[i] = p.Subtract(prev).Multiply(scale).Pt()
		prev = p
	}
	return dps
}

// Cache pre-computes the Tangent and returns a BezierCache
func (b Bezier) Cache() BezierCache {
	return BezierCache{
		Bezier:  b,
		Tangent: b.Tangent(),
	}
}

// BezierCache holds both a Bezier curve and it's Tangent and is more efficient if
// both Pts and Vs will be needed from a curve.
type BezierCache struct {
	Bezier
	Tangent
}

// V1 fulfills d2.V1
func (bc BezierCache) V1(t float64) d2.V {
	return bc.Tangent.V1(t)
}

// V1c0 fulfills d2.V1c0
func (bc BezierCache) V1c0() d2.V1 {
	return bc
}
