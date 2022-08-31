package boxmodel

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/affine"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/box"
)

type frame struct {
	node  uint32
	child byte
}

// the cursor manages the operation of moving through the tree. It is used both
// to construct the tree
type cursor struct {
	*tree
	idx        uint32
	parent     frame
	stack      []frame
	x, y, size float64
	match      uint32
}

func (c *cursor) moveTo(child byte) {
	c.stack = append(c.stack, c.parent)
	c.stacklessMoveTo(child)
}

func (c *cursor) stacklessMoveTo(child byte) {
	c.parent = frame{
		node:  c.idx,
		child: child,
	}
	c.size /= 2
	if child&1 == 1 {
		c.x += c.size
	}
	if child&2 == 2 {
		c.y += c.size
	}
	c.idx = c.nodes[c.idx-firstParent][child]
}

func (c *cursor) pop() (child byte) {
	ln := len(c.stack) - 1
	c.idx = c.parent.node
	child = c.parent.child
	if child&1 == 1 {
		c.x -= c.size
	}
	if child&2 == 2 {
		c.y -= c.size
	}
	c.parent = c.stack[ln]
	c.stack = c.stack[:ln]
	c.size *= 2
	return
}

func (c *cursor) reset() {
	c.idx = c.start
	c.stack = c.stack[:0]
	c.x = 0
	c.y = 0
	c.size = 1
}

func (c *cursor) set(tag uint32) {
	c.nodes[c.parent.node-firstParent][c.parent.child] = tag
}

func (c *cursor) center() d2.Pt {
	s := c.size / 2
	return d2.Pt{
		X: c.h.Pt1(c.x + s).X,
		Y: c.v.Pt1(c.y + s).Y,
	}
}

func (c *cursor) insert(x, y float64, depth int) {
	for i := 0; i < depth; i++ {
		if c.idx == unknownLeaf {
			c.idx = uint32(len(c.nodes)) + firstParent
			c.nodes = append(c.nodes, children{})
			c.set(c.idx)
		}
		s := c.size / 2
		var child byte
		if x > c.x+s {
			child = 1
		}
		if y > c.y+s {
			child += 2
		}
		c.moveTo(child)
	}
	c.set(perimeterLeaf)
	c.reset()
}

func (c *cursor) tag(s shape.Shape) {
	if c.idx == perimeterLeaf {
		c.tree.perimeter++
	} else if c.idx == unknownLeaf {
		if s.Contains(c.center()) {
			c.tree.inside++
			c.set(insideLeaf)
		} else {
			c.tree.outside++
			c.set(outsideLeaf)
		}
	} else {
		for _, child := range cIter {
			c.moveTo(child)
			c.tag(s)
			c.pop()
		}
	}
}

type sum struct {
	centroid *affine.Weighted
	area     float64
}

func (c *cursor) sum(s *sum) {
	if c.idx == perimeterLeaf {
		a := c.size * c.size / 2.0
		s.centroid.Weight(c.center(), a)
		s.area += a
	} else if c.idx == insideLeaf {
		a := c.size * c.size
		s.centroid.Weight(c.center(), a)
		s.area += a
	} else if c.idx >= firstParent {
		for _, child := range cIter {
			c.moveTo(child)
			c.sum(s)
			c.pop()
		}
	}
}

func (c *cursor) Next() (b *box.Box, done bool) {
	for c.nextLeaf(0) {
		if c.idx == c.match {
			return c.box(), false
		}
	}
	return &box.Box{}, true
}

func (c *cursor) nextLeaf(child byte) bool {
	if c.idx < firstParent || child > 3 {
		if c.parent.child == 255 {
			return false
		}
		return c.nextLeaf(c.pop() + 1)
	}
	c.moveTo(child)
	for c.idx >= firstParent {
		c.moveTo(0)
	}
	return true
}

func (c *cursor) box() *box.Box {
	return &box.Box{
		d2.Pt{c.h.Pt1(c.x).X, c.v.Pt1(c.y).Y},
		d2.Pt{c.h.Pt1(c.x + c.size).X, c.v.Pt1(c.y + c.size).Y},
	}
}

func (c *cursor) scan(s shape.Shape, depth int) {
	step := 1.0 / float64(uint(1)<<uint(depth))
	h := line.Line{
		D: c.h.D,
	}
	v := line.Line{
		D: c.v.D,
	}
	var buf []float64
	for t := step / 2; t <= 1.0; t += step {
		h.T0 = c.v.Pt1(t)
		buf = s.LineIntersections(h, buf[:0])
		for _, ht := range buf {
			c.insert(ht, t, depth)
		}

		v.T0 = c.h.Pt1(t)
		buf = s.LineIntersections(v, buf[:0])
		for _, vt := range buf {
			c.insert(t, vt, depth)
		}
	}
}

func (c *cursor) at(pt d2.Pt) uint32 {
	if c.idx < firstParent {
		return c.idx
	}
	b := c.box()
	center := b.Centroid()
	var child byte
	if pt.X > center.X {
		child += 1
	}
	if pt.Y > center.Y {
		child += 2
	}
	c.stacklessMoveTo(child)
	return c.at(pt)
}

type intersection struct {
	t    float64
	kind byte
}

func (c *cursor) lineIntersections(l line.Line, buf []float64, hits []intersection) []intersection {
	bx := c.box()
	buf = bx.LineIntersections(l, buf)
	ln := len(buf)
	if ln == 0 {
		return hits
	}
	if c.idx < firstParent {
		if ln >= 2 {
			buf[0] = (buf[0] + buf[1]) / 2.0
		}
		return append(hits, intersection{
			t:    buf[0],
			kind: byte(c.idx),
		})
	}

	buf = buf[:0]
	for child := byte(0); child < 4; child++ {
		c.moveTo(child)
		hits = c.lineIntersections(l, buf, hits)
		c.pop()
	}
	return hits
}
