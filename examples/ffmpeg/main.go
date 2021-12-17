package main

import (
	"image"
	"os"
	"runtime/pprof"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/d2/draw"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/ffmpeg"
)

var (
	profile = true
	size    = 500
	frames  = 200
	aspect  = grid.Widescreen
	curve   = []d2.Pt{
		{.1, .1},
		{0.32, 2},
		{0.64, -1},
		{.9, .9},
	}
)

func main() {
	if profile {
		f, err := os.Create("out.prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer func() {
			pprof.StopCPUProfile()
			f.Close()
		}()
	}

	s := ffmpeg.NewByAspect("test", size, aspect)
	f := &Framer{
		Settings: s,
		ContextGenerator: draw.ContextGenerator{
			Size:  s.Size,
			Clear: draw.Color(1, 1, 1),
			Set:   draw.Color(1, 0, 0),
		},
		FrameCount: frames,
		Path: bezier.Bezier(
			d2.Scale(s.Size.D2()).T().Slice(curve),
		),
		ArrowScale: 0.1,
	}

	err := f.Settings.Framer(f)
	if err != nil {
		panic(err)
	}
}

type Framer struct {
	*ffmpeg.Settings
	draw.ContextGenerator
	FrameCount int
	Path       d2.Pt1V1
	ArrowScale float64
}

func (f *Framer) Frame(idx int, img image.Image) (image.Image, error) {
	ctx := f.GenerateForImage(img)

	ctx.Pt1(f.Path)
	ctx.SetRGB(0, 0, 1)
	t := float64(idx) / float64(f.FrameCount)
	pt := f.Path.Pt1(t)
	v := f.Path.V1(t).Multiply(f.ArrowScale)
	ctx.Arrow(pt, v)

	return ctx.Image(), nil
}

func (f *Framer) Frames() int {
	return f.FrameCount
}
