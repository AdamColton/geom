package cc

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFacetMesh(t *testing.T) {
	f := Facet{
		Vertex{
			Pt:   d3.Pt{0, 0, 0},
			Next: 0.1,
			Prev: 0.1,
		},
		Vertex{
			Pt:   d3.Pt{0, 1, 0},
			Next: 0.1,
			Prev: 0.1,
		},
		Vertex{
			Pt:   d3.Pt{1, 1, 0},
			Next: 0.1,
			Prev: 0.1,
		},
		Vertex{
			Pt:   d3.Pt{1, 0, 0},
			Next: 0.1,
			Prev: 0.1,
		},
	}
	fm := f.polygons()
	assert.NotNil(t, fm)

	expected := facetMesh{
		pts: []d3.Pt{
			{0.0000, 0.0000, 0.0000},
			{0.0000, 0.1000, 0.0000},
			{0.0000, 0.9000, 0.0000},
			{0.0000, 1.0000, 0.0000},
			{0.1000, 1.0000, 0.0000},
			{0.9000, 1.0000, 0.0000},
			{1.0000, 1.0000, 0.0000},
			{1.0000, 0.9000, 0.0000},
			{1.0000, 0.1000, 0.0000},
			{1.0000, 0.0000, 0.0000},
			{0.9000, 0.0000, 0.0000},
			{0.1000, 0.0000, 0.0000},
			{0.1000, 0.1000, 0.0000},
			{0.1000, 0.9000, 0.0000},
			{0.9000, 0.9000, 0.0000},
			{0.9000, 0.1000, 0.0000},
		},
		facets: [][]int{
			{0, 1, 12, 11},
			{1, 2, 13, 12},
			{3, 4, 13, 2},
			{4, 5, 14, 13},
			{6, 7, 14, 5},
			{7, 8, 15, 14},
			{9, 10, 15, 8},
			{10, 11, 12, 15},
			{12, 13, 14, 15},
		},
	}

	assert.Equal(t, expected.facets, fm.facets)

	for i, p := range fm.pts {
		assert.InDelta(t, 0, p.Distance(expected.pts[i]), 1E-6)
	}
}

func TestStuff(t *testing.T) {
	m := mesh.Extrude([]d3.Pt{
		{0, 0, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
	}, d3.V{0, 0, 1})
	cc := ccMesh{
		Mesh: m,
	}

	cc.populateEdges()
	assert.Len(t, cc.pt2edge, len(m.Pts))

	cc.setFacePoints()
	assert.Equal(t, len(cc.Mesh.Polygons), len(cc.facePoints))
	assert.Equal(t, d3.Pt{0.5, 0.5, 0}, cc.facePoints[0])

	cc.setEdgePoints()
	assert.Len(t, cc.edgePoints, len(cc.edge2face))

	cc.setBaryPoints()
	assert.Equal(t, `Pt(0.2222, 0.2222, 0.2222)`, cc.baryPoints[0].String())

	m2 := cc.subdivide()
	assert.True(t, len(m2.Pts) > 0)
}
