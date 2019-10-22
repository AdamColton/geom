package solid

import (
	"sort"
	"testing"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
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

func TestLine(t *testing.T) {
	a, b := d3.Pt{2, 3, 1}, d3.Pt{3, 0, 4}
	e := NewEdge(a, b)
	geomtest.Equal(t, a, e.Pt1(0))
	geomtest.Equal(t, b, e.Pt1(1))
}

func TestEdgeMeshErrors(t *testing.T) {
	a, b := d3.Pt{2, 3, 1}, d3.Pt{3, 0, 4}
	em := NewEdgeMesh()
	err := em.Add(a, b)
	assert.NoError(t, err)
	err = em.Add(a, b)
	assert.NoError(t, err)
	err = em.Add(a, b)
	assert.Error(t, err)
	assert.Equal(t, ErrEdgeOverUsed{}.Error(), err.Error())
	err = em.Add(a, b, d3.Pt{})
	assert.Error(t, err)
	err = em.Add(a)
	assert.Error(t, err)
	assert.Equal(t, ErrTwoPoints{}.Error(), err.Error())
}

func TestIdxEdgeStr(t *testing.T) {
	e := IdxEdge{2, 3}
	assert.Equal(t, "[2, 3]", e.String())
}

func TestIdxEdgeMeshErrors(t *testing.T) {
	var (
		a uint32 = 3
		b uint32 = 4
	)
	em := NewIdxEdgeMesh()
	err := em.Add(a, b)
	assert.NoError(t, err)
	err = em.Add(a, b)
	assert.NoError(t, err)
	err = em.Add(a, b)
	assert.Error(t, err)
	assert.Equal(t, ErrEdgeOverUsed{}.Error(), err.Error())
	err = em.Add(a, b, 0)
	assert.Error(t, err)
	err = em.Add(a)
	assert.Error(t, err)
	assert.Equal(t, ErrTwoPoints{}.Error(), err.Error())
	err = em.Add(5, 6, 7)
	assert.NoError(t, err)
}

func TestIdxMeshEdges(t *testing.T) {
	em := NewIdxEdgeMesh()
	em.Add(1, 2, 3)
	got := em.Edges()
	sortEdges(got)
	expected := []IdxEdge{
		{1, 2},
		{1, 3},
		{2, 3},
	}
	assert.Equal(t, expected, got)
}

func sortEdges(es []IdxEdge) {
	sort.Slice(es, func(i, j int) bool {
		return es[i][0] < es[j][0] || (es[i][0] == es[j][0] && es[i][1] <= es[j][1])
	})
}

func TestIdxMeshSingleDouble(t *testing.T) {
	em := NewIdxEdgeMesh()
	em.Add(1, 2, 3)
	em.Add(1, 3, 4)
	gots, gotd := em.SingleDouble()
	sortEdges(gots)
	sortEdges(gotd)
	expecteds := []IdxEdge{
		{1, 2},
		{1, 4},
		{2, 3},
		{3, 4},
	}
	expectedd := []IdxEdge{
		{1, 3},
	}
	assert.Equal(t, expecteds, gots)
	assert.Equal(t, expectedd, gotd)
}
