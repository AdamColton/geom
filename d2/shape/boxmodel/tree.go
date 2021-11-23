package boxmodel

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/box"
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
	if t.nodes == nil && t.updateNodes != nil {
		t.updateNodes()
		t.updateNodes = nil
	}
	return &cursor{
		idx:   t.start,
		tree:  t,
		stack: make([]frame, 0, t.depth),
		size:  1.0,
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

func (t *tree) InsideCursor() (Iterator, box.Box, bool) {
	return t.match(insideLeaf)
}

func (t *tree) PerimeterCursor() (Iterator, box.Box, bool) {
	return t.match(perimeterLeaf)
}

func (t *tree) OutsideCursor() (Iterator, box.Box, bool) {
	return t.match(outsideLeaf)
}

func (t *tree) tree() *tree {
	return t
}

func (t *tree) match(tag uint32) (Iterator, box.Box, bool) {
	c := t.root()
	c.match = tag
	b, done := c.Next()
	return c, b, done
}
