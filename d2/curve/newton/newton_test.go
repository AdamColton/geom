package newton

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"

	//"github.com/adamcolton/geom/d2/curve/poly"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/stretchr/testify/assert"
)

func TestPolyLineIntersection(t *testing.T) {
	l1 := line.New(d2.Pt{0, 0}, d2.Pt{1, 1})
	l2 := line.New(d2.Pt{1, 0}, d2.Pt{0, 1})

	p1 := poly.Poly{l1.T0.V(), l1.D}
	p2 := poly.Poly{l2.T0.V(), l2.D}

	t1, t2 := NewDistance(p1, p2).Min(0, 0)
	assert.InDelta(t, 0, p1.Pt1(t1).Distance(p2.Pt1(t2)), 1e-5)

	t1, t2 = NewDistance(p1, p2).Min(10, 10)
	assert.InDelta(t, 0, p1.Pt1(t1).Distance(p2.Pt1(t2)), 1e-5)

	t1, t2 = NewDistance(p1, p2).Min(0.5, 0.5)
	assert.InDelta(t, 0, p1.Pt1(t1).Distance(p2.Pt1(t2)), 1e-5)
}

func TestPolyCubicIntersection(t *testing.T) {
	p0 := poly.Bezier([]d2.Pt{d2.Pt{10, 10}, d2.Pt{50, 400}, d2.Pt{400, 300}, d2.Pt{200, 20}})
	p1 := poly.Bezier([]d2.Pt{d2.Pt{20, 300}, d2.Pt{10, 10}, d2.Pt{300, 30}, d2.Pt{400, 350}})

	dst := NewDistance(p0, p1)
	t0, t1 := dst.Min(0, 0)
	assert.InDelta(t, 0, p0.Pt1(t0).Distance(p1.Pt1(t1)), 0.002)

	t0, t1 = dst.Min(-5, 3)
	assert.InDelta(t, 0, p0.Pt1(t0).Distance(p1.Pt1(t1)), 0.002)
}

func TestPolyCubicNoIntersection(t *testing.T) {
	p0 := poly.Bezier([]d2.Pt{d2.Pt{10, 10}, d2.Pt{50, 200}, d2.Pt{400, 200}, d2.Pt{200, 20}})
	p1 := poly.Bezier([]d2.Pt{d2.Pt{20, 450}, d2.Pt{60, 300}, d2.Pt{450, 300}, d2.Pt{450, 470}})

	dst := NewDistance(p0, p1)
	t0a, t1a := dst.Min(0, 0)
	t0b, t1b := dst.Min(-5, 3)
	// Accuracy isn't as good when curves don't intersect.
	assert.InDelta(t, t0a, t0b, 0.05)
	assert.InDelta(t, t1a, t1b, 0.05)
}

func TestDistancePartialDerivatives(t *testing.T) {
	p0 := poly.Bezier([]d2.Pt{d2.Pt{10, 10}, d2.Pt{50, 400}, d2.Pt{400, 300}, d2.Pt{200, 20}})
	p1 := poly.Bezier([]d2.Pt{d2.Pt{20, 300}, d2.Pt{10, 10}, d2.Pt{300, 30}, d2.Pt{400, 350}})
	dst := NewDistance(p0, p1)

	small := 1e-6
	s := grid.Scale{0.1, 0.1, 0, 0}
	grid.Pt{11, 11}.Iter().Each(func(idx int, gpt grid.Pt) {
		t0, t1 := s.T(gpt)
		d, dt0, dt1 := dst.D(t0, t1)
		bd, bdt0, bdt1 := dst.BruteD(t0, t1, small)

		assert.Equal(t, dst.At(t0, t1), d)
		assert.InDelta(t, bd, d, 1e-4)
		assert.InDelta(t, bdt0, dt0, 1e-4)
		assert.InDelta(t, bdt1, dt1, 1e-4)
	})
}

func TestDG(t *testing.T) {
	p0 := poly.Bezier([]d2.Pt{d2.Pt{10, 10}, d2.Pt{50, 400}, d2.Pt{400, 300}, d2.Pt{200, 20}})
	p1 := poly.Bezier([]d2.Pt{d2.Pt{20, 300}, d2.Pt{10, 10}, d2.Pt{300, 30}, d2.Pt{400, 350}})
	dst := NewDistance(p0, p1)

	for t0 := -1.0; t0 < 2.0; t0 += 0.1 {
		for t1 := -1.0; t1 < 2.0; t1 += 0.1 {
			_, dt0, dt1 := dst.D(t0, t1)
			dgf := dst.dg(t0, t1, dt0, dt1)
			dgbf := dst.dgBrute(t0, t1, dt0, dt1, 1e-15)
			for g := -0.1; g < 0.1; g += 0.01 {
				d, dg := dgf(g)
				db, dgb := dgbf(g)
				pdg := math.Abs(dg-dgb) / math.Abs(dgb)
				pd := math.Abs(d-db) / math.Abs(db)
				assert.InDelta(t, 0, pdg, 0.022)
				assert.InDelta(t, 0, pd, 1e-2)
			}
		}
	}
}

func TestG(t *testing.T) {
	p0 := poly.Bezier([]d2.Pt{d2.Pt{10, 10}, d2.Pt{50, 400}, d2.Pt{400, 300}, d2.Pt{200, 20}})
	p1 := poly.Bezier([]d2.Pt{d2.Pt{20, 300}, d2.Pt{10, 10}, d2.Pt{300, 30}, d2.Pt{400, 350}})
	dst := NewDistance(p0, p1)
	// small needs to be at least one order larger than the step size we're
	// checking against in g.
	small := 1e-2

	for t0 := -1.0; t0 < 2.0; t0 += 0.1 {
		for t1 := -1.0; t1 < 2.0; t1 += 0.1 {
			before, dt0, dt1 := dst.D(t0, t1)
			g := dst.g(t0, t1, dt0, dt1)

			best, _, _ := dst.D(t0-g*dt0, t1-g*dt1)
			bestSub, _, _ := dst.D(t0-(g-small)*dt0, t1-(g-small)*dt1)
			bestPlus, _, _ := dst.D(t0-(g+small)*dt0, t1-(g+small)*dt1)

			assert.True(t, best <= before)
			assert.True(t, best <= bestSub)
			assert.True(t, best <= bestPlus)
		}
	}
}

func (dst Distance) dgBrute(t0, t1, dt0, dt1, small float64) dgfn {
	return func(g float64) (d float64, dd_dg float64) {
		g -= small
		t0g, t1g := t0-g*dt0, t1-g*dt1
		pt0, pt1 := dst.pt0(t0g), dst.pt1(t1g)
		d1 := pt0.Subtract(pt1).Mag2()
		g += 2 * small
		t0g, t1g = t0-g*dt0, t1-g*dt1
		pt0, pt1 = dst.pt0(t0g), dst.pt1(t1g)
		d2 := pt0.Subtract(pt1).Mag2()
		d = d1
		dd_dg = (d2 - d1) / (2 * small)
		return
	}
}

func (dst Distance) BruteD(t0, t1, small float64) (d float64, dt0 float64, dt1 float64) {
	pt0, pt1 := dst.pt0(t0-small), dst.pt1(t1)
	v := pt0.Subtract(pt1)
	da0 := v.Mag2()

	pt0, pt1 = dst.pt0(t0+small), dst.pt1(t1)
	v = pt0.Subtract(pt1)
	db0 := v.Mag2()

	pt0, pt1 = dst.pt0(t0), dst.pt1(t1-small)
	v = pt0.Subtract(pt1)
	da1 := v.Mag2()
	pt0, pt1 = dst.pt0(t0), dst.pt1(t1+small)
	v = pt0.Subtract(pt1)
	db1 := v.Mag2()

	d = (da0 + db0 + da1 + db1) / 4.0
	dt0 = (db0 - da0) / (2 * small)
	dt1 = (db1 - da1) / (2 * small)
	return d, dt0, dt1
}
