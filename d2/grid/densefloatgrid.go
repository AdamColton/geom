package grid

import "errors"

type FloatGenerator func(pt Pt) float64

type DenseFloatGrid struct {
	Iterator
	Data []float64
}

func NewDenseFloatGrid(i IteratorFactory, generator FloatGenerator) DenseFloatGrid {
	iter := i.Iter()
	g := DenseFloatGrid{
		Iterator: iter,
		Data:     make([]float64, iter.Size().Area()),
	}
	if generator != nil {
		for done := g.Reset(); !done; done = g.Next() {
			g.Data[g.Idx()] = generator(g.Pt())
		}
	}
	return g
}

// Get a value, fulfills Grid
func (g DenseFloatGrid) Get(pt Pt) (interface{}, error) {
	return g.GetFloat(pt)
}

// Set a value, fulfills Grid
func (g DenseFloatGrid) Set(pt Pt, val interface{}) error {
	switch v := val.(type) {
	case float64:
		return g.SetFloat(pt, v)
	case float32:
		return g.SetFloat(pt, float64(v))
	case int:
		return g.SetFloat(pt, float64(v))
	case int8:
		return g.SetFloat(pt, float64(v))
	case int16:
		return g.SetFloat(pt, float64(v))
	case int32:
		return g.SetFloat(pt, float64(v))
	case int64:
		return g.SetFloat(pt, float64(v))
	case uint:
		return g.SetFloat(pt, float64(v))
	case uint8:
		return g.SetFloat(pt, float64(v))
	case uint16:
		return g.SetFloat(pt, float64(v))
	case uint32:
		return g.SetFloat(pt, float64(v))
	case uint64:
		return g.SetFloat(pt, float64(v))
	}
	return errors.New("val must be numeric type")
}

// GetFloat from DenseFloatGrid
func (g DenseFloatGrid) GetFloat(pt Pt) (float64, error) {
	idx := g.PtIdx(pt)
	if idx < 0 || idx >= len(g.Data) {
		return 0, errors.New("Pt outside of grid")
	}
	return g.Data[idx], nil
}

// SetFloat in DenseFloatGrid
func (g DenseFloatGrid) SetFloat(pt Pt, val float64) error {
	idx := g.PtIdx(pt)
	if idx < 0 || idx >= len(g.Data) {
		return errors.New("Pt outside of grid")
	}
	g.Data[g.PtIdx(pt)] = val
	return nil
}

func (g DenseFloatGrid) Normalize() {
	min := g.Data[0]
	max := g.Data[0]
	for _, v := range g.Data[1:] {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}

	d := max - min
	for i, v := range g.Data {
		g.Data[i] = (v - min) / d
	}
}
