package mesh

import (
	"bytes"
	"math"
	"sort"
	"testing"

	"github.com/adamcolton/geom/d3/solid"
	"github.com/adamcolton/geom/geomtest"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d3"
	"github.com/stretchr/testify/assert"
)

func TestFace(t *testing.T) {
	f := []d3.Pt{
		{0, 0, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
	}

	m := Extrude(f, d3.V{0, 0, 1})

	assert.Equal(t, f, m.Face(0))
}

func TestSingleDouble(t *testing.T) {
	tm := &TriangleMesh{
		Polygons: [][][3]uint32{
			{
				{1, 2, 3},
				{0, 1, 3},
				{0, 3, 4},
			},
		},
	}

	s, d := tm.Edges()
	assert.Len(t, s, 5)
	assert.Len(t, d, 2)
}

func TestMeshTransform(t *testing.T) {
	f := []d3.Pt{
		{2, 0, 0},
		{0, 2, 0},
		{-2, 0, 0},
		{0, -2, 0},
	}

	m := Extrude(f, d3.V{0, 0, 1})
	m = m.T(d3.Rotation{angle.Deg(45), d3.XY}.T())
	sq2 := math.Sqrt2
	expected := []d3.Pt{
		{sq2, sq2, 0},
		{-sq2, sq2, 0},
		{-sq2, -sq2, 0},
		{sq2, -math.Sqrt2, 0},
	}
	geomtest.Equal(t, expected, m.Face(0))
}

func TestMeshString(t *testing.T) {
	f := []d3.Pt{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
	}

	m := Extrude(f, d3.V{0, 0, 1})

	expected := `v 0 0 0
v 0 1 0
v 1 0 0
v 0 0 1
v 0 1 1
v 1 0 1
f 1 2 3
f 4 5 6
f 3 1 4 6
f 1 2 5 4
f 2 3 6 5
`
	assert.Equal(t, expected, m.String())
}

func TestMeshEdges(t *testing.T) {
	f := []d3.Pt{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
	}

	m := Extrude(f, d3.V{0, 0, 1})
	got := m.Edges()
	expected := []solid.IdxEdge{
		{0, 1},
		{0, 2},
		{0, 3},
		{1, 2},
		{1, 4},
		{2, 5},
		{3, 4},
		{3, 5},
		{4, 5},
	}
	sort.Slice(got, func(i, j int) bool {
		return got[i][0] < got[j][0] || (got[i][0] == got[j][0] && got[i][1] < got[j][1])
	})
	assert.Equal(t, expected, got)
}

func TestTriangleMesh(t *testing.T) {
	f := []d3.Pt{
		{0, 0, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
	}

	m := Extrude(f, d3.V{0, 0, 1})
	tm, err := m.TriangleMesh()
	assert.NoError(t, err)
	assert.Len(t, tm.Polygons, 6)
	for _, p := range tm.Polygons {
		assert.Len(t, p, 2)
	}
}

func TestTriangleMeshTransform(t *testing.T) {
	f := []d3.Pt{
		{2, 0, 0},
		{0, 2, 0},
		{-2, 0, 0},
		{0, -2, 0},
	}

	m := Extrude(f, d3.V{0, 0, 1})
	tm, err := m.TriangleMesh()
	assert.NoError(t, err)
	tm = tm.T(d3.Rotation{angle.Deg(45), d3.XY}.T())
	sq2 := math.Sqrt2
	expected := [][3]d3.Pt{
		{
			{sq2, sq2, 0},
			{-sq2, sq2, 0},
			{-sq2, -sq2, 0},
		}, {
			{sq2, sq2, 0},
			{-sq2, -sq2, 0},
			{sq2, -sq2, 0},
		},
	}
	got := tm.Face(0)
	assert.Len(t, got, 2)
	geomtest.Equal(t, expected[0][:], got[0][:])
	geomtest.Equal(t, expected[1][:], got[1][:])
}

func TestRoundXY(t *testing.T) {
	tm := TriangleMesh{
		Pts: []d3.Pt{
			{1.1, 2.2, 3.9},
		},
	}
	tm.RoundXY()
	geomtest.Equal(t, d3.Pt{1, 2, 3.9}, tm.Pts[0])
}

func TestReadObj(t *testing.T) {
	buf := bytes.NewBufferString(`v 0 0 0
v 0 1 0
v 1 0 0
v 0 0 1
v 0 1 1
v 1 0 1
f 1 2 3
f 4 5 6
f 3 1 4 6
f 1 2 5 4
f 2 3 6 5
	`)

	m, err := ReadObj(buf)
	assert.NoError(t, err)
	assert.Len(t, m.Pts, 6)
	assert.Len(t, m.Polygons, 5)
}
