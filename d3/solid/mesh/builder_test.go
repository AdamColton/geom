package mesh

import (
	"bytes"
	"github.com/adamcolton/geom/d3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilder(t *testing.T) {
	b := NewBuilder()

	p0 := d3.Pt{0, 0, 0}
	p1 := d3.Pt{1.1, 0, 0}
	p2 := d3.Pt{0, 1, 0}
	p3 := d3.Pt{0, 0, 1}

	assert.NoError(t, b.Add([]d3.Pt{p0, p1, p2}))
	assert.NoError(t, b.Add([]d3.Pt{p0, p1, p3}))
	assert.NoError(t, b.Add([]d3.Pt{p1, p2, p3}))
	assert.False(t, b.Solid())

	assert.NoError(t, b.Add([]d3.Pt{p0, p2, p3}))
	assert.True(t, b.Solid())

	assert.Len(t, b.Pts, 4)

	buf := bytes.NewBuffer(nil)
	b.WriteObj(buf)

	//	m2, _ := ReadObj(buf)

	//assert.Equal(t, b.Mesh, *m2)
}

func TestExtrude(t *testing.T) {
	f := []d3.Pt{
		{0, 0, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
	}

	m := Extrude(f, d3.V{0, 0, 1})
	assert.Len(t, m.Pts, 8)
}
