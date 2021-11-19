package main

import (
	"image/color"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/calc/descent"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/d2/curve/distance"
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

	draw.Call(gen.Generate,
		Descent,
	)
}

func convertFn(fn distance.Fn) descent.Fn {
	return func(x []float64) float64 {
		return fn(x[0], x[1])
	}
}

func convertDFn(dfn distance.DFn) descent.DFn {
	return func(x, buf []float64) []float64 {
		var out []float64
		if len(buf) >= 2 {
			out = buf[:2]
		} else {
			out = make([]float64, 2)
		}
		out[0], out[1] = dfn(x[0], x[1])
		return out
	}
}

func Descent(ctx *draw.Context) {
	b := bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	l := line.New(d2.Pt{0, 200}, d2.Pt{500, 300})
	fn, dfn := distance.New2(b, l)

	d := &DescentPlotter{
		solver: &descent.Solver{
			Ln:  2,
			Fn:  convertFn(fn),
			DFn: convertDFn(dfn),
		},
		pmm: &cmpr.MinMax{
			Max: 338881,
		},
		vmm: &cmpr.MinMax{
			Max: 2246393,
		},
		scale: grid.Scale{
			X: 1.0 / 500.0,
			Y: 1.0 / 500.0,
		},
		arrowFactor: 500.0 / 10.0,
		startingPoints: []grid.Pt{
			{20, 20},
			{365, 365},
			{363, 363},
			{20, 480},
		},
	}

	d.drawfn(ctx)

	// x*t + dx
	// x1*500/10 = x2
	//sf := 500.0 / 10.0
	//s.X *= sf
	//s.Y *= sf
	d.drawArrows(ctx)

	d.drawDescent(ctx)
}

type DescentPlotter struct {
	solver         *descent.Solver
	pmm, vmm       *cmpr.MinMax
	scale          grid.Scale
	arrowFactor    float64
	startingPoints []grid.Pt
}

func (d *DescentPlotter) Draw(ctx *draw.Context) {
	d.drawfn(ctx)
}

func (d *DescentPlotter) drawfn(ctx *draw.Context) {
	mm := cmpr.NewMinMax()
	img := ctx.Image().(Setter)
	buf := make([]float64, 2)
	for i := range size.Iter().Chan() {
		buf[0], buf[1] = d.scale.T(i)
		z := d.solver.Fn(buf)
		mm.Update(z)
		c := draw.Color(z/d.pmm.Max, 0, 0)
		img.Set(i.X, i.Y, c)
	}
}

func (d *DescentPlotter) drawArrows(ctx *draw.Context) {
	if d.solver.DFn == nil {
		d.solver.SetDFn(nil)
	}

	s := grid.Scale{
		X: d.arrowFactor * d.scale.X,
		Y: d.arrowFactor * d.scale.Y,
	}

	mm := cmpr.NewMinMax()
	ctx.SetLineWidth(3)
	ctx.SetRGB(0, 0, 1)
	buf := make([]float64, 2)
	x := make([]float64, 2)
	for i := range (grid.Pt{10, 10}).Iter().Chan() {
		x[0], x[1] = s.T(i)
		vs := d.solver.DFn(x, buf)
		v := (d2.V{vs[0], vs[1]}).Multiply(200 / d.vmm.Max)
		mm.Update(v.Mag())
		ctx.Arrow(i.D2().Pt().Multiply(50), v)
	}
}

func (d *DescentPlotter) drawDescent(ctx *draw.Context) {
	if d.solver.DFn == nil {
		d.solver.SetDFn(nil)
	}

	si := d.scale.Inv()

	p := []float64{0, 0}
	for _, sp := range d.startingPoints {
		p[0], p[1] = d.scale.T(sp)
		px, py := p[0], p[1]
		ctx.SetLineWidth(3)
		steps := 20
		buf1 := make([]float64, 2)
		buf2 := make([]float64, 2)
		for i := 0; i < steps; i++ {
			d.solver.Step(p, buf1, buf2)
			x0, y0 := si.F(px, py)
			x1, y1 := si.F(p[0], p[1])
			ctx.SetRGB(float64(i%2), 0.5+float64(i)/float64(2*steps), 0)
			ctx.DrawLine(x0, y0, x1, y1)
			ctx.Stroke()
			px, py = p[0], p[1]
		}
	}

}
