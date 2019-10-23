package grid

import "errors"

// DenseGrid is backed by a slice and should have an entry for every node in
// the grid.
type DenseGrid struct {
	Iterator
	Data []interface{}
}

// NewDenseGrid creates a DenseGrid with Data allocated at the correct size. If
// a Generator is given, the data will be populated.
func NewDenseGrid(i IteratorFactory, generator Generator) *DenseGrid {
	iter := i.Iter()
	g := &DenseGrid{
		Iterator: iter,
		Data:     make([]interface{}, iter.Size().Area()),
	}
	if generator != nil {
		for done := g.Reset(); !done; done = g.Next() {
			g.Data[g.Idx()] = generator(g.Pt())
		}
	}
	return g
}

// Get a value, fulfills Grid
func (g *DenseGrid) Get(pt Pt) (interface{}, error) {
	idx := g.PtIdx(pt)
	if idx < 0 || idx >= len(g.Data) {
		return nil, errors.New("Pt outside of grid")
	}
	return g.Data[idx], nil
}

// Set a value, fulfills Grid
func (g *DenseGrid) Set(pt Pt, val interface{}) error {
	idx := g.PtIdx(pt)
	if idx < 0 || idx >= len(g.Data) {
		return errors.New("Pt outside of grid")
	}
	g.Data[idx] = val
	return nil
}
