package material

import "image/color"

type Color struct {
	R, G, B float64
}

func (c *Color) RGBA() *color.RGBA {
	return &color.RGBA{uint8(255 * c.R), uint8(255 * c.G), uint8(255 * c.B), 255}
}

func RGBAColor(c *color.RGBA) *Color {
	return &Color{
		R: float64(c.R) / 255,
		G: float64(c.G) / 255,
		B: float64(c.B) / 255,
	}
}

func (c *Color) Scale(scale float64) *Color {
	return &Color{
		R: c.R * scale,
		G: c.G * scale,
		B: c.B * scale,
	}
}

func Reflect(colors ...*Color) *Color {
	out := &Color{1, 1, 1}
	for _, c := range colors {
		out.R *= c.R
		out.G *= c.G
		out.B *= c.B
	}
	return out
}

func Radiate(colors ...*Color) *Color {
	out := &Color{1, 1, 1}
	for _, c := range colors {
		out.R *= (1 - c.R)
		out.G *= (1 - c.G)
		out.B *= (1 - c.B)
	}
	out.R = 1 - out.R
	out.G = 1 - out.G
	out.B = 1 - out.B
	return out
}

func Avg(colors ...*Color) *Color {
	out := &Color{0, 0, 0}
	if len(colors) == 0 {
		return out
	}
	for _, c := range colors {
		out.R += c.R
		out.G += c.G
		out.B += c.B
	}
	ln64 := float64(len(colors))
	out.R /= ln64
	out.G /= ln64
	out.B /= ln64
	return out
}
