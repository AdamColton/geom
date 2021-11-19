package distance

import (
	"math"

	"github.com/adamcolton/geom/d2"
)

type Fn func(t0, t1 float64) float64

type DFn func(t0, t1 float64) (float64, float64)

// this should probably get moved to curve
func New(a, b d2.Pt1) (Fn, DFn) {
	var (
		fn = func(t0, t1 float64) float64 {
			return a.Pt1(t0).Distance(b.Pt1(t1))
		}
	)

	return fn, NewDFn(a, b)
}

func NewDFn(a, b d2.Pt1) DFn {
	da := d2.GetV1(a)
	db := d2.GetV1(b)
	return func(t0, t1 float64) (float64, float64) {
		// d = ( (bx-lx)^2 + (by-ly)^2 )^0.5
		// f(g) = g^0.5 :. f'(g) = ((g^-0.5)/2) * g'
		// g(h,i) = h^2 + i^2 :. g'(h,i) = 2h*h' + 2i*i'
		// h(bx) = bx-lx :. h' = bx'
		// i(by) = by-ly :. i' = by'
		// dbx/dt0 = db(t0).x
		// dby/dt0 = db(t0).y
		// d = f(g(h(bx(t0)),i(by(t0))))

		da_t0 := da.V1(t0)
		hi := a.Pt1(t0).Subtract(b.Pt1(t1))
		hi2 := hi.Multiply(2)
		dg_t0 := hi2.X*da_t0.X + hi2.Y*da_t0.Y
		g := hi.X*hi.X + hi.Y*hi.Y
		g = (math.Pow(g, -0.5) / 2)

		db_t1 := db.V1(t1)
		dg_t1 := -hi2.X*db_t1.X - hi2.Y*db_t1.Y

		return g * dg_t0, g * dg_t1
	}
}

func New2(a, b d2.Pt1) (Fn, DFn) {
	var (
		fn = func(t0, t1 float64) float64 {
			d := a.Pt1(t0).Subtract(b.Pt1(t1))
			return d.X*d.X + d.Y*d.Y
		}
	)

	return fn, NewDFn2(a, b)
}

func NewDFn2(a, b d2.Pt1) DFn {
	da := d2.GetV1(a)
	db := d2.GetV1(b)
	return func(t0, t1 float64) (float64, float64) {
		// d = (bx-lx)^2 + (by-ly)^2

		// g(h,i) = h^2 + i^2 :. g'(h,i) = 2h*h' + 2i*i'
		// h(bx) = bx-lx :. h' = bx'
		// i(by) = by-ly :. i' = by'
		// dbx/dt0 = db(t0).x
		// dby/dt0 = db(t0).y
		// d = f(g(h(bx(t0)),i(by(t0))))

		da_t0 := da.V1(t0)
		hi := a.Pt1(t0).Subtract(b.Pt1(t1))
		hi2 := hi.Multiply(2)
		dg_t0 := hi2.X*da_t0.X + hi2.Y*da_t0.Y

		db_t1 := db.V1(t1)
		dg_t1 := -hi2.X*db_t1.X - hi2.Y*db_t1.Y

		return dg_t0, dg_t1
	}
}
