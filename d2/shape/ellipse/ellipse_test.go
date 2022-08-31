package ellipse

import (
	"fmt"
	"math"
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestEllipse(t *testing.T) {
	f1 := d2.Pt{0, 0}
	f2 := d2.Pt{2, 0}
	r := 1.0

	e := New(f1, f2, r)
	expected := d2.Pt{1, 0}
	actual := e.Pt2(0.5, 0.5)
	geomtest.Equal(t, expected, actual)
	actual = e.Pt2c1(0.5).Pt1(0.5)
	geomtest.Equal(t, expected, actual)
	assert.Equal(t, math.Sqrt2*math.Pi, e.Area())

	// Test Pt1 by definition of ellipse
	p0 := e.Pt1(0)
	d0 := f1.Distance(p0) + f2.Distance(p0)
	for i := 0.0; i <= 1.0; i += 0.2 {
		p := e.Pt1(i)
		d := f1.Distance(p) + f2.Distance(p)
		assert.InDelta(t, d0, d, 1e-10)
	}

	geomtest.Equal(t, d2.Pt{1, 0}, e.Centroid())

	expectedIntersections := []float64{0.58333, 0.25}
	gotIntersections := e.LineIntersections(line.New(d2.Pt{0, -2}, d2.Pt{4, 2}), nil)
	assert.InDeltaSlice(t, expectedIntersections, gotIntersections, 1e-5)

	assert.Equal(t, e.perimeter, e.Arc())

	e = New(f2, f1, -r)
	assert.Equal(t, math.Sqrt2*math.Pi, e.Area())

	e = New(f1, f1, r)
	assert.Equal(t, 2*math.Pi*r, e.Perimeter())
}

func TestCircle(t *testing.T) {
	r := 20.0
	c := NewCircle(d2.Pt{0, 0}, r)
	assert.Equal(t, r, c.Radius())
	assert.Equal(t, 2*r*math.Pi, c.Perimeter())
}

func TestCircumscribeCircle(t *testing.T) {
	tri := triangle.Triangle{
		{0, 0},
		{0, 1},
		{1, 0},
	}
	c := CircumscribeCircle(tri)
	r := c.Radius()
	center := c.Centroid()
	geomtest.Equal(t, d2.Pt{0.5, 0.5}, center)
	assert.InDelta(t, math.Sqrt(0.5), r, 1e-5)

	for _, pt := range tri {
		assert.InDelta(t, r, pt.Distance(center), 1e-5)
	}
}

func TestPerimeter(t *testing.T) {
	for r := 3.0; r < 7.0; r++ {
		for x := 1.0; x < 7.0; x++ {
			e := New(d2.Pt{0, 0}, d2.Pt{x, 0}, r)
			p := e.Perimeter()
			ps := e.PerimeterSeries(10)
			err := (ps - p) / ps
			assert.InDelta(t, 0, err, 1e-3, fmt.Sprint(x, r))
		}
	}
}

func TestContains(t *testing.T) {
	f1 := d2.Pt{2, 2}
	f2 := d2.Pt{4, 2}
	r := 1.0
	e := New(f1, f2, r)
	m, M := e.BoundingBox()

	pt0 := e.Pt1(0)
	d0 := f1.Distance(pt0) + f2.Distance(pt0)
	scale := grid.Scale{1.0 / 10.0, 1.0 / 10.0, 0, 0}
	grid.Pt{60, 40}.Iter().Each(func(idx int, gpt grid.Pt) {
		x, y := scale.T(gpt)
		p := d2.Pt{x, y}
		d := f1.Distance(p) + f2.Distance(p)
		c := e.Contains(p)
		assert.Equal(t, d-1e-5 < d0, c, fmt.Sprint(p, d, d0))
		if c {
			assert.True(t, p.X > m.X && p.Y > m.Y && p.X < M.X && p.Y < M.Y)
		}
	})
}

func TestConvexHullPoints(t *testing.T) {
	f1 := d2.Pt{2, 2}
	f2 := d2.Pt{4, 2}
	r := 1.0
	e := New(f1, f2, r)

	a := make(polygon.AssertConvexHuller, 10, 13)
	for i := range a {
		a[i] = e.Pt1(float64(i) / 10)
	}
	a = append(a, e.Centroid())
	a = append(a, f1, f2)

	geomtest.Equal(t, a, e)
}
