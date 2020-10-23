package raytrace

import (
	"math/rand"
	"testing"

	"github.com/adamcolton/geom/angle"

	"github.com/adamcolton/geom/geomtest"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/stretchr/testify/assert"
)

func TestReflect(t *testing.T) {
	// Reflection doesn't care which way the normal vector is pointing.

	tri := &triangle.Triangle{
		{-1, -1, 0},
		{1, -1, 0},
		{-1, 1, 0},
	}
	l := line.New(d3.Pt{-1, -1, -1}, d3.Pt{0, 0, 0})
	n := tri.Normal().Normal()
	geomtest.Equal(t, d3.V{0, 0, 1}, n)
	i := tri.Intersection(l)
	assert.True(t, i.Does)
	r := line.Line{
		T0: l.Pt1(i.T),
		D:  reflect(l.D, n),
	}
	geomtest.Equal(t, d3.Pt{1, 1, -1}, r.Pt1(1.0))

	tri = &triangle.Triangle{
		{-1, -1, 0},
		{-1, 1, 0},
		{1, -1, 0},
	}
	l = line.New(d3.Pt{-1, -1, -1}, d3.Pt{0, 0, 0})
	n = tri.Normal().Normal()
	geomtest.Equal(t, d3.V{0, 0, -1}, n)
	i = tri.Intersection(l)
	assert.True(t, i.Does)
	r = line.Line{
		T0: l.Pt1(i.T),
		D:  reflect(l.D, n),
	}
	geomtest.Equal(t, d3.Pt{1, 1, -1}, r.Pt1(1.0))
}

func TestRandomAngle(t *testing.T) {
	v := randomAngle(angle.Rad(0)).T().V(d3.V{1, 0, 0})
	geomtest.Equal(t, d3.V{1, 0, 0}, v)

	for ang := angle.Rad(0); ang.Deg() < 45; ang += angle.Deg(1) {
		for i := 0; i < 100; i++ {
			vIn := d3.V{rand.Float64(), rand.Float64(), rand.Float64()}
			vIn = vIn.Normal()
			q := randomAngle(ang)
			// the angle should be normal
			assert.InDelta(t, 1.0, q.A*q.A+q.B*q.B+q.C*q.C, 1e-5)

			vOut := q.T().V(vIn)
			out := float64(vIn.Ang(vOut))
			rad := ang.Rad()
			assert.True(t, -rad <= out && out <= rad)
		}
	}
}
