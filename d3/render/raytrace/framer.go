package raytrace

import (
	"image"

	"github.com/adamcolton/geom/d3/render/scene"
)

type Framer struct {
	Count int
	Scene *scene.Scene
	*RayFrame
}

func (f *Framer) Frames() int {
	return f.Count
}

func (f *Framer) Frame(idx int, img image.Image) (image.Image, error) {
	s := &SceneFrame{
		SceneFrame: f.Scene.Frame(idx),
		RayFrame:   f.RayFrame,
	}

	if f.Shaders == nil {
		s.PopulateShaders()
	}

	rgbaImg, _ := img.(*image.RGBA)
	return s.Image(rgbaImg), nil
}
