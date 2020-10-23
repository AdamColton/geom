package zbuf

import (
	"image"

	"github.com/adamcolton/geom/d3/render/scene"
)

type Scene struct {
	*scene.Scene
	*ZBufFrame
	// ImageScale will render the image larger than the final size then scale it
	// down. This helps eliminate artifacts.
	ImageScale float64
	buf        *ZBuffer
}

func (s *Scene) Frame(idx int, img image.Image) (image.Image, error) {
	f := &Frame{
		Frame:     s.Scene.Frame(idx),
		ZBufFrame: s.ZBufFrame,
	}

	if s.Shaders == nil {
		f.PopulateShaders()
	}

	if s.buf == nil {
		size := f.Camera.Size.Multiply(s.ImageScale)
		if s.buf == nil {
			s.buf = New(size.X, size.Y)
		}
	}

	rgbaImg, _ := img.(*image.RGBA)
	return s.buf.Draw(f, rgbaImg), nil
}
