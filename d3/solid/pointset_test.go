package solid

import (
	"testing"

	"github.com/adamcolton/geom/d3"
	"github.com/stretchr/testify/assert"
)

func TestPointSet(t *testing.T) {
	ps := NewPointSet()

	var p d3.Pt
	_, found := ps.Has(p)
	assert.False(t, found)

	idx := ps.Add(p)
	idx2, found := ps.Has(p)
	assert.True(t, found)
	assert.Equal(t, idx, idx2)
	assert.Equal(t, idx, ps.Add(p))

	p.X = 123
	idx = ps.Add(p)
	idx2, found = ps.Has(p)
	assert.True(t, found)
	assert.Equal(t, idx, idx2)
}
