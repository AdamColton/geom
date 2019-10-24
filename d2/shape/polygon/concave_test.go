package polygon

import (
	"testing"

	"github.com/adamcolton/geom/geomtest"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/stretchr/testify/assert"
)

func TestConcave(t *testing.T) {
	pp := PolarPolygon{
		{3, angle.Rot(0.0 / 8.0)},
		{3, angle.Rot(2.0 / 8.0)},
		{3, angle.Rot(4.0 / 8.0)},
		{3, angle.Rot(6.0 / 8.0)},
		{1, angle.Rot(1.0 / 8.0)},
		{1, angle.Rot(3.0 / 8.0)},
		{1, angle.Rot(5.0 / 8.0)},
		{1, angle.Rot(7.0 / 8.0)},
	}
	pp.Sort()
	p := pp.Polygon(d2.Pt{})
	assert.False(t, p.Convex())

	ccp := NewConcavePolygon(p)
	geomtest.Equal(t, d2.Pt{3, 0}, ccp.Pt2(0, 0))
	geomtest.Equal(t, d2.Pt{3, 0}, ccp.Pt1(0))
	geomtest.Equal(t, d2.Pt{0, 0}, ccp.Centroid())
	assert.InDelta(t, 8.48528, ccp.Area(), 1e-4)
	assert.InDelta(t, 8.48528, ccp.SignedArea(), 1e-4)
	assert.InDelta(t, 19.195598, ccp.Perimeter(), 1e-4)
	assert.True(t, ccp.Contains(d2.Pt{0, 0}))
	assert.False(t, ccp.Contains(d2.Pt{1, 1}))

}
