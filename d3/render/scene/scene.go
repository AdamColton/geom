package scene

import (
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d3/solid/mesh"
	"github.com/adamcolton/geom/work"
)

// Scene is used to generate the camera and meshest for a sequence of frames.
// Size will populate the size on the camera if the CameraFactory does not set
// the size.
type Scene struct {
	CameraFactory
	Meshes     []*Mesh
	FrameCount int
	ImgSize    grid.Pt
}

// NewScene creates Scene and allocates the Mesh slice.
func NewScene(camera CameraFactory, meshes int) *Scene {
	return &Scene{
		CameraFactory: camera,
		Meshes:        make([]*Mesh, 0, meshes),
	}
}

// Size returns the scene size and fulfills Sizer.
func (s *Scene) Size() grid.Pt {
	return s.ImgSize
}

// Frames is the length of the scene and fullfils Framer.
func (s *Scene) Frames() int {
	return s.FrameCount
}

// AddMesh appends a Mesh to the scene.
func (s *Scene) AddMesh(original *mesh.TriangleMesh, tf TransformFactory, shader interface{}) {
	s.Meshes = append(s.Meshes, &Mesh{
		Original:         original,
		TransformFactory: tf,
		Shader:           shader,
	})
}

// Frame generates a specific SceneFrame.
func (s *Scene) Frame(frameIdx int) *Frame {
	ln := len(s.Meshes)
	f := &Frame{
		Meshes: make([]*FrameMesh, ln),
		Camera: s.Camera(frameIdx),
	}

	if f.Camera.Size.X == 0 {
		f.Camera.Size = s.ImgSize
	}

	work.RunRange(ln, func(mIdx, _ int) {
		f.Meshes[mIdx] = s.Meshes[mIdx].Frame(frameIdx)
	})

	return f
}
