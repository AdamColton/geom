package line

import (
	"github.com/adamcolton/geom/d3"
	"strings"
)

type Line struct {
	T0 d3.Pt
	D  d3.V
}

func (l Line) Pt(t float64) d3.Pt {
	return l.T0.Add(l.D.Multiply(t))
}

func (l Line) V(t float64) d3.V {
	return l.D
}

func New(start, end d3.Pt) Line {
	return Line{
		T0: start,
		D:  end.Subtract(start),
	}
}

func (l Line) String() string {
	return strings.Join([]string{
		"Line[ ",
		l.T0.String(),
		" + ",
		l.D.String(),
		"*t ]",
	}, "")
}

/*
given lines A, B the closest point will be perpendicular to both.
If the solutions are m and n, this means
s := Line(A(m), B(n))
s.V.Dot(A.V) == 0 && s.V.Dot(B.V) == 0

s.V.X = A.P.X + A.V.X*m - B.P.X - B.V.X*n
s.V.Y = A.P.Y + A.V.Y*m - B.P.Y - B.V.Y*n
s.V.Z = A.P.Z + A.V.Z*m - B.P.Z - B.V.Z*n

::: s.V.Dot(A.V) == 0 :::
(A.P.X + A.V.X*m - B.P.X - B.V.X*n) * A.V.X +
(A.P.Y + A.V.Y*m - B.P.Y - B.V.Y*n) * A.V.Y +
(A.P.Z + A.V.Z*m - B.P.Z - B.V.Z*n) * A.V.Z = 0

A.P.X*A.V.X + A.V.X*A.V.X*m - B.P.X*A.V.X - B.V.X*A.V.X*n +
A.P.Y*A.V.Y + A.V.Y*A.V.Y*m - B.P.Y*A.V.Y - B.V.Y*A.V.Y*n +
A.P.Z*A.V.Z + A.V.Z*A.V.Z*m - B.P.Z*A.V.Z - B.V.Z*A.V.Z*n = 0

-(B.V.X*A.V.X*n + B.V.Y*A.V.Y*n + B.V.Z*A.V.Z*n) +
(A.V.X*A.V.X*m + A.V.Y*A.V.Y*m + A.V.Z*A.V.Z*m) +
A.P.X*A.V.X - B.P.X*A.V.X +
A.P.Y*A.V.Y - B.P.Y*A.V.Y +
A.P.Z*A.V.Z - B.P.Z*A.V.Z = 0

m * (A.V.X*A.V.X + A.V.Y*A.V.Y + A.V.Z*A.V.Z) +
A.P.X*A.V.X - B.P.X*A.V.X +
A.P.Y*A.V.Y - B.P.Y*A.V.Y +
A.P.Z*A.V.Z - B.P.Z*A.V.Z /
(B.V.X*A.V.X + B.V.Y*A.V.Y + B.V.Z*A.V.Z)
= n

::: s.V.Dot(B.V) == 0 :::
m(A.V.X*B.V.X + A.V.Y*B.V.Y + A.V.Z*B.V.Z) +
A.P.X*B.V.X - B.P.X*B.V.X +
A.P.Y*B.V.Y - B.P.Y*B.V.Y +
A.P.Z*B.V.Z - B.P.Z*B.V.Z /
(B.V.X*B.V.X + B.V.Y*B.V.Y + B.V.Z*B.V.Z)
= n

..:: Sub ::..
M1 := A.V.X*A.V.X + A.V.Y*A.V.Y + A.V.Z*A.V.Z
    = A.V.Mag2()
C1 := A.P.X*A.V.X - B.P.X*A.V.X + A.P.Y*A.V.Y - B.P.Y*A.V.Y + A.P.Z*A.V.Z - B.P.Z*A.V.Z
    = A.V.X(A.P.X - B.P.X) + A.V.Y(A.P.Y - B.P.Y) + A.V.Z(A.P.Z - B.P.Z)
    = A.P.Subtract(B.P).Dot(A.V)
D1 := B.V.X*A.V.X + B.V.Y*A.V.Y + B.V.Z*A.V.Z
    = A.V.Dot(B.V)
M2 := A.V.X*B.V.X + A.V.Y*B.V.Y + A.V.Z*B.V.Z
    = D1
C2 := A.P.X*B.V.X - B.P.X*B.V.X + A.P.Y*B.V.Y - B.P.Y*B.V.Y + A.P.Z*B.V.Z - B.P.Z*B.V.Z
    = B.V.X(A.P.X - B.P.X) + B.V.Y(A.P.Y - B.P.Y + B.V.Z(A.P.Z - B.P.Z)
    = A.P.Subtract(B.P).Dot(B.V)
D2 := B.V.X*B.V.X + B.V.Y*B.V.Y + B.V.Z*B.V.Z
    = B.V.Mag2()

(m*M1 + C1) / D1 = (m*M2 + C2) / D2
m*M1*D2 + C1*D2 = m*M2*D1 + C2*D1
m*(M1*D2 - M2*D1) = C2*D1 - C1*D2
m = (C2*D1 - C1*D2) / (M1*D2 - D1*D1)

m(A.V.X*B.V.X + A.V.Y*B.V.Y + A.V.Z*B.V.Z) +
A.P.X*B.V.X - B.P.X*B.V.X + A.P.Y*B.V.Y - B.P.Y*B.V.Y + A.P.Z*B.V.Z - B.P.Z*B.V.Z /
(B.V.X*B.V.X + B.V.Y*B.V.Y + B.V.Z*B.V.Z)
= n
= (m*D1 + C2) / D2

*/

// Closest returns the parametric points on both lines at the point they pass
// closest to eachother
func (l Line) Closest(l2 Line) (float64, float64) {
	d2 := l2.D.Mag2()
	if d2 == 0 {
		return 0, 0
	}
	m1 := l.D.Mag2()
	d1 := l.D.Dot(l2.D)
	c := l.T0.Subtract(l2.T0)
	c1 := c.Dot(l.D)
	c2 := c.Dot(l2.D)

	d := (m1*d2 - d1*d1)
	if d == 0 {
		return 0, 0
	}
	t0 := (c2*d1 - c1*d2) / d
	t1 := (t0*d1 + c2) / d2
	return t0, t1
}
