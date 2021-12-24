package newton

import (
	"math"

	"github.com/adamcolton/geom/d2"
)

type Distance struct {
	pt0, pt1 func(float64) d2.Pt
	v0, v1   func(float64) d2.V
}

func NewDistance(p0, p1 d2.Pt1) Distance {
	return Distance{
		pt0: p0.Pt1,
		pt1: p1.Pt1,
		v0:  d2.GetV1(p0).V1,
		v1:  d2.GetV1(p1).V1,
	}
}

func (dst Distance) At(t0, t1 float64) float64 {
	return dst.pt0(t0).Subtract(dst.pt1(t1)).Mag2()
}

func (dst Distance) D(t0, t1 float64) (d float64, dt0 float64, dt1 float64) {
	pt0, pt1 := dst.pt0(t0), dst.pt1(t1)
	v0, v1 := dst.v0(t0), dst.v1(t1)
	v := pt0.Subtract(pt1)
	d = v.Mag2()
	dt0 = 2 * (v.X*v0.X + v.Y*v0.Y)
	dt1 = 2 * (v.X*-v1.X + v.Y*-v1.Y)
	return
}

/*
// not taking sqrt - not necessary, simplifies math
D(t0,t1) = (p0X(t0) - p1X(t1))^2 + (p0Y(t0) - p1Y(t1))^2

p0X(t0) = ∑ p0[i].X*t0^i
p1X(t1) = ∑ p1[i].X*t1^i
p0Y(t0) = ∑ p0[i].Y*t0^i
p1Y(t1) = ∑ p1[i].Y*t1^i

p0X'(t0) = ∑ i*p0[i].X*t0^(i-1)
p1X'(t1) = ∑ i*p1[i].X*t1^(i-1)
p0y'(t0) = ∑ i*p0[i].Y*t0^(i-1)
p1y'(t1) = ∑ i*p1[i].Y*t1^(i-1)

d_D/d_t0 (t0,t1) =  2(p0X(t0) - p1X(t1))*p0X'(t0) + 2(p0Y(t0) - p1Y(t1))*p0Y'(t0)
d_D/d_t1 (t0,t1) =  2(p0X(t0) - p1X(t1))*-p1X'(t1) + 2(p0Y(t0) - p1Y(t1))*-p1Y'(t1)

STEP
g(amma) : how far to step in gradient direction
0 = (p0X(t0-g*dt0) - p1X(t1-g*dt1))^2 + (p0Y(t0-g*dt0) - p1Y(t1-g*dt1))^2
We can reduce this. Since we've got
  0 = A^2 + B^2
then
  A == 0 && B == 0

so
0 = p0X(t0-g*dt0) - p1X(t1-g*dt1)
p0X(t0) = ∑ p0[i].X*(t0-g*dt0)^i
p1X(t1) = ∑ p1[i].X*(t1-g*dt1)^i
excpet, this is no longer being written for poly, so just
0 = p0X(t0-g*dt0) - p1X(t1-g*dt1)

t := t0-g*dt0
f := p0x(t)
df/dg = (df/dt)(dt/dg)

maybe that doesn't work because we can hit an t0 --> 0 where t1 does not

*/

func (dst Distance) Min(t0, t1 float64) (float64, float64) {
	found := false
	max := 1000
	var d float64
	for i := 0; !found && i < max && !math.IsNaN(t0) && !math.IsNaN(t1); i++ {
		t0, t1, d = dst.minStep(t0, t1)
		found = d < 1e-5
	}
	return t0, t1
}

func (dst Distance) minStep(t0, t1 float64) (float64, float64, float64) {
	d, dt0, dt1 := dst.D(t0, t1)
	g := dst.g(t0, t1, dt0, dt1)
	t0 -= dt0 * g
	t1 -= dt1 * g
	return t0, t1, d
}

/*
t0, t1, dt0 and dt1 are all fixed in this calculation, only g is moving.
Newton's method is used to find the minimum for g.

d(g) = (p0X(t0-g*dt0) - p1X(t1-g*dt1))^2 + (p0Y(t0-g*dt0) - p1Y(t1-g*dt1))^2
t0(g) = t0-g*dt0
t1(g) = t1-g*dt1
A = p0X(t0) - p1X(t1)
B = p0Y(t0) - p1Y(t1)
v = V{A,B}
d(g) = A^2 + B^2 = v.Mag2()
d'(g) = 2A*A' + 2B*B' = 2(A*A' + B*B')
A'(g) = p0X'*t0' - p1X'*t1'
B'(g) = p0Y'*t0' - p1Y'*t1'
t0'(g) = -dt0
t1'(g) = -dt1
and then just Newtons Method:
g = g - d(g)/d'(g)
*/

type dgfn func(float64) (float64, float64)

func (dst Distance) dg(t0, t1, dt0, dt1 float64) dgfn {
	return func(g float64) (d float64, dd_dg float64) {
		t0g, t1g := t0-g*dt0, t1-g*dt1
		pt0, pt1 := dst.pt0(t0g), dst.pt1(t1g)
		v := pt0.Subtract(pt1)
		dv0, dv1 := dst.v0(t0g).Multiply(-dt0), dst.v1(t1g).Multiply(-dt1)
		dv := dv0.Subtract(dv1)
		d = v.Mag2()
		//fmt.Printf("t0g:%.2e t1g:%.2e d:%.2e\n", t0g, t1g, d)
		dd_dg = 2 * v.Dot(dv)
		return
	}
}

func (dst Distance) g(t0, t1, dt0, dt1 float64) float64 {
	max := 15.0
	g := 0.0
	dg := dst.dg(t0, t1, dt0, dt1)
	bestg := g
	bestd := dst.At(t0, t1)
	var step, d float64
	var done bool
	i := 0.0
	for ; i < max && !done; i++ {
		step, d, done = dg.step(g)
		if bestd > d {
			bestg, bestd = g, d
		}
		g -= step
	}
	return bestg
}

func (dg dgfn) step(g float64) (float64, float64, bool) {
	d, dd_dg := dg(g)
	if dd_dg == 0 {
		return g, d, true
	}
	step := d / dd_dg
	return step, d, false
}
