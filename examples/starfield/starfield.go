package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"runtime/pprof"

	"github.com/adamcolton/geom/d3/shape/plane"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/grid"
	d2poly "github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render"
	"github.com/adamcolton/geom/d3/render/ffmpeg"
	triangle3 "github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

// Stop using gg
// Setup scene
// Setup a pipeline
// * Init scene frame
// * Scene transforms
// * Camera transforms
// * Render mesh to zbuf
// * merge zbufs
// * draw

func main() {
	f, err := os.Create("profile.out")
	if err != nil {
		panic(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	width := 546
	height := 1 + (width*9)/16
	c := setupCamera(width, height)
	m := getMesh()
	//es, _ := m.Edges()

	stars := defineStarField()

	out := &ffmpeg.Proc{
		Framerate:          15,
		Name:               "stars",
		ConstantRateFactor: 25,
		InputFormat:        "bmp",
		Width:              width,
		Height:             height,
	}
	if err := out.Start(); err != nil {
		panic(err)
	}

	buf := render.New(width, height)

	for frame := 0; frame < 1000; frame++ {
		c.Pt.Z = float64(frame) * (-150.0 / 1000.0)
		ct := c.T()
		buf.Reset()
		for _, s := range stars {
			if s.Z+0.2 > c.Z {
				continue
			}
			space := d3.Rotation{
				s.Rad + s.speed*angle.Rad(frame),
				d3.XZ,
			}.
				T().
				T(
					d3.Translate(s.V).T(),
				)
			rm := render.NewRenderMesh(&m, space, ct, starShader)
			buf.Add(rm)
			//buf.Edge(es, mt, &([3]float64{0, 0, 0}))
		}
		img := getImage(width, height)
		buf.Draw(img)
		out.AddPng(img)
		fmt.Println("Frame ", frame)
	}
	out.Close()
}

func getMesh() mesh.TriangleMesh {
	outer := d2poly.RegularPolygonRadius(d2.Pt{0, 0}, 2, angle.Rot(0.25), 5)
	inner := d2poly.RegularPolygonRadius(d2.Pt{0, 0}, 1, angle.Rot(0.35), 5)
	star2d := make([]d2.Pt, 10)
	for i, oPt := range outer {
		star2d[2*i] = oPt
		star2d[2*i+1] = inner[i]
	}
	f := plane.New(d3.Pt{0, 0, 0}, d3.Pt{1, 0, 0}, d3.Pt{0, 1, 0}).ConvertMany(star2d)
	tm := mesh.TriangleMesh{
		Pts: append(f, d3.Pt{0, 0, 0.5}, d3.Pt{0, 0, -0.5}),
	}
	for i := uint32(0); i < 5; i++ {
		tm.Polygons = append(tm.Polygons, [][3]uint32{
			{i * 2, i*2 + 1, 10},
		}, [][3]uint32{
			{i * 2, i*2 + 1, 11},
		}, [][3]uint32{
			{i*2 + 1, (i*2 + 2) % 10, 10},
		}, [][3]uint32{
			{i*2 + 1, (i*2 + 2) % 10, 11},
		})
	}

	return tm
}

func setupCamera(w, h int) render.Camera {
	return render.Camera{
		Pt:    d3.Pt{0, 0, 0},
		Q:     d3.Q{1, 0, 0, 0},
		Near:  0.1,
		Far:   200,
		Angle: 3.1415 / 2.0,
		W:     w,
		H:     h,
	}
}

type star struct {
	d3.V
	angle.Rad
	speed angle.Rad
}

func defineStarField() []star {
	out := make([]star, 100)
	for i := range out {
		v2 := d2.Polar{
			M: rand.Float64()*10 + 3,
			A: angle.Rot(rand.Float64()),
		}.V()

		out[i] = star{
			V:     d3.V{v2.X, v2.Y, rand.Float64() * -200},
			Rad:   angle.Rot(rand.Float64()),
			speed: angle.Rot(rand.Float64()*.05) + .05,
		}
		if rand.Intn(2) == 0 {
			out[i].speed = -out[i].speed
		}
	}
	return out
}

var black = color.RGBA{0, 0, 0, 255}

func starShader(ctx *render.Context) *color.RGBA {
	if ctx.B.U < 0.03 || ctx.B.V < 0.03 || ctx.B.U+ctx.B.V > 0.97 {
		return &black
	}
	tIdxs := ctx.Original.Polygons[ctx.PolygonIdx][ctx.TriangleIdx]
	n := (&triangle3.Triangle{
		ctx.Space[tIdxs[0]],
		ctx.Space[tIdxs[1]],
		ctx.Space[tIdxs[2]],
	}).Normal().Normal()
	r := (n.X*0.25 + 0.75) * 255
	g := (n.Y*0.25 + 0.75) * 255

	return &(color.RGBA{uint8(r), uint8(g), 0, 255})

}

var white = color.RGBA{255, 255, 255, 255}
var baseImg *image.RGBA

func getImage(width, height int) *image.RGBA {
	if baseImg == nil {
		baseImg = image.NewRGBA(image.Rect(0, 0, width, height))
		for iter, done := (grid.Pt{width, height}.Iter()).Start(); !done; done = iter.Next() {
			pt := iter.Pt()
			baseImg.SetRGBA(pt.X, pt.Y, white)
		}
	}
	img := &image.RGBA{
		Pix:    make([]uint8, len(baseImg.Pix)),
		Stride: baseImg.Stride,
		Rect:   baseImg.Rect,
	}
	copy(img.Pix, baseImg.Pix)
	return img
}
