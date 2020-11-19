package draw

import (
	"math"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape/boxmodel"
	"github.com/adamcolton/geom/iter"
)

// Pt1 draws a curve that fulfills d2.Pt1.
func Pt1(ctx Ctx, pt1 d2.Pt1, r iter.FloatRange) {
	prev := pt1.Pt1(r.Start)
	r.Start += r.Step
	r.Each(func(t float64) {
		cur := pt1.Pt1(t)
		ctx.DrawLine(prev.X, prev.Y, cur.X, cur.Y)
		prev = cur
	})
	ctx.Stroke()
}

// BoxModel is rendered as a solid
func BoxModel(ctx Ctx, model boxmodel.BoxModel) {
	for c, b, done := model.InsideCursor(); !done; b, done = c.Next() {
		v := b.V()
		ctx.DrawRectangle(b[0].X, b[0].Y, v.X, v.Y)
	}
	ctx.Fill()
}

var (
	t45cw = d2.Chain{
		d2.Rotate(3.25 * math.Pi / 4.0),
		d2.Scale(d2.V{0.2, 0.2}),
	}.T()

	t45ccw = d2.Chain{
		d2.Rotate(-3.25 * math.Pi / 4.0),
		d2.Scale(d2.V{0.2, 0.2}),
	}.T()
)

func Arrow(ctx Ctx, pt d2.Pt, v d2.V) {
	ctx.SetLineWidth(v.Mag() * 0.04)
	end := pt.Add(v)
	ctx.DrawLine(pt.X, pt.Y, end.X, end.Y)
	va := t45cw.V(v)
	pta := end.Add(va)
	ctx.DrawLine(end.X, end.Y, pta.X, pta.Y)
	va = t45ccw.V(v)
	pta = end.Add(va)
	ctx.DrawLine(end.X, end.Y, pta.X, pta.Y)
	ctx.Stroke()
}

// V1 draws a curve with arrows indicating the derivative
func V1(ctx Ctx, pt1 d2.Pt1, r iter.FloatRange, scale float64) {
	v1 := d2.GetV1(pt1)
	r.Each(func(t float64) {
		pt := pt1.Pt1(t)
		v := v1.V1(t)
		v = v.Multiply(scale / math.Sqrt(v.Mag()))
		Arrow(ctx, pt, v)
	})
}

// Pt2 draws a d2.Pt2 as stripes
func Pt2(ctx Ctx, pt2 d2.Pt2, r0, r1 iter.FloatRange) {
	pt2c1 := d2.GetPt2c1(pt2)
	r0.Each(func(t0 float64) {
		Pt1(ctx, pt2c1.Pt2c1(t0), r1)
	})
}

// OnPt1 draws circles on the d2.Pt1 for each value in ts.
func OnPt1(ctx Ctx, pt1 d2.Pt1, ts []float64, radius float64) {
	for _, t := range ts {
		pt := pt1.Pt1(t)
		ctx.DrawCircle(pt.X, pt.Y, radius)
	}
	ctx.Stroke()
}
