package zbuf

import (
	"image/color"

	"github.com/adamcolton/geom/d3"

	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3/shape/triangle"
)

// TriangleRef is a reference to a triangle in the scene. The MeshIdx,
// PolygonIdx and TriangleIdx provide the indexes to pull the value out of a
// SceneFrame. Space represents the triangle in scene space and Camera
// represents the triangle in Camera space (after the camera transform).
// NSpace and NCamera provide the respective normals.
type TriangleRef struct {
	MeshIdx         int
	PolygonIdx      int
	TriangleIdx     int
	Space, Camera   *triangle.Triangle
	NSpace, NCamera d3.V
	BSpace, BCamera *triangle.BT
	Orgin, U        int
}

type Context struct {
	*SceneFrame
	*TriangleRef
	barycentric.B
}

type Shader func(ctx *Context) *color.RGBA

type ZBufShader interface {
	ZBufShader(ctx *Context) *color.RGBA
}
