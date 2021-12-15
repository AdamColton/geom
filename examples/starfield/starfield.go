package main

import (
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime/pprof"

	"github.com/adamcolton/geom/d3/render/material"
	"github.com/adamcolton/geom/d3/render/raytrace"
	"github.com/adamcolton/geom/d3/render/scene"
	"github.com/adamcolton/geom/d3/render/zbuf"
	"github.com/adamcolton/geom/d3/shape/plane"
	"github.com/adamcolton/geom/ffmpeg"
	"github.com/nfnt/resize"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	d2poly "github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

// TODO
// * Look at profiling - speed up raytrace rendering
// * set up shaders to look good.
// * pass out into sceneframe.intersect - this is a lot of allocation and gc

// For this to run, ffmpeg must be installed

const (
	// Frames of video to render
	frames = 500
	// Number of stars to render
	stars = 50

	// width sets the size, the aspect ratio is always widescreen
	width = 1000

	// Set to between 1.0 and 2.0
	// 1.0 is low quality
	// 2.0 is high quality
	imageScale = 2.0

	// enable the profiler
	profile = false

	doRay  = false
	doZbuf = true
)

func main() {
	if profile {
		f, err := os.Create("profile.out")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer func() {
			pprof.StopCPUProfile()
			f.Close()
		}()
	}

	proc := ffmpeg.New("stars")

	s := &scene.Scene{
		CameraFactory: cameraFactory{
			w:      proc.Size.X,
			h:      proc.Size.Y,
			frames: frames,
		},
		FrameCount: frames,
	}

	starMesh := getStarMesh()
	for _, starTransform := range defineStarField() {
		s.AddMesh(starMesh, starTransform, starMaterial)
	}

	sf := s.Frame(0)
	ray := &raytrace.Frame{
		Frame: sf,
		RayFrame: &raytrace.RayFrame{
			Background: raytrace.NewMaterialWrapper(backgroundMaterial).RayShader,
			Depth:      3,
			RayMult:    2,
			ImageScale: 1.25,
		},
	}
	if doRay {
		img := ray.Image(nil)
		f, _ := os.Create("ray.png")
		png.Encode(f, resize.Resize(uint(ray.Camera.Size.X), 0, img, resize.Bilinear))
		f.Close()
	}

	zb := &zbuf.Scene{
		Scene:      s,
		ImageScale: imageScale,
		ZBufFrame: &zbuf.ZBufFrame{
			Near:       0.1,
			Far:        200,
			Background: color.RGBA{255, 255, 255, 255},
		},
	}
	if doZbuf {
		zimg, _ := zb.Frame(0, nil)
		f, _ := os.Create("zbuf.png")
		png.Encode(f, resize.Resize(uint(ray.Camera.Size.X), 0, zimg, resize.Bilinear))
		f.Close()
	}

	if doZbuf {
		proc.Name = "zbuf"
		proc.Framer(
			&zbuf.Scene{
				Scene:      s,
				ImageScale: imageScale,
				ZBufFrame: &zbuf.ZBufFrame{
					Near:       0.1,
					Far:        200,
					Background: color.RGBA{255, 255, 255, 255},
				},
			},
		)
	}

	if doRay {
		proc.Name = "ray"
		proc.Framer(
			&raytrace.Scene{
				Scene: s,
				RayFrame: &raytrace.RayFrame{
					ImageScale: imageScale,
					Depth:      3,
					RayMult:    2,
					Background: raytrace.NewMaterialWrapper(backgroundMaterial).RayShader,
				},
			},
		)
	}
}

type cameraFactory struct {
	w, h   int
	frames int
}

func (cf cameraFactory) Camera(frameIdx int) *scene.Camera {
	d := float64(frameIdx) / float64(frames)
	rot := angle.Rot(d)
	var q d3.Q
	q.A, q.D = rot.Sincos()

	return scene.NewCamera(d3.Pt{0, 0, d * -150.0}, angle.Deg(45)).
		SetSize(cf.w, cf.h).
		SetRot(q)
}

func getStarMesh() *mesh.TriangleMesh {
	var points = 5
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
		tm.Polygons = append(tm.Polygons,
			[][3]uint32{
				{i * 2, i*2 + 1, (up * 2)},
			},
			[][3]uint32{
				{i*2 + 1, i * 2, (up * 2) + 1},
			},
			[][3]uint32{
				{i*2 + 1, (i*2 + 2) % (up * 2), (up * 2)},
			},
			[][3]uint32{
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

var starMaterial = &material.Material{
	Specular:    angle.Deg(5),
	Diffuse:     angle.Deg(1),
	Color:       &material.Color{0.95, 0.95, 0},
	Border:      0.03,
	BorderColor: &material.Color{0, 0, 0},
}

var backgroundMaterial = material.Material{
	Specular:    angle.Deg(5),
	Diffuse:     angle.Deg(1),
	Color:       &material.Color{1, 1, 1},
	BorderColor: &material.Color{1, 1, 1},
}
