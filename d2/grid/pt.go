package grid

import (
	"github.com/adamcolton/geom/d2"
)

type Pt struct {
	X, Y int
}

func (pt Pt) Area() int {
	a := pt.X * pt.Y
	if a < 0 {
		return -a
	}
	return a
}

func (pt Pt) D2() d2.D2 {
	return d2.D2{float64(pt.X), float64(pt.Y)}
}

func (pt Pt) Abs() Pt {
	if pt.X < 0 {
		pt.X = -pt.X
	}
	if pt.Y < 0 {
		pt.Y = -pt.Y
	}
	return pt
}

func (pt Pt) Add(pt2 Pt) Pt {
	return Pt{
		X: pt.X + pt2.X,
		Y: pt.Y + pt2.Y,
	}
}

func (pt Pt) Subtract(pt2 Pt) Pt {
	return Pt{
		X: pt.X - pt2.X,
		Y: pt.Y - pt2.Y,
	}
}

func (pt Pt) Multiply(scale int) Pt {
	return Pt{
		X: pt.X * scale,
		Y: pt.Y * scale,
	}
}

func (pt Pt) To(pt2 Pt) Iterator {
	return Range{pt, pt2}.Iter()
}

func (pt Pt) Iter() Iterator {
	return Pt{}.To(pt)
}

type Scale struct {
	X, Y, DX, DY float64
}

func (s Scale) T(pt Pt) (float64, float64) {
	return float64(pt.X)*s.X + s.DX, float64(pt.Y)*s.Y + s.DY
}

func Origin() Pt{
	return Pt{0,0}
}