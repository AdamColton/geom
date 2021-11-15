package bezier

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/stretchr/testify/assert"
)

func TestNewRelativeCompositeBezier(t *testing.T) {
	c := NewRelativeCompositeBezier([]Segment{
		{{0, 1}, {0, 1}, {1, 0}},
		{{0, -1}, {0, -1}, {1, 0}},
	}, d2.IndentityTransform())

	assert.Equal(t, d2.Pt{1, -2.25}, c.Pt1(-0.25))
	assert.Equal(t, d2.Pt{0, 0}, c.Pt1(0))
	assert.Equal(t, d2.Pt{0.5000, 0.7500}, c.Pt1(0.25))
	assert.Equal(t, d2.Pt{1, 0}, c.Pt1(0.5))
	assert.Equal(t, d2.Pt{1.5, -0.75}, c.Pt1(0.75))
	assert.Equal(t, d2.Pt{2, 0}, c.Pt1(1))
	assert.Equal(t, d2.Pt{1, 2.25}, c.Pt1(1.25))
}
