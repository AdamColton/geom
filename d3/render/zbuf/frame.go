package zbuf

import (
	"image/color"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/material"
	"github.com/adamcolton/geom/d3/render/scene"
	"github.com/adamcolton/geom/work"
)

type ZBufFrame struct {
	Near, Far  float64
	Background color.RGBA
	Shaders    []Shader
}

type Frame struct {
	*scene.Frame
	*ZBufFrame
	CameraMeshes [][]d3.Pt
}

func (f *Frame) PopulateCameraMeshes() {
	ln := len(f.Meshes)
	f.CameraMeshes = make([][]d3.Pt, ln)

	t := (&Camera{
		Camera: f.Camera,
		Near:   f.Near,
		Far:    f.Far,
	}).T()

	work.RunRange(ln, func(mIdx, _ int) {
		f.CameraMeshes[mIdx] = t.PtsScl(f.Meshes[mIdx].Space)
	})
}

func (f *Frame) PopulateShaders() {
	ln := len(f.Meshes)
	f.Shaders = make([]Shader, ln)
	work.RunRange(ln, func(idx, _ int) {
		f.innerPopulateShaders(idx)
	})
}

func (f *Frame) innerPopulateShaders(idx int) {
	s := f.Meshes[idx].Shader
	z, ok := s.(ZBufShader)
	if ok {
		f.Shaders[idx] = z.ZBufShader
		return
	}
	m, ok := s.(*material.Material)
	if ok {
		mw := &MaterialWrapper{*m}
		f.Shaders[idx] = mw.ZBufShader
	}
}
