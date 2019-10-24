package grid

// SparseGrid is backed by a map so it can hold an arbitrary amount of data.
// If a generator is defined, it will be invoked lazily.
type SparseGrid struct {
	Iterator
	Data map[Pt]interface{}
	Generator
}

// NewSparseGrid creates a SparseGrid. If an IteratorFactory is defined, that
// will be used as an iterator. This is for convience only as a sparse grid is
// not actually limited to any range. If a generator is defined, it will be
// invoked lazily when accessing points that have not been set.
func NewSparseGrid(i IteratorFactory, generator Generator) *SparseGrid {
	g := &SparseGrid{
		Data:      make(map[Pt]interface{}),
		Generator: generator,
	}
	if i != nil {
		g.Iterator = i.Iter()
	}
	return g
}

// Get the value at the defined Pt in the grid. If the point value has not been
// set and a generator has been set the value will be generated and stored.
func (g *SparseGrid) Get(pt Pt) (interface{}, error) {
	v, ok := g.Data[pt]
	if !ok && g.Generator != nil {
		v = g.Generator(pt)
		if v != nil {
			g.Data[pt] = v
		}
	}
	return v, nil
}

// Set the value at the point. Setting a value to nil will remove it from the
// underlying map. This means that if there is a generator, the next Get will
// update that point with the value from the generator.
func (g *SparseGrid) Set(pt Pt, val interface{}) error {
	if val == nil {
		delete(g.Data, pt)
	} else {
		g.Data[pt] = val
	}
	return nil
}
