// Package box provids a 3D bounding box.
package box

import (
	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/adamcolton/geom/geomerr"
)

// Box is defined by 2 points, a min and a max.
type Box [2]d3.Pt

// New creates a Box by finding the relative Min and Max across all the
// coordinates given.
func New(pts ...d3.Pt) *Box {
	b := &Box{}
	b[0], b[1] = d3.MinMax(pts...)
	return b
}

// Add points to expand the bounds if necessary.
func (b *Box) Add(pts ...d3.Pt) {
	m, M := d3.MinMax(pts...)
	b[0] = d3.Min(b[0], m)
	b[1] = d3.Max(b[1], M)
}

// LineIntersection returns the lowest point on the line greater than 0 that
// intersects the cube. If there is no intersection it will return -1 and false.
func (b *Box) LineIntersection(l line.Line) (float64, bool) {
	best := -1.0
	if l.D.X != 0 {
		t := l.AtX(b[0].X)
		x0 := l.Pt1(t)
		if x0.Y >= b[0].Y && x0.Y <= b[1].Y &&
			x0.Z >= b[0].Z && x0.Z <= b[1].Z {
			best = t
		}

		t = l.AtX(b[1].X)
		x1 := l.Pt1(t)
		if x1.Y >= b[0].Y && x1.Y <= b[1].Y &&
			x1.Z >= b[0].Z && x1.Z <= b[1].Z &&
			t > 0 && (t < best || best < 0) {
			best = t
		}
	}

	if l.D.Y != 0 {
		t := l.AtY(b[0].Y)
		y0 := l.Pt1(t)
		if y0.X >= b[0].X && y0.X <= b[1].X &&
			y0.Z >= b[0].Z && y0.Z <= b[1].Z &&
			t > 0 && (t < best || best < 0) {
			best = t
		}

		t = l.AtY(b[1].Y)
		y1 := l.Pt1(t)
		if y1.X >= b[0].X && y1.X <= b[1].X &&
			y1.Z >= b[0].Z && y1.Z <= b[1].Z &&
			t > 0 && (t < best || best < 0) {
			best = t
		}
	}

	if l.D.Z != 0 {
		t := l.AtZ(b[0].Z)
		z0 := l.Pt1(t)
		if z0.X >= b[0].X && z0.X <= b[1].X &&
			z0.Y >= b[0].Y && z0.Y <= b[1].Y &&
			t > 0 && (t < best || best < 0) {
			best = t
		}

		t = l.AtZ(b[1].Z)
		z1 := l.Pt1(t)
		if z1.X >= b[0].X && z1.X <= b[1].X &&
			z1.Y >= b[0].Y && z1.Y <= b[1].Y &&
			t > 0 && (t < best || best < 0) {
			best = t
		}
	}

	return best, best >= 0
}

// AssertEqual fulfills geomtest.AssertEqualizer. It will check if two boxes
// are equal.
func (b *Box) AssertEqual(to interface{}, t cmpr.Tolerance) error {
	if err := geomerr.NewTypeMismatch(b, to); err != nil {
		return err
	}
	b2 := to.(*Box)

	return geomerr.NewSliceErrs(2, -1, func(i int) error {
		return b[i].AssertEqual(b2[i], t)
	})
}
