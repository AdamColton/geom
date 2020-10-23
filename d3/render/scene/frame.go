package scene

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

type SceneFrame struct {
	Camera *Camera
	Meshes []*FrameMesh
}

func (s *SceneFrame) TriangleCount() int {
	out := 0
	for _, m := range s.Meshes {
		out += m.Original.GetTriangleCount()
	}
	return out
}

type TransformFactory interface {
	T(frameIdx int) *d3.T
}

func (sf *SceneFrame) AddMesh(original *mesh.TriangleMesh, t *d3.T, shader interface{}) {
	sf.Meshes = append(sf.Meshes, &FrameMesh{
		Original: original,
		Shader:   shader,
		Space:    t.PtsScl(original.Pts),
	})
}
