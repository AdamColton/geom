package ellipse

import (
	"math"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/ellipsearc"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/triangle"
)

// Ellipse fulfills Shape
type Ellipse struct {
	perimeter *ellipsearc.EllipseArc
}

// New returns an Ellipse with foci f1 and f2 and a minor radius of r.
// The perimeter point that corresponds to an angle of 0 will be 1/4 rotation
// going from f1 to f2, which will lie along the minor axis. So an ellipse with
// foci (0,0) and (0,2) with a minor radius of 1 will have angle 0 at point
// (1,1).
func New(f1, f2 d2.Pt, r float64) Ellipse {
	return Ellipse{
		perimeter: ellipsearc.New(f1, f2, r),
	}
}

// FillCurve a curve that lies inside of the ellipse. It is returned by
// Ellipse.Pt2c1.
type FillCurve struct {
	*d2.T
	*ellipsearc.EllipseArc
}

// Pt1 fulfills d2.Pt1. Given a parametric value it returns a Pt on the curve.
func (fc FillCurve) Pt1(t0 float64) d2.Pt {
	t0 = (t0 / 4.0) + (1.0 / 8.0)
	return fc.T.Pt(fc.EllipseArc.Pt1(t0))
}

func (e Ellipse) fillT(t0 float64) *d2.T {
	f0, f1 := e.perimeter.Foci()
	b := line.Bisect(f0, f1)
	tFrom := &triangle.Triangle{
		e.perimeter.Pt1(1.0 / 8.0),
		e.perimeter.Pt1(3.0 / 8.0),
		b.Pt1(0.5),
	}

	//0 ==> 1/8  0.5 ==> 0  1 ==> -1/8
	t0 = t0*-0.25 + 0.125

	tTo := &triangle.Triangle{
		e.perimeter.Pt1(t0),
		e.perimeter.Pt1(0.5 - t0),
		b.Pt1(t0 * 4),
	}
	tfrm, _ := triangle.Transform(tFrom, tTo)
	return &tfrm
}

// Pt2c1 returns a curve that lies inside of the ellipse. All curves in the
// range 0.0 to 1.0 will fill the entire ellipse.
func (e Ellipse) Pt2c1(t0 float64) d2.Pt1 {
	return FillCurve{
		T:          e.fillT(t0),
		EllipseArc: e.perimeter,
	}
}

// Pt2 fulfils d2.Pt2 taking two parametric points and returning a point in the
// ellipse. It does so by mapping the ellipse onto a unit square with points on
// the perimeter every 1/4 rotation corresponding to a corner of the square.
func (e Ellipse) Pt2(t0, t1 float64) d2.Pt {
	t1 = (t1 / 4.0) + (1.0 / 8.0)
	return e.fillT(t0).Pt(e.Pt1(t1))
}

// Pt1 fulfils d2.Pt1 taking one parametric point and returning a point on the
// perimeter.
func (e Ellipse) Pt1(t0 float64) d2.Pt {
	return e.perimeter.Pt1(t0)
}

// Area returns the area of the Ellipse
func (e Ellipse) Area() float64 {
	a := e.SignedArea()
	if a < 0 {
		return -a
	}
	return a
}

// SignedArea returns the area of the ellipse, though the value may be negative
// depending on polarity.
func (e Ellipse) SignedArea() float64 {
	M, m := e.perimeter.Axis()
	return m * M * math.Pi
}

// Perimeter returns the length of the perimeter of the ellipse. This value may
// have a slight error that will grow as the ellipse is elongated.
func (e Ellipse) Perimeter() float64 {
	// https://www.youtube.com/watch?v=5nW3nJhBHL0
	M, m := e.perimeter.Axis()
	d, s := (M - m), (M + m)
	h := (d * d) / (s * s)
	p := 1 + ((3 * h) / (10 + math.Sqrt(4-3*h)))
	p *= math.Pi * s
	return p
}

// PerimeterSeries approximates the perimeter using a series
func (e Ellipse) PerimeterSeries(terms int) float64 {
	M, m := e.perimeter.Axis()
	d, s := (M - m), (M + m)
	h := (d * d) / (s * s)

	// https://en.wikipedia.org/wiki/Ellipse#Circumference
	dn, dd2 := 1.0, 1.0
	// n --> Π i*2+3
	// 3*5*7*9....
	// d1 --> 2^i
	// d2 --> i!
	// hi --> h^i
	n, d1, d2 := 1.0, 1.0, 1.0
	p := 1.0 + h/4.0
	hi := h
	for i := 0; i < terms; i++ {
		dn += 2.0
		dd2 += 1.0
		n *= dn
		d1 *= 2.0
		d2 *= dd2
		hi *= h
		p += (n * hi) / (d1 * d2)
	}
	p *= math.Pi * s
	return p
}

// Centroid returns the center of the ellipse
func (e Ellipse) Centroid() d2.Pt {
	return e.perimeter.Centroid()
}

// Contains returns true if the point f is contained in the ellipse
func (e Ellipse) Contains(pt d2.Pt) bool {
	d := pt.Subtract(e.perimeter.Centroid())

	_, as, ac := e.perimeter.Angle()
	M, m := e.perimeter.Axis()

	h := (d.X*ac + d.Y*as) / M
	v := (d.Y*ac - d.X*as) / m
	return h*h+v*v <= 1
}

// Arc returns the EllipseArc defined by the perimeter
func (e Ellipse) Arc() *ellipsearc.EllipseArc {
	return e.perimeter
}

// BoundingBox fulfills shape.BoundingBoxer. It returns the corners of a Box
// that will exactly contain the Ellipse. It is a wrapper around
// EllipseArc.BoundingBox.
func (e Ellipse) BoundingBox() (min d2.Pt, max d2.Pt) {
	return e.perimeter.BoundingBox()
}

// LineIntersections fulfills line.LineIntersector.
func (e Ellipse) LineIntersections(l line.Line, buf []float64) []float64 {
	return e.perimeter.LineIntersections(l, buf)
}
