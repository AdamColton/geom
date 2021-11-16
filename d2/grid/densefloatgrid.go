package grid

import "errors"

// FloatGenerator returns a float64 given a pt
type FloatGenerator func(pt Pt) float64

// DenseFloatGrid stores a fixed size grid of float64s.
type DenseFloatGrid struct {
	Iterator
	Data []float64
}

// NewDenseFloatGrid with a size determined by the Iterator generated by the
// IteratorFactory. If generator is not nil, the values will be prepopulated.
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

// Get a value, fulfills Grid. The underlying type will always be float64
func (g DenseFloatGrid) Get(pt Pt) (interface{}, error) {
	return g.GetFloat(pt)
}

// Set a value, fulfills Grid. Any numeric value can be passed in and will be
// cast to float64.
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

// Normalize sets the lowest value to 0 and the highest value to 1 and
// distributes all other value proportionally.
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
	if d == 0 {
		for i := range g.Data {
			g.Data[i] = 0
		}
	} else {
		for i, v := range g.Data {
			g.Data[i] = (v - min) / d
		}
	}
}
