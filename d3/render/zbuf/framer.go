package zbuf

import (
	"image"

	"github.com/adamcolton/geom/d3/render/scene"
)

type Framer struct {
	Count int
	Scene *scene.Scene
	*ZBufFrame
	// ImageScale will render the image larger than the final size then scale it
	// down. This helps eliminate artifacts.
	ImageScale float64
	buf        *ZBuffer
}

func (f *Framer) Frame(idx int, img image.Image) (image.Image, error) {
	s := &SceneFrame{
		SceneFrame: f.Scene.Frame(idx),
		ZBufFrame:  f.ZBufFrame,
	}

	if f.Shaders == nil {
		s.PopulateShaders()
	}

	if f.buf == nil {
		size := s.Camera.Size.Multiply(f.ImageScale)
		if f.buf == nil {
			f.buf = New(size.X, size.Y)
		}
	}

	rgbaImg, _ := img.(*image.RGBA)
	return f.buf.Draw(s, rgbaImg), nil
}

func (f *Framer) Frames() int {
	return f.Count
}
