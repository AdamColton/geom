package d2

import (
	"math"
	"strconv"
	"strings"
)

/*

| a b c |   | x |
| d e f | * | y | = | ax+by+c dx+ey+f gx+hy+i|
| g h i |   | 1 |

| (0,0) (1,0) (2,0) |   | [0][0] [0][1] [0][2] |
| (0,1) (1,1) (2,1) | = | [1][0] [1][1] [1][2] |
| (0,2) (1,2) (2,2) |   | [2][0] [2][1] [2][2] |

| a b c |   | j k l |   | aj+bm+cp ak+bn+cq al+bo+cr |
| d e f | * | m n o | = | dj+em+fp dk+en+fq dl+eo+fr |
| g h i |   | p q r |   | gj+hm+ip gk+hn+iq gl+ho+ir |

*/

type T [3][3]float64

func IndentityTransform() T {
	return T{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
}

func (t T) Pt(pt Pt) (Pt, float64) {
	return Pt{
		pt.X*t[0][0] + pt.Y*t[0][1] + t[0][2],
		pt.X*t[1][0] + pt.Y*t[1][1] + t[1][2],
	}, pt.X*t[2][0] + pt.Y*t[2][1] + t[2][2]
}

func (t T) V(v V) (V, float64) {
	return V{
		v.X*t[0][0] + v.Y*t[0][1] + t[0][2],
		v.X*t[1][0] + v.Y*t[1][1] + t[1][2],
	}, v.X*t[2][0] + v.Y*t[2][1] + t[2][2]
}

func (t T) T(t2 T) T {
	return T{
		{
			t[0][0]*t2[0][0] + t[1][0]*t2[0][1] + t[2][0]*t2[0][2],
			t[0][1]*t2[0][0] + t[1][1]*t2[0][1] + t[2][1]*t2[0][2],
			t[0][2]*t2[0][0] + t[1][2]*t2[0][1] + t[2][2]*t2[0][2],
		}, {
			t[0][0]*t2[1][0] + t[1][0]*t2[1][1] + t[2][0]*t2[1][2],
			t[0][1]*t2[1][0] + t[1][1]*t2[1][1] + t[2][1]*t2[1][2],
			t[0][2]*t2[1][0] + t[1][2]*t2[1][1] + t[2][2]*t2[1][2],
		}, {
			t[0][0]*t2[2][0] + t[1][0]*t2[2][1] + t[2][0]*t2[2][2],
			t[0][1]*t2[2][0] + t[1][1]*t2[2][1] + t[2][1]*t2[2][2],
			t[0][2]*t2[2][0] + t[1][2]*t2[2][1] + t[2][2]*t2[2][2],
		},
	}
}

func Scale(v V) T {
	return T{
		{v.X, 0, 0},
		{0, v.Y, 0},
		{0, 0, 1},
	}
}

func Rotate(angle float64) T {
	s, c := math.Sincos(angle)
	return T{
		{c, -s, 0},
		{s, c, 0},
		{0, 0, 1},
	}
}

func Translate(v V) T {
	return T{
		{1, 0, v.X},
		{0, 1, v.Y},
		{0, 0, 1},
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
		"), (",
		strconv.FormatFloat(t[1][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[1][2], 'f', Prec, 64),
		"), (",
		strconv.FormatFloat(t[2][0], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][1], 'f', Prec, 64),
		", ",
		strconv.FormatFloat(t[2][2], 'f', Prec, 64),
		") ]",
	}, "")
}
