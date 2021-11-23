package boxmodel

import (
	"errors"
)

func (c *cursor) compressionMap(nextIdx uint32, m map[children]uint32) (uint32, uint32) {
	if c.idx < firstParent {
		return nextIdx, c.idx
	}
	var cs children
	for _, child := range cIter {
		c.moveTo(child)
		nextIdx, cs[child] = c.compressionMap(nextIdx, m)
		c.pop()
	}
	if idx, found := m[cs]; found {
		return nextIdx, idx + firstParent
	}
	m[cs] = nextIdx
	return nextIdx + 1, nextIdx + firstParent
}

// Compressor compacts the quad-trees and reuses segments between models.
type Compressor interface {
	Add(name string, model BoxModel) (BoxModel, error)
	Get(name string) BoxModel
	All() map[string]BoxModel
	Stats() (int, int)
	private()
}

type compressor struct {
	nodes   []children
	trees   map[string]*tree
	cblocks map[children]uint32
	nextIdx uint32
	sum     int
}

func (c *compressor) private() {}

// NewCompressor initilizes a Compressor
func NewCompressor() Compressor {
	return &compressor{
		trees:   make(map[string]*tree),
		cblocks: make(map[children]uint32),
	}
}

// Stats allows the compression to be calculated. The first value is the
// number of boxes the compressor represents and the second is the number it
// actually has to store.
func (c *compressor) Stats() (int, int) {
	return c.sum, len(c.cblocks)
}

func (c *compressor) updateNodes() {
	if ln := len(c.cblocks); len(c.nodes) <= ln {
		c.nodes = make([]children, ln)
		for cs, cb := range c.cblocks {
			c.nodes[cb] = cs
		}
	}
	for _, t := range c.trees {
		t.nodes = c.nodes
	}
}

// Add a model to the compressor by name
func (c *compressor) Add(name string, model BoxModel) (BoxModel, error) {
	if _, alreadyExists := c.trees[name]; alreadyExists {
		return nil, errors.New("Model '" + name + "' already exists")
	}
	t := model.tree()
	c.sum += t.inside + t.outside + t.perimeter
	var start uint32
	c.nextIdx, start = t.root().compressionMap(c.nextIdx, c.cblocks)

	out := &tree{
		start:       start,
		depth:       t.depth,
		h:           t.h,
		v:           t.v,
		inside:      t.inside,
		outside:     t.outside,
		perimeter:   t.perimeter,
		area:        t.area,
		centroid:    t.centroid,
		updateNodes: c.updateNodes,
	}
	c.trees[name] = out
	return out, nil
}

func (c *compressor) Get(name string) BoxModel {
	return c.trees[name]
}

func (c *compressor) All() map[string]BoxModel {
	all := make(map[string]BoxModel, len(c.trees))
	for name, t := range c.trees {
		all[name] = t
	}
	return all
}
