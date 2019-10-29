package plane

import (
	"strings"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d3"
)

// Plane is defined by a point on the plane and a vector perpendicular to plan.
type Plane struct {
	Origin  d3.Pt
	X, Y, N d3.V
}

// New creates a plane from 3 points. The first point is the origin on the
// plane, the second indicates the direction that will be treated as the X axis
// in the plane. The distance between origin and X is normalized. Ref provides a
// third point that lies in the plane.
func New(origin, x, ref d3.Pt) Plane {
	dx := x.Subtract(origin).Normal()
	n := origin.Subtract(ref).Cross(dx).Normal()
	dy := n.Cross(dx)
	return Plane{
		Origin: origin,
		X:      dx,
		Y:      dy,
		N:      n,
	}
}

// Project a point onto the plane. Returns the closes point in the plane and a
// vector from that point to the original 3d point.
func (p Plane) Project(pt d3.Pt) (d2.Pt, d3.V) {
	v := pt.Subtract(p.Origin)
	vOut := p.N.Project(v)
	vIn := v.Subtract(vOut)
	return d2.Pt{vIn.Dot(p.X), vIn.Dot(p.Y)}, vOut
}

// Convert a 2D point in the plane to a 3D point
func (p Plane) Convert(pt d2.Pt) d3.Pt {
	return p.Origin.Add(p.X.Multiply(pt.X)).Add(p.Y.Multiply(pt.Y))
}

func (p Plane) ConvertMany(pts []d2.Pt) []d3.Pt {
	out := make([]d3.Pt, len(pts))
	for i, pt := range pts {
		out[i] = p.Origin.Add(p.X.Multiply(pt.X)).Add(p.Y.Multiply(pt.Y))
	}
	return out
}

// String fulfills Stringer on Plane
func (p Plane) String() string {
	return strings.Join([]string{
		"Plane( Origin:",
		p.Origin.String(),
		" X:",
		p.X.String(),
		" Y:",
		p.Y.String(),
		" N:",
		p.N.String(),
		")",
	}, "")
}
