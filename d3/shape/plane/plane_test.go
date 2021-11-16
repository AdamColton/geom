package plane

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestPlane(t *testing.T) {
	p := New(d3.Pt{0, 0, 0}, d3.Pt{2, 0, 0}, d3.Pt{5, 5, 0})
	geomtest.Equal(t, d3.Pt{0, 0, 0}, p.Origin)
	geomtest.Equal(t, d3.V{1, 0, 0}, p.X)
	geomtest.Equal(t, d3.V{0, 1, 0}, p.Y)
	geomtest.Equal(t, d3.V{0, 0, 1}, p.N)

	p = New(d3.Pt{0, 0, 0}, d3.Pt{2, 0, 0}, d3.Pt{5, -5, 0})
	geomtest.Equal(t, d3.Pt{0, 0, 0}, p.Origin)
	geomtest.Equal(t, d3.V{1, 0, 0}, p.X)
	geomtest.Equal(t, d3.V{0, -1, 0}, p.Y)
	geomtest.Equal(t, d3.V{0, 0, -1}, p.N)
}

func TestProject(t *testing.T) {
	p := New(d3.Pt{0, 0, 0}, d3.Pt{2, 0, 0}, d3.Pt{5, 5, 0})
	d2Pt, v := p.Project(d3.Pt{1, 2, 3})
	geomtest.Equal(t, d2.Pt{1, 2}, d2Pt)
	geomtest.Equal(t, d3.V{0, 0, 3}, v)

	d2Pt, v = p.Project(d3.Pt{-1, 2, 0})
	geomtest.Equal(t, d2.Pt{-1, 2}, d2Pt)
	geomtest.Equal(t, d3.V{0, 0, 0}, v)

	p = New(d3.Pt{0, 0, 0}, d3.Pt{2, 0, 0}, d3.Pt{5, -5, 0})
	d2Pt, v = p.Project(d3.Pt{1, 2, 3})
	geomtest.Equal(t, d2.Pt{1, -2}, d2Pt)
	geomtest.Equal(t, d3.V{0, 0, 3}, v)
	geomtest.Equal(t, d3.Pt{1, 2, 0}, p.Convert(d2Pt))

	d2Pt, v = p.Project(d3.Pt{-1, 2, 0})
	geomtest.Equal(t, d2.Pt{-1, -2}, d2Pt)
	geomtest.Equal(t, d3.V{0, 0, 0}, v)
}

func TestString(t *testing.T) {
	p := New(d3.Pt{0, 0, 0}, d3.Pt{2, 0, 0}, d3.Pt{5, 5, 0})
	assert.Equal(t, "Plane( Origin:Pt(0.0000, 0.0000, 0.0000) X:V(1.0000, 0.0000, 0.0000) Y:V(0.0000, 1.0000, -0.0000) N:V(-0.0000, 0.0000, 1.0000))", p.String())
}
