package d3

import (
	"math"
	"strconv"
	"strings"

	"github.com/adamcolton/geom/angle"
)

// Q is a quaternion used for rotations. B, C and D correspond to the X, Y and Z
// axis.
type Q struct {
	A, B, C, D float64
}

func QX(ang angle.Rad) Q {
	s, c := (ang / 2.0).Sincos()
	return Q{c, -s, 0, 0}
}

func QY(ang angle.Rad) Q {
	s, c := (ang / 2.0).Sincos()
	return Q{c, 0, -s, 0}
}

func QZ(ang angle.Rad) Q {
	s, c := (ang / 2.0).Sincos()
	return Q{c, 0, 0, -s}
}

func QV(v V) Q {
	s, c := (angle.Atan(v.Y, v.X) / 2.0).Sincos()
	qy := Q{c, 0, 0, -s}
	if v.Y < 1e-5 && v.Y > -1e-1 {
		qy = Q{1, 0, 0, 0}
	}
	v = qy.TInv().V(v)
	s, c = (angle.Atan(v.Z, v.X) / 2.0).Sincos()
	qz := Q{c, 0, s, 0}
	if v.Z < 1e-5 && v.Z > -1e-1 {
		qz = Q{1, 0, 0, 0}
	}
	out := qz.Product(qy)
	return out
}

func (q Q) Normalize() Q {
	d := q.A*q.A + q.B*q.B + q.C*q.C + q.D*q.D
	if d == 0 {
		return Q{}
	}
	d = math.Sqrt(d)
	return Q{
		A: q.A / d,
		B: q.B / d,
		C: q.C / d,
		D: q.D / d,
	}
}

func (q Q) Product(q2 Q) Q {
	return Q{
		A: q.A*q2.A - q.B*q2.B - q.C*q2.C - q.D*q2.D,
		B: q.A*q2.B + q.B*q2.A + q.C*q2.D - q.D*q2.C,
		C: q.A*q2.C - q.B*q2.D + q.C*q2.A + q.D*q2.B,
		D: q.A*q2.D + q.B*q2.C - q.C*q2.B + q.D*q2.A,
	}
}

func (q Q) T() *T {
	return &T{
		{
			1 - 2*q.C*q.C - 2*q.D*q.D,
			2*q.B*q.C + 2*q.A*q.D,
			2*q.B*q.D - 2*q.A*q.C,
			0,
		}, {
			2*q.B*q.C - 2*q.A*q.D,
			1 - 2*q.B*q.B - 2*q.D*q.D,
			2*q.C*q.D + 2*q.A*q.B,
			0,
		}, {
			2*q.B*q.D + 2*q.A*q.C,
			2*q.C*q.D - 2*q.A*q.B,
			1 - 2*q.B*q.B - 2*q.C*q.C,
			0,
		}, {
			0,
			0,
			0,
			1,
		},
	}
}

func (q Q) TInv() *T {
	return Q{q.A, -q.B, -q.C, -q.D}.T()
}

func (q Q) String() string {
	return strings.Join([]string{
		"Q(",
		strconv.FormatFloat(q.A, 'f', Prec, 64),
		" + ",
		strconv.FormatFloat(q.B, 'f', Prec, 64),
		"i + ",
		strconv.FormatFloat(q.C, 'f', Prec, 64),
		"j + ",
		strconv.FormatFloat(q.D, 'f', Prec, 64),
		"k)",
	}, "")
}
