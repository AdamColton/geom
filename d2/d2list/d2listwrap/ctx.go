package d2listwrap

import (
	"github.com/adamcolton/geom/d2/d2list"
	"github.com/adamcolton/geom/d2/draw"
)

type DrawPoint struct {
	R   float64
	RGB [3]float64
}

var DefaultDrawPoint = &DrawPoint{
	R:   3,
	RGB: [3]float64{1, 0, 0},
}

type PointContext struct {
	d2list.PointList
	Listeners []Updater
	*DrawPoint
}

func NewPointContext(pts d2list.PointList) *PointContext {
	return &PointContext{
		PointList: pts,
	}
}

func (pc *PointContext) Update() {
	for _, l := range pc.Listeners {
		l.Update()
	}
}

func (pc *PointContext) SubList(pts ...int) PointSubList {
	return PointSubList{
		PointList: pc,
		Idxs:      pts,
	}
}
func (pc *PointContext) AddListeners(listeners ...Updater) {
	pc.Listeners = append(pc.Listeners, listeners...)
}

func (pc *PointContext) Draw(ctx *draw.Context) {
	d := pc.DrawPoint
	if d == nil {
		d = DefaultDrawPoint
	}
	ctx.SetRGB(d.RGB[0], d.RGB[1], d.RGB[2])
	DrawPointList(pc, d.R, ctx)
}

func (pc *PointContext) Line(start, end int) *Line {
	l := &Line{
		PointList: pc.SubList(),
	}
	l.Update()
	pc.AddListeners(l)
	return l
}

func (pc *PointContext) Bezier(pts ...int) *Bezier {
	b := &Bezier{
		PointList: pc.SubList(pts...),
	}
	b.Update()
	pc.AddListeners(b)
	return b
}

func (pc *PointContext) Triangle(a, b, c int) *Triangle {
	t := &Triangle{
		PointList: pc.SubList(a, b, c),
	}
	t.Update()
	pc.AddListeners(t)
	return t
}
