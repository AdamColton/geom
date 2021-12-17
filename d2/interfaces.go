package d2

import (
	"math"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
)

// Point can return a Pt
type Point interface {
	Pt() Pt
}

// Vector can return a V
type Vector interface {
	V() V
}

// Pt1 takes one parametric value and returns a Pt
type Pt1 interface {
	Pt1(t0 float64) Pt
}

// Pt1c0 returns a Pt1, it may cache some computation to optimize future calls
// to Pt1
type Pt1c0 interface {
	Pt1c0() Pt1
}

// Pt2 takes two parametric values and returns a Pt
type Pt2 interface {
	Pt2(t0, t1 float64) Pt
}

// Pt2c1 curries one of two parametric values returning a Pt1
type Pt2c1 interface {
	Pt2c1(t0 float64) Pt1
}

// V1 takes one parametric value and returns a V
type V1 interface {
	V1(t0 float64) V
}

// V1c0 returns a V1, it may cache some computation to optimize future callse to
// V1
type V1c0 interface {
	V1c0() V1
}

// Pt1V1 has one argument parametric methods for both Pt1 and V1, typically this
// represents a curve and it's derivative.
type Pt1V1 interface {
	Pt1
	V1
}

// Limit is used to indicate if a parametric method is bounded to [0,1] or
// unbounded
type Limit byte

const (
	// LimitUndefined indicates that the requested parametric method is not
	// defined.
	LimitUndefined Limit = iota
	// LimitBounded indicates that the behavior of a parametric funnction
	// outside the range [0,1] is not defined.
	LimitBounded
	// LimitUnbounded indicates that passing parametric values outside the range
	// [0,1] should behave predictibly.
	LimitUnbounded
)

// Limiter can describe the behavior of it's parametric methods that return a Pt
type Limiter interface {
	L(t, c int) Limit
}

// VLimiter can describe the behavior of it's parametric methods that return a V
type VLimiter interface {
	VL(t, c int) Limit
}

const (
	small = 1e-7
	big   = 1.0 / small
)

// V1Wrapper takes any Pt1 and approximates V1
type V1Wrapper struct {
	P     Pt1
	Small float64
}

// V1 approximates V1 from two points close together
func (v1 V1Wrapper) V1(t0 float64) V {
	if v1.Small == 0 {
		v1.Small = small
	}
	big := 1 / v1.Small
	return v1.P.Pt1(t0 + v1.Small).Subtract(v1.P.Pt1(t0)).Multiply(big)
}

// Pt1 calls underlying Pt1 to fulfill Pt1V1 interface
func (v1 V1Wrapper) Pt1(t0 float64) Pt {
	return v1.P.Pt1(t0)
}

// GetV1 takes any Pt1 and returns the optimal V1.
func GetV1(of Pt1) V1 {
	if v1c0, ok := of.(V1c0); ok {
		return v1c0.V1c0()
	}
	if v1, ok := of.(V1); ok {
		return v1
	}
	return V1Wrapper{of, 0}
}

// Pt2c1Wrapper wraps any Pt2 to convert it to a Pt2c1
type Pt2c1Wrapper struct {
	Pt2
}

// T0Wrapper curries the t0 parametric value to a Pt2
type T0Wrapper struct {
	T0 float64
	Pt2
}

// GetPt2c1 returns a Pt2c1 from any Pt2 prefering Pt2c1 on the interface if it
// is defined
func GetPt2c1(of Pt2) Pt2c1 {
	if pt2c1, ok := of.(Pt2c1); ok {
		return pt2c1
	}
	return Pt2c1Wrapper{of}
}

// Pt2c1 curries the first argument using a T0Wrapper
func (pt2c1 Pt2c1Wrapper) Pt2c1(t0 float64) Pt1 {
	return T0Wrapper{t0, pt2c1.Pt2}
}

// Pt1 passes the curried t0 and the method argument t1 into the underlying Pt2
func (t0w T0Wrapper) Pt1(t1 float64) Pt {
	return t0w.Pt2.Pt2(t0w.T0, t1)
}

// TGen is a Transform generator. If both the T and it's inverse are needed,
// Pair may reduce some duplicate calculations.
type TGen interface {
	T() *T
	TInv() *T
	Pair() [2]*T
}

// AssertV1 checks that the derivative is close the derivative approximation. A
// good base check that an algorithm isn't completely wrong.
type AssertV1 struct{}

// AssertEqual fulfils geomtest.AssertEqualizer.
func (AssertV1) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	// V1Wrapper is not very accurate, so we reduce the required tolerance.
	t = cmpr.Tolerance(math.Sqrt(float64(t)))

	pv, ok := actual.(Pt1V1)
	if !ok {
		return geomerr.TypeMismatch(Pt1V1(nil), actual)
	}
	fn := make([]func(float64) V, 2, 3)
	fn[0] = V1Wrapper{pv, 0}.V1
	fn[1] = pv.V1
	if v1c0, ok := pv.(V1c0); ok {
		fn = append(fn, v1c0.V1c0().V1)
	}

	var badPoints geomerr.SliceErrs
	v := make([]V, len(fn))
	for i := 0.0; i <= 1.0; i += 0.01 {
		for j, f := range fn {
			v[j] = f(i)
			if j > 0 {
				if v[0].AssertEqual(v[j], t) != nil {
					r := V{
						X: v[0].X / v[j].X,
						Y: v[0].Y / v[j].Y,
					}
					str := "V1"
					if j == 2 {
						str = "V1c0"
					}
					badPoints = badPoints.AppendF(int(i/0.01), "Bad derivative at %s(%0.2f) Expected: %s Got:%s Ratio:%s", str, i, v[0], v[j], r)
				}
			}
		}
	}

	if len(badPoints) == 0 {
		return nil
	}
	return badPoints
}
