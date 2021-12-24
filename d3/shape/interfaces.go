package shape

import "github.com/adamcolton/geom/d3"

type Area interface {
	Area() float64
	SignedArea() float64
}

type Container interface {
	Contains(d3.Pt) bool
}

type Centroid interface {
	Centroid() d3.Pt
}

type Perimeter interface {
	Perimeter() float64
}

type BoundingCube interface {
	BoundingCube() (min d3.Pt, max d3.Pt)
}
