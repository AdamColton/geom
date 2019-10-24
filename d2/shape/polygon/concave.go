package polygon

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/triangle"
)

// ConcavePolygon represents a Polygon with at least one concave angle.
type ConcavePolygon struct {
	concave   Polygon
	regular   Polygon
	triangles [][2]*triangle.Triangle
}

// GetTriangles takes triangle indexes from FindTriangles and returns a slice
// of triangles. This can be used to map one polygon to another with the same
// number of sides.
func GetTriangles(triangles [][3]uint32, p Polygon) []*triangle.Triangle {
	ts := make([]*triangle.Triangle, len(triangles))
	for i, t := range triangles {
		ts[i] = &triangle.Triangle{p[t[0]], p[t[1]], p[t[2]]}
	}
	return ts
}

// NewConcavePolygon converts a Polygon to a ConcavePolygon
func NewConcavePolygon(concave Polygon) ConcavePolygon {
	regular := RegularPolygonRadius(d2.Pt{}, 1, 0, len(concave))
	tIdxs := concave.FindTriangles()
	cts := GetTriangles(tIdxs, concave)
	rts := GetTriangles(tIdxs, regular)
	ts := make([][2]*triangle.Triangle, len(tIdxs))
	for i := range cts {
		ts[i][0] = rts[i]
		ts[i][1] = cts[i]
	}

	return ConcavePolygon{
		concave:   concave,
		regular:   regular,
		triangles: ts,
	}
}

// Pt2 returns a point in the ConcavePolygon adhereing to the shape rules
func (c ConcavePolygon) Pt2(t0, t1 float64) d2.Pt {
	pt := c.regular.Pt2(t0, t1)

	for _, ts := range c.triangles {
		if !ts[0].Contains(pt) {
			continue
		}
		tfrm, _ := triangle.Transform(ts[0], ts[1])
		return tfrm.Pt(pt)
	}

	// point is on perimeter
	for _, ts := range c.triangles {
		if line.New(ts[0][0], ts[0][1]).Closest(pt).Distance(pt) < 1E-5 ||
			line.New(ts[0][1], ts[0][2]).Closest(pt).Distance(pt) < 1E-5 ||
			line.New(ts[0][2], ts[0][0]).Closest(pt).Distance(pt) < 1E-5 {
			tfrm, _ := triangle.Transform(ts[0], ts[1])
			return tfrm.Pt(pt)
		}
	}

	return d2.Pt{}
}

// Pt1 returns a point on the perimeter
func (c ConcavePolygon) Pt1(t0 float64) d2.Pt { return c.concave.Pt1(t0) }

// Area of the polygon
func (c ConcavePolygon) Area() float64 { return c.concave.Area() }

// SignedArea returns the Area and may be negative depending on the polarity.
func (c ConcavePolygon) SignedArea() float64 { return c.concave.SignedArea() }

// Perimeter returns the total length of the perimeter
func (c ConcavePolygon) Perimeter() float64 { return c.concave.Perimeter() }

// Contains returns true of the point f is inside of the polygon
func (c ConcavePolygon) Contains(f d2.Pt) bool { return c.concave.Contains(f) }

// Centroid returns the center of mass of the polygon
func (c ConcavePolygon) Centroid() d2.Pt { return c.concave.Centroid() }
