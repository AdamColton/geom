package main

import (
	"fmt"
	"image/color"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/calc/descent"
	"github.com/adamcolton/geom/d2"
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

func Descent(ctx *draw.Context) {
	fn := func(x, y float64) float64 {
		x -= 3
		y -= 4
		x *= x
		y *= y
		return x + y
	}

	sf := 0.1
	s := grid.Scale{
		X:  sf,
		DX: -250 * sf,
		Y:  sf,
		DY: -250 * sf,
	}
	si := s.Inv()
	// x*t+dx = f
	// inf --> 1/xf-dx

	px, py := s.T(grid.Pt{10, 10})

	mm := cmpr.NewMinMax()
	img := ctx.Image().(Setter)
	for i := range size.Iter().Chan() {
		x, y := s.T(i)
		z := fn(x, y)
		mm.Update(z)
		c := draw.Color(z/1625, 0, 0)
		img.Set(i.X, i.Y, c)
	}

	solver := &descent.Solver{
		Ln: 2,
		Fn: func(v []float64) float64 {
			return fn(v[0], v[1])
		},
	}
	solver.SetDFn()
	// x*t + dx
	// x1*500/10 = x2
	sf = 500.0 / 10.0
	s.X *= sf
	s.Y *= sf

	ctx.SetLineWidth(3)
	ctx.SetRGB(0, 0, 1)
	for i := range (grid.Pt{10, 10}).Iter().Chan() {
		x, y := s.T(i)
		vs := solver.Derivative([]float64{x, y})
		v := d2.V{vs[0], vs[1]}
		ctx.Arrow(i.D2().Pt().Multiply(50), v)
	}

	fmt.Println(mm)

	p := []float64{px, py}
	fmt.Println(px, py)
	for i := 0; i < 20; i++ {
		solver.Step(p)
		x0, y0 := si.F(px, py)
		x1, y1 := si.F(p[0], p[1])
		ctx.SetRGB(0, 0.5+float64(i)/40, 0)
		ctx.DrawLine(x0, y0, x1, y1)
		ctx.Stroke()
		fmt.Println(p)
		fmt.Println(x0, y0, x1, y1)
		px, py = p[0], p[1]
	}

}
