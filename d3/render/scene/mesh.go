package scene

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

// TriangleIndex provides an index to a specific triangle in a SceneFrame.
type TriangleIndex struct {
	MeshIdx     int
	PolygonIdx  int
	TriangleIdx int
}

// Mesh is used to generate a FrameMesh.
type Mesh struct {
	Original *mesh.TriangleMesh
	TransformFactory
	Shader interface{}
}

// FrameMesh represents a mesh transformed for a specific frame. The Space slice
// holds the points of the mesh on the specific frame.
type FrameMesh struct {
	Original *mesh.TriangleMesh
	Shader   interface{}
	Space    []d3.Pt
}

// Frame generates the FrameMesh for a specific frame.
func (m *Mesh) Frame(idx int) *FrameMesh {
	return &FrameMesh{
		Original: m.Original,
		Shader:   m.Shader,
		Space:    m.TransformFactory.T(idx).Pts(m.Original.Pts),
	}
}
