package affine

import (
	"github.com/adamcolton/geom/d2"
)

// Weighted finds the center of a set of weighted points.
type Weighted struct {
	Pts     []d2.Pt
	Weights []float64
	Sum     float64
}

// NewWeighted pre-allocates the capacity for the calculation.
func NewWeighted(capacity int) *Weighted {
	return &Weighted{
		Pts:     make([]d2.Pt, 0, capacity),
		Weights: make([]float64, 0, capacity),
	}
}

// Add points to the weighted set, each with a weight of 1.
func (w *Weighted) Add(pts ...d2.Pt) {
	for _, pt := range pts {
		w.Pts = append(w.Pts, pt)
		w.Weights = append(w.Weights, 1)
		w.Sum++
	}
}

// Weight adds a point with a weight
func (w *Weighted) Weight(pt d2.Pt, weight float64) {
	w.Pts = append(w.Pts, pt)
	w.Weights = append(w.Weights, weight)
	w.Sum += weight
}

// Centroid finds the center of the weighted points.
func (w *Weighted) Centroid() d2.Pt {
	var sum d2.Pt
	for i, p := range w.Pts {
		wi := w.Weights[i] / w.Sum
		sum.X += p.X * wi
		sum.Y += p.Y * wi
	}
	return sum
}

// Center of a set of points.
type Center []d2.Pt

// Centroid of the point set.
func (c Center) Centroid() d2.Pt {
	var cp d2.Pt
	for _, p := range c {
		cp.X += p.X
		cp.Y += p.Y
	}
	l := float64(len(c))
	cp.X /= l
	cp.Y /= l
	return cp
}
