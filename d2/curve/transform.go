package curve

import "github.com/adamcolton/geom/d2"

type TransformCurveWrapper struct {
	Curve d2.Pt1
	T     *d2.T
}

func (tcw TransformCurveWrapper) Pt1(t float64) d2.Pt {
	return tcw.T.Pt(tcw.Curve.Pt1(t))
}
