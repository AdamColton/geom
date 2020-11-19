package boxmodel

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/ellipse"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	c := NewCompressor()
	tr, err := c.Add("triangles", New(shape.Subtract{
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
	}, 12))
	assert.NoError(t, err)

	_, err = c.Add("ellipse", New(ellipse.New(d2.Pt{100, 350}, d2.Pt{400, 110}, 170), 12))
	assert.NoError(t, err)

	_, err = c.Add("intersection", New(shape.Subtract{
		shape.Intersection{
			ellipse.NewCircle(d2.Pt{250, 250}, 230),
			&triangle.Triangle{
				{100, 250}, {490, 100}, {490, 400},
			},
		},
		ellipse.NewCircle(d2.Pt{350, 250}, 40),
	}, 12))
	assert.NoError(t, err)

	bm := tr
	outside := 0
	for c, _, done := bm.OutsideCursor(); !done; _, done = c.Next() {
		outside++
	}
	assert.Equal(t, outside, bm.Outside())

	inside := 0
	for c, _, done := bm.InsideCursor(); !done; _, done = c.Next() {
		inside++
	}
	assert.Equal(t, inside, bm.Inside())

	perimeter := 0
	for c, _, done := bm.PerimeterCursor(); !done; _, done = c.Next() {
		perimeter++
	}
	assert.Equal(t, perimeter, bm.Perimeter())
}
