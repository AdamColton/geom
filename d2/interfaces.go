package d2

type Curve func(t float64) Pt

type Curver interface {
	Pt(t float64) Pt
}

type DerivativeCurver interface {
	Curver
	V(t float64) V
}

// Surface functions describe 2D area parametrically. They should have the
// following properties
// * All points on the perimeter should have either t0==0 or t1==0
// * The surface should have no creases
// * Pt(ta0,ta1)==Pt(tb0,tb1) --> ta0==tb0 && ta1==tb1
type Surface func(t0, t1 float64) Pt

type ParametricShape interface {
	Pt(t0, t1 float64) Pt
}

type ParametricShapeCurve interface {
	ParametricShape
	Curve(t0 float64) Curve
}

type AreaShape interface {
	Area() float64
	SignedArea() float64
}

type ContainerShape interface {
	Contains(Pt) bool
}

type CentroidShape interface {
	Centroid() Pt
}

type PerimeterShape interface {
	Perimeter() float64
}

type BoundingBoxShape interface {
	BoundingBox() (min Pt, max Pt)
}
