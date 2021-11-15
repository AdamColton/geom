package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/d2/curve/ellipsearc"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/draw"
	"github.com/adamcolton/geom/d2/generate"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/boxmodel"
	"github.com/adamcolton/geom/d2/shape/ellipse"
	"github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/adamcolton/geom/iter"
	"github.com/fogleman/gg"
)

func main() {
	Clear()

	gen := draw.ContextGenerator{
		Size:  grid.Pt{500, 500},
		Clear: draw.Color(1, 1, 1),
		Set:   draw.Color(1, 0, 0),
	}

	draw.Call(gen.Generate,
		EllipseArc,
		BoxModel,
		Arrow,
		CurveWithArrows,
		EllipseFill,
		EllipseIntersection,
		EllipseContains,
		CircumscribeCircle,
		TriangleFill,
		ConcaveFill,
		BezierIntersection,
		Bezier,
		Blossom,
		BezierSegment,
		PtNorm,
		BezBezIntersection,
	)
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

func EllipseArc(ctx *draw.Context) {
	ctx.Pt1(ellipsearc.New(d2.Pt{100, 100}, d2.Pt{400, 400}, 100))
}

func BoxModel(ctx *draw.Context) {
	g := gg.NewRadialGradient(250, 250, 20, 250, 250, 250)
	g.AddColorStop(0, color.RGBA{255, 0, 0, 255})
	g.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	ctx.Ctx.(*gg.Context).SetFillStyle(g)

	ctx.BoxModel(boxmodel.New(shape.Subtract{
		shape.Union{
			ellipse.NewCircle(d2.Pt{250, 250}, 230),
			&triangle.Triangle{
				{100, 250}, {490, 100}, {490, 400},
			},
		},
		ellipse.NewCircle(d2.Pt{350, 250}, 40),
	}, 10))
}

func Arrow(ctx *draw.Context) {
	ctx.Arrow(d2.Pt{250, 250}, d2.V{100, 100})
}

var largeStep = iter.Include(1, 0.05)

func CurveWithArrows(ctx *draw.Context) {
	c := bezier.Bezier{{0, 0}, {160, 1000}, {320, -500}, {500, 500}}
	ctx.Pt1(c)
	ctx.SetRGB(0, 0, 1)
	ctx.V1(c, largeStep, 1.5)
}

func EllipseFill(ctx *draw.Context) {
	e := ellipse.New(d2.Pt{100, 100}, d2.Pt{400, 400}, 100)
	ctx.Pt1(e)
	ctx.SetRGB(0, 0, 1)
	ctx.Pt2(e, largeStep)
}

func EllipseIntersection(ctx *draw.Context) {
	e := ellipse.New(d2.Pt{100, 100}, d2.Pt{400, 400}, 100)
	ctx.Pt1(e)

	fn := func(p1, p2 d2.Pt) {
		l := line.New(p1, p2)
		ctx.Pt1(l)
		ctx.OnPt1(l, e.LineIntersections(l, nil), 3)
	}

	ctx.SetRGB(0, 0, 1)
	for f := range iter.FloatChan(0, 500, 20) {
		fn(d2.Pt{f, 0}, d2.Pt{f, 500})

	}

	ctx.SetRGB(0, 0.75, 0)
	for f := range iter.FloatChan(0, 180, 10) {
		l := line.Line{
			T0: d2.Pt{250, 250},
			D:  d2.Polar{250, angle.Deg(f)}.V(),
		}
		fn(l.Pt1(-1), l.Pt1(1))
	}

	ctx.SetRGB(0, 0.5, 1)
	for y := range iter.FloatChan(-150, 500, 30) {
		l := line.Line{
			T0: d2.Pt{0, y},
			D:  d2.V{500, 150},
		}
		fn(l.Pt1(0), l.Pt1(1))
	}
}

func EllipseContains(ctx *draw.Context) {
	e := ellipse.New(d2.Pt{100, 100}, d2.Pt{400, 400}, 100)
	ctx.SetRGB(0, 0, 1)
	ctx.Pt1(e)

	m, M := e.BoundingBox()
	d := M.Subtract(m)
	ctx.DrawRectangle(m.X, m.Y, d.X, d.Y)
	ctx.Stroke()

	size := d2.Pt{500, 500}
	for i := 0; i < 3000; i++ {
		pt := generate.PtIn(size)
		if e.Contains(pt) {
			ctx.SetRGB(0, 1, 0)
		} else {
			ctx.SetRGB(1, 0, 0)
		}
		ctx.DrawCircle(pt.X, pt.Y, 3)
		ctx.Stroke()
	}
}

func CircumscribeCircle(ctx *draw.Context) {
	t := triangle.Triangle{
		{100, 100},
		{400, 100},
		{100, 400},
	}
	ctx.Pt1(ellipse.CircumscribeCircle(t))
	ctx.Pt1(&t)
}

func TriangleFill(ctx *draw.Context) {
	t := triangle.Triangle{
		{100, 100},
		{400, 100},
		{100, 400},
	}
	ctx.Pt1(&t)
	ctx.SetRGB(0, 0, 1)
	ctx.Pt2(&t, largeStep)
}

func ConcaveFill(ctx *draw.Context) {
	p := polygon.NewConcavePolygon(polygon.Polygon{
		{100, 250},
		{400, 100},
		{250, 250},
		{400, 400},
	})
	ctx.Pt1(p)
	ctx.SetRGB(0, 0, 1)
	ctx.Pt2(p, largeStep)
}

func BezBezIntersection(ctx *draw.Context) {
	b1 := bezier.Bezier{
		{0, 0},
		{133, 1000},
		{266, -500},
		{450, 500},
	}
	b2 := bezier.Bezier{
		{0, 500},
		{300, -100},
		{500, 500},
	}
	ctx.Pt1(b1)
	ctx.Stroke()
	ctx.SetRGB(0, 0, 1)
	ctx.Pt1(b2)
	ctx.Stroke()
	ctx.SetRGB(0, 0.5, 0)
	is := b1.Intersections(b2)
	for _, t := range is {
		pt1 := b1.Pt1(t[0])
		pt2 := b2.Pt1(t[1])
		ctx.DrawCircle(pt1.X, pt1.Y, 10)
		ctx.DrawCircle(pt2.X, pt2.Y, 10)
	}
	ctx.Stroke()
}

func BezierIntersection(ctx *draw.Context) {
	b := bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	l := line.New(d2.Pt{0, 200}, d2.Pt{500, 300})
	ctx.Pt1(b)
	ctx.SetRGB(0, 0, 1)
	ctx.Pt1(l)
	ctx.SetRGB(0, 0.5, 0)
	ctx.OnPt1(l, b.LineIntersections(l, nil), 10)
}

func Bezier(ctx *draw.Context) {
	b := bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	ctx.Pt1(b)
}

func Blossom(ctx *draw.Context) {
	b := bezier.Bezier{
		{0, 0},
		{166, 500},
		{333, 0},
		{500, 500},
	}
	var out bezier.Bezier
	largeStep.Each(func(f1 float64) {
		largeStep.Each(func(f2 float64) {
			out = out[:0]
			largeStep.Each(func(f3 float64) {
				out = append(out, b.Blossom(f1, f2, f3))
			})
			ctx.SetRGB(f1, f2, 0.5)
			ctx.Pt1(out)
		})
	})
}

func BezierSegment(ctx *draw.Context) {
	b := bezier.Bezier{
		{0, 0},
		{166, 1000},
		{333, -500},
		{500, 500},
	}
	ctx.Pt1(b)
	ctx.SetRGB(0, 0, 1)
	ctx.Pt1(b.Segment(0.25, .75))
}

func PtNorm(ctx *draw.Context) {
	v := d2.V{250, 250}
	for i := 0; i < 2000; i++ {
		pt := generate.PtNorm().Multiply(75).Add(v)
		ctx.Circle(pt, 3)
	}
}
