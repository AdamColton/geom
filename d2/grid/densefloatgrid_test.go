package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDenseFloatGrid(t *testing.T) {
	fn := func(pt Pt) float64 {
		return float64(pt.X * pt.Y)
	}
	dfg := NewDenseFloatGrid(Pt{5, 5}, fn)

	c := 0
	for i, done := dfg.Start(); !done; done = dfg.Next() {
		f, err := dfg.Get(i.Pt())
		assert.NoError(t, err)
		assert.Equal(t, fn(i.Pt()), f)
		c++
	}

	assert.Equal(t, dfg.Size().Area(), c)

	var _ Grid = dfg
}

func TestSetTypes(t *testing.T) {
	tt := map[string]interface{}{
		"int":     int(10),
		"int8":    int8(10),
		"int16":   int16(10),
		"int32":   int32(10),
		"int64":   int64(10),
		"uint":    uint(10),
		"uint8":   uint8(10),
		"uint16":  uint16(10),
		"uint32":  uint32(10),
		"uint64":  uint64(10),
		"float32": float32(10),
		"float64": float64(10),
	}

	dfg := NewDenseFloatGrid(Pt{1, 1}, nil)
	pt := Pt{0, 0}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.NoError(t, dfg.SetFloat(pt, 0))
			f, err := dfg.GetFloat(pt)
			assert.NoError(t, err)
			assert.Equal(t, float64(0), f)

			assert.NoError(t, dfg.Set(pt, tc))
			f, err = dfg.GetFloat(pt)
			assert.NoError(t, err)
			assert.Equal(t, float64(10), f)
		})
	}

	assert.Error(t, dfg.Set(pt, "not valid"))
}

func TestDenseFloatGridRangeCheck(t *testing.T) {
	dfg := NewDenseFloatGrid(Pt{1, 1}, nil)
	pt := Pt{1, 1}

	assert.Error(t, dfg.Set(pt, 1))
	_, err := dfg.Get(pt)
	assert.Error(t, err)
}

func TestNormalize(t *testing.T) {
	dfg := &DenseFloatGrid{
		Iterator: Pt{2, 2}.Iter(),
		Data: []float64{
			2, 1,
			4, 5,
		},
	}
	dfg.Normalize()
	expected := []float64{
		0.25, 0,
		0.75, 1,
	}
	assert.Equal(t, expected, dfg.Data)

	dfg = &DenseFloatGrid{
		Iterator: Pt{2, 2}.Iter(),
		Data: []float64{
			2, 2,
			2, 2,
		},
	}
	dfg.Normalize()
	expected = []float64{
		0, 0,
		0, 0,
	}
	assert.Equal(t, expected, dfg.Data)
}
