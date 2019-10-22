package line

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestLineSegments(t *testing.T) {
	ls := Segments{
		d2.Pt{0, 0},
		d2.Pt{1, 1},
		d2.Pt{2, 0},
		d2.Pt{3, 1},
		d2.Pt{4, 0},
	}
	geomtest.Equal(t, d2.Pt{-8, -8}, ls.Pt1(-2))
	geomtest.Equal(t, d2.Pt{-0.5, -0.5}, ls.Pt1(-0.125))
	geomtest.Equal(t, ls[0], ls.Pt1(0))
	geomtest.Equal(t, d2.Pt{0.5, 0.5}, ls.Pt1(0.125))
	geomtest.Equal(t, ls[1], ls.Pt1(0.25))
	geomtest.Equal(t, d2.Pt{1.5, 0.5}, ls.Pt1(0.375))
	geomtest.Equal(t, ls[2], ls.Pt1(0.5))
	geomtest.Equal(t, d2.Pt{2.5, 0.5}, ls.Pt1(0.625))
	geomtest.Equal(t, ls[3], ls.Pt1(0.75))
	geomtest.Equal(t, d2.Pt{3.5, 0.5}, ls.Pt1(0.875))
	geomtest.Equal(t, ls[4], ls.Pt1(1))
	geomtest.Equal(t, d2.Pt{4.5, -0.5}, ls.Pt1(1.125))

	is := ls.Intersections(New(d2.Pt{-1, 0.5}, d2.Pt{5, 0.5}))
	if assert.Len(t, is, 4) {
		assert.InDelta(t, 0.25, is[0], 1e-4)
		assert.InDelta(t, 0.4166, is[1], 1e-4)
		assert.InDelta(t, 0.5833, is[2], 1e-4)
		assert.InDelta(t, 0.75, is[3], 1e-4)
	}

	ls = Segments{
		d2.Pt{0, 0},
		d2.Pt{1, 1},
		d2.Pt{2, 0},
	}

	geomtest.Equal(t, ls[0], ls.Pt1(0))
	geomtest.Equal(t, d2.Pt{0.5, 0.5}, ls.Pt1(0.25))
	geomtest.Equal(t, ls[1], ls.Pt1(0.5))
	geomtest.Equal(t, d2.Pt{1.5, 0.5}, ls.Pt1(0.75))
	geomtest.Equal(t, ls[2], ls.Pt1(1))
}

func TestSegmentsEdgeCase(t *testing.T) {
	var s Segments
	geomtest.Equal(t, d2.Pt{}, s.Pt1(0))
	assert.Nil(t, s.Intersections(New(d2.Pt{1, 2}, d2.Pt{3, 5})))
	s = Segments{d2.Pt{6, 7}}
	geomtest.Equal(t, d2.Pt{6, 7}, s.Pt1(0))

	s = Segments{d2.Pt{0, 0}, d2.Pt{1, 1}}
	i := s.Intersections(New(d2.Pt{1, 1.5}, d2.Pt{2, 1.5}))
	assert.Nil(t, i)
}
