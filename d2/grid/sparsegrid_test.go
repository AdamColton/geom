package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSparseGrid(t *testing.T) {
	fn := func(pt Pt) interface{} {
		return pt.X * pt.Y
	}

	g := NewSparseGrid(Pt{5, 5}, fn)

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

	v, err = g.Get(Pt{20, 20})
	assert.NoError(t, err)
	assert.Equal(t, 400, v)

	assert.NoError(t, g.Set(Pt{20, 20}, nil))
	_, found := g.Data[Pt{20, 20}]
	assert.False(t, found)

	var _ Grid = g
}
