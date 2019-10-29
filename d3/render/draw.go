package render

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/adamcolton/geom/d3/solid"
	"github.com/adamcolton/geom/d3/solid/mesh"
	"github.com/fogleman/gg"
)

func Wireframe(ctx *gg.Context, m mesh.TriangleMesh, internalColor, edgeColor *[3]float64) {
	h := float64(ctx.Height() - 1)
	s, d := m.Edges()
	ctx.SetLineWidth(1)
	ctx.SetRGB(internalColor[0], internalColor[1], internalColor[2])
	for _, e := range d {
		p0, p1 := m.Pts[e[0]], m.Pts[e[1]]
		ctx.DrawLine(p0.X, h-p0.Y, p1.X, h-p1.Y)
	}
	ctx.Stroke()
	ctx.SetLineWidth(2)
	ctx.SetRGB(edgeColor[0], edgeColor[1], edgeColor[2])
	for _, e := range s {
		p0, p1 := m.Pts[e[0]], m.Pts[e[1]]
		ctx.DrawLine(p0.X, h-p0.Y, p1.X, h-p1.Y)
	}
	ctx.Stroke()
}

func Solid(ctx *gg.Context, m mesh.TriangleMesh, colors []*[3]float64) {
	buf := newZbuf(ctx.Width(), ctx.Height())

	buf.Add(m, colors)
	buf.Draw(ctx)

}

type ZBuffer struct {
	w, h int
	buf  []bufEntry
}

type bufEntry struct {
	Color *[3]float64
	Z     float64
	Set   bool
}

func newZbuf(w, h int) ZBuffer {
	return ZBuffer{
		w:   w,
		h:   h,
		buf: make([]bufEntry, w*h),
	}
}

func (z ZBuffer) Insert(pt d3.Pt, color *[3]float64) {
	if pt.X < 0 || pt.X >= float64(z.w) || pt.Y < 0 || pt.Y >= float64(z.h) || pt.Z < -1 || pt.Z > 1 {
		return
	}
	idx := getIdx(z.w, &pt)
	if idx < 0 || idx >= len(z.buf) || (z.buf[idx].Set && z.buf[idx].Z < pt.Z) {
		return
	}
	z.buf[idx].Color = color
	z.buf[idx].Z = pt.Z
	z.buf[idx].Set = true
}

func getIdx(w int, pt *d3.Pt) int {
	return w*int(pt.Y) + int(pt.X)
}

func New(w, h int) ZBuffer {
	return newZbuf(w, h)
}

func (buf ZBuffer) Add(m mesh.TriangleMesh, colors []*[3]float64) {
	w64 := float64(buf.w)
	h64 := float64(buf.h)
	for i := range m.Polygons {
		p := m.Face(i)
		c := colors[i]
		for _, t := range p {
			ok := false
			for _, pt := range t {
				x := pt.X >= 0 && pt.X <= w64
				y := pt.Y >= 0 && pt.Y <= h64
				z := pt.Z <= 1 && pt.Z >= -1
				ok = x && y && z
				if ok {
					break
				}
			}
			if !ok {
				continue
			}
			bi, bt := Scan(triangle.Triangle(t), 0.75)
			for b, done := bi.Start(); !done; b, done = bi.Next() {
				buf.Insert(bt.PtB(b), c)
			}
		}
	}
}

const zfix = -0.0001

var cards = [5]d3.V{
	{0, 0, zfix},
	{1, 0, zfix},
	{0, 1, zfix},
	{-1, 0, zfix},
	{0, -1, zfix},
}

func (buf ZBuffer) Edge(es []solid.IdxEdge, m mesh.TriangleMesh, c *[3]float64) {
	w64 := float64(buf.w)
	h64 := float64(buf.h)
	for _, e := range es {
		ok := false
		for _, ptIdx := range e {
			pt := m.Pts[ptIdx]
			x := pt.X >= 0 && pt.X <= w64
			y := pt.Y >= 0 && pt.Y <= h64
			z := pt.Z <= 1 && pt.Z >= -1
			ok = x && y && z
			if ok {
				break
			}
		}
		if !ok {
			continue
		}
		l := line.New(m.Pts[e[0]], m.Pts[e[1]])
		d := 1.0 / (l.D.Mag() * 2)
		for i := 0.0; i < 1.0; i += d {
			pt := l.Pt1(i)
			for _, d := range cards {
				buf.Insert(pt.Add(d), c)
			}
		}
	}
}

func (buf ZBuffer) Draw(ctx *gg.Context) {
	for idx, be := range buf.buf {
		if !be.Set {
			continue
		}
		x, y := idx%buf.w, idx/buf.w
		ctx.SetRGB(be.Color[0], be.Color[1], be.Color[2])
		ctx.SetPixel(x, buf.h-y-1)
		ctx.Stroke()
	}
}
