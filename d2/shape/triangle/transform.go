package triangle

import (
	"github.com/adamcolton/geom/d2"
)

/*
Given triangles A and B, the matrix transformation is:

| k l m |   | Ax |   | Bx |   | k*Ax + l*Ay + m |
| n p q | * | Ay | = | By | = | n*Ax + p*Ay + q |
            | 1  |

---- Finding k,l,m ----
Bx = k*Ax + l*Ay + m
Bx - k*Ax - l*Ay = m
B0x - k*A0x - l*A0y = B1x - k*A1x - l*A1y = B2x - k*A2x - l*A2y

B0x - k*A0x - l*A0y = B1x - k*A1x - l*A1y
k*A1x - k*A0x = B1x - B0x + l*A0y - l*A1y
k(A1x - A0x) = (B1x - B0x) + l(A0y - A1y)
k = (B1x - B0x)/(A1x - A0x) + l(A0y - A1y)/(A1x - A0x)
D1 = (A1x - A0x)
R = (B1x - B0x) / D1
S = (A0y - A1y) / D1
k = R + lS

B1x - k*A1x - l*A1y = B2x - k*A2x - l*A2y
k(A2x - A1x) = (B2x - B1x) + l(A1y - A2y)
l = k(A2x - A1x)/(A1y - A2y) - (B2x - B1x)/(A1y - A2y)
D2 = (A1y - A2y)
T = (A2x - A1x) / D2
U = (B2x - B1x) / D2
l = kT - U
l = RT + lST - U
l - lST = RT - U
l(1 - ST)= RT - U
D3 = 1 - ST
l  = (RT - U) / D3

---- Finding n,p,q ----
By = n*Ax + p*Ay + q
By - n*Ax - p*Ay = q
B0y - n*A0x - p*A0y = B1y - n*A1x - p*A1y = B2y - n*A2x - p*A2y

B0y - n*A0x - p*A0y = B1y - n*A1x - p*A1y
n(A1x - A0x) = (B1y - B0y) + p(A0y - A1y)
n = (B1y - B0y)/(A1x - A0x) + p(A0y - A1y)/(A1x - A0x)
D1 = A1x - A0x
R = (B1y - B0y) / D1
S = (A0y - A1y) / D1
n = R + pS

B1y - n*A1x - p*A1y = B2y - n*A2x - p*A2y
n(A2x - A1x) = (B2y - B1y) + p(A1y - A2y)
p = n(A2x - A1x)/(A1y - A2y) - (B2y - B1y)/(A1y - A2y)
D2 = (A1y - A2y)
T = (A2x - A1x) / D2
U = (B2y - B1y) / D2
p = nT - U
p = RT + pST - U
p - pST = RT - U
p(1 - ST) = RT - U
D3 = 1 - ST
p = (RT - U) / D3
*/

// CouldNotResolveErr is returned if a TriangleTransform cannot resolve one of
// the axis.
type CouldNotResolveErr string

// Error fulfils the error interface
func (c CouldNotResolveErr) Error() string {
	return string(c)
}

// Transform takes 2 triangles and returns a Transformation that when
// applied to the vertexes of A produces the vertex of B.
func Transform(a, b Triangle) (d2.T, error) {
	tfrm := d2.T{}

	// find k, l, m
	var found bool
	var k, l, m float64
	for i0 := 0; !found && i0 < 3; i0++ {
		i1 := (i0 + 1) % 3
		i2 := (i0 + 2) % 3
		d1 := a[i1].X - a[i0].X
		if d1 == 0 {
			continue
		}
		d2 := a[i1].Y - a[i2].Y
		if d2 == 0 {
			continue
		}
		s := (a[i0].Y - a[i1].Y) / d1
		t := (a[i2].X - a[i1].X) / d2
		d3 := 1 - s*t
		if d3 == 0 {
			continue
		}
		found = true

		r := (b[i1].X - b[i0].X) / d1
		u := (b[i2].X - b[i1].X) / d2

		l = (r*t - u) / d3
		k = r + l*s
		m = b[i0].X - k*a[i0].X - l*a[i0].Y
	}

	if !found {
		return tfrm, CouldNotResolveErr("Could not resolve X transform")
	}

	// find n, p, q
	var n, p, q float64
	found = false
	for i0 := 0; !found && i0 < 3; i0++ {
		i1 := (i0 + 1) % 3
		i2 := (i0 + 2) % 3
		d1 := a[i1].X - a[i0].X
		if d1 == 0 {
			continue
		}
		d2 := a[i1].Y - a[i2].Y
		if d2 == 0 {
			continue
		}
		s := (a[i0].Y - a[i1].Y) / d1
		t := (a[i2].X - a[i1].X) / d2
		d3 := 1 - s*t
		if d3 == 0 {
			continue
		}
		found = true

		r := (b[i1].Y - b[i0].Y) / d1
		u := (b[i2].Y - b[i1].Y) / d2

		p = (r*t - u) / d3
		n = r + p*s
		q = b[i0].Y - n*a[i0].X - p*a[i0].Y
	}

	if !found {
		return tfrm, CouldNotResolveErr("Could not resolve Y transform")
	}

	tfrm[0][0] = k
	tfrm[1][0] = n
	tfrm[0][1] = l
	tfrm[1][1] = p
	tfrm[0][2] = m
	tfrm[1][2] = q
	tfrm[2][2] = 1

	return tfrm, nil
}
