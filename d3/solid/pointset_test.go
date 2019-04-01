package solid

import (
	"github.com/adamcolton/geom/d3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPointSet(t *testing.T) {
	ps := NewPointSet()

	var p d3.Pt
	_, found := ps.Has(p)
	assert.False(t, found)

	idx := ps.Add(p)
	idx2, found := ps.Has(p)
	assert.False(t, found)
	assert.Equal(t, idx, idx2)

	p.X = 123
	idx = ps.Add(p)
	idx2, found = ps.Has(p)
	assert.False(t, found)
	assert.Equal(t, idx, idx2)
}
