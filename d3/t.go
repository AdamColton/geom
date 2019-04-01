package d3

import (
	"strconv"
	"strings"
)

type T [4][4]float64

/*

| a b c d |   | x |
| e f g h | * | y | = | ax+by+cz+d ex+fy+gz+h ix+jy+kz+l mx+ny+pz+q |
| i j k l |   | z |
| m n p q |   | 1 |

| (0,0) (1,0) (2,0) |   | [0][0] [0][1] [0][2] |
| (0,1) (1,1) (2,1) | = | [1][0] [1][1] [1][2] |
| (0,2) (1,2) (2,2) |   | [2][0] [2][1] [2][2] |

| a b c |   | j k l |   | aj+bm+cp ak+bn+cq al+bo+cr |
| d e f | * | m n o | = | dj+em+fp dk+en+fq dl+eo+fr |
| g h i |   | p q r |   | gj+hm+ip gk+hn+iq gl+ho+ir |


| a1 b1 c1 d1 |   | a2 b2 c2 d2 |   | a1a2+b1e2+c1i2+d1m2 a1b2+b1f2+c1j2+d1n2 a1c2+b1g2+c1k2+d1p2 |
| e1 f1 g1 h1 | * | e2 f2 g2 h2 | = |
| i1 j1 k1 l1 |   | i2 j2 k2 l2 |   |
| m1 n1 p1 q1 |   | m2 n2 p2 q2 |   |
*/

func IndentityTransform() T {
	return T{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (t T) Pt(pt Pt) (Pt, float64) {
	return Pt{
		pt.X*t[0][0] + pt.Y*t[0][1] + pt.Z*t[0][2] + t[0][3],
		pt.X*t[1][0] + pt.Y*t[1][1] + pt.Z*t[1][2] + t[1][3],
		pt.X*t[2][0] + pt.Y*t[2][1] + pt.Z*t[2][2] + t[2][3],
	}, pt.X*t[3][0] + pt.Y*t[3][1] + pt.Z*t[3][2] + t[3][3]
}

func (t T) V(v V) (V, float64) {
	return V{
		v.X*t[0][0] + v.Y*t[0][1] + v.Z*t[0][2] + t[0][3],
		v.X*t[1][0] + v.Y*t[1][1] + v.Z*t[1][2] + t[1][3],
		v.X*t[2][0] + v.Y*t[2][1] + v.Z*t[2][2] + t[2][3],
	}, v.X*t[3][0] + v.Y*t[3][1] + v.Z*t[3][2] + t[3][3]
}

func (t T) T(t2 T) T {
	return T{
		{
			t[0][0]*t2[0][0] + t[1][0]*t2[0][1] + t[2][0]*t2[0][2] + t[3][0]*t2[0][3],
			t[0][1]*t2[0][0] + t[1][1]*t2[0][1] + t[2][1]*t2[0][2] + t[3][1]*t2[0][3],
			t[0][2]*t2[0][0] + t[1][2]*t2[0][1] + t[2][2]*t2[0][2] + t[3][2]*t2[0][3],
			t[0][3]*t2[0][0] + t[1][3]*t2[0][1] + t[2][3]*t2[0][2] + t[3][3]*t2[0][3],
		}, {
			t[0][0]*t2[1][0] + t[1][0]*t2[1][1] + t[2][0]*t2[1][2] + t[3][0]*t2[1][3],
			t[0][1]*t2[1][0] + t[1][1]*t2[1][1] + t[2][1]*t2[1][2] + t[3][1]*t2[1][3],
			t[0][2]*t2[1][0] + t[1][2]*t2[1][1] + t[2][2]*t2[1][2] + t[3][2]*t2[1][3],
			t[0][3]*t2[1][0] + t[1][3]*t2[1][1] + t[2][2]*t2[1][2] + t[3][3]*t2[1][3],
		}, {
			t[0][0]*t2[2][0] + t[1][0]*t2[2][1] + t[2][0]*t2[2][2] + t[3][0]*t2[2][3],
			t[0][1]*t2[2][0] + t[1][1]*t2[2][1] + t[2][1]*t2[2][2] + t[3][1]*t2[2][3],
			t[0][2]*t2[2][0] + t[1][2]*t2[2][1] + t[2][2]*t2[2][2] + t[3][2]*t2[2][3],
			t[0][3]*t2[2][0] + t[1][3]*t2[2][1] + t[2][3]*t2[2][2] + t[3][3]*t2[2][3],
		}, {
			t[0][0]*t2[3][0] + t[1][0]*t2[3][1] + t[2][0]*t2[3][2] + t[3][0]*t2[3][3],
			t[0][1]*t2[3][0] + t[1][1]*t2[3][1] + t[2][1]*t2[3][2] + t[3][1]*t2[3][3],
			t[0][2]*t2[3][0] + t[1][2]*t2[3][1] + t[2][2]*t2[3][2] + t[3][2]*t2[3][3],
			t[0][3]*t2[3][0] + t[1][3]*t2[3][1] + t[2][3]*t2[3][2] + t[3][3]*t2[3][3],
		},
	}
}

func Scale(v V) T {
	return T{
		{v.X, 0, 0, 0},
		{0, v.Y, 0, 0},
		{0, 0, v.Z, 0},
		{0, 0, 0, 1},
	}
}

func Translate(v V) T {
	return T{
		{1, 0, 0, v.X},
		{0, 1, 0, v.Y},
		{0, 0, 1, v.Z},
		{0, 0, 0, 1},
	}
}

func (t T) String() string {
	return strings.Join([]string{
		"T[ (",
		strconv.FormatFloat(t[0][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[0][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[0][2], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[0][3], 'f', Prec, 64),
		"), (",
		strconv.FormatFloat(t[1][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][2], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][3], 'f', Prec, 64),
		"), (",
		strconv.FormatFloat(t[2][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][2], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][3], 'f', Prec, 64),
		"), (",
		strconv.FormatFloat(t[3][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[3][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[3][2], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[3][3], 'f', Prec, 64),
		") ]",
	}, "")
}
