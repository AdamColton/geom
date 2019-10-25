package mesh

import (
	"github.com/adamcolton/geom/d3"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestExrudeExtrusion(t *testing.T) {
	f := []d3.Pt{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
	}

	tr := d3.Translate(d3.V{0, 0, 1}).T()
	m := NewExtrusion(f).Extrude(tr, tr).Close()
	assert.Len(t, m.Polygons, 8)
}

func TestEdgeExtrudeExtrusion(t *testing.T) {
	f := []d3.Pt{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
	}

	m := NewExtrusion(f).
		EdgeExtrude(d3.Scale(d3.V{2, 2, 0}).T()).
		Extrude(d3.Translate(d3.V{0, 0, 1}).T()).
		EdgeMerge(d3.Scale(d3.V{0.5, 0.5, 1}).T()).
		Close()

	assert.Len(t, m.Polygons, 23)
}
