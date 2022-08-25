package box

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

// Box is a rectangle that lies orthoganal to the plane. The first point should
// be the min point and the second should be the max.
type Box [2]d2.Pt

// New Box containing all the points passed in.
func New(pts ...d2.Pt) *Box {
	b := &Box{}
	b[0], b[1] = d2.MinMax(pts...)
	return b
}

func (b *Box) Add(pts ...d2.Pt) {
	for _, pt := range pts {
		if pt.X < b[0].X {
			b[0].X = pt.X
		}
		if pt.Y < b[0].Y {
			b[0].Y = pt.Y
		}
		if pt.X > b[1].X {
			b[1].X = pt.X
		}
		if pt.Y > b[1].Y {
			b[1].Y = pt.Y
		}
	}
}

// V is the vector from the min point to the max point
func (b *Box) V() d2.V {
	return b[1].Subtract(b[0])
}

// Vertex returns one of the 4 verticies of the box proceeding counter
// clockwise.
func (b *Box) Vertex(n int) d2.Pt {
	n %= 4
	if n == 2 {
		return b[1]
	}
	p := b[0]
	if n == 1 {
		p.X = b[1].X
	} else if n == 3 {
		p.Y = b[1].Y
	}
	return p
}

// Side returns one of the sides proceeding counter clockwise.
func (b *Box) Side(n int) line.Line {
	return line.New(b.Vertex(n), b.Vertex(n+1))
}

// Sides returns an array of the 4 sides in counter clockwise order.
func (b *Box) Sides() *[4]line.Line {
	v := b.V()
	return &[4]line.Line{
		{b[0], d2.V{v.X, 0}},
		{d2.Pt{b[1].X, b[0].Y}, d2.V{0, v.Y}},
		{b[1], d2.V{-v.X, 0}},
		{d2.Pt{b[0].X, b[1].Y}, d2.V{0, -v.Y}},
	}
}

// LineIntersections returning the sides of the box that are intersected.
// Fulfills line.Intersector and shape.Shape.
func (b *Box) LineIntersections(l line.Line, buf []float64) []float64 {
	max := len(buf)
	buf = buf[:0]
	for _, s := range b.Sides() {
		t0, t1, ok := l.Intersection(s)
		if ok && t0 >= 0 && t0 < 1 {
			buf = append(buf, t1)
			if max == 1 {
				return buf
			}
		}
	}
	return buf
}

// Centroid returns the center of the box, fulfilling shape.Centroid
func (b *Box) Centroid() d2.Pt {
	return line.New(b[0], b[1]).Pt1(0.5)
}

// Area of the box fulling shape.Area.
func (b *Box) Area() float64 {
	a := b.SignedArea()
	if a < 0 {
		return -a
	}
	return a
}

// SignedArea of the box fulling shape.Area. If the box is well formed the area
// should always be positive.
func (b *Box) SignedArea() float64 {
	v := b.V()
	return v.X * v.Y
}

// Contains returns true of the box contains the point. Fulfills shape.Container
// and shape.Shape.
func (b *Box) Contains(p d2.Pt) bool {
	return b[0].X <= p.X &&
		b[1].X >= p.X &&
		b[0].Y <= p.Y &&
		b[1].Y >= p.Y
}

// Perimeter of the box. Fulfills shape.Perimeter.
func (b *Box) Perimeter() float64 {
	v := b.V()
	return 2 * (v.X + v.Y)
}

// BoundingBox fulfils shape.BoundingBoxer and shape.Shape.
func (b *Box) BoundingBox() (min, max d2.Pt) {
	return b[0], b[1]
}
