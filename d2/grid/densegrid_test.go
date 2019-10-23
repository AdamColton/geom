package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDenseGrid(t *testing.T) {
	fn := func(pt Pt) interface{} {
		return pt.X * pt.Y
	}

	g := NewDenseGrid(Pt{5, 5}, fn)

	c := 0
	for i, done := g.Start(); !done; done = g.Next() {
		v, err := g.Get(i.Pt())
		assert.NoError(t, err)
		assert.Equal(t, fn(i.Pt()), v)
		c++
	}
	assert.Equal(t, c, g.Size().Area())

	assert.NoError(t, g.Set(Pt{2, 2}, 100))
	v, err := g.Get(Pt{2, 2})
	assert.NoError(t, err)
	assert.Equal(t, 100, v)

	var _ Grid = g
}

func TestDenseGridRangeCheck(t *testing.T) {
	dfg := NewDenseGrid(Pt{1, 1}, nil)
	pt := Pt{1, 1}

	assert.Error(t, dfg.Set(pt, 1))
	_, err := dfg.Get(pt)
	assert.Error(t, err)
}
