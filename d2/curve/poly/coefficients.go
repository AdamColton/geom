package poly

import (
	"github.com/adamcolton/geom/d2"
)

// Coefficients wraps the concept of a list of d2.V. It must return the length
// and be able to return the coefficient at any index.
type Coefficients interface {
	Coefficient(idx int) d2.V
	Len() int
}

// X converts and instance of Coefficients into a set of 1D Coefficients using
// their X components.
type X struct{ Coefficients }

// Coefficient returns the X value of the underlying Coefficients.
func (x X) Coefficient(idx int) float64 {
	if idx >= x.Len() || idx < 0 {
		return 0
	}
	return x.Coefficients.Coefficient(idx).X
}

// Len returns the Len of the underlying Coefficients.
func (x X) Len() int {
	return x.Coefficients.Len()
}

// Y converts and instance of Coefficients into a set of 1D Coefficients using
// their Y components.
type Y struct{ Coefficients }

// Coefficient returns the Y value of the underlying Coefficients.
func (y Y) Coefficient(idx int) float64 {
	if idx >= y.Len() || idx < 0 {
		return 0
	}
	return y.Coefficients.Coefficient(idx).Y
}

// Len returns the Len of the underlying Coefficients.
func (y Y) Len() int {
	return y.Coefficients.Len()
}

// Slice fulfills Coefficients using a Slice.
type Slice []d2.V

// Coefficient returns the d2.V at the given index if it is in range, otherwise
// it returns d2.V{0,0}.
func (s Slice) Coefficient(idx int) d2.V {
	if idx >= s.Len() || idx < 0 {
		return d2.V{}
	}
	return s[idx]
}

// Len returns the length of the underlying slice.
func (s Slice) Len() int {
	return len(s)
}
