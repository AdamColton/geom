package main

import (
	"fmt"
	"math/rand"

	"github.com/adamcolton/geom/d3/shape/plane"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	d2poly "github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render"
	"github.com/adamcolton/geom/d3/render/ffmpeg"
	triangle3 "github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/adamcolton/geom/d3/solid/mesh"
	"github.com/adamcolton/geom/examples/ggctx"
)

func main() {
	scale := 512.0
	c := setupCamera(scale)
	m := getMesh()
	//es, _ := m.Edges()

	size := int(scale * 2)
	stars := defineStarField()

	out := &ffmpeg.Proc{
		Framerate:          15,
		Name:               "stars",
		ConstantRateFactor: 25,
		Width:              size,
		Height:             size,
	}
	if err := out.Start(); err != nil {
		panic(err)
	}

	post := d3.Translate(d3.V{1, 1, 0}).T().T(d3.Scale(d3.V{scale, scale, -1}).T())

	for frame := 0; frame < 1000; frame++ {
		c.Pt.Z = float64(frame) * (-150.0 / 1000.0)
		ct := c.T().T(post)
		buf := render.New(size, size)
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
		ctx := ggctx.New(size, size)
		buf.Draw(ctx)
		out.AddFrame(ctx)
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

func setupCamera(scale float64) render.Camera {
	return render.Camera{
		Pt:    d3.Pt{0, 0, 0},
		Q:     d3.Q{1, 0, 0, 0},
		Near:  0.1,
		Far:   200,
		Angle: 3.1415 / 2.0,
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

func starShader(ctx *render.Context) *[3]float64 {
	if ctx.B.U < 0.03 || ctx.B.V < 0.03 || ctx.B.U+ctx.B.V > 0.97 {
		return &([3]float64{0, 0, 0})
	}
	tIdxs := ctx.Original.Polygons[ctx.PolygonIdx][ctx.TriangleIdx]
	n := (&triangle3.Triangle{
		ctx.Space[tIdxs[0]],
		ctx.Space[tIdxs[1]],
		ctx.Space[tIdxs[2]],
	}).Normal().Normal()
	r := n.X*0.25 + 0.75
	g := n.Y*0.25 + 0.75

	return &([3]float64{r, g, 0})

}
