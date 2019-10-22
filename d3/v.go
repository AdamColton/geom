package d3

import (
	"math"
	"strconv"
	"strings"
)

type V D3

// Mag2 Returns the sqaure of the magnitude of the vector.
func (v V) Mag2() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v V) V() V { return v }

// Mag returns the magnitude (distance to origin) of the vector
func (v V) Mag() float64 {
	return math.Sqrt(v.Mag2())
}

func (v V) Add(v2 V) V {
	return V{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v V) Subtract(v2 V) V {
	return V{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

// Cross returns the cross product of the two vectors
func (v V) Cross(v2 V) V {
	return V{
		v.Y*v2.Z - v.Z*v2.Y,
		v.Z*v2.X - v.X*v2.Z,
		v.X*v2.Y - v.Y*v2.X,
	}
}

// Dot returns the dot product of two vectors
func (v V) Dot(v2 V) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

// Returns the projection of v2 onto v.
func (v V) Project(v2 V) V {
	return v.Multiply(v.Dot(v2))
}

func (v V) Multiply(scale float64) V {
	return V{
		v.X * scale,
		v.Y * scale,
		v.Z * scale,
	}
}

func (v V) Normal() V {
	m := v.Mag()
	if m == 0 {
		return v
	}
	return v.Multiply(1 / m)
}

func (v V) Abs() V {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	if v.Z < 0 {
		v.Z = -v.Z
	}
	return v
}

// String fulfills Stringer, returns the vector as "(X, Y)"
func (v V) String() string {
	return strings.Join([]string{
		"V(",
		strconv.FormatFloat(v.X, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(v.Y, 'f', Prec, 64),
		", ",
		strconv.FormatFloat(v.Z, 'f', Prec, 64),
		")",
	}, "")
}
