package list

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
	"github.com/adamcolton/geom/geomtest"
)

type Slice[T any] []T

type Slicer[T any] interface {
	ToSlice() Slice[T]
}

func (s Slice[T]) Len() int {
	return len(s)
}

func (s Slice[T]) Idx(idx int) (t T) {
	if idx < 0 || idx >= len(s) {
		return
	}
	return s[idx]
}

func (s Slice[T]) ToSlice() Slice[T] {
	return s
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (s Slice[T]) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	l, ok := actual.(List[T])
	if !ok {
		return fmt.Errorf("type %s does not fulfill List", reflect.TypeOf(actual))
	}
	ln := s.Len()
	if _, _, err := geomerr.NewLenMismatch(ln, l.Len()); err != nil {
		return err
	}

	var errs geomerr.SliceErrs
	for i, si := range s {
		var x interface{} = si
		if ta, ok := x.(geomtest.AssertEqualizer); ok {
			errs = errs.Append(i, ta.AssertEqual(l.Idx(i), t))
		}

	}
	return errs.Ret()
}

func NewSlice[T any](list List[T]) Slice[T] {
	if ps, ok := list.(Slicer[T]); ok {
		return ps.ToSlice()
	}
	out := make(Slice[T], list.Len())
	for i := range out {
		out[i] = list.Idx(i)
	}
	return out
}
