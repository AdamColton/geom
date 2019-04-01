package mesh

import (
	"github.com/adamcolton/geom/d3"
	"github.com/stretchr/testify/assert"
	"testing"
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
