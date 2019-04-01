package solid

import (
	"github.com/adamcolton/geom/d3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPairs(t *testing.T) {
	p0 := d3.Pt{0, 0, 0}
	p1 := d3.Pt{1.1, 0, 0}
	p2 := d3.Pt{0, 1, 0}
	p3 := d3.Pt{0, 0, 1}

	edges := [][2]d3.Pt{
		{p0, p1},
		{p1, p2},
		{p2, p0},
		{p1, p0},
		{p1, p3},
		{p3, p0},
		{p1, p2},
		{p2, p3},
		{p3, p1},
		{p0, p2},
		{p2, p3},
		{p3, p0},
	}

	em := NewEdgeMesh()
	for _, pts := range edges {
		assert.False(t, em.Solid())
		assert.NoError(t, em.Add(pts[0], pts[1]))
	}
	assert.True(t, em.Solid())
}

func TestTriangles(t *testing.T) {
	p0 := d3.Pt{0, 0, 0}
	p1 := d3.Pt{1.1, 0, 0}
	p2 := d3.Pt{0, 1, 0}
	p3 := d3.Pt{0, 0, 1}

	triangles := [][3]d3.Pt{
		{p0, p1, p2},
		{p0, p1, p3},
		{p1, p2, p3},
		{p0, p2, p3},
	}

	em := NewEdgeMesh()
	for _, tr := range triangles {
		assert.False(t, em.Solid())
		assert.NoError(t, em.Add(tr[0], tr[1], tr[2]))
	}
	assert.True(t, em.Solid())
}

func TestPairsIdx(t *testing.T) {
	edges := [][2]uint32{
		{0, 1},
		{1, 2},
		{2, 0},
		{1, 0},
		{1, 3},
		{3, 0},
		{1, 2},
		{2, 3},
		{3, 1},
		{0, 2},
		{2, 3},
		{3, 0},
	}

	em := NewIdxEdgeMesh()
	for _, pts := range edges {
		assert.False(t, em.Solid())
		assert.NoError(t, em.Add(pts[0], pts[1]))
	}
	assert.True(t, em.Solid())
}
