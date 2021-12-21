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

type SceneFrame struct {
	*scene.SceneFrame
	*ZBufFrame
	CameraMeshes [][]d3.Pt
}

func (sf *SceneFrame) PopulateCameraMeshes() {
	ln := len(sf.Meshes)
	sf.CameraMeshes = make([][]d3.Pt, ln)

	t := (&Camera{
		Camera: sf.Camera,
		Near:   sf.Near,
		Far:    sf.Far,
	}).T()

	work.RunRange(ln, func(mIdx, _ int) {
		sf.CameraMeshes[mIdx] = t.PtsScl(sf.Meshes[mIdx].Space)
	})
}

func (sf *SceneFrame) PopulateShaders() {
	ln := len(sf.Meshes)
	sf.Shaders = make([]Shader, ln)
	work.RunRange(ln, func(idx, _ int) {
		sf.innerPopulateShaders(idx)
	})
}

func (sf *SceneFrame) innerPopulateShaders(idx int) {
	s := sf.Meshes[idx].Shader
	z, ok := s.(ZBufShader)
	if ok {
		sf.Shaders[idx] = z.ZBufShader
		return
	}
	m, ok := s.(*material.Material)
	if ok {
		mw := &MaterialWrapper{*m}
		sf.Shaders[idx] = mw.ZBufShader
	}
}
