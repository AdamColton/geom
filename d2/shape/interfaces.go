package shape

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

// Area of a shape, signed area may take procession into account. For instance,
// a triangle with points defined proceeding counter-clockwise wil have a
// positive area and proceeding clockwise will have a negative area.
type Area interface {
	Area() float64
	SignedArea() float64
}

// Container checks if a shape contains a point.
type Container interface {
	Contains(d2.Pt) bool
}

// Centroid is the center of mass of a shape.
type Centroid interface {
	Centroid() d2.Pt
}

// Perimeter length of the shape
type Perimeter interface {
	Perimeter() float64
}

// BoundingBoxer returns the corners of a bounding box.
type BoundingBoxer interface {
	BoundingBox() (min, max d2.Pt)
}

// Closest point on the perimeter to given point.
type Closest interface {
	Closest(pt d2.Pt) d2.Pt
}

// Shape is an interface that is easy to implement for most primitives but
// allows for complex generic shape functions.
type Shape interface {
	Container
	line.Intersector
	BoundingBoxer
}
