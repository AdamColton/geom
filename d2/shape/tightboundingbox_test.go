package shape_test

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/box"
	"github.com/adamcolton/geom/d2/shape/ellipse"
	"github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/stretchr/testify/assert"
)

func TestTightCircle(t *testing.T) {
	start := &box.Box{{-2, -2}, {2, 2}}
	s := shape.Rebound{
		Shape: ellipse.NewCircle(d2.Pt{0, 0}, 1.0),
		Box:   start,
	}

	m, M, b := shape.TightBoundingBox(s, 10)
	assert.True(t, b)

	// get the perfect bounding box from the underlying circle
	perfect := box.New(s.Shape.BoundingBox())

	// expect that solution has moved >99% of the way from the starting
	// solution to the perfect solution
	expected := box.Box{
		line.Line{
			T0: start[0],
			D:  start[0].Subtract(perfect[0]),
		}.Pt1(0.01),
		line.Line{
			T0: start[1],
			D:  start[1].Subtract(perfect[1]),
		}.Pt1(0.01),
	}

	assert.False(t, perfect.Contains(m))
	assert.False(t, perfect.Contains(M))
	assert.True(t, expected.Contains(m))
	assert.True(t, expected.Contains(M))
}

func TestNewBoundingBox(t *testing.T) {
	s := shape.Subtract{
		polygon.RectangleTwoPoints(d2.Pt{0, 0}, d2.Pt{1, 1}),
		ellipse.NewCircle(d2.Pt{0, 0.5}, 1),
	}
	r := shape.NewBoundingBox(s)
	m, M := r.BoundingBox()
	sbox := box.New(s.BoundingBox())

	assert.True(t, sbox.Contains(m))
	assert.True(t, sbox.Contains(M))
	assert.True(t, sbox[0].X < m.X)
}
