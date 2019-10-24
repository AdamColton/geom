package mesh

import (
	"github.com/adamcolton/geom/d3"
	//"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExtrusion(t *testing.T) {
	f := []d3.Pt{
		{0, 0, 0},
		{0.5, -0.25, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0.5, 1.25, 0},
		{0, 1, 0},
	}

	m := NewExtrusion(f).
		EdgeExtrude(d3.Scale(d3.V{2, 2, 1}).T()).
		Extrude(
			d3.Translate(d3.V{0, 0, 1}).T(),
			d3.Translate(d3.V{0, 0, 2}).T(),
			d3.Translate(d3.V{0, 0, 1}).T(),
		).
		EdgeMerge(d3.Scale(d3.V{0.5, 0.5, 1}).T()).
		Close()

	file, _ := os.Create("temp.obj")
	m.WriteObj(file)
	file.Close()
}
