package shape_test

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestIntersection(t *testing.T) {
	s := shape.Intersection{
		&triangle.Triangle{{0, 1}, {2, 1}, {1, 3}},
		&triangle.Triangle{{0, 2}, {2, 2}, {1, 0}},
	}

	assert.True(t, s.Contains(d2.Pt{1, 1}))
	assert.False(t, s.Contains(d2.Pt{1, 0.1}))

	expected := []float64{1.0 / 3.0, 2.0 / 3.0}
	l := line.New(d2.Pt{0, 0}, d2.Pt{2, 3})
	assert.Equal(t, expected, s.LineIntersections(l, nil))
	assert.Equal(t, expected[:1], s.LineIntersections(l, []float64{0}))
	assert.Equal(t, expected[:2], s.LineIntersections(l, []float64{0, 0}))

	m, M := s.BoundingBox()
	geomtest.Equal(t, d2.Pt{0, 1}, m)
	geomtest.Equal(t, d2.Pt{2, 2}, M)

	geomtest.Equal(t, []d2.Pt{{0.5, 1.0}, {1.5, 1}, {0.5, 2}, {0.25, 1.5}}, s.ConvexHull())
}

func TestUnion(t *testing.T) {
	s := shape.Union{
		&triangle.Triangle{{0, 1}, {2, 1}, {1, 3}},
		&triangle.Triangle{{0, 2}, {2, 2}, {1, 0}},
	}

	assert.True(t, s.Contains(d2.Pt{1, 1}))
	assert.True(t, s.Contains(d2.Pt{1, 0.1}))

	expected := []float64{5.0 / 7.0, 2.0 / 7.0}
	l := line.New(d2.Pt{0, 0}, d2.Pt{2, 3})
	assert.Equal(t, expected, s.LineIntersections(l, nil))
	assert.Equal(t, expected[:1], s.LineIntersections(l, []float64{0}))
	assert.Equal(t, expected[:2], s.LineIntersections(l, []float64{0, 0}))

	m, M := s.BoundingBox()
	geomtest.Equal(t, d2.Pt{0, 0}, m)
	geomtest.Equal(t, d2.Pt{2, 3}, M)
}

func TestSubtract(t *testing.T) {
	s := shape.Subtract{
		&triangle.Triangle{{0, 1}, {2, 1}, {1, 3}},
		&triangle.Triangle{{0, 2}, {2, 2}, {1, 0}},
	}

	assert.False(t, s.Contains(d2.Pt{1, 1}))
	assert.False(t, s.Contains(d2.Pt{1, 0.1}))
	assert.True(t, s.Contains(d2.Pt{1, 2.9}))

	expected := []float64{5.0 / 7.0, 2.0 / 3.0}
	l := line.New(d2.Pt{0, 0}, d2.Pt{2, 3})
	assert.Equal(t, expected, s.LineIntersections(l, nil))
	assert.Equal(t, expected[:1], s.LineIntersections(l, []float64{0}))
	assert.Equal(t, expected[:2], s.LineIntersections(l, []float64{0, 0}))

	m, M := s.BoundingBox()
	geomtest.Equal(t, d2.Pt{0, 1}, m)
	geomtest.Equal(t, d2.Pt{2, 3}, M)

	geomtest.Equal(t, []d2.Pt{{0, 1}, {2, 1}, {1, 3}}, s.ConvexHull())
}
