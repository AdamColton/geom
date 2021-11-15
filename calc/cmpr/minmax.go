package cmpr

import (
	"fmt"
	"math"
)

type MinMax struct {
	Min, Max float64
}

func NewMinMax() *MinMax {
	return &MinMax{
		Min: math.Inf(1),
		Max: math.Inf(-1),
	}
}

func (mm *MinMax) Update(f float64) bool {
	out := false
	if f < mm.Min {
		mm.Min = f
		out = true
	}
	if f > mm.Max {
		mm.Max = f
		out = true
	}
	return out
}

func (mm *MinMax) String() string {
	return fmt.Sprintf("(%f, %f)", mm.Min, mm.Max)
}
