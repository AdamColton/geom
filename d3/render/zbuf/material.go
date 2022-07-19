package zbuf

import (
	"image/color"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d3/render/material"
)

type MaterialWrapper struct {
	material.Material
}

func (m *MaterialWrapper) ZBufShader(ctx *Context) *color.RGBA {
	if ctx.B.U < m.Border || ctx.B.V < m.Border || 1-ctx.B.U-ctx.B.V < m.Border {
		return m.BorderColor.RGBA()
	}

	n := ctx.TriangleRef.NSpace

	scale := ((angle.Rot(n.Z).Cos() + 1) / 4) + 0.5
	c := m.Color.Scale(scale)

	ptSpc := ctx.BSpace.PtB(ctx.B)
	v := ctx.Frame.Camera.Pt.Subtract(ptSpc).Normal()
	ang := v.Ang(ctx.NSpace)
	if ang < m.Specular {
		d := float64(ang / (4 * m.Specular))
		l := angle.Rot(d).Cos()
		li := 1 - l
		bc := material.RGBAColor(&(ctx.Frame.ZBufFrame.Background)).Scale(l)
		c = c.Scale(li)
		c.R += bc.R
		c.G += bc.G
		c.B += bc.B
	}

	return c.RGBA()
}
