package drawgrid

import (
	"github.com/adamcolton/geom/d2/grid"
	"image"
	"image/color"
	"image/png"
	"os"
)

func FloatToGrey(i interface{}) color.Color {
	f := i.(float64)
	g := uint8(f * 255)
	return color.RGBA{g, g, g, 255}
}

func Image(g grid.Grid, toColor func(interface{}) color.Color) *image.RGBA {
	s := g.Size()
	img := image.NewRGBA(image.Rect(0, 0, s.X, s.Y))
	for i, ok := g.Start(); ok; ok = i.Next() {
		pt := i.Pt()
		img.Set(pt.X, pt.Y, toColor(g.Get(pt)))
	}
	return img
}

func SavePNG(img image.Image, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	err = png.Encode(f, img)
	if err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
