package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/zbuf"
	"github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/adamcolton/geom/examples/ggctx"
	"github.com/adamcolton/geom/work"
)

// Demonstrates the triangle scan operation by scanning a triangle after several
// transformations.

func main() {
	Clear()
	rotate(triangle.Triangle{{-0.5, -0.5, 0}, {.5, -0.5, 0}, {-0.5, .5, 0}}, "1")
	rotate(triangle.Triangle{{-0.1, -0.1, 0}, {.5, -0.5, 0}, {-0.5, .5, 0}}, "2")
}

func rotate(t triangle.Triangle, id string) {
	scale := 300.0
	scaledOffset := d3.V{250, 250, 0}

	step := 0.05
	r := int(1.0 / step)
	work.RunRange(r, func(rangeIdx, coreIdx int) {
		rot := float64(rangeIdx) * step
		ctx := ggctx.New(500, 500)
		tr := d3.Rotation{
			Angle: angle.Rot(rot),
			Plane: d3.XY}.T()
		tr = tr.T(d3.ScaleF(scale).T())
		tr = tr.T(d3.Translate(scaledOffset).T())
		t := &triangle.Triangle{
			tr.Pt(t[0]),
			tr.Pt(t[1]),
			tr.Pt(t[2]),
		}
		bi, bt := zbuf.Scan(t, 1.0, 1.0, nil, nil)
		for b, done := bi.Start(); !done; b, done = bi.Next() {
			ctx.SetRGB(b.U, b.V, 0)
			pt := bt.PtB(b)
			ctx.DrawPoint(pt.X, pt.Y, 1)
			ctx.Stroke()
		}

		ctx.SetRGB(0, 1, 0)
		ctx.DrawLine(t[0].X, t[0].Y, t[1].X, t[1].Y)
		ctx.DrawLine(t[0].X, t[0].Y, t[2].X, t[2].Y)
		ctx.DrawLine(t[2].X, t[2].Y, t[1].X, t[1].Y)
		ctx.Stroke()
		ctx.SavePNG(fmt.Sprintf("trianglescan-%s-%.2f.png", id, rot))
	})
}

func Clear() {
	files, err := filepath.Glob("*.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		os.Remove(f)
	}
}
