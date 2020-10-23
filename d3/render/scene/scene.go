package scene

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

type CameraFactory interface {
	Camera(frameIdx int) *Camera
}

type Scene struct {
	GetCamera CameraFactory
	Meshes    []*Mesh
}

func NewScene(camera CameraFactory, meshes int) *Scene {
	return &Scene{
		GetCamera: camera,
		Meshes:    make([]*Mesh, 0, meshes),
	}
}

func (s *Scene) AddMesh(original *mesh.TriangleMesh, tf TransformFactory, shader interface{}) {
	s.Meshes = append(s.Meshes, &Mesh{
		Original:         original,
		TransformFactory: tf,
		Shader:           shader,
	})
}

type TransformFactory interface {
	T(frameIdx int) *d3.T
}

type Mesh struct {
	Original *mesh.TriangleMesh
	TransformFactory
	Shader interface{}
}

type FrameMesh struct {
	Original *mesh.TriangleMesh
	Shader   interface{}
	Space    []d3.Pt
}

type SceneFrame struct {
	Camera *Camera
	Meshes []*FrameMesh
}

func (sf *SceneFrame) AddMesh(original *mesh.TriangleMesh, t *d3.T, shader interface{}) {
	sf.Meshes = append(sf.Meshes, &FrameMesh{
		Original: original,
		Shader:   shader,
		Space:    t.PtsScl(original.Pts),
	})
}

func (s *Scene) Frame(frameIdx int) *SceneFrame {
	ln := len(s.Meshes)
	sf := &SceneFrame{
		Meshes: make([]*FrameMesh, ln),
	}
	sf.Camera = s.GetCamera.Camera(frameIdx)

	RunRange(ln, func(mIdx, _ int) {
		m := s.Meshes[mIdx]
		fm := &FrameMesh{
			Original: m.Original,
			Shader:   m.Shader,
		}
		fm.Space = m.T(frameIdx).PtsScl(fm.Original.Pts)
		sf.Meshes[mIdx] = fm
	})

	return sf
}
