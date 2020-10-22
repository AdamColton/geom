package zbuf

import (
	"image"
	"image/color"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/ffmpeg"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

type Scene struct {
	Camera
	Framerate          byte
	Name               string
	ConstantRateFactor byte
	Background         color.RGBA
	// ImageScale will render the image larger than the final size then scale it
	// down. This helps eliminate artifacts.
	ImageScale                float64
	wg                        sync.WaitGroup
	proc                      *ffmpeg.Proc
	toCameraTransform, toZbuf chan *Frame
	toFF, recycleImg          chan *image.RGBA
	done                      chan bool
}

type Frame struct {
	Camera
	Meshes    []*RenderMesh
	s         *Scene
	Wireframe bool
}

func (s *Scene) NewFrame(meshes int) *Frame {
	return &Frame{
		s:      s,
		Meshes: make([]*RenderMesh, 0, meshes),
		Camera: s.Camera,
	}
}

func (sf *Frame) AddMesh(m *mesh.TriangleMesh, shader Shader, space *d3.T) {
	sf.Meshes = append(sf.Meshes, &RenderMesh{
		Original: m,
		Space:    space.Pts(m.Pts),
		Shader:   shader,
	})
}

func (s *Scene) init() {
	s.toCameraTransform = make(chan *Frame)
	s.toZbuf = make(chan *Frame)
	s.toFF = make(chan *image.RGBA)
	s.recycleImg = make(chan *image.RGBA)

	go s.cameraTransform()
	go s.zbuf()
	go s.ff()
}

func (sf *Frame) Render() {
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
	cpus := runtime.NumCPU()
	for sf := range s.toCameraTransform {
		if sf == nil {
			break
		}
		ct := sf.Camera.T()
		scale := d3.Scale(d3.V{float64(sf.Camera.Width), float64(sf.Camera.Height), 1}).T()
		ct = ct.T(scale)

		var idx32 int32 = -1
		ln := len(sf.Meshes)
		wg := &sync.WaitGroup{}
		wg.Add(cpus)
		fn := func() {
			for {
				idx := int(atomic.AddInt32(&idx32, 1))
				if idx >= ln {
					break
				}
				m := sf.Meshes[idx]
				m.Camera = ct.PtsScl(m.Space)
			}
			wg.Add(-1)
		}
		for i := 0; i < cpus; i++ {
			go fn()
		}
		wg.Wait()
		s.toZbuf <- sf
	}
	close(s.toZbuf)
}

func (s *Scene) zbuf() {
	// Todo: break this up into parallel renders
	buf := newZbuf(int(float64(s.Camera.Width)*s.ImageScale), int(float64(s.Camera.Height)*s.ImageScale), &s.Background)
	for sf := range s.toZbuf {
		for _, m := range sf.Meshes {
			buf.Add(m)
		}
		img := <-s.recycleImg
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
		Width:              s.Camera.Width,
		Height:             s.Camera.Height,
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
	return image.NewRGBA(image.Rect(0, 0, int(float64(s.Camera.Width)*s.ImageScale), int(float64(s.Camera.Height)*s.ImageScale)))
}
