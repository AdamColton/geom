package mesh

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/affine"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/adamcolton/geom/d3/solid"
)

type Extrusion struct {
	cur    []uint32
	points *solid.PointSet
	faces  [][]uint32
}

func NewExtrusion(face []d3.Pt) *Extrusion {
	ln := len(face)
	e := &Extrusion{
		points: solid.NewPointSet(),
		cur:    make([]uint32, ln),
	}
	for i, pt := range face {
		e.cur[i] = e.points.Add(pt)
	}
	e.faces = append(e.faces, e.cur)
	return e
}

func (e *Extrusion) applyTRelativeToCenter(t *d3.T) *d3.T {
	w := &affine.Weighted{}
	for _, ptIdx := range e.cur {
		w.Add(e.points.Pts[ptIdx])
	}
	center := d3.Pt{}.Subtract(w.Get())
	return d3.NewTSet().
		AddBoth(d3.Translate(center).Pair()).
		Add(t).
		Get()
}

func (e *Extrusion) Extrude(ts ...*d3.T) *Extrusion {
	ln := len(e.cur)
	for _, t := range ts {
		t = e.applyTRelativeToCenter(t)
		nxt := make([]uint32, len(e.cur))
		for i, ptIdx := range e.cur {
			nxt[i] = e.points.Add(t.Pt(e.points.Pts[ptIdx]))
		}

		prev := ln - 1
		for i := range e.cur {
			e.faces = append(e.faces, []uint32{
				e.cur[prev],
				e.cur[i],
				nxt[i],
				nxt[prev],
			})
			prev = i
		}

		e.cur = nxt
	}
	return e
}

func (e *Extrusion) EdgeExtrude(t *d3.T) *Extrusion {
	lnCur := len(e.cur)
	lnNxt := 3 * lnCur
	nxt := make([]uint32, lnNxt)

	prev := e.points.Pts[e.cur[lnCur-1]]
	t = e.applyTRelativeToCenter(t)
	for i, cIdx := range e.cur {
		cur := e.points.Pts[cIdx]
		l := line.New(prev, cur)
		for j := 0; j < 3; j++ {
			f := float64(j) / 3.0
			nxt[(i*3+j-3+lnNxt)%lnNxt] = e.points.Add(t.Pt(l.Pt1(f)))
		}
		prev = cur
	}

	for i, cIdx := range e.cur {
		e.faces = append(e.faces, []uint32{
			nxt[i*3],
			nxt[i*3+1],
			cIdx,
			nxt[(i*3+lnNxt-1)%lnNxt],
		}, []uint32{
			nxt[i*3+1],
			nxt[i*3+2],
			e.cur[(i+1)%lnCur],
			cIdx,
		})
	}

	e.cur = nxt

	return e
}

func (e *Extrusion) EdgeMerge(t *d3.T) *Extrusion {
	lnCur := len(e.cur)
	if lnCur%3 != 0 {
		return e
	}
	lnNxt := lnCur / 3
	nxt := make([]uint32, lnNxt)
	t = e.applyTRelativeToCenter(t)
	for i := range nxt {
		nxt[i] = e.points.Add(t.Pt(e.points.Pts[e.cur[i*3]]))
	}

	for i, nIdx := range nxt {
		e.faces = append(e.faces, []uint32{
			nIdx,
			e.cur[(i*3-1+lnCur)%lnCur],
			e.cur[i*3],
			e.cur[i*3+1],
		}, []uint32{
			nIdx,
			e.cur[i*3+1],
			e.cur[i*3+2],
			nxt[(i+1)%lnNxt],
		})
	}

	e.cur = nxt
	return e
}

func (e *Extrusion) Close() Mesh {
	e.faces = append(e.faces, e.cur)
	return Mesh{
		Polygons: e.faces,
		Pts:      e.points.Pts,
	}
}
