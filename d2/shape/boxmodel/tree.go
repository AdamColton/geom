package boxmodel

import (
	"sort"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/box"
	"github.com/adamcolton/geom/d2/shape/polygon"
)

type tree struct {
	start                      uint32
	nodes                      []children
	depth                      int
	h, v                       line.Line
	inside, outside, perimeter int
	compressed                 bool
	updateNodes                func()
	area                       float64
	centroid                   d2.Pt
}

func (t *tree) root() *cursor {
	c := t.stacklessRoot()
	c.stack = make([]frame, 0, t.depth)
	return c
}

func (t *tree) stacklessRoot() *cursor {
	if t.nodes == nil && t.updateNodes != nil {
		t.updateNodes()
		t.updateNodes = nil
	}
	return &cursor{
		idx:  t.start,
		tree: t,
		size: 1.0,
		parent: frame{
			node:  0,
			child: 255,
		},
	}
}

func (t *tree) Inside() int    { return t.inside }
func (t *tree) Perimeter() int { return t.perimeter }
func (t *tree) Outside() int   { return t.outside }

func (t *tree) Area() float64       { return t.area }
func (t *tree) SignedArea() float64 { return t.area }
func (t *tree) Centroid() d2.Pt     { return t.centroid }

func (t *tree) InsideCursor() (Iterator, *box.Box, bool) {
	return t.match(insideLeaf)
}

func (t *tree) PerimeterCursor() (Iterator, *box.Box, bool) {
	return t.match(perimeterLeaf)
}

func (t *tree) OutsideCursor() (Iterator, *box.Box, bool) {
	return t.match(outsideLeaf)
}

func (t *tree) ConvexHull() []d2.Pt {
	pts := make([]d2.Pt, 0, t.inside*4)
	for i, b, done := t.InsideCursor(); !done; b, done = i.Next() {
		pts = append(pts, b.ConvexHull()...)
	}
	return polygon.ConvexHull(pts...)
}

func (t *tree) tree() *tree {
	return t
}

func (t *tree) match(tag uint32) (Iterator, *box.Box, bool) {
	c := t.root()
	c.match = tag
	b, done := c.Next()
	return c, b, done
}

func (t *tree) Contains(pt d2.Pt) bool {
	// TODO: cursor that doesn't store stack
	c := t.stacklessRoot()
	if !c.box().Contains(pt) {
		return false
	}
	v := c.at(pt)
	return v == insideLeaf || v == perimeterLeaf
}

func (t *tree) LineIntersections(l line.Line, buf []float64) []float64 {
	max := len(buf)
	buf = buf[:0]
	hits := t.root().lineIntersections(l, buf, nil)
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].t < hits[j].t
	})

	start := -1
	for i, h := range hits {
		if h.kind == byte(perimeterLeaf) {
			if start == -1 {
				start = i
			}
		} else if start != -1 {
			end := i - 1
			if start == end {
				buf = append(buf, hits[start].t)
			} else {
				buf = append(buf, (hits[start].t+hits[end].t)/2)
			}
			if len(buf) == max {
				break
			}
			start = -1
		}
	}
	return buf
}
