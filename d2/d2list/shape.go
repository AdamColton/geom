package d2list

import (
	"errors"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape"
	"github.com/adamcolton/geom/d2/shape/box"
	"github.com/adamcolton/geom/list"
)

type ShapeList list.List[shape.Shape]

type ShapeSlice = list.Slice[shape.Shape]

type ShapeGenerator interface {
	GenerateShape(pts PointList) shape.Shape
}

type ShapeGeneratorList list.List[ShapeGenerator]

type ShapeGeneratorSlice = list.Slice[shape.Shape]

func GenerateShapeSlice(pts PointList, l ShapeGeneratorList) ShapeSlice {
	out := make(ShapeSlice, l.Len())
	for i := range out {
		out[i] = l.Idx(i).GenerateShape(pts)
	}
	return out
}

type ShapePerimeter interface {
	d2.Pt1
	line.Intersector
	shape.BoundingBoxer
}

type Shape struct {
	Curves list.List[ShapePerimeter]
}

func (s Shape) Contains(pt d2.Pt) bool {
	ln := s.Curves.Len()
	count := 0
	ray := line.New(pt, pt.Add(d2.V{1, 1}))
	var intersections []float64
	for i := 0; i < ln; i++ {
		c := s.Curves.Idx(i)
		intersections = c.LineIntersections(ray, intersections[:0])
		for _, x := range intersections {
			ipt := c.Pt1(x)
			t := ipt.X - pt.X
			if t > 0 {
				count++
			}
		}
	}
	return count%2 == 0
}

func (s Shape) BoundingBox() (min, max d2.Pt) {
	ln := s.Curves.Len()
	var out *box.Box
	for i := 0; i < ln; i++ {
		c := s.Curves.Idx(i)
		out = out.Expand(c.BoundingBox())
	}
	return out.BoundingBox()
}

func (s Shape) LineIntersections(l line.Line, buf []float64) []float64 {
	ln := s.Curves.Len()
	buf = buf[:0]
	for i := 0; i < ln; i++ {
		c := s.Curves.Idx(i)
		buf = append(buf, c.LineIntersections(l, buf[len(buf):])...)
	}
	return buf
}

func (s Shape) Validate(t cmpr.Tolerance) error {
	ln := s.Curves.Len()
	if ln == 0 {
		return errors.New("shape contains no curves")
	}
	d := float64(t * t)
	prev := s.Curves.Idx(ln - 1)
	for i := 0; i < ln; i++ {
		c := s.Curves.Idx(ln - 1)
		if c.Pt1(0).Subtract(prev.Pt1(1)).Mag2() > d {
			return errors.New("ends don't line up")
		}
		// TODO: after descent branch, check for intersections
		prev = c
	}
	return nil
}
