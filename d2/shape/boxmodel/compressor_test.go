package boxmodel

import (
	"testing"

	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	s := shape.Subtract{
		shape.Union{
			shape.Subtract{
				&triangle.Triangle{
					{20, 20}, {480, 40}, {250, 470},
				},
				&triangle.Triangle{
					{200, 200}, {400, 400}, {250, 50},
				},
			},
			&triangle.Triangle{
				{20, 200}, {50, 30}, {200, 200},
			},
		},
		&triangle.Triangle{
			{100, 100}, {100, 150}, {150, 100},
		},
	}
	bm := New(s, 12)
	ctr, err := NewCompressor().Add("triangles", bm)
	assert.NoError(t, err)

	outside := 0
	for c, b, done := ctr.OutsideCursor(); !done; b, done = c.Next() {
		assert.False(t, s.Contains(b.Centroid()))
		outside++
	}
	assert.Equal(t, outside, ctr.Outside())

	inside := 0
	for c, b, done := ctr.InsideCursor(); !done; b, done = c.Next() {
		assert.True(t, s.Contains(b.Centroid()))
		inside++
	}
	assert.Equal(t, inside, ctr.Inside())

	perimeter := 0
	for c, _, done := ctr.PerimeterCursor(); !done; _, done = c.Next() {
		perimeter++
	}
	assert.Equal(t, perimeter, ctr.Perimeter())

	// reduces node count by more than 10x
	assert.Less(t, len(ctr.tree().nodes), len(bm.tree().nodes)/10)
}
