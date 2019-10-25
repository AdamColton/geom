package cc

import (
	"testing"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
	"github.com/stretchr/testify/assert"
)

func TestMethods(t *testing.T) {
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

func TestSubdivide(t *testing.T) {
	m := mesh.Extrude([]d3.Pt{
		{0, 0, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
	}, d3.V{0, 0, 1})

	out := Subdivide(m, 2)
	assert.Len(t, out.Polygons, 96)
}
