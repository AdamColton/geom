// Package affine allows for combinations of points where the sum of the
// coefficients is zero.
package affine

import (
	"github.com/adamcolton/geom/d3"
)

// Weighted allows for each point to be individually weighted before taking the
// sum.
type Weighted struct {
	Pts     []d3.Pt
	Weights []float64
	Sum     float64
}

// Add many points, all with a weight of 1
func (af *Weighted) Add(pts ...d3.Pt) {
	for _, pt := range pts {
		af.Pts = append(af.Pts, pt)
		af.Weights = append(af.Weights, 1)
		af.Sum++
	}
}

// Weight an individual point
func (af *Weighted) Weight(pt d3.Pt, w float64) {
	af.Pts = append(af.Pts, pt)
	af.Weights = append(af.Weights, w)
	af.Sum += w
}

// Get the current weighted sum.
func (af *Weighted) Get() d3.Pt {
	var p d3.Pt
	for i, wp := range af.Pts {
		w := af.Weights[i] / af.Sum
		p.X += wp.X * w
		p.Y += wp.Y * w
		p.Z += wp.Z * w
	}
	return p
}

// Center represents the center of a set of points, all equally weighted
type Center []d3.Pt

// Centroid returns the average point.
func (c Center) Centroid() d3.Pt {
	var cp d3.Pt
	for _, p := range c {
		cp.X += p.X
		cp.Y += p.Y
		cp.Z += p.Z
	}
	l := float64(len(c))
	cp.X /= l
	cp.Y /= l
	cp.Z /= l
	return cp
}
