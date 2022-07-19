package scene

import (
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/geomtest"
)

func TestCamera(t *testing.T) {
	c := NewCamera(d3.Pt{0, 0, 0}, angle.Rad(0))

	c.SetSize(200, 300)
	geomtest.Equal(t, grid.Pt{X: 200, Y: 300}, c.Size)

	c.Widescreen(500)
	geomtest.Equal(t, grid.Pt{X: 500, Y: 281}, c.Size)

	c.Square(450)
	geomtest.Equal(t, grid.Pt{X: 450, Y: 450}, c.Size)

	c.Fullscreen(400)
	geomtest.Equal(t, grid.Pt{X: 400, Y: 300}, c.Size)
}
