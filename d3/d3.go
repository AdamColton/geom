package d3

import (
	"math"
	"strconv"
	"strings"
)

type D3 struct {
	X, Y, Z float64
}

func (d D3) Mag2() float64 {
	return d.X*d.X + d.Y*d.Y + d.Z*d.Z
}

func (d D3) Mag() float64 {
	return math.Sqrt(d.Mag2())
}

func (d D3) String() string {
	return strings.Join([]string{
		"D3(",
		strconv.FormatFloat(d.X, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(d.Y, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(d.Z, 'f', Prec, 64),
		")",
	}, "")
}
