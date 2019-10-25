// Package geomtest provides helpers for testing the geom packages.
package geomtest

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d3"
)

const (
	small = 1e-10
	big   = 1.0 / small
)

// Equal can compare many types from various packages and check that they are
// considered equal. This allows for small variations in floating point values
// to be ignored. Supported types are d2.Pt, d2.V, d3.Pt, d3.V
func Equal(t *testing.T, expected, actual interface{}) bool {
	switch a := expected.(type) {
	case d2.Pt:
		b, ok := actual.(d2.Pt)
		if !ok {
			t.Error("Types do not match")
			return false
		}
		v := a.Subtract(b).Abs()
		if v.X > small || v.Y > small {
			t.Errorf("Expected %s got %s", a, b)
			return false
		}
		return true
	case []d2.Pt:
		b, ok := actual.([]d2.Pt)
		if !ok {
			t.Error("Types do not match")
			return false
		}
		if len(a) != len(b) {
			t.Error("Lengths do not match")
			return false
		}
		equal := true
		for i, a := range a {
			v := a.Subtract(b[i]).Abs()
			if v.X > small || v.Y > small {
				t.Errorf("At %d expected %s got %s", i, a, b[i])
				equal = false
			}
		}
		return equal
	case []d3.Pt:
		b, ok := actual.([]d3.Pt)
		if !ok {
			t.Error("Types do not match")
			return false
		}
		if len(a) != len(b) {
			t.Error("Lengths do not match")
			return false
		}
		equal := true
		for i, a := range a {
			v := a.Subtract(b[i]).Abs()
			if v.X > small || v.Y > small || v.Z > small {
				t.Errorf("At %d expected %s got %s", i, a, b[i])
				equal = false
			}
		}
		return equal
	case d2.V:
		b, ok := actual.(d2.V)
		if !ok {
			t.Error("Types do not match")
			return false
		}
		v := d2.V{
			X: a.X - b.X,
			Y: a.Y - b.Y,
		}.Abs()
		if v.X > small || v.Y > small {
			t.Errorf("Expected %s got %s", a, b)
			return false
		}
		return true
	case d3.Pt:
		b, ok := actual.(d3.Pt)
		if !ok {
			t.Error("Types do not match")
			return false
		}
		v := a.Subtract(b).Abs()
		if v.X > small || v.Y > small || v.Z > small {
			t.Errorf("Expected %s got %s", a, b)
			return false
		}
		return true
	case d3.V:
		b, ok := actual.(d3.V)
		if !ok {
			t.Error("Types do not match")
			return false
		}
		v := a.Subtract(b).Abs()
		if v.X > small || v.Y > small || v.Z > small {
			t.Errorf("Expected %s got %s", a, b)
			return false
		}
		return true
	}
	t.Error("Unsupported type")
	return false
}

// V1 checks that the derivative is close the derivative approximation. A good
// base check that an algorithm isn't completely wrong.
func V1(t *testing.T, pv d2.Pt1V1) bool {
	fn := make([]func(float64) d2.V, 2, 3)
	fn[0] = func(t0 float64) d2.V {
		return pv.Pt1(t0 + small).Subtract(pv.Pt1(t0)).Multiply(big)
	}
	fn[1] = pv.V1
	if v1c0, ok := pv.(d2.V1c0); ok {
		fn = append(fn, v1c0.V1c0().V1)
	}

	ok := true
	v := make([]d2.V, len(fn))
	for i := 0.0; i <= 1.0; i += 0.01 {
		for j, f := range fn {
			v[j] = f(i)
			if j > 0 {
				vd := d2.V{
					X: v[0].X - v[j].X,
					Y: v[0].Y - v[j].Y,
				}.Abs()
				if vd.X > 1e-4 || vd.Y > 1e-4 {
					r := d2.V{
						X: v[0].X / v[j].X,
						Y: v[0].Y / v[j].Y,
					}
					t.Errorf("Bad derivative %0.2f %d %s %s %s", i, j, v[0], v[j], r)
					ok = false
				}
			}
		}
	}
	return ok
}
