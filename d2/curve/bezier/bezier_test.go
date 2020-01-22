package bezier

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestBezier(t *testing.T) {
	b := Bezier{{0, 0}, {0.5, 1}, {1, 0}}
	assert.Equal(t, d2.Pt{0, 0}, b.Pt1(0))
	assert.Equal(t, d2.Pt{1, 0}, b.Pt1(1))
	assert.Equal(t, d2.Pt{0.5, 0.5}, b.Pt1(0.5))
	assert.Equal(t, d2.Pt{-1, -4}, b.Pt1(-1))

	b = Bezier{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	assert.Equal(t, d2.Pt{0, 0}, b.Pt1(0))
	assert.Equal(t, d2.Pt{1, 0}, b.Pt1(1))
	assert.Equal(t, d2.Pt{0.5, 0.75}, b.Pt1(0.5))

	b = Bezier{{0, 0}, {0, 1}, {1, -1}, {1, 0}}
	assert.Equal(t, d2.Pt{0, 0}, b.Pt1(0))
	assert.Equal(t, d2.Pt{1, 0}, b.Pt1(1))
	assert.Equal(t, d2.Pt{0.5, 0}, b.Pt1(0.5))
	assert.Equal(t, b.Pt1(0.75).X, 1-b.Pt1(0.25).X)
	assert.Equal(t, b.Pt1(0.75).Y, -b.Pt1(0.25).Y)
}

func TestLimits(t *testing.T) {
	assert.Equal(t, d2.LimitUnbounded, Bezier{}.L(1, 1))
	assert.Equal(t, d2.LimitUndefined, Bezier{}.L(2, 1))
	assert.Equal(t, d2.LimitUnbounded, Bezier{}.VL(1, 0))
	assert.Equal(t, d2.LimitUnbounded, Bezier{}.VL(1, 1))
	assert.Equal(t, d2.LimitUndefined, Bezier{}.VL(2, 1))
}

func TestV1(t *testing.T) {
	geomtest.V1(t, Bezier{{0, 0}, {0.5, 1}, {1, 0}})
	geomtest.V1(t, Bezier{{0, 0}, {0, 1}, {1, -1}, {1, 0}})
}

func TestIntersection(t *testing.T) {
	b := Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	l := line.New(d2.Pt{0, 200}, d2.Pt{500, 300})

	li := b.LineIntersections(l)
	bi := b.BezierIntersections(l)
	for i, lt := range li {
		ptl := l.Pt1(lt)
		ptb := b.Pt1(bi[i])
		geomtest.Equal(t, ptl, ptb)
	}
}

func TestCache(t *testing.T) {
	b := Bezier{{0, 0}, {0.5, 1}, {1, 0}}.Cache()
	geomtest.Equal(t, d2.V{1, 1}, b.V1(.25))
	geomtest.Equal(t, d2.V{1, 1}, b.V1c0().V1(.25))
}

func TestBlossomAndSegment(t *testing.T) {
	b := Bezier{{0, 0}, {0.5, 1}, {1, 0}}
	pt := b.Blossom(.2, .3, .4)
	geomtest.Equal(t, d2.Pt{.25, .38}, pt)

	s := b.Segment(0.25, 0.75)
	geomtest.Equal(t, b.Pt1(0.25), s.Pt1(0))
	geomtest.Equal(t, b.Pt1(0.5), s.Pt1(0.5))
	geomtest.Equal(t, b.Pt1(0.75), s.Pt1(1))

	s = b.SegmentBuf(0.25, 0.75, make([]d2.Pt, 10))
	geomtest.Equal(t, b.Pt1(0.25), s.Pt1(0))
	geomtest.Equal(t, b.Pt1(0.5), s.Pt1(0.5))
	geomtest.Equal(t, b.Pt1(0.75), s.Pt1(1))
}
