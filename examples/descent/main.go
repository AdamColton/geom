package main

import (
	"image/color"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/calc/descent"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/draw"
	"github.com/adamcolton/geom/d2/grid"
)

var size = grid.Pt{500, 500}

type Setter interface {
	Set(x, y int, c color.Color)
}

func main() {
	gen := draw.ContextGenerator{
		Size:  size,
		Clear: draw.Color(1, 1, 1),
		Set:   draw.Color(1, 0, 0),
	}

	draw.Call(gen.Generate, Descent)
}

var (
	b = bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	l  = line.New(d2.Pt{0, 200}, d2.Pt{500, 300})
	fn = func(t0, t1 float64) float64 {
		return b.Pt1(t0).Distance(l.Pt1(t1))
	}
)

func Descent(ctx *draw.Context) {
	s := grid.Scale{
		X: 1.0 / 500.0,
		Y: 1.0 / 500.0,
	}
	si := s.Inv()
	// x*t+dx = f
	// inf --> 1/xf-dx

	// start at a point corresponding to pixel (10,10)
	startingPoints := []struct {
		grid.Pt
		x, y float64
	}{
		{
			Pt: grid.Pt{10, 10},
		}, {
			Pt: grid.Pt{365, 365},
		}, {
			Pt: grid.Pt{363, 363},
		}, {
			Pt: grid.Pt{20, 480},
		},
	}
	for i := range startingPoints {
		startingPoints[i].x, startingPoints[i].y = s.T(startingPoints[i].Pt)
	}

	drawfn(s, ctx)

	solver := &descent.Solver{
		Ln: 2,
		Fn: func(v []float64) float64 {
			return fn(v[0], v[1])
		},
	}
	solver.SetDFn()
	// x*t + dx
	// x1*500/10 = x2
	sf := 500.0 / 10.0
	s.X *= sf
	s.Y *= sf
	drawArrows(solver, s, ctx)

	for _, sp := range startingPoints {
		drawDescent(sp.x, sp.y, solver, s, si, ctx)
	}
}

func drawfn(s grid.Scale, ctx *draw.Context) {
	mm := cmpr.NewMinMax()
	img := ctx.Image().(Setter)
	for i := range size.Iter().Chan() {
		x, y := s.T(i)
		z := fn(x, y)
		mm.Update(z)
		c := draw.Color(z/582, 0, 0)
		img.Set(i.X, i.Y, c)
	}
	//fmt.Println(mm)
}

func drawArrows(solver *descent.Solver, s grid.Scale, ctx *draw.Context) {
	mm := cmpr.NewMinMax()
	ctx.SetLineWidth(3)
	ctx.SetRGB(0, 0, 1)
	buf := make([]float64, 2)
	x := make([]float64, 2)
	for i := range (grid.Pt{10, 10}).Iter().Chan() {
		x[0], x[1] = s.T(i)
		vs := solver.Derivative(x, buf)
		v := (d2.V{vs[0], vs[1]}).Multiply(0.06)
		mm.Update(v.Mag())
		ctx.Arrow(i.D2().Pt().Multiply(50), v)
	}
	//fmt.Println(mm)
}

func drawDescent(px, py float64, solver *descent.Solver, s, si grid.Scale, ctx *draw.Context) {
	p := []float64{px, py}
	ctx.SetLineWidth(3)
	steps := 20
	buf1 := make([]float64, 2)
	buf2 := make([]float64, 2)
	for i := 0; i < steps; i++ {
		solver.Step(p, buf1, buf2)
		x0, y0 := si.F(px, py)
		x1, y1 := si.F(p[0], p[1])
		ctx.SetRGB(float64(i%2), 0.5+float64(i)/float64(2*steps), 0)
		ctx.DrawLine(x0, y0, x1, y1)
		ctx.Stroke()
		px, py = p[0], p[1]
	}
}
