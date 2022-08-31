package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/draw"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/boxmodel"
	"github.com/adamcolton/geom/d2/shape/ellipse"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/fogleman/gg"
)

const (
	profile = true
)

var (
	triangles = shape.Subtract{
		shape.Union{
			shape.Subtract{
				&triangle.Triangle{
					{20, 20}, {480, 40}, {250, 470},
				},
				&triangle.Triangle{
					{200, 200}, {400, 400}, {250, 50},
				},
			},
			&triangle.Triangle{
				{20, 200}, {50, 30}, {200, 200},
			},
		},
		&triangle.Triangle{
			{100, 100}, {100, 150}, {150, 100},
		},
	}
	ell       = ellipse.New(d2.Pt{100, 350}, d2.Pt{400, 110}, 170)
	intersect = shape.Subtract{
		shape.Intersection{
			ellipse.NewCircle(d2.Pt{250, 250}, 230),
			&triangle.Triangle{
				{100, 250}, {490, 100}, {490, 400},
			},
		},
		ellipse.NewCircle(d2.Pt{350, 250}, 40),
	}
	compressed = boxmodel.NewCompressor()
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

	clear()
	Compress()

	gen := draw.ContextGenerator{
		Size:  grid.Pt{500, 500},
		Clear: draw.Color(1, 1, 1),
		Set:   draw.Color(1, 0, 0),
	}

	draw.Call(gen.Generate,
		Triangles,
		Ellipse,
		Intersection,
		Compressed0,
	)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func clear() {
	files, err := filepath.Glob("*.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		os.Remove(f)
	}
}

func drawModel(ctx *draw.Context, bm boxmodel.BoxModel) {
	ctx.SetRGB(0, 0.5, 0)
	for c, b, done := bm.OutsideCursor(); !done; b, done = c.Next() {
		d := b.V()
		ctx.DrawRectangle(b[0].X, b[0].Y, d.X, d.Y)
	}
	ctx.Stroke()

	ctx.SetRGB(0, 0, 1)
	for c, b, done := bm.InsideCursor(); !done; b, done = c.Next() {
		d := b.V()
		ctx.DrawRectangle(b[0].X, b[0].Y, d.X, d.Y)
	}
	ctx.Stroke()

	ctx.SetRGB(1, 0, 0)
	for c, b, done := bm.PerimeterCursor(); !done; b, done = c.Next() {
		d := b.V()
		ctx.DrawRectangle(b[0].X, b[0].Y, d.X, d.Y)
	}
	ctx.Stroke()

	size := ctx.Image().Bounds().Max
	ctx.SetRGB(0, 0, 0)
	l := line.New(d2.Pt{0.4 * float64(size.X), 0}, d2.Pt{0.6 * float64(size.X), float64(size.Y)})
	ctx.CurvePts(l, bm.LineIntersections(l, nil))
	ctx.Pt1(l)
}

func solid(bm boxmodel.BoxModel, name string) {
	ctx := gg.NewContext(500, 500)
	ctx.SetRGB(1, 1, 1)
	ctx.DrawRectangle(0, 0, 500, 500)
	ctx.Fill()
	ctx.Stroke()

	ctx.SetRGB(0, 0, 1)
	for c, b, ok := bm.InsideCursor(); ok; b, ok = c.Next() {
		d := b.V()
		ctx.DrawRectangle(b[0].X, b[0].Y, d.X+.5, d.Y+.5)
		ctx.Fill()
	}

	ctx.SavePNG(name)
}

func Triangles(ctx *draw.Context) {
	drawModel(ctx, boxmodel.New(triangles, 8))
}

func Ellipse(ctx *draw.Context) {
	drawModel(ctx, boxmodel.New(ell, 8))
}
func Intersection(ctx *draw.Context) {
	drawModel(ctx, boxmodel.New(intersect, 8))
}

func Compress() {
	_, err := compressed.Add("triangles", boxmodel.New(triangles, 12))
	checkErr(err)

	_, err = compressed.Add("ellipse", boxmodel.New(ell, 12))
	checkErr(err)

	_, err = compressed.Add("intersection", boxmodel.New(intersect, 12))
	checkErr(err)
}

func Compressed0(ctx *draw.Context) {
	drawModel(ctx, compressed.Get("triangles"))
}

func Compressed1(ctx *draw.Context) {
	drawModel(ctx, compressed.Get("ellipse"))
}

func Compressed2(ctx *draw.Context) {
	drawModel(ctx, compressed.Get("intersection"))
}
