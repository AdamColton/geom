package d2set

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomerr"
)

type PointSlicer interface {
	ToPointSlice() PointSlice
}

type PointSlice []d2.Pt

func (ps PointSlice) Len() int {
	return len(ps)
}

func (ps PointSlice) Get(n int) d2.Pt {
	return ps[n]
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (ps PointSlice) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	set, ok := actual.(PointList)
	if !ok {
		return fmt.Errorf("Type %s does not fulfill PointSet", reflect.TypeOf(actual))
	}
	ln := ps.Len()
	if _, _, err := geomerr.NewLenMismatch(ln, set.Len()); err != nil {
		return err
	}

	var errs geomerr.SliceErrs
	for i, p := range ps {
		errs = errs.Append(i, p.AssertEqual(set.Get(i), t))
	}
	return errs.Ret()
}

func NewPointSlice(ps PointList) PointSlice {
	if ps, ok := ps.(PointSlicer); ok {
		return ps.ToPointSlice()
	}
	out := make(PointSlice, ps.Len())
	for i := range out {
		out[i] = ps.Get(i)
	}
	return out
}
