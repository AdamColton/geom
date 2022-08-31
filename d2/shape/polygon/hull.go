package polygon

import (
	"fmt"
	"math"
	"sort"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomerr"
)

// GrahamTolerance is the deviation from 0 of a cross product that will still
// be treated as 0.
var GrahamTolerance = cmpr.Tolerance(1e-7)

// ConvexHull is currently a wrapper around GrahamScan. A more efficient
// algorithm may be added later.
func ConvexHull(pts ...d2.Pt) []d2.Pt {
	return GrahamScan(pts...)
}

// GrahamScan: https://en.wikipedia.org/wiki/Graham_scan
// Finds the convex hull of a set of points in O(n log n).
func GrahamScan(pts ...d2.Pt) []d2.Pt {
	ln := len(pts)
	if ln < 4 {
		return pts
	}
	least := 0
	for i, pt := range pts[1:] {
		if pt.Y < pts[least].Y {
			least = i + 1
		}
	}

	type ptRef struct {
		idx int
		ang angle.Rad
	}
	order := make([]ptRef, ln-1)
	for i := 0; i < ln-1; i++ {
		idx := (least + i + 1) % ln
		order[i].idx = idx
		order[i].ang = pts[idx].Subtract(pts[least]).Angle()
	}
	sort.Slice(order, func(i, j int) bool {
		return order[i].ang < order[j].ang
	})

	stack := make([]d2.Pt, 0, int(math.Log(float64(ln)))+1)
	stack = append(stack, pts[least], pts[order[0].idx])
	sln := 2

	for _, o := range order[1:] {
		pt := pts[o.idx]
		for {
			d1, d2 := pt.Subtract(stack[sln-1]), pt.Subtract(stack[sln-2])
			if (GrahamTolerance.Zero(d1.X) && GrahamTolerance.Zero(d1.Y)) ||
				(GrahamTolerance.Zero(d2.X) && GrahamTolerance.Zero(d2.Y)) {
				break
			}
			c := d1.Cross(d2)
			if c < -float64(GrahamTolerance) {
				stack = append(stack, pt)
				sln++
				break
			} else if c < float64(GrahamTolerance) {
				break
			}
			sln--
			stack = stack[:sln]
		}
	}
	return stack
}

// ErrNotConvex is return if a convex hull is not actually convex.
type ErrNotConvex struct{}

// Error fulfills error and returns the string "hull is not convex"
func (ErrNotConvex) Error() string {
	return "hull is not convex"
}

// AssertConvexHull fulfills geomtest.AssertEqualizer. It is a list of points.
// Given a convex hull, it checks that the hull is convex and contains all the
// points in it's list.
type AssertConvexHull []d2.Pt

// AssertEqual fulfills geomtest.AssertEqualizer. It checks that actual is a
// slice of points that form a convex hull containing all the points in
// AssertConvexHull.
func (a AssertConvexHull) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	pts, ok := actual.([]d2.Pt)
	if !ok {
		return geomerr.TypeMismatch([]d2.Pt(nil), actual)
	}
	p := Polygon(pts)
	if !p.Convex() {
		return ErrNotConvex{}
	}
	return geomerr.NewSliceErrs(len(a), -1, func(i int) error {
		if !p.Contains(a[i]) {
			return fmt.Errorf("%s", a[i])
		}
		return nil
	})
}

// AssertConvexHuller fulfills geomtest.AssertEqualizer. It is a list of points.
// Given an interface that fulfills ConvexHuller, it checks that the hull is
// convex and contains all the points in it's list.
type AssertConvexHuller []d2.Pt

// ConvexHuller is taken from shape.ConvexHuller to avoid cyclic imports.
type ConvexHuller interface {
	ConvexHull() []d2.Pt
}

// AssertEqual fulfills geomtest.AssertEqualizer. It check that actual is an an
// instance of ConvexHuller and the ConvexHull returned is convex and contains
// all the points in AssertConvexHuller.
func (a AssertConvexHuller) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	ch, ok := actual.(ConvexHuller)
	if !ok {
		return geomerr.TypeMismatch([]d2.Pt(nil), actual)
	}
	return AssertConvexHull(a).AssertEqual(ch.ConvexHull(), t)
}
