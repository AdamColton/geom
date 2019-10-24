package polygon

import (
	"sort"

	"github.com/adamcolton/geom/d2"
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
