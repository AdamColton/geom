package render

import (
	"image/color"

	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

type RenderMesh struct {
	Original *mesh.TriangleMesh
	Space    []d3.Pt
	Camera   []d3.Pt
	Shader
}

func NewRenderMesh(tm *mesh.TriangleMesh, space, camera *d3.T, shader Shader) *RenderMesh {
	s := space.Pts(tm.Pts)
	c := camera.PtsScl(s)
	return &RenderMesh{
		Original: tm,
		Space:    s,
		Camera:   c,
		Shader:   shader,
	}
}

func (rm *RenderMesh) ApplyT(space, camera *d3.T) {
	rm.Space = space.Pts(rm.Original.Pts)
	rm.Camera = camera.PtsScl(rm.Space)
}

type Context struct {
	barycentric.B
	*RenderMesh
	PolygonIdx, TriangleIdx int
	d3.Pt
}

type Shader func(ctx *Context) *color.RGBA
