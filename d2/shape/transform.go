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

func Transform(s Shape, t *d2.T) Shape {
	if ts, ok := s.(TransformShaper); ok {
		return ts.TransformShape(t)
	}
	v := reflect.ValueOf(s)
	m := v.MethodByName("T")
	mt := m.Type()
	if mt.Kind() == reflect.Func && mt.NumIn() == 1 && mt.In(0) == ttype && mt.NumOut() == 1 && mt.Out(0).Implements(shapeType) {
		return m.Call([]reflect.Value{reflect.ValueOf(t)})[0].Interface().(Shape)
	}
	return nil
}

type TransformShapeWrapper struct {
	Shape
	d2.Pair
}

func NewTransformShapeWrapper(s Shape, t d2.TGen) TransformShapeWrapper {
	return TransformShapeWrapper{
		Shape: s,
		Pair:  t.Pair(),
	}
}

func (tsw TransformShapeWrapper) Contains(pt d2.Pt) bool {
	return tsw.Shape.Contains(tsw.Pair[1].Pt(pt))
}

func (tsw TransformShapeWrapper) LineIntersections(l line.Line, buf []float64) []float64 {
	return tsw.Shape.LineIntersections(l.T(tsw.Pair[1]), buf)
}

func (tsw TransformShapeWrapper) ConvexHull() []d2.Pt {
	return tsw.Pair[0].Slice(tsw.Shape.ConvexHull())
}
