package main

import (
	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render"
	"github.com/adamcolton/geom/d3/render/ffmpeg"
	"github.com/adamcolton/geom/d3/solid/mesh"
	"github.com/adamcolton/geom/examples/ggctx"
)

func main() {
	scale := 250.0
	c := setupCamera(scale)
	m := getMesh()
	colors := []*[3]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 0},
		{1, 0, 1},
		{0, 1, 1},

		{0.5, 0, 0},
		{0, 0.5, 0},
		{0, 0, 0.5},
	}

	size := int(scale * 2)
	out := &ffmpeg.Proc{
		Framerate:          15,
		Name:               "rendermesh",
		ConstantRateFactor: 25,
		Width:              size,
		Height:             size,
	}
	if err := out.Start(); err != nil {
		panic(err)
	}

	post := d3.Translate(d3.V{1, 1, 0}).T().T(d3.Scale(d3.V{scale, scale, -1}).T())

	step := angle.Rot(0.05)
	for rot := angle.Rot(0); rot < angle.Rot(1); rot += step {
		r := d3.Rotation{rot, d3.XY}.T()
		t := r.T(d3.Translate(d3.V{0, 0, -5}).T()).T(c.T()).T(post)
		m1 := m.T(t)

		ctx := ggctx.New(size, size)
		render.Wireframe(ctx, m1, colors[0], colors[8])
		out.AddFrame(ctx)
	}
	for rot := angle.Rot(0); rot < angle.Rot(1); rot += step {
		r := d3.Rotation{rot, d3.XY}.T()
		t := r.T(d3.Translate(d3.V{0, 0, -5}).T()).T(c.T()).T(post)
		m1 := m.T(t)
		render.RoundXY(m1)

		ctx := ggctx.New(size, size)
		render.Solid(ctx, m1, colors)
		out.AddFrame(ctx)
	}
	for rot := angle.Rot(0); rot < angle.Rot(1); rot += step {
		r := d3.Rotation{rot, d3.XZ}.T()
		t := r.T(d3.Translate(d3.V{0, 0, -5}).T()).T(c.T()).T(post)
		m1 := m.T(t)
		render.RoundXY(m1)

		ctx := ggctx.New(size, size)
		render.Wireframe(ctx, m1, colors[0], colors[8])
		out.AddFrame(ctx)
	}
	for rot := angle.Rot(0); rot < angle.Rot(1); rot += step {
		r := d3.Rotation{rot, d3.XZ}.T()
		t := r.T(d3.Translate(d3.V{0, 0, -5}).T()).T(c.T()).T(post)
		m1 := m.T(t)
		render.RoundXY(m1)

		ctx := ggctx.New(size, size)
		render.Solid(ctx, m1, colors)
		out.AddFrame(ctx)
	}

	out.Close()
}

func getMesh() mesh.TriangleMesh {
	f := []d3.Pt{
		{0, 2, 0},
		{1.5, 3.5, 0},
		{3, 2, 0},
		{2, 2, 0},
		{2, 0, 0},
		{1, 0, 0},
		{1, 2, 0},
	}
	f = d3.Translate(d3.V{-1.5, 0, 0}).T().Pts(f)

	m := mesh.NewExtrusion(f).
		Extrude(d3.Translate(d3.V{0, 0, -1}).T()).
		Close()

	tm, err := m.TriangleMesh()
	if err != nil {
		panic(err)
	}
	return tm
}

func setupCamera(scale float64) render.Camera {
	return render.Camera{
		Pt:    d3.Pt{0, 0, 0},
		Q:     d3.Q{1, 0, 0, 0},
		Near:  0.1,
		Far:   10,
		Angle: 3.1415 / 2.0,
	}
}
