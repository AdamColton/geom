package ellipsearc

import (
	"math"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

// EllipseArc fulfills Path and describes an elliptic arc. Start defines where the
// arc starts and Length defines the arc length, in radians. Start defaults to 0
// which is the point on the arc that intersects the ray from foci-1 to foci-2
type EllipseArc struct {
	Start, Length float64
	c             d2.Pt     // center
	sMa, sma      float64   // semi-major axis, semi-minor axis
	a             angle.Rad // angle
	as, ac        float64   // angle sin, cos
}

// New returns an EllipseArc with foci f1 and f2 and a minor radius of r. The
// perimeter point that corresponds to an angle of 0 will be 1/4 rotation going
// from f1 to f2, which will lie along the minor axis. So an ellipse with foci
// (0,0) and (0,2) with a minor radius of 1 will have angle 0 at point (1,1).
func New(pt1, pt2 d2.Pt, r float64) EllipseArc {
	e := EllipseArc{
		c:      line.New(pt1, pt2).Pt1(0.5),
		Length: math.Pi * 2,
	}
	d := pt2.Subtract(pt1)
	e.a = d.Angle()
	e.as, e.ac = e.a.Sincos()

	e.sma = r
	e.sMa = d2.Pt{d.Mag(), 2 * r}.Mag() / 2

	return e
}

// Pt1 returns the float64 vector at t.
func (e EllipseArc) Pt1(t float64) d2.Pt {
	return e.c.Add(e.ByAngle((e.Length)*t + e.Start))
}

// V1 returns a tangent vector at t.
func (e EllipseArc) V1(t float64) d2.V {
	t = (e.Length)*t + e.Start
	return e.ByAngle(t + math.Pi/2).V().Multiply(2 * math.Pi)
}

// ByAngle returns the vector at the given angle relative to the center
func (e EllipseArc) ByAngle(a float64) d2.V {
	// https://en.wikipedia.org/wiki/Parametric_equation#Ellipse
	s, c := math.Sincos(a)
	return d2.V{
		X: e.sMa*e.ac*c - e.sma*e.as*s,
		Y: e.sMa*e.as*c + e.sma*e.ac*s,
	}
}

// Foci of the ellipse
func (e EllipseArc) Foci() (d2.Pt, d2.Pt) {
	fociLen := math.Sqrt(e.sMa*e.sMa - e.sma*e.sma)
	v1 := d2.V{fociLen * e.ac, fociLen * e.as}
	v2 := d2.V{-v1.X, -v1.Y}
	return e.c.Add(v2), e.c.Add(v1)
}

// Centroid of the ellipse
func (e EllipseArc) Centroid() d2.Pt {
	return e.c
}

// Axis returns the lengths of the major and minor axis of the ellipse
func (e EllipseArc) Axis() (major, minor float64) {
	return e.sMa, e.sma
}

// Angle returns the information about the ellipse rotation angle.
func (e EllipseArc) Angle() (ang angle.Rad, sin, cos float64) {
	return e.a, e.as, e.ac
}

var padding = 1e-10

// BoundingBox fulfills shape.BoundingBoxer. It returns a min and max that are
// the corners of a bounding rectangle that will contain the ellipse arc. This
// ignores the start and end of the arc so it may not be the minimal bounding
// box.
func (e EllipseArc) BoundingBox() (min, max d2.Pt) {
	// https://math.stackexchange.com/questions/91132/how-to-get-the-limits-of-rotated-ellipse
	a2 := e.sMa * e.sMa
	b2 := e.sma * e.sma
	c2 := e.ac * e.ac
	s2 := e.as * e.as

	x := math.Sqrt(a2*c2 + b2*s2)
	y := math.Sqrt(a2*s2 + b2*c2)

	return e.c.Add(d2.V{-x - padding, -y - padding}), e.c.Add(d2.V{x + padding, y + padding})
}

// LineIntersections fulfills line.LineIntersector. Returns the points on the
// line that intersect the EllipseArc relative to the line.
func (e EllipseArc) LineIntersections(l line.Line) []float64 {
	// TODO: this is not correct, works for horizontal or vertical lines
	// but the further off from horizontal it is, the more error there is.
	//
	// http://quickcalcbasic.com/ellipse%20line%20intersection.pdf
	// Intersection of Rotated Ellipse with Sloping Line(s)
	v, h := e.sma, e.sMa
	v2, h2 := v*v, h*h
	s, c := e.as, e.ac
	s2, c2 := s*s, c*c
	cs := c * s
	l.T0.X -= e.c.X
	l.T0.Y -= e.c.Y
	if l.D.X == 0 {
		if l.D.Y == 0 {
			return nil
		}
		x := l.T0.X

		A := v2*s2 + h2*c2
		B := 2 * x * cs * (v2 - h2)
		C := x*x*(v2*c2+h2*s2) - h2*v2

		sqrt := B*B - 4*A*C
		if sqrt < 0 {
			return nil
		}
		sqrt = math.Sqrt(sqrt)

		y0 := (-B + sqrt) / (2 * A)
		y1 := (-B - sqrt) / (2 * A)
		t0 := (y0 - l.T0.Y) / l.D.Y
		t1 := (y1 - l.T0.Y) / l.D.Y
		return []float64{t0, t1}
	}
	m := l.M()
	m2 := m * m
	b := l.B()
	b2 := b * b

	A := v2*(c2+2*m*cs+m2*s2) + h2*(m2*c2-2*m*cs+s2)
	B := b * (2*v2*(cs+m*s2) + 2*h2*(m*c2-cs))
	C := b2*(v2*s2+h2*c2) - h2*v2

	sqrt := B*B - 4*A*C
	if sqrt < 0 {
		return nil
	}
	sqrt = math.Sqrt(sqrt)
	x1 := (-B + sqrt) / (2 * A)
	x2 := (-B - sqrt) / (2 * A)
	t1 := (x1 - l.T0.X) / l.D.X
	t2 := (x2 - l.T0.X) / l.D.X

	// Todo - check if this is inside start/end

	return []float64{t1, t2}
}

// ABC returns the values so the A*x^2 + B*x + C = 0 is true for the ellipse at
// y.
func (e EllipseArc) ABC(y float64) (float64, float64, float64) {
	v, h := e.sma, e.sMa
	v2, h2 := v*v, h*h
	s, c := e.as, e.ac
	s2, c2 := s*s, c*c
	y -= e.c.Y
	return v2*c2 + h2*s2,
		2 * y * c * s * (v2 - h2),
		y*y*(v2*s2+h2*c2) - h2*v2
}
