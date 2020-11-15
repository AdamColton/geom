// Package generate helps generate random data from the types in the d2 package.
package generate

import (
	"math"
	"math/rand"

	"github.com/adamcolton/geom/d2"
)

// Pt generates a random point in the range (0,0) to (1,1)
func Pt() d2.Pt {
	return d2.Pt{
		X: rand.Float64(),
		Y: rand.Float64(),
	}
}

// PtIn generates a random point from (0,0) to the given point.
func PtIn(pt d2.Pt) d2.Pt {
	return d2.Pt{
		X: pt.X * rand.Float64(),
		Y: pt.Y * rand.Float64(),
	}
}

// V generates a random Vector in the range (0,0) to (1,1)
func V() d2.V {
	return d2.V{
		X: rand.Float64(),
		Y: rand.Float64(),
	}
}

// Polar returns a random vector with length 1.
func Polar() d2.V {
	s, c := math.Sincos(rand.Float64() * 2 * math.Pi)
	return d2.V{
		X: s,
		Y: c,
	}
}

// PtNorm returns a random point with a normal distribution.
func PtNorm() d2.Pt {
	return d2.Pt{
		X: rand.NormFloat64(),
		Y: rand.NormFloat64(),
	}
}

// VNorm returns a random Vector with a normal distribution.
func VNorm() d2.Pt {
	return d2.Pt{
		X: rand.NormFloat64(),
		Y: rand.NormFloat64(),
	}
}
