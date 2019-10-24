package polygon

import (
	"errors"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/d2"
	d2polygon "github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/shape/plane"
)

// Polygon should be a set of coplanar points
type Polygon []d3.Pt

// D2 converts the polygon to a 2D polygon returning the polygon and the plane.
// If there are not at least 3 points or the points are not coplanar, an error
// is returned.
func (p Polygon) D2() (d2polygon.Polygon, plane.Plane, error) {
	if len(p) < 3 {
		return nil, plane.Plane{}, errors.New("Polygon must have at least 3 points")
	}
	pln := plane.New(p[0], p[1], p[2])
	out := make([]d2.Pt, len(p))
	out[0], _ = pln.Project(p[0])
	out[1], _ = pln.Project(p[1])
	out[2], _ = pln.Project(p[2])
	var v d3.V
	for i, pt := range p[3:] {
		out[i+3], v = pln.Project(pt)
		if !cmpr.Zero(v.Mag2()) {
			return nil, pln, errors.New("Not planar")
		}
	}
	return d2polygon.Polygon(out), pln, nil
}

// From2D takes a 2D polygon and a Plane and creates a 3D polygon.
func From2D(poly2 d2polygon.Polygon, pln plane.Plane) Polygon {
	out := make(Polygon, len(poly2))
	for i, pt := range poly2 {
		out[i] = pln.Convert(pt)
	}
	return out
}
