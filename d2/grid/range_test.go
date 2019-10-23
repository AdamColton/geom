package grid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRange(t *testing.T) {
	r := Range{
		{0, 0},
		{3, 3},
	}
	assert.True(t, r.Contains(Pt{1, 1}))
	assert.True(t, r.Contains(Pt{0, 0}))
	assert.False(t, r.Contains(Pt{3, 3}))
	assert.Equal(t, Pt{0, 0}, r.Min())
	assert.Equal(t, Pt{2, 2}, r.Max())

	r = Range{
		{0, 0},
		{-3, -3},
	}
	assert.True(t, r.Contains(Pt{-1, -1}))
	assert.True(t, r.Contains(Pt{0, 0}))
	assert.False(t, r.Contains(Pt{-3, -3}))
	assert.Equal(t, Pt{-2, -2}, r.Min())
	assert.Equal(t, Pt{0, 0}, r.Max())

	r = Range{
		{-3, -3},
		{0, 0},
	}
	assert.Equal(t, Pt{-3, -3}, r.Min())
	assert.Equal(t, Pt{-1, -1}, r.Max())

	r = Range{
		{3, 3},
		{0, 0},
	}
	assert.Equal(t, Pt{1, 1}, r.Min())
	assert.Equal(t, Pt{3, 3}, r.Max())
}
