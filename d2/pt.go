package d2

import (
	"strconv"
	"strings"
)

type Point interface {
	Pt() Pt
}

type Pt struct {
	X, Y float64
}

func (pt Pt) Subtract(pt2 Pt) V {
	return V{
		pt.X - pt2.X,
		pt.Y - pt2.Y,
	}
}

func (pt Pt) Add(v V) Pt {
	return Pt{
		pt.X + v.X,
		pt.Y + v.Y,
	}
}

func (pt Pt) Pt() Pt {
	return pt
}

// Distance returns the distance between to points
func (pt Pt) Distance(pt2 Pt) float64 {
	return pt.Subtract(pt2).Mag()
}

func (pt Pt) Multiply(scale float64) Pt {
	return Pt{pt.X * scale, pt.Y * scale}
}

// Prec is the precision for the String method on F
var Prec = 4

// String fulfills Stringer, returns the vector as "(X, Y)"
func (pt Pt) String() string {
	return strings.Join([]string{
		"Pt(",
		strconv.FormatFloat(pt.X, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(pt.Y, 'f', Prec, 64),
		")",
	}, "")
}
