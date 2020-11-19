package draw

import (
	"image"
	"image/color"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d2/shape/boxmodel"
	"github.com/adamcolton/geom/iter"
	"github.com/fogleman/gg"
)

// Ctx is meant to represent *gg.Context.
type Ctx interface {
	Stroke()
	Fill()
	DrawLine(x1, y1, x2, y2 float64)
	DrawRectangle(x, y, w, h float64)
	SetLineWidth(lineWidth float64)
	DrawCircle(x, y, r float64)
	SavePNG(string) error
	SetRGB(r, g, b float64)
	Image() image.Image
}

// Context passes in the Ctx and FloatRange to simplify calling the draw
// functions.
type Context struct {
	Ctx
	iter.FloatRange
}

func (c *Context) Pt1(pt1 d2.Pt1) {
	Pt1(c, pt1, c.FloatRange)
}

func (c *Context) BoxModel(model boxmodel.BoxModel) {
	BoxModel(c, model)
}

func (c *Context) Arrow(pt d2.Pt, v d2.V) {
	Arrow(c, pt, v)
}

func (c *Context) V1(pt1 d2.Pt1, r iter.FloatRange, scale float64) {
	V1(c, pt1, r, scale)
}

func (c *Context) Pt2(pt2 d2.Pt2, r iter.FloatRange) {
	Pt2(c, pt2, r, r)
}

func (c *Context) OnPt1(pt1 d2.Pt1, ts []float64, radius float64) {
	OnPt1(c, pt1, ts, radius)
}

func (c *Context) Circle(pt d2.Pt, radius float64) {
	c.Ctx.DrawCircle(pt.X, pt.Y, radius)
	c.Ctx.Stroke()
}

// ContextGenerator is a helper for generating Contexts
type ContextGenerator struct {
	Size       grid.Pt
	Clear, Set color.Color
}

// Generate a Context of the defined size, set the background color and set the
// active color.
func (g *ContextGenerator) Generate() *Context {
	return g.GenerateForCtx(gg.NewContext(g.Size.X, g.Size.Y))
}

// GenerateForImage takes in an image and returns a Context. If img is nil,
// a call to Generate is returned.
func (g *ContextGenerator) GenerateForImage(img image.Image) *Context {
	if img == nil {
		return g.Generate()
	}
	return g.GenerateForCtx(gg.NewContextForImage(img))
}

// GenerateForCtx clears the gg.Context and uses it as Ctx in the returned
// Context.
func (g *ContextGenerator) GenerateForCtx(ctx *gg.Context) *Context {
	ctx.SetColor(g.Clear)
	ctx.Clear()
	ctx.SetColor(g.Set)
	return &Context{
		Ctx:        ctx,
		FloatRange: iter.Include(1, 0.0025),
	}
}
