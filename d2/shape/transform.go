package shape

import (
	"reflect"

	"github.com/adamcolton/geom/d2"
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
