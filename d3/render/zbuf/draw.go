package zbuf

import (
	"image"

	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/shape/triangle"
)

type ZBuffer struct {
	w, h int
	buf  []bufEntry
	set  []bool
}

type bufEntry struct {
	*RenderMesh
	PolygonIdx, TriangleIdx int
	barycentric.B
	Z float64
}

func newZbuf(w, h int) ZBuffer {
	size := w * h
	return ZBuffer{
		w:   w,
		h:   h,
		buf: make([]bufEntry, size),
		set: make([]bool, size),
	}
}

func (z ZBuffer) Insert(pt d3.Pt, b barycentric.B, pIdx, tIdx int, rm *RenderMesh) {
	if pt.X < 0 || pt.X >= float64(z.w) || pt.Y < 0 || pt.Y >= float64(z.h) || pt.Z < -1 || pt.Z > 1 {
		return
	}
	idx := getIdx(z.w, &pt)
	if idx < 0 || idx >= len(z.buf) || (z.set[idx] && z.buf[idx].Z < pt.Z) {
		return
	}
	buf := &z.buf[idx]
	buf.RenderMesh = rm
	buf.B = b
	buf.Z = pt.Z
	buf.PolygonIdx = pIdx
	buf.TriangleIdx = tIdx
	z.set[idx] = true
}

func (z ZBuffer) Reset() {
	for i := range z.set {
		z.set[i] = false
	}
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
			bi, bt := Scan(triangle.Triangle(t), 0.9)
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

func (buf ZBuffer) Draw(img *image.RGBA) {
	for idx, be := range buf.buf {
		if !buf.set[idx] {
			continue
		}
		x, y := idx%buf.w, idx/buf.w
		c := be.Shader(&Context{
			B:           be.B,
			RenderMesh:  be.RenderMesh,
			PolygonIdx:  be.PolygonIdx,
			TriangleIdx: be.TriangleIdx,
		})
		img.SetRGBA(x, buf.h-y-1, *c)
	}
}
