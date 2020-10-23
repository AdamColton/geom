package scene

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

// Frame holds the camera to produce the frame and the meshes to render.
type Frame struct {
	Camera *Camera
	Meshes []*FrameMesh
}

// TriangleCount is the total number of triangles in the frame.
func (f *Frame) TriangleCount() int {
	out := 0
	for _, m := range f.Meshes {
		out += m.Original.GetTriangleCount()
	}
	return out
}

// TransformFactory produces a transform per frame index.
type TransformFactory interface {
	T(frameIdx int) *d3.T
}

// AddMesh to a sceneframe. The transform will be applied to all the points in
// the original mesh.
func (f *Frame) AddMesh(original *mesh.TriangleMesh, t *d3.T, shader interface{}) {
	f.Meshes = append(f.Meshes, &FrameMesh{
		Original: original,
		Shader:   shader,
		Space:    t.PtsScl(original.Pts),
	})
}
