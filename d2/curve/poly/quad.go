package poly

import (
	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/newton"
)

/*

> Quad for 1D
f = a*x2 + b*x + c
f' = 2a*x + b
f' = 0 @
2a*x + b = 0
x = -b/2a
this is the middle, so
m = -b/2a
f(m) = a*m2+b*m+c
a*m2 = b2/2a
b*m = -b2/2a
f(m) = b2/2a -b2/2a
// ha!
f(m) = c



{ax*t0^2, ay*t0^2} + {bx*t0, by*t0} + {cx, cy} = {0 , 0}

*/

const small cmpr.Tolerance = 1e-6

func Quad(a, b, c d2.V, buf []float64) []float64 {
	if buf != nil {
		buf = buf[:0]
	}
	// if small.Zero(a.X) {
	// 	t0, t1 := quad(a.Y, b.Y, c.Y)

	// 	x := b.X*t0 + c.X
	// 	if small.Zero(x) {
	// 		buf = append(buf, t0)
	// 	}

	// 	if !small.Zero(t0 - t1) {
	// 		x = b.X*t1 + c.X
	// 		if small.Zero(x) {
	// 			buf = append(buf, t1)
	// 		}
	// 	}
	// 	return buf
	// }

	// t0, t1 := quad(a.X, b.X, c.X)
	// y := a.Y*t0*t0 + b.Y*t0 + c.Y
	// if small.Zero(y) {
	// 	buf = append(buf, t0)
	// }

	// if !small.Zero(t0 - t1) {
	// 	y = a.Y*t1*t1 + b.Y*t1 + c.Y
	// 	if small.Zero(y) {
	// 		buf = append(buf, t1)
	// 	}
	// }

	return buf
}

func Newton(p1, p2 Poly) (t0, t1 float64) {
	return newton.NewDistance(p1, p2).Min(0, 0)
}
