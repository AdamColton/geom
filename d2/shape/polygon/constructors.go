package polygon

import (
	"math"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
)

// RectangleTwoPoints takes two points and returns a Polygon representing a
// rectangle.
func RectangleTwoPoints(p1, p2 d2.Pt) Polygon {
	return Polygon{
		p1,
		d2.Pt{p2.X, p1.Y},
		p2,
		d2.Pt{p1.X, p2.Y},
	}
}

// RectanglePointWidthLength takes a point and a vector and returns a rectangle
func RectanglePointWidthLength(pt d2.Pt, v d2.V) Polygon {
	pt2 := pt.Add(v)
	return Polygon{
		pt,
		d2.Pt{pt2.X, pt.Y},
		pt2,
		d2.Pt{pt.X, pt2.Y},
	}
}

// RegularPolygonRadius constructs a regular polygon. The radius is measured
// from the center of each side.
func RegularPolygonRadius(center d2.Pt, radius float64, a angle.Rad, sides int) Polygon {
	ps := make(Polygon, sides)
	p := d2.Polar{radius, a}
	da := Tau / float64(sides)
	for i := range ps {
		ps[i] = center.Add(p.V())
		p.A += angle.Rad(da)
	}
	return ps
}

const (
	// Tau constant
	Tau = math.Pi * 2
	// Pi constant
	Pi = math.Pi
)

var rpclC = math.Sin(Pi/2) / (2)

// RegularPolygonSideLength constructs a regular polygon defined by the length
// of the sides.
func RegularPolygonSideLength(center d2.Pt, sideLength float64, a angle.Rad, sides int) Polygon {
	// A right triangle is formed with the hypotenuse being length r (which we
	// want to find), one angle being 360°/(2n) and the opposite side being length
	// (s/2). So the sine law gives us:
	// r / sin(90°) = (s/2) / sin(180°/n)
	// which is gives us
	// r = (s*sin(90°)) / (2*sin(180°/n))
	// r = (sin(90°)/2) * (r/sin(180°/n))
	// so (sin(90°)/2) comes out as a constant, rpclC
	ao := Pi / float64(sides)
	r := (rpclC * sideLength) / math.Sin(ao)
	// rotate backwards so the first line is tangent to the X axis
	a -= angle.Rad(ao)
	return RegularPolygonRadius(center, r, a, sides)
}
