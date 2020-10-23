package zbuf

import (
	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3/shape/triangle"
)

/*

Step = [
m n
p 0
]

m*U.X + n*V.X = dx
m*U.Y + n*V.Y = 0
p*U.Y = 1 :. p = 1/U.Y

m = (dx - n*V.X) / U.X

[(dx - n*V.X) / U.X]*U.Y + n*V.Y = 0
r = U.Y/U.X
dx*r - n*r*V.X + n*V.Y = 0
n(V.Y - r*V.X) = -dx*r
n = (-dx*r) / (V.Y - r*V.X)

? U.X == 0 :
n*V.X = dx
m*U.Y + n*V.Y = 0

n = dx/V.X [ V.X == 0 --> trimngle is vertical line]

m*U.Y + n*V.Y = 0
m = (-n*V.Y)/ U.Y
*/

func Scan(t *triangle.Triangle, dx, dy float64) (*barycentric.BIterator, *triangle.BT) {
	bi := scanU(t)
	bt := t.BT(bi.Origin, bi.U)
	if bt == nil || bt.U.Y == 0 {
		return nil, bt // triangle is horizontal line
	}
	bi.Step[1] = barycentric.B{U: dy / bt.U.Y, V: 0}

	c := bt.U.Cross(bt.V)
	if c.Z > 0 {
		dx = -dx
	}

	var m, n float64
	if bt.U.X == 0 {
		if bt.V.X == 0 {
			return nil, bt // triangle is vertical line
		}
		n = dx / bt.V.X
		m = (-n * bt.V.Y) / bt.U.Y
	} else {
		r := bt.U.Y / bt.U.X
		d := bt.V.Y - r*bt.V.X
		if d == 0 {
			return nil, bt // ???
		}
		n = (-dx * r) / d
		m = (dx - n*bt.V.X) / bt.U.X
	}

	bi.Step[0] = barycentric.B{U: m, V: n}

	return bi, bt
}

func scanU(t *triangle.Triangle) *barycentric.BIterator {
	bi := &barycentric.BIterator{}

	// Choose Origin and U so that they span the height of the triangle
	// So Origin has the lowest Y and U has the highest Y
	for i, p := range t[1:] {
		if p.Y < t[bi.Origin].Y || (p.Y == t[bi.Origin].Y && p.X < t[bi.Origin].X) {
			bi.Origin = i + 1
		}
		if p.Y > t[bi.U].Y || (p.Y == t[bi.U].Y && p.X > t[bi.U].X) {
			bi.U = i + 1
		}
	}
	return bi
}
