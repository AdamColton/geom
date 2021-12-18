package grid

import (
	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomerr"
)

// Pt is a cell in a grid
type Pt struct {
	X, Y int
}

// Area always returns a positive value.
func (pt Pt) Area() int {
	a := pt.X * pt.Y
	if a < 0 {
		return -a
	}
	return a
}

// D2 converts a Pt to a d2.D2.
func (pt Pt) D2() d2.D2 {
	return d2.D2{float64(pt.X), float64(pt.Y)}
}

// Abs returns a Pt where both X and Y are positive
func (pt Pt) Abs() Pt {
	if pt.X < 0 {
		pt.X = -pt.X
	}
	if pt.Y < 0 {
		pt.Y = -pt.Y
	}
	return pt
}

// Add two Pts
func (pt Pt) Add(pt2 Pt) Pt {
	return Pt{
		X: pt.X + pt2.X,
		Y: pt.Y + pt2.Y,
	}
}

// Subtract two points
func (pt Pt) Subtract(pt2 Pt) Pt {
	return Pt{
		X: pt.X - pt2.X,
		Y: pt.Y - pt2.Y,
	}
}

// Multiply a Pt by a scale value
func (pt Pt) Multiply(scale float64) Pt {
	return Pt{
		X: int(float64(pt.X) * scale),
		Y: int(float64(pt.Y) * scale),
	}
}

// To creates an Iterator between two points
func (pt Pt) To(pt2 Pt) Iterator {
	return Range{pt, pt2}.Iter()
}

// Iter creates an Iterator from the origin to this Pt
func (pt Pt) Iter() Iterator {
	return Pt{}.To(pt)
}

// Scale is used to convert a Grid Pt to two float64 values, often
type Scale struct {
	X, Y, DX, DY float64
}

// T returns the scaled values corresponding to the point. Typically these are
// used as parametric values.
func (s Scale) T(pt Pt) (float64, float64) {
	return float64(pt.X)*s.X + s.DX, float64(pt.Y)*s.Y + s.DY
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (pt Pt) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	if err := geomerr.NewTypeMismatch(pt, actual); err != nil {
		return err
	}
	pt2 := actual.(Pt)
	if pt.X != pt2.X || pt.Y != pt2.Y {
		return geomerr.NotEqual(pt, pt2)
	}
	return nil
}

const (
	// TwoMask sets the one bit to zero, the result is always divisible by 2.
	TwoMask int = (^int(0)) ^ 1
)

// Mask performs and AND operation with the mask on the X and Y value of the
// given point.
func (pt Pt) Mask(and int) Pt {
	return Pt{
		X: (pt.X) & and,
		Y: (pt.Y) & and,
	}
}
