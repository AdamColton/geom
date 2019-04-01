package cc

import (
	"github.com/adamcolton/geom/d3"
)

type affinePoint struct {
	pts     []d3.Pt
	weights []float64
	sum     float64
}

func (af *affinePoint) add(pt d3.Pt) {
	af.pts = append(af.pts, pt)
	af.weights = append(af.weights, 1)
	af.sum++
}

func (af *affinePoint) weight(pt d3.Pt, w float64) {
	af.pts = append(af.pts, pt)
	af.weights = append(af.weights, w)
	af.sum += w
}

func (af *affinePoint) Get() d3.Pt {
	var p d3.Pt
	for i, wp := range af.pts {
		w := af.weights[i] / af.sum
		p.X += wp.X * w
		p.Y += wp.Y * w
		p.Z += wp.Z * w
	}
	return p
}
