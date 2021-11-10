package raytrace

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d3/render/material"
)

type MaterialWrapper struct {
	border                       float64
	borderMaterial, mainMaterial *Material
}

func NewMaterialWrapper(m material.Material) *MaterialWrapper {
	return &MaterialWrapper{
		border: m.Border,
		borderMaterial: &Material{
			Color:    m.BorderColor,
			Luminous: m.Luminous,
			Diffuse:  m.Diffuse,
		},
		mainMaterial: &Material{
			Color:    m.Color,
			Luminous: m.Luminous,
			Diffuse:  m.Diffuse,
		},
	}
}

func (m *MaterialWrapper) RayShader(ctx *Context) *Material {
	// if ctx.B.U < m.border || ctx.B.V < m.border || 1-ctx.B.U-ctx.B.V < m.border {
	// 	return m.BorderColor.RGBA()
	// }
	return m.mainMaterial
}

// func (m *MaterialWrapper) RayShader(ctx *Context) *Material {
// 	m :=
// 	if ctx.B.U < m.Border || ctx.B.V < m.Border || 1-ctx.B.U-ctx.B.V < m.Border {
// 		return m.BorderColor.RGBA()
// 	}

// 	n := ctx.TriangleRef.NSpace

// 	scale := ((angle.Rot(n.Z).Cos() + 1) / 4) + 0.5
// 	c := m.Color.Scale(scale)

// 	ptSpc := ctx.BSpace.PtB(ctx.B)
// 	v := ctx.SceneFrame.Camera.Pt.Subtract(ptSpc).Normal()
// 	ang := v.Ang(ctx.NSpace)
// 	if ang < m.Specular {
// 		d := float64(ang / (4 * m.Specular))
// 		l := angle.Rot(d).Cos()
// 		li := 1 - l
// 		bc := material.RGBAColor(&(ctx.SceneFrame.ZBufFrame.Background)).Scale(l)
// 		c = c.Scale(li)
// 		c.R += bc.R
// 		c.G += bc.G
// 		c.B += bc.B
// 	}

// 	return c.RGBA()
// }

type Material struct {
	Color    *material.Color
	Luminous float64
	Diffuse  angle.Rad
}

type Context struct {
	*Intersection
}

type Shader func(*Context) *Material

type RayShader interface {
	RayShader(*Context) *Material
}
