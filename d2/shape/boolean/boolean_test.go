package boolean_test

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/boolean"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestIntersection(t *testing.T) {
	s := boolean.Intersection{
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

	geomtest.Equal(t, []d2.Pt{{1.5, 1}, {1.75, 1.5}, {1.5, 2}, {0.5, 2}, {0.25, 1.5}, {0.5, 1}}, s.ConvexHull())

}

func TestUnion(t *testing.T) {
	s := boolean.Union{
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

	geomtest.Equal(t, []d2.Pt{{1, 0}, {2, 1}, {2, 2}, {1, 3}, {0, 2}, {0, 1}}, s.ConvexHull())
}

func TestSubtract(t *testing.T) {
	s := boolean.Subtract{
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

	geomtest.Equal(t, []d2.Pt{{0, 1}, {2, 1}, {1, 3}}, s.ConvexHull())
}
