package boxmodel

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/affine"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/box"
)

// BoxModel represents a shape that has been encoded as a set of boxes.
type BoxModel interface {
	// Inside returns the number of boxes inside the shape.
	Inside() int
	// InsideCursor returns a cursor that will iterate over all the boxes
	// inside the shape.
	InsideCursor() (Iterator, box.Box, bool)

	// Perimeter returns the number of boxes that contain the perimeter.
	Perimeter() int
	// PerimeterCursor returns a cursor that will iterate over all the boxes
	// on the perimeter of the shape.
	PerimeterCursor() (Iterator, box.Box, bool)

	// Outside returns the number of boxes outside the shape.
	Outside() int
	// OutsideCursor returns a cursor that will iterate over all the boxes
	// inside the shape.
	OutsideCursor() (Iterator, box.Box, bool)

	// Area is an approximation of the area of the shape. It is the sum of all
	// the boxes inside the shape and half the area of the boxes on the
	// perimeter.
	Area() float64

	// SignedArea is the same as Area.
	SignedArea() float64

	ConvexHull() []d2.Pt

	// Centroid is the center of mass of the shape.
	Centroid() d2.Pt
	tree() *tree
}

// Iterator iterates over a collection of boxes
type Iterator interface {
	Next() (b box.Box, done bool)
}

// New BoxModel representing the shape.
func New(s shape.Shape, depth int) BoxModel {
	b := box.New(s.ConvexHull()...)
	t := &tree{
		start: firstParent,
		nodes: make([]children, 1, 1<<(depth+2)),
		depth: depth,
		h: line.Line{
			T0: b[0],
			D:  d2.V{b[1].X - b[0].X, 0},
		},
		v: line.Line{
			T0: b[0],
			D:  d2.V{0, b[1].Y - b[0].Y},
		},
	}
	root := t.root()
	root.scan(s, depth)
	root.tag(s)
	sm := &sum{
		centroid: affine.NewWeighted(root.inside + root.perimeter),
	}
	root.sum(sm)
	root.area = sm.area * b.Area()
	root.centroid = sm.centroid.Centroid()

	return root.tree
}
