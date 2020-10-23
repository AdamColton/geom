package scene

import (
	"github.com/adamcolton/geom/d3/solid/mesh"
	"github.com/adamcolton/geom/work"
)

type Scene struct {
	CameraFactory
	Meshes []*Mesh
}

func NewScene(camera CameraFactory, meshes int) *Scene {
	return &Scene{
		CameraFactory: camera,
		Meshes:        make([]*Mesh, 0, meshes),
	}
}

func (s *Scene) AddMesh(original *mesh.TriangleMesh, tf TransformFactory, shader interface{}) {
	s.Meshes = append(s.Meshes, &Mesh{
		Original:         original,
		TransformFactory: tf,
		Shader:           shader,
	})
}

func (s *Scene) Frame(frameIdx int) *SceneFrame {
	ln := len(s.Meshes)
	sf := &SceneFrame{
		Meshes: make([]*FrameMesh, ln),
		Camera: s.Camera(frameIdx),
	}

	work.RunRange(ln, func(mIdx, _ int) {
		m := s.Meshes[mIdx]
		sf.Meshes[mIdx] = &FrameMesh{
			Original: m.Original,
			Shader:   m.Shader,
			Space:    m.T(frameIdx).PtsScl(m.Original.Pts),
		}
	})

	return sf
}
