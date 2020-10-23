package main

import (
	"image/color"
	"math/rand"
	"os"
	"runtime/pprof"

	"github.com/adamcolton/geom/d3/render/ffmpeg"
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
	frames = 1000
	// Number of stars to render
	stars = 200

	// width sets the size, the aspect ratio is always widescreen
	width = 1000

	// Set to between 1.0 and 2.0
	// 1.0 is low quality
	// 2.0 is high quality
	imageScale = 1.5

	// enable the profiler
	profile = false
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
		defer f.Close()
	}

	proc := ffmpeg.NewWidescreen("stars", width)

	s := &scene.Scene{
		GetCamera: cameraFactory{
			w:      proc.Width,
			h:      proc.Height,
			frames: frames,
		},
	}

	starMesh := getStarMesh()
	for _, starTransform := range defineStarField() {
		s.AddMesh(starMesh, starTransform, starShader{})
	}

	proc.Framer(&zbuf.Framer{
		Count:      frames,
		Scene:      s,
		Near:       0.1,
		Far:        200,
		Background: color.RGBA{255, 255, 255, 255},
		ImageScale: 1.25,
	})

}

type cameraFactory struct {
	w, h   int
	frames int
}

func (cf cameraFactory) Camera(frameIdx int) *scene.Camera {
	d := float64(frameIdx) / float64(frames)
	c := &scene.Camera{
		Width:  cf.w,
		Height: cf.h,
		Angle:  angle.Deg(45),
	}
	c.Pt.Z = d * -150.0
	rot := angle.Rot(d)
	c.Q.A, c.Q.D = rot.Sincos()
	return c
}

func getStarMesh() *mesh.TriangleMesh {
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
			{(i*2 + 2) % (up * 2), i*2 + 1, (up * 2)},
		}, [][3]uint32{
			{(i*2 + 2) % (up * 2), i*2 + 1, (up * 2) + 1},
		})
	}

	return &tm
}

type star struct {
	d3.V
	angle.Rad
	speed  angle.Rad
	offset *d3.T
}

func (s *star) T(frame int) *d3.T {
	xz := d3.Rotation{
		s.Rad + s.speed*angle.Rad(float64(frame)),
		d3.XZ,
	}.T()
	t := d3.Translate(s.V).T()
	return xz.T(s.offset).T(t)

}

func defineStarField() []*star {
	out := make([]*star, stars)
	for i := range out {
		v2 := d2.Polar{
			M: rand.Float64()*10 + 3,
			A: angle.Rot(rand.Float64()),
		}.V()

		out[i] = &star{
			V:     d3.V{v2.X, v2.Y, rand.Float64() * -200},
			Rad:   angle.Rot(rand.Float64()),
			speed: angle.Rot(rand.Float64()*.05) + .05,
			offset: d3.Rotation{
				angle.Rot(rand.Float64()),
				d3.XY,
			}.T(),
		}
		if rand.Intn(2) == 0 {
			out[i].speed = -out[i].speed
		}
	}
	return out
}

type starShader struct{}

var black = color.RGBA{0, 0, 0, 255}

func (starShader) ZBufShader(ctx *zbuf.Context) *color.RGBA {
	if ctx.B.U < 0.03 || ctx.B.V < 0.03 || ctx.B.U+ctx.B.V > 0.97 {
		return &black
	}
	m := ctx.SceneFrame.Meshes[ctx.MeshIdx]
	tIdxs := m.Original.Polygons[ctx.PolygonIdx][ctx.TriangleIdx]
	n := (&triangle3.Triangle{
		m.Space[tIdxs[0]],
		m.Space[tIdxs[1]],
		m.Space[tIdxs[2]],
	}).Normal().Normal()

	r := ((angle.Rot(n.X).Cos() + 1) / 4) + 0.5
	g := ((angle.Rot(n.Y).Cos() + 1) / 4) + 0.5

	return &(color.RGBA{uint8(r * 255), uint8(g * 255), 0, 255})
}
