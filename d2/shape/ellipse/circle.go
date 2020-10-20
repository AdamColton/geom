package ellipse

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape/triangle"
)

// Circle fulfills Shape.
type Circle struct {
	Ellipse
}

// NewCircle returns a circle defined by a center and radius
func NewCircle(center d2.Pt, radius float64) Circle {
	return Circle{Ellipse: New(center, center, radius)}
}

// CircumscribeCircle creates a circle where all three verticies lie on the
// perimeter.
func CircumscribeCircle(t triangle.Triangle) Circle {
	c := t.CircumCenter()
	r := t[0].Distance(c)
	return NewCircle(c, r)
}

// Radius returns the radius of the circle
func (c Circle) Radius() float64 {
	r, _ := c.perimeter.Axis()
	return r
}
