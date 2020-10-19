package zbuf

import (
	"image"
	"image/color"
	"sync"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/ffmpeg"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

type Scene struct {
	W, H               int
	A                  angle.Rad
	Near, Far          float64
	Framerate          byte
	Name               string
	ConstantRateFactor byte
	Background         color.RGBA
	// ImageScale will render the image larger than the final size then scale it
	// down. This helps eliminate artifacts.
	ImageScale                float64
	wg                        sync.WaitGroup
	proc                      *ffmpeg.Proc
	toCameraTransform, toZbuf chan *SceneFrame
	toFF, recycleImg          chan *image.RGBA
	done                      chan bool
}

type SceneFrame struct {
	Camera
	Meshes    []*RenderMesh
	s         *Scene
	Wireframe bool
}

func (s *Scene) NewFrame(pt d3.Pt, q d3.Q, meshes int) *SceneFrame {
	return &SceneFrame{
		s:      s,
		Meshes: make([]*RenderMesh, 0, meshes),
		Camera: Camera{
			Pt:    pt,
			Q:     q,
			Near:  s.Near,
			Far:   s.Far,
			Angle: s.A,
			W:     int(float64(s.W) * s.ImageScale),
			H:     int(float64(s.H) * s.ImageScale),
		},
	}
}

// It would be more efficient to store the space transform then combine it with
// the camera transform and not save the space points. If the shader needs it,
// it can be computed then.

func (sf *SceneFrame) AddMesh(m *mesh.TriangleMesh, shader Shader, space *d3.T) {
	sf.Meshes = append(sf.Meshes, &RenderMesh{
		Original: m,
		Space:    space.Pts(m.Pts),
		Shader:   shader,
	})
}

func (s *Scene) init() {
	s.toCameraTransform = make(chan *SceneFrame)
	s.toZbuf = make(chan *SceneFrame)
	s.toFF = make(chan *image.RGBA)
	s.recycleImg = make(chan *image.RGBA)

	go s.cameraTransform()
	go s.zbuf()
	go s.ff()
}

func (sf *SceneFrame) Render() {
	if sf.s.toCameraTransform == nil {
		sf.s.init()
	}
	sf.s.toCameraTransform <- sf
}

func (s *Scene) Done() {
	s.done = make(chan bool)
	s.toCameraTransform <- nil
	<-s.done
}

func (s *Scene) cameraTransform() {
	for sf := range s.toCameraTransform {
		if sf == nil {
			break
		}
		ct := sf.Camera.T()
		scale := d3.Scale(d3.V{float64(sf.Camera.W), float64(sf.Camera.H), 1}).T()
		ct = ct.T(scale)
		for _, m := range sf.Meshes {
			m.Camera = ct.PtsScl(m.Space)
		}
		s.toZbuf <- sf
	}
	close(s.toZbuf)
}

func (s *Scene) zbuf() {
	// Todo: break this up into parallel renders
	buf := newZbuf(int(float64(s.W)*s.ImageScale), int(float64(s.H)*s.ImageScale))
	for sf := range s.toZbuf {
		buf.Reset()
		for _, m := range sf.Meshes {
			buf.Add(m)
		}
		img := <-s.recycleImg
		s.reset(img)
		buf.Draw(img)
		s.toFF <- img
	}
	<-s.recycleImg
	close(s.toFF)
}

func (s *Scene) ff() {
	proc := &ffmpeg.Proc{
		Framerate:          s.Framerate,
		Name:               s.Name,
		ConstantRateFactor: s.ConstantRateFactor,
		Width:              s.W,
		Height:             s.H,
	}
	proc.Start()
	s.recycleImg <- s.makeImg()
	for img := range s.toFF {
		proc.AddFrame(img)
		s.recycleImg <- img
	}
	proc.Close()
	s.done <- true
}

func (s *Scene) makeImg() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, int(float64(s.W)*s.ImageScale), int(float64(s.H)*s.ImageScale)))
}

func (s *Scene) reset(img *image.RGBA) {
	grid.Pt{img.Rect.Max.X, img.Rect.Max.Y}.Iter().Each(func(idx int, pt grid.Pt) {
		img.SetRGBA(pt.X, pt.Y, s.Background)
	})
}
