package grid

type SparseGrid struct {
	Iterator
	Data map[Pt]interface{}
	Generator
}

func NewSparseGrid(i IteratorFactory, generator Generator) *SparseGrid {
	g := &SparseGrid{
		Iterator:  i.Iter(),
		Data:      make(map[Pt]interface{}),
		Generator: generator,
	}
	return g
}

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

func (g *SparseGrid) Set(pt Pt, val interface{}) error {
	if val == nil {
		delete(g.Data, pt)
	} else {
		g.Data[pt] = val
	}
	return nil
}
