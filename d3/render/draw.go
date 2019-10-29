package render

import (
	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/fogleman/gg"
)

type ZBuffer struct {
	w, h int
	buf  []bufEntry
}

type bufEntry struct {
	*RenderMesh
	PolygonIdx, TriangleIdx int
	barycentric.B
	Z   float64
	Set bool
}

func newZbuf(w, h int) ZBuffer {
	return ZBuffer{
		w:   w,
		h:   h,
		buf: make([]bufEntry, w*h),
	}
}

func (z ZBuffer) Insert(pt d3.Pt, b barycentric.B, pIdx, tIdx int, rm *RenderMesh) {
	if pt.X < 0 || pt.X >= float64(z.w) || pt.Y < 0 || pt.Y >= float64(z.h) || pt.Z < -1 || pt.Z > 1 {
		return
	}
	idx := getIdx(z.w, &pt)
	if idx < 0 || idx >= len(z.buf) || (z.buf[idx].Set && z.buf[idx].Z < pt.Z) {
		return
	}
	z.buf[idx].RenderMesh = rm
	z.buf[idx].B = b
	z.buf[idx].Z = pt.Z
	z.buf[idx].PolygonIdx = pIdx
	z.buf[idx].TriangleIdx = tIdx
	z.buf[idx].Set = true
}

func getIdx(w int, pt *d3.Pt) int {
	return w*int(pt.Y) + int(pt.X)
}

func New(w, h int) ZBuffer {
	return newZbuf(w, h)
}

func (buf ZBuffer) Add(rm *RenderMesh) {
	w64 := float64(buf.w)
	h64 := float64(buf.h)
	for pIdx, p := range rm.Original.Polygons {
		for tIdx, ptIdxs := range p {
			t := [3]d3.Pt{
				rm.Camera[ptIdxs[0]],
				rm.Camera[ptIdxs[1]],
				rm.Camera[ptIdxs[2]],
			}
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
				buf.Insert(bt.PtB(b), b, pIdx, tIdx, rm)
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

func (buf ZBuffer) Draw(ctx *gg.Context) {
	for idx, be := range buf.buf {
		if !be.Set {
			continue
		}
		x, y := idx%buf.w, idx/buf.w
		c := be.Shader(&Context{
			B:           be.B,
			RenderMesh:  be.RenderMesh,
			PolygonIdx:  be.PolygonIdx,
			TriangleIdx: be.TriangleIdx,
		})
		ctx.SetRGB(c[0], c[1], c[2])
		ctx.SetPixel(x, buf.h-y-1)
		ctx.Stroke()
	}
}
