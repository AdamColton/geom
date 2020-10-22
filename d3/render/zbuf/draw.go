package zbuf

import (
	"image"
	"image/color"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/shape/triangle"
)

type ZBuffer struct {
	w, h       int
	background *color.RGBA
	buf        []bufEntry
	set        []bool
	cpus       int
}

type bufEntry struct {
	*RenderMesh
	PolygonIdx, TriangleIdx int
	barycentric.B
	Z float64
}

func newZbuf(w, h int, background *color.RGBA) ZBuffer {
	size := w * h
	return ZBuffer{
		w:          w,
		h:          h,
		background: background,
		buf:        make([]bufEntry, size),
		set:        make([]bool, size),
		cpus:       runtime.NumCPU(),
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

func getIdx(w int, pt *d3.Pt) int {
	return w*int(pt.Y) + int(pt.X)
}

func New(w, h int, background *color.RGBA) ZBuffer {
	return newZbuf(w, h, background)
}

func (buf ZBuffer) Add(rm *RenderMesh) {
	w64 := float64(buf.w)
	h64 := float64(buf.h)

	ln := len(rm.Original.Polygons)
	wg := &sync.WaitGroup{}
	wg.Add(buf.cpus)
	var idx32 int32 = -1
	fn := func() {
		for {
			pIdx := int(atomic.AddInt32(&idx32, 1))
			if pIdx >= ln {
				break
			}
			p := rm.Original.Polygons[pIdx]
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
				bi, bt := Scan(triangle.Triangle(t), 0.8)
				for b, done := bi.Start(); !done; b, done = bi.Next() {
					buf.Insert(bt.PtB(b), b, pIdx, tIdx, rm)
				}
			}
		}
		wg.Add(-1)
	}
	for i := 0; i < buf.cpus; i++ {
		go fn()
	}
	wg.Wait()
}

func (buf ZBuffer) Draw(img *image.RGBA) {
	wg := &sync.WaitGroup{}
	wg.Add(buf.cpus)
	ln := len(buf.buf)
	var idx32 int32 = -1
	fn := func(offset int) {
		var c *color.RGBA
		for {
			idx := int(atomic.AddInt32(&idx32, 1))
			if idx >= ln {
				break
			}
			x, y := idx%buf.w, idx/buf.w
			be := buf.buf[idx]
			if buf.set[idx] {
				c = be.Shader(&Context{
					B:           be.B,
					RenderMesh:  be.RenderMesh,
					PolygonIdx:  be.PolygonIdx,
					TriangleIdx: be.TriangleIdx,
				})
				buf.set[idx] = false
			} else {
				c = buf.background
			}

			img.SetRGBA(x, buf.h-y-1, *c)
		}
		wg.Add(-1)
	}
	for i := 0; i < buf.cpus; i++ {
		go fn(i)
	}
	wg.Wait()
}
