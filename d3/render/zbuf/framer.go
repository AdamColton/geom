package zbuf

import (
	"image"
	"image/color"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/scene"
)

type Framer struct {
	Count      int
	Scene      *scene.Scene
	Near, Far  float64
	buf        *ZBuffer
	Background color.RGBA
	// ImageScale will render the image larger than the final size then scale it
	// down. This helps eliminate artifacts.
	ImageScale float64
	Shaders    []Shader
}

type SceneFrame struct {
	*scene.SceneFrame
	Near, Far    float64
	CameraMeshes [][]d3.Pt
	Background   color.RGBA
	Shaders      []Shader
}

func (f *Framer) Frame(idx int, img image.Image) (image.Image, error) {
	s := &SceneFrame{
		SceneFrame: f.Scene.Frame(idx),
		Near:       f.Near,
		Far:        f.Far,
		Background: f.Background,
		Shaders:    f.Shaders,
	}

	if f.Shaders == nil {
		s.PopulateShaders()
		f.Shaders = s.Shaders
	} else {
		s.Shaders = f.Shaders
	}

	if f.buf == nil {
		w := int(float64(s.Camera.Width) * f.ImageScale)
		h := int(float64(s.Camera.Height) * f.ImageScale)
		if f.buf == nil {
			f.buf = New(w, h)
		}
	}

	rgbaImg, _ := img.(*image.RGBA)
	return f.buf.Draw(s, rgbaImg), nil
}

func (f *Framer) Frames() int {
	return f.Count
}

func (sf *SceneFrame) PopulateCameraMeshes() {
	ln := len(sf.Meshes)
	sf.CameraMeshes = make([][]d3.Pt, ln)

	t := Camera{
		Camera: *sf.Camera,
		Near:   sf.Near,
		Far:    sf.Far,
	}.T()

	scene.RunRange(ln, func(mIdx, _ int) {
		sf.CameraMeshes[mIdx] = t.PtsScl(sf.Meshes[mIdx].Space)
	})
}

func (sf *SceneFrame) PopulateShaders() {
	ln := len(sf.Meshes)
	sf.Shaders = make([]Shader, ln)
	scene.RunRange(ln, func(idx, _ int) {
		s, ok := sf.Meshes[idx].Shader.(ZBufShader)
		if ok {
			sf.Shaders[idx] = s.ZBufShader
		}
	})
}
