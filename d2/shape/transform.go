package shape

import (
	"reflect"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

var (
	ttype     = reflect.TypeOf((*d2.T)(nil))
	shapeType = reflect.TypeOf((*Shape)(nil)).Elem()
)

// Transform a shape. If the shape fulfills TransformShaper, that will be
// invoked. If it follow the pattern of having a method named T that returns a
// new copy of the underlying shape, that will be invoked. The final fallback
// is to create a TransformShapeWrapper.
func Transform(s Shape, t d2.TGen) Shape {
	if ts, ok := s.(TransformShaper); ok {
		return ts.TransformShape(t.GetT())
	}
	v := reflect.ValueOf(s)
	m := v.MethodByName("T")
	mt := m.Type()
	if mt.Kind() == reflect.Func && mt.NumIn() == 1 && mt.In(0) == ttype && mt.NumOut() == 1 && mt.Out(0).Implements(shapeType) {
		return m.Call([]reflect.Value{reflect.ValueOf(t.GetT())})[0].Interface().(Shape)
	}
	return TransformShapeWrapper{
		P:     t.Pair(),
		Shape: s,
	}
}

// TransformShapeWrapper fulfills Shape and applies a transform to a Shape.
type TransformShapeWrapper struct {
	P d2.Pair
	Shape
}

// Contains returns true if the transform applied to the underlying shape
// contains the point.
func (tsw TransformShapeWrapper) Contains(pt d2.Pt) bool {
	return tsw.Shape.Contains(tsw.P[1].Pt(pt))
}

// LineIntersections returns the parametric values relative to the line that
// intersect the transformed shape.
func (tsw TransformShapeWrapper) LineIntersections(l line.Line, buf []float64) []float64 {
	return tsw.Shape.LineIntersections(l.T(tsw.P[1]), buf)
}

// ConvexHull of the transformed shape.
func (tsw TransformShapeWrapper) ConvexHull() []d2.Pt {
	return tsw.P[0].Slice(tsw.Shape.ConvexHull())
}
