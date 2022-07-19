package material

import "image/color"

// Color is an RGB color where each value is in the range 0 to 1.
type Color struct {
	R, G, B float64
}

// RGBA converts a color to an RGBA color with A set to 255.
func (c *Color) RGBA() *color.RGBA {
	return &color.RGBA{uint8(255 * c.R), uint8(255 * c.G), uint8(255 * c.B), 255}
}

// RGBAColor converts an RGBA color to Color.
func RGBAColor(c *color.RGBA) *Color {
	return &Color{
		R: float64(c.R) / 255,
		G: float64(c.G) / 255,
		B: float64(c.B) / 255,
	}
}

// Scale multiplies all values of a color by a scale factor.
func (c *Color) Scale(scale float64) *Color {
	return &Color{
		R: c.R * scale,
		G: c.G * scale,
		B: c.B * scale,
	}
}

// Reflect finds the product of a set of colors - this is for Raytracing
// reflections.
func Reflect(colors ...*Color) *Color {
	out := &Color{1, 1, 1}
	for _, c := range colors {
		out.R *= c.R
		out.G *= c.G
		out.B *= c.B
	}
	return out
}

// Radiate compounds the brightness of a set of colors.
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

// Avg finds the average of a set of colors.
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
