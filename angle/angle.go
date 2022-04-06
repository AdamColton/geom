package angle

import (
	"math"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
)

const (
	Tau Rad = math.Pi * 2
)

// Rad is angle in radians
type Rad float64

// Deg is angle in degrees
func Deg(d float64) Rad {
	return Rad(d * math.Pi / 180)
}

// Rot is angle as percent of a rotation
func Rot(r float64) Rad {
	return Rad(r * 2 * math.Pi)
}

// Atan returns the angle in radians formed by x and y.
func Atan(y, x float64) Rad {
	return Rad(math.Atan2(y, x))
}

// Acos returns the arccosine in radians of x.
func Acos(x float64) Rad {
	return Rad(math.Acos(x))
}

// Rad returns the angle as float64 reprenting radians
func (r Rad) Rad() float64 {
	return float64(r)
}

// Deg returns the angle as float64 reprenting degrees
func (r Rad) Deg() float64 {
	return float64(r) * 180 / math.Pi
}

// Rot returns the angle as float64 reprenting percent rotations
func (r Rad) Rot() float64 {
	return float64(r) / (2 * math.Pi)
}

// Sin returns the sine of the angle
func (r Rad) Sin() float64 {
	return math.Sin(float64(r))
}

// Cos returns the cosine of the angle
func (r Rad) Cos() float64 {
	return math.Cos(float64(r))
}

// Sincos returns both the sine and cosine of the angle
func (r Rad) Sincos() (float64, float64) {
	return math.Sincos(float64(r))
}

// Normal returns a value between 0 and 2*Pi.
func (r Rad) Normal() Rad {
	if r <= -Tau || r >= Tau {
		x := float64(r) / float64(Tau)
		x -= math.Floor(x)
		r = Rad(x) * Tau
	}
	if r < 0 {
		r += Tau
	}
	return r
}

// AssertEqual fulfils geomtest.AssertEqualizer. It checks that the normals are
// equal. This means that Rad(1) == Rad(1+2*Pi).
func (r Rad) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	if err := geomerr.NewTypeMismatch(r, actual); err != nil {
		return err
	}
	r2 := actual.(Rad).Normal()
	r = r.Normal()

	if !t.Zero(float64(r - r2)) {
		return geomerr.NotEqual(r, r2)
	}
	return nil
}
