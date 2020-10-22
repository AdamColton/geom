package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"runtime/pprof"

	"github.com/adamcolton/geom/d3/render/scene"
	"github.com/adamcolton/geom/d3/render/zbuf"
	"github.com/adamcolton/geom/d3/shape/plane"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	d2poly "github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d3"
	triangle3 "github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

// For this to run, ffmpeg must be installed

const (
	// Frames of video to render
	frames = 100
	// Number of stars to render
	stars = 200

	// width sets the size, the aspect ratio is always widescreen
	width = 500

	// Set to between 1.0 and 2.0
	// 1.0 is low quality
	// 2.0 is high quality
	imageScale = 1.5

	// enable the profiler
	profile = true
)

var cr = string([]byte{13})

func main() {
	if profile {
		f, err := os.Create("profile.out")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	m := getMesh()
	w := width
	if w%2 == 1 {
		w++
	}
	h := (w * 9) / 16
	if h%2 == 1 {
		h++
	}

	s := &zbuf.Scene{
		Camera: zbuf.Camera{
			Camera: scene.Camera{
				Width:  w,
				Height: h,
				Angle:  angle.Deg(45),
				Q:      d3.Q{1, 0, 0, 0},
			},
			Near: 0.1,
			Far:  200,
		},
		Framerate:          15,
		Name:               "stars",
		ConstantRateFactor: 25,
		Background:         color.RGBA{255, 255, 255, 255},
		ImageScale:         1.25,
	}

	stars := defineStarField()

	for frame := 0; frame < frames; frame++ {
		d := float64(frame) / float64(frames)
		s.Camera.Pt.Z = d * -150.0
		rot := angle.Rot(d)
		s.Camera.Q.A, s.Camera.Q.D = rot.Sincos()
		f := s.NewFrame(len(stars))
		for _, star := range stars {
			if star.Z+0.2 > s.Camera.Pt.Z {
				continue
			}
			f.AddMesh(&m, starShader, star.T(float64(frame)))
		}
		f.Render()
		fmt.Print(cr, "Frame ", frame, "         ")
	}
	s.Done()
}

func getMesh() mesh.TriangleMesh {
	var points = 8
	outer := d2poly.RegularPolygonRadius(d2.Pt{0, 0}, 2, angle.Rot(0.25), points)
	inner := d2poly.RegularPolygonRadius(d2.Pt{0, 0}, 1, angle.Rot(0.25+1.0/float64(points*2)), points)
	star2d := make([]d2.Pt, points*2)
	for i, oPt := range outer {
		star2d[2*i] = oPt
		star2d[2*i+1] = inner[i]
	}
	f := plane.New(d3.Pt{0, 0, 0}, d3.Pt{1, 0, 0}, d3.Pt{0, 1, 0}).ConvertMany(star2d)
	tm := mesh.TriangleMesh{
		Pts: append(f, d3.Pt{0, 0, 0.5}, d3.Pt{0, 0, -0.5}),
	}

	up := uint32(points)
	for i := uint32(0); i < up; i++ {
		tm.Polygons = append(tm.Polygons, [][3]uint32{
			{i * 2, i*2 + 1, (up * 2)},
		}, [][3]uint32{
			{i * 2, i*2 + 1, (up * 2) + 1},
		}, [][3]uint32{
			{i*2 + 1, (i*2 + 2) % (up * 2), (up * 2)},
		}, [][3]uint32{
			{i*2 + 1, (i*2 + 2) % (up * 2), (up * 2) + 1},
		})
	}

	return tm
}

type star struct {
	d3.V
	angle.Rad
	speed angle.Rad
}

func (s *star) T(frame float64) *d3.T {
	return d3.Rotation{
		s.Rad + s.speed*angle.Rad(frame),
		d3.XZ,
	}.T().
		T(d3.Translate(s.V).T())
}

func defineStarField() []star {
	out := make([]star, stars)
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

func starShader(ctx *zbuf.Context) *color.RGBA {
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
