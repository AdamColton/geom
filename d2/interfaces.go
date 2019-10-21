package d2

type Point interface {
	Pt() Pt
}

type Vector interface {
	V() V
}

type Pt1 interface {
	Pt1(t0 float64) Pt
}

type Pt1c0 interface {
	Pt1c0() Pt1
}

type Pt2 interface {
	Pt2(t0, t1 float64) Pt
}

type Pt2c1 interface {
	Pt2c1(t0 float64) Pt1
}

type V1 interface {
	V1(t0 float64) V
}

type V1c0 interface {
	V1c0() V1
}

type Pt1V1 interface {
	Pt1
	V1
}

type Limit byte

const (
	LimitUndefined Limit = iota
	LimitBounded
	LimitUnbounded
)

type Limiter interface {
	L(t, c int) Limit
}

type VLimiter interface {
	VL(t ...int) Limit
}

type Area interface {
	Area() float64
	SignedArea() float64
}

type Container interface {
	Contains(Pt) bool
}

type Centroid interface {
	Centroid() Pt
}

type Perimeter interface {
	Perimeter() float64
}

type BoundingBoxer interface {
	BoundingBox() (min Pt, max Pt)
}

type Closest interface {
	Closest(pt Pt) Pt
}

const (
	small = 1e-5
	big   = 1.0 / small
)

type V1Wrapper struct {
	Pt1
}

func (v1 V1Wrapper) V1(t0 float64) V {
	return v1.Pt1.Pt1(t0 + small).Subtract(v1.Pt1.Pt1(t0)).Multiply(big)
}

func GetV1(of Pt1) V1 {
	if v1c0, ok := of.(V1c0); ok {
		return v1c0.V1c0()
	}
	if v1, ok := of.(V1); ok {
		return v1
	}
	return V1Wrapper{of}
}

type Pt2c1Wrapper struct {
	Pt2
}

type T0Wrapper struct {
	T0 float64
	Pt2
}

func GetPt2c1(of Pt2) Pt2c1 {
	if pt2c1, ok := of.(Pt2c1); ok {
		return pt2c1
	}
	return Pt2c1Wrapper{of}
}

func (pt2c1 Pt2c1Wrapper) Pt2c1(t0 float64) Pt1 {
	return T0Wrapper{t0, pt2c1.Pt2}
}

func (t0w T0Wrapper) Pt1(t1 float64) Pt {
	return t0w.Pt2.Pt2(t0w.T0, t1)
}

type TGen interface {
	T() T
	TInv() T
	Pair() [2]T
}
