package raytrace

import (
	"image"

	"github.com/adamcolton/geom/d3/render/scene"
)

type Scene struct {
	*scene.Scene
	*RayFrame
}

func (f *Scene) Frame(idx int, img image.Image) (image.Image, error) {
	s := &Frame{
		Frame:    f.Scene.Frame(idx),
		RayFrame: f.RayFrame,
	}

	if f.Shaders == nil {
		s.PopulateShaders()
	}

	rgbaImg, _ := img.(*image.RGBA)
	return s.Image(rgbaImg), nil
}
