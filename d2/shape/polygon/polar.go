package polygon

import (
	"sort"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/affine"
)

// PolarPolygon is useful in constructing a Polygon when the order of the
// vertexes is not known. It does not fulfill shape.
type PolarPolygon []d2.Polar

// Len returns the number of points, fulfills sort.Interface
func (p PolarPolygon) Len() int { return len(p) }

// Less fulfills sort.Interface
func (p PolarPolygon) Less(i, j int) bool { return p[i].A < p[j].A }

// Swap fulfills sort.Interface
func (p PolarPolygon) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Sort the PolarPolygon by angle
func (p PolarPolygon) Sort() {
	sort.Sort(p)
}

// Polygon converts the PolarPolygon to a Polygon
func (p PolarPolygon) Polygon(center d2.Pt) Polygon {
	ply := make(Polygon, len(p))
	for i, plr := range p {
		ply[i] = center.Add(plr.V())
	}
	return ply
}

// Create a new PolarPolygon from a set of points. The result will be translated
// so the center is (0,0) and the previous center returned.
func NewPolar(pts []d2.Pt) (polarPolygon PolarPolygon, center d2.Pt) {
	center = affine.Center(pts).Centroid()
	polarPolygon = make(PolarPolygon, 0, len(pts))
	for _, pt := range pts {
		polarPolygon = append(polarPolygon, pt.Subtract(center).Polar())
	}
	polarPolygon.Sort()
	return polarPolygon, center
}
