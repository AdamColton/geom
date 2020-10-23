package zbuf

import (
	"image"
	"image/color"
	"runtime"

	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/scene"
	"github.com/adamcolton/geom/d3/shape/triangle"
)

type ZBuffer struct {
	w, h     int
	w64, h64 float64
	buf      []bufEntry
	set      []bool
	cpus     int
}

type triRef struct {
	meshIdx     int
	polygonIdx  int
	triangleIdx int
}

type bufEntry struct {
	*triRef
	barycentric.B
	Z float64
}

func New(w, h int) *ZBuffer {
	size := w * h
	return &ZBuffer{
		w:    w,
		h:    h,
		w64:  float64(w),
		h64:  float64(h),
		buf:  make([]bufEntry, size),
		set:  make([]bool, size),
		cpus: runtime.NumCPU(),
	}
}

func (buf *ZBuffer) Draw(sf *SceneFrame, img *image.RGBA) *image.RGBA {
	if sf.CameraMeshes == nil {
		sf.PopulateCameraMeshes()
	}
	if sf.Shaders == nil {
		sf.PopulateShaders()
	}
	buf.scanScene(sf)

	if img == nil {
		r := image.Rect(0, 0, buf.w, buf.h)
		img = image.NewRGBA(r)
	}
	buf.draw(sf, img)
	return img
}

func getTriRefs(sf *SceneFrame) []triRef {
	ln := 0
	for _, m := range sf.Meshes {
		for _, p := range m.Original.Polygons {
			ln += len(p)
		}
	}
	out := make([]triRef, ln)
	idx := 0
	for mIdx, m := range sf.Meshes {
		for pIdx, p := range m.Original.Polygons {
			for tIdx := range p {
				out[idx] = triRef{
					meshIdx:     mIdx,
					polygonIdx:  pIdx,
					triangleIdx: tIdx,
				}
				idx++
			}
		}
	}
	return out
}

func (buf *ZBuffer) scanScene(sf *SceneFrame) {
	trs := getTriRefs(sf)
	ln := len(trs)

	dx := 0.8 / buf.w64
	dy := 0.8 / buf.h64

	scene.RunRange(ln, func(idx, _ int) {
		tr := trs[idx]
		m := sf.Meshes[tr.meshIdx]
		cm := sf.CameraMeshes[tr.meshIdx]
		p := m.Original.Polygons[tr.polygonIdx]
		ptIdxs := p[tr.triangleIdx]
		t := &triangle.Triangle{
			cm[ptIdxs[0]],
			cm[ptIdxs[1]],
			cm[ptIdxs[2]],
		}

		if !triangleVisible(t) {
			return
		}

		bi, bt := Scan(t, dx, dy)
		for b, done := bi.Start(); !done; b, done = bi.Next() {
			buf.insert(bt.PtB(b), b, &tr)
		}
	})
}

func triangleVisible(t *triangle.Triangle) bool {
	for _, pt := range t {
		if !(pt.X < 0 || pt.X > 1 ||
			pt.Y < 0 || pt.Y > 1 ||
			pt.Z < 0 || pt.Z > 1) {
			return true
		}
	}
	return false
}

func (buf *ZBuffer) insert(pt d3.Pt, b barycentric.B, tr *triRef) {
	if pt.X < 0 || pt.X > 1 || pt.Y < 0 || pt.Y > 1 || pt.Z < 0 || pt.Z > 1 {
		return
	}
	pt.X *= buf.w64
	pt.Y *= buf.h64
	idx := getIdx(buf.w, &pt)
	if idx < 0 || idx >= len(buf.buf) || (buf.set[idx] && buf.buf[idx].Z < pt.Z) {
		return
	}
	be := &buf.buf[idx]
	be.triRef = tr
	be.B = b
	be.Z = pt.Z
	buf.set[idx] = true
}

func getIdx(w int, pt *d3.Pt) int {
	return w*int(pt.Y) + int(pt.X)
}

func (buf *ZBuffer) draw(sf *SceneFrame, img *image.RGBA) {
	scene.RunRange(len(buf.buf), func(idx, _ int) {
		var c *color.RGBA
		x, y := idx%buf.w, idx/buf.w
		be := buf.buf[idx]
		if buf.set[idx] {
			c = sf.Shaders[be.meshIdx](&Context{
				SceneFrame:  sf,
				MeshIdx:     be.meshIdx,
				PolygonIdx:  be.polygonIdx,
				TriangleIdx: be.triangleIdx,
				B:           be.B,
			})
			buf.set[idx] = false
		} else {
			c = &sf.Background
		}
		img.SetRGBA(x, buf.h-y-1, *c)
	})
}
