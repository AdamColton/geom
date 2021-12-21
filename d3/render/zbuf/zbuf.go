package zbuf

import (
	"image"
	"image/color"
	"runtime"

	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/scene"
	"github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/adamcolton/geom/work"
)

type ZBuffer struct {
	w, h     int
	w64, h64 float64
	buf      []bufEntry
	cpus     int
}

type bufEntry struct {
	*TriangleRef
	barycentric.B
	Z   float64
	set bool
}

func New(w, h int) *ZBuffer {
	size := w * h
	return &ZBuffer{
		w:    w,
		h:    h,
		w64:  float64(w),
		h64:  float64(h),
		buf:  make([]bufEntry, size),
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

func getTriRefs(sf *SceneFrame) []TriangleRef {
	out := make([]TriangleRef, sf.SceneFrame.TriangleCount())
	idx := 0
	for mIdx, m := range sf.Meshes {
		for pIdx, p := range m.Original.Polygons {
			for tIdx := range p {
				out[idx] = TriangleRef{
					TriangleIndex: scene.TriangleIndex{
						MeshIdx:     mIdx,
						PolygonIdx:  pIdx,
						TriangleIdx: tIdx,
					},
				}
				idx++
			}
		}
	}
	return out
}

type cacheRecord struct {
	*barycentric.BIterator
	*triangle.BT
	Space *triangle.BT
}

func (buf *ZBuffer) scanScene(sf *SceneFrame) {
	trs := getTriRefs(sf)
	ln := len(trs)

	dx := 0.8 / buf.w64
	dy := 0.8 / buf.h64

	cache := make([]*cacheRecord, runtime.NumCPU())
	for i := range cache {
		cache[i] = &cacheRecord{
			BIterator: &barycentric.BIterator{},
			BT:        &triangle.BT{},
			Space:     &triangle.BT{},
		}
	}

	work.RunRange(ln, func(idx, threadIdx int) {
		csh := cache[threadIdx]
		tr := &trs[idx]
		m := sf.Meshes[tr.MeshIdx]
		cm := sf.CameraMeshes[tr.MeshIdx]
		p := m.Original.Polygons[tr.PolygonIdx]
		ptIdxs := p[tr.TriangleIdx]
		tr.Camera = &triangle.Triangle{
			cm[ptIdxs[0]],
			cm[ptIdxs[1]],
			cm[ptIdxs[2]],
		}

		if !triangleVisible(tr.Camera) {
			return
		}

		tr.Space = &triangle.Triangle{
			m.Space[ptIdxs[0]],
			m.Space[ptIdxs[1]],
			m.Space[ptIdxs[2]],
		}
		tr.NSpace = tr.Space.Normal().Normal()
		tr.NCamera = tr.Camera.Normal().Normal()

		bi, bCam := Scan(tr.Camera, dx, dy, csh.BIterator, csh.BT)
		tr.BCamera = bCam
		tr.BSpace = tr.Space.BT(bi.Origin, bi.U, csh.Space)
		for b, done := bi.Start(); !done; b, done = bi.Next() {
			buf.insert(bCam.PtB(b), b, tr)
		}
	})
}

func triangleVisible(t *triangle.Triangle) bool {
	// backside culling
	n := t.Normal()
	if n.Z < 0 {
		return false
	}
	// is triangle in frame
	for _, pt := range t {
		if !(pt.X < 0 || pt.X > 1 ||
			pt.Y < 0 || pt.Y > 1 ||
			pt.Z < 0 || pt.Z > 1) {
			return true
		}
	}
	return false
}

func (buf *ZBuffer) insert(pt d3.Pt, b barycentric.B, tr *TriangleRef) {
	if pt.X < 0 || pt.X > 1 || pt.Y < 0 || pt.Y > 1 || pt.Z < 0 || pt.Z > 1 {
		return
	}
	pt.X *= buf.w64
	pt.Y *= buf.h64
	idx := getIdx(buf.w, &pt)
	if idx < 0 || idx >= len(buf.buf) {
		return
	}
	be := &buf.buf[idx]
	if be.set && be.Z < pt.Z {
		return
	}
	be.TriangleRef = tr
	be.B = b
	be.Z = pt.Z
	be.set = true
}

func getIdx(w int, pt *d3.Pt) int {
	return w*int(pt.Y) + int(pt.X)
}

func (buf *ZBuffer) draw(sf *SceneFrame, img *image.RGBA) {
	work.RunRange(buf.h, func(y, _ int) {
		idx := y * buf.w
		var c *color.RGBA
		ctx := &Context{}
		for x := 0; x < buf.w; x++ {
			be := &buf.buf[idx]
			if be.set {
				ctx.SceneFrame = sf
				ctx.TriangleRef = be.TriangleRef
				ctx.B = be.B
				c = sf.Shaders[be.MeshIdx](ctx)
				be.set = false
			} else {
				c = &sf.Background
			}
			img.SetRGBA(x, buf.h-y-1, *c)
			idx++
		}
	})
}
