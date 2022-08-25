package shape

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/box"
)

var DefaultNewBoundingBoxSteps = 8

func NewBoundingBox(s Shape) Rebound {
	r := Rebound{
		Shape: s,
	}

	m, M, ok := TightBoundingBox(s, DefaultNewBoundingBoxSteps)
	if ok {
		r.Box = &box.Box{m, M}
	}
	return r
}

type Rebound struct {
	Shape
	*box.Box
}

func (rs Rebound) BoundingBox() (min, max d2.Pt) {
	if rs.Box != nil {
		return rs.Box.BoundingBox()
	}
	return rs.Shape.BoundingBox()
}

func (rs Rebound) Contains(pt d2.Pt) bool {
	return rs.Shape.Contains(pt)
}

func (rs Rebound) LineIntersections(l line.Line, buf []float64) []float64 {
	return rs.Shape.LineIntersections(l, buf)
}

func TightBoundingBox(s Shape, steps int) (min, max d2.Pt, found bool) {
	m, M := s.BoundingBox()
	op := &tightBoxOp{
		outer: box.New(m, M),
		steps: steps,
		shape: s,
	}
	op.initScan()
	found = op.inner != nil
	if !found {
		return
	}
	for i := 0; i < steps; i++ {
		op.expand()
	}
	min, max = op.outer.BoundingBox()
	return
}

type tightBoxOp struct {
	outer *box.Box
	inner *box.Box
	steps int
	shape Shape
}

func (op *tightBoxOp) expand() {
	mid := &box.Box{
		line.New(op.outer[0], op.inner[0]).Pt1(0.5),
		line.New(op.outer[1], op.inner[1]).Pt1(0.5),
	}
	var ts []float64
	for i, s := range mid.Sides() {
		ts = op.shape.LineIntersections(s, ts[:0])
		if len(ts) > 0 {
			op.linePoints(s, ts)
		} else {
			op.clip(i, s)
		}
	}
}

func (op *tightBoxOp) clip(i int, l line.Line) {
	switch i {
	case 0:
		op.outer[0].Y = l.T0.Y
	case 1:
		op.outer[1].X = l.T0.X
	case 2:
		op.outer[1].Y = l.T0.Y
	case 3:
		op.outer[0].X = l.T0.X
	}
}

func (op *tightBoxOp) initScan() {
	d := op.outer.V()
	scanx := line.Line{
		T0: op.outer[0],
		D:  d2.V{d.X, 0},
	}
	scany := line.Line{
		T0: op.outer[0],
		D:  d2.V{0, d.Y},
	}
	x := line.Line{
		D: scany.D,
	}
	y := line.Line{
		D: scanx.D,
	}
	for i := 0; i < op.steps; i++ {
		t := subdiv(i)
		x.T0 = scanx.Pt1(t)
		if op.linePoints(x, op.shape.LineIntersections(x, nil)) {
			return
		}
		y.T0 = scany.Pt1(t)
		if op.linePoints(y, op.shape.LineIntersections(y, nil)) {
			return
		}
	}
}

func (op *tightBoxOp) linePoints(l line.Line, ts []float64) bool {
	if len(ts) == 0 {
		return false
	}
	if op.inner == nil {
		op.inner = box.New(l.Pt1(ts[0]))
	} else {
		op.inner.Add(l.Pt1(ts[0]))
	}
	for _, t := range ts[1:] {
		op.inner.Add(l.Pt1(t))
	}
	return true
}

var subdivMemo = []float64{0.5, 0, 0, 0, 0, 0, 0, 0}

// subdiv returns a fraction A/B such that
// A < B && A % 0 ==1
// B == 2^x (where x is some positive integer)
// This allows scanning the range 0:1
// so that each successive pass works out from the center
// but never overlaps a previous pass
// The first few values will be: 1/2, 1/4, 3/4, 3/8, 5/8, 1/8, 7/8, 7/16, 9/16,
// 5/16, 11/16, 3/16, 13/16, 1/16, 15/16
func subdiv(n int) float64 {
	ln := len(subdivMemo)
	if n >= ln {
		for ; n >= ln; ln *= 2 {
		}
		cp := make([]float64, ln)
		copy(cp, subdivMemo)
		subdivMemo = cp
	}
	if subdivMemo[n] > 0 {
		return subdivMemo[n]
	}
	n++
	base := subdivbase(n)
	num := base / 2
	idx := n - num
	if idx%2 == 0 {
		num -= idx + 1
	} else {
		num += idx
	}
	out := float64(num) / float64(base)
	subdivMemo[n-1] = out
	return out
}

var subdivbasememo = make([]int, 8)

// subdivbase is only called by subdiv
// it finds the denominator of the fraction
// This will be the smallest 2^x greater than n.
func subdivbase(n int) int {
	ln := len(subdivbasememo)
	if n > ln {
		for ; n > ln; ln *= 2 {
		}
		cp := make([]int, ln)
		copy(cp, subdivbasememo)
		subdivbasememo = cp
	}
	if subdivbasememo[n-1] > 0 {
		return subdivbasememo[n-1]
	}
	base := 2
	for {
		if n < base {
			subdivbasememo[n-1] = base
			return base
		}
		base <<= 1
	}
}
