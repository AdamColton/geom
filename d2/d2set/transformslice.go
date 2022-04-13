package d2set

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomerr"
)

type TransformSlicer interface {
	ToTransformSlice() TransformSlice
}

type TransformSlice []*d2.T

func (ts TransformSlice) Len() int {
	return len(ts)
}

func (ts TransformSlice) Get(n int) *d2.T {
	return ts[n]
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (ts TransformSlice) AssertEqual(actual interface{}, tol cmpr.Tolerance) error {
	lst, ok := actual.(TransformList)
	if !ok {
		return fmt.Errorf("Type %s does not fulfill PointSet", reflect.TypeOf(actual))
	}
	ln := ts.Len()
	if _, _, err := geomerr.NewLenMismatch(ln, lst.Len()); err != nil {
		return err
	}

	var errs geomerr.SliceErrs
	for i, t := range ts {
		errs = errs.Append(i, t.AssertEqual(lst.Get(i), tol))
	}
	return errs.Ret()
}

func NewTransformSlice(tlst TransformList) TransformSlice {
	if ts, ok := tlst.(TransformSlicer); ok {
		return ts.ToTransformSlice()
	}
	out := make(TransformSlice, tlst.Len())
	for i := range out {
		out[i] = tlst.Get(i)
	}
	return out
}
