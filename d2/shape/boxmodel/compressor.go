package boxmodel

import (
	"errors"
	"hash/crc32"
)

type cblock struct {
	idx uint32
	children
}

func (c *cursor) compressionMap(nextIdx uint32, m map[uint32]cblock, b []byte) (uint32, uint32) {
	if c.idx < firstParent {
		return nextIdx, c.idx
	}
	var cs children
	for _, child := range cIter {
		c.moveTo(child)
		nextIdx, cs[child] = c.compressionMap(nextIdx, m, b)
		c.pop()
	}
	cs.bytes(b)
	hash := crc32.ChecksumIEEE(b)
	if cb, found := m[hash]; found {
		return nextIdx, cb.idx + firstParent
	}
	m[hash] = cblock{
		idx:      nextIdx,
		children: cs,
	}
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
	cblocks map[uint32]cblock
	nextIdx uint32
	sum     int
}

func (c *compressor) private() {}

// NewCompressor initilizes a Compressor
func NewCompressor() Compressor {
	return &compressor{
		trees:   make(map[string]*tree),
		cblocks: make(map[uint32]cblock),
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
		for _, cb := range c.cblocks {
			c.nodes[cb.idx] = cb.children
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
	b := make([]byte, 16)
	c.nextIdx, start = t.root().compressionMap(c.nextIdx, c.cblocks, b)

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
