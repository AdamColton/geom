package main

import (
	"fmt"
	"image/color"

	"github.com/adamcolton/geom/d3/render/zbuf"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d3"
	triangle3 "github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/adamcolton/geom/d3/solid/cc"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

const (
	frames     = 50
	width      = 1920
	imageScale = 1.25
	maxDivs    = 5
	framerate  = 25
)

var cr = string([]byte{13})

func main() {
	w := width
	if w%2 == 1 {
		w++
	}
	h := (w * 9) / 16
	if h%2 == 1 {
		h++
	}

	scene := &zbuf.Scene{
		W:    w,
		H:    h,
		A:    angle.Deg(45),
		Near: 0.1, Far: 200,
		Framerate:          framerate,
		Name:               "cc",
		ConstantRateFactor: 25,
		Background:         color.RGBA{100, 100, 100, 255},
		ImageScale:         imageScale,
	}

	cPt := d3.Pt{}
	q := d3.Q{1, 0, 0, 0}

	meshFacts := []func(int) mesh.TriangleMesh{
		getCube,
		getHex,
	}

	var shader zbuf.Shader
	for mfIdx, mf := range meshFacts {
		for divs := 0; divs <= maxDivs; divs++ {
			m := mf(divs)
			if divs == maxDivs {
				shader = smoothShader
			} else {
				shader = ccShader
			}
			tr := d3.Translate(d3.V{0, 0, 10}).Pair()
			for frame := 0; frame < frames; frame++ {
				d := float64(frame) / float64(frames)
				rot := d3.Rotation{angle.Rot(d), d3.YZ}.T()
				f := scene.NewFrame(cPt, q, 1)
				t := tr[0].T(rot).T(tr[1])
				f.AddMesh(&m, shader, t)
				f.Render()
				fmt.Print(cr, mfIdx, divs, " YZ Frame ", frame, "         ")
			}
			for frame := 0; frame < frames; frame++ {
				d := float64(frame) / float64(frames)
				rot := d3.Rotation{angle.Rot(d), d3.XZ}.T()
				f := scene.NewFrame(cPt, q, 1)
				t := tr[0].T(rot).T(tr[1])
				f.AddMesh(&m, shader, t)
				f.Render()
				fmt.Print(cr, mfIdx, divs, " XZ Frame ", frame, "         ")
			}
		}
	}
	scene.Done()
}

func getHex(divs int) mesh.TriangleMesh {
	f := []d3.Pt{
		{0, 0, 0},
		{0.5, -0.25, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0.5, 1.25, 0},
		{0, 1, 0},
	}

	m := mesh.NewExtrusion(f).
		EdgeExtrude(d3.Scale(d3.V{2, 2, 1}).T()).
		Extrude(
			d3.Translate(d3.V{0, 0, 1}).T(),
			d3.Translate(d3.V{0, 0, 2}).T(),
			d3.Translate(d3.V{0, 0, 1}).T(),
		).
		EdgeMerge(d3.Scale(d3.V{0.5, 0.5, 1}).T()).
		Close()

	m = m.T(d3.Translate(d3.V{-.5, -.5, -12}).T())
	m = cc.Subdivide(m, divs)

	tm, err := m.TriangleMesh()
	if err != nil {
		panic(err)
	}
	return tm
}

func getCube(divs int) mesh.TriangleMesh {
	f := []d3.Pt{
		{-1, -1, 0},
		{1, -1, 0},
		{1, 1, 0},
		{-1, 1, 0},
	}

	m := mesh.NewExtrusion(f).
		Extrude(
			d3.Translate(d3.V{0, 0, 2}).T(),
		).
		Close()

	m = m.T(d3.Translate(d3.V{0, 0, -11}).T())
	m = cc.Subdivide(m, divs)

	tm, err := m.TriangleMesh()
	if err != nil {
		panic(err)
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
	}.
		T().
		T(
			d3.Translate(s.V).T(),
		)
}

var black = color.RGBA{0, 0, 0, 255}

func ccShader(ctx *zbuf.Context) *color.RGBA {
	if ctx.B.U < 0.03 || ctx.B.V < 0.03 || ctx.B.U+ctx.B.V > 0.97 {
		return &black
	}
	tIdxs := ctx.Original.Polygons[ctx.PolygonIdx][ctx.TriangleIdx]
	n := (&triangle3.Triangle{
		ctx.Space[tIdxs[0]],
		ctx.Space[tIdxs[1]],
		ctx.Space[tIdxs[2]],
	}).Normal().Abs().Normal()
	r := (n.X*0.25 + 0.75) * 255
	g := (n.Y*0.25 + 0.75) * 255
	b := (n.Z*0.25 + 0.75) * 255

	return &(color.RGBA{uint8(r), uint8(g), uint8(b), 255})
}

var light = d3.V{0.1, 0.2, 1}.Normal()

func smoothShader(ctx *zbuf.Context) *color.RGBA {
	tIdxs := ctx.Original.Polygons[ctx.PolygonIdx][ctx.TriangleIdx]
	on := (&triangle3.Triangle{
		ctx.Original.Pts[tIdxs[0]],
		ctx.Original.Pts[tIdxs[1]],
		ctx.Original.Pts[tIdxs[2]],
	}).Normal()
	bright := light.Dot(on) * 0.5
	no := on.Abs().Normal()
	ns := (&triangle3.Triangle{
		ctx.Space[tIdxs[0]],
		ctx.Space[tIdxs[1]],
		ctx.Space[tIdxs[2]],
	}).Normal().Abs().Normal()
	r := (ns.X*0.25 + no.X*0.25 + bright) * 255
	g := (ns.Y*0.25 + no.Y*0.25 + bright) * 255
	b := (ns.Z*0.25 + no.Z*0.25 + bright) * 255

	return &(color.RGBA{uint8(r), uint8(g), uint8(b), 255})
}
