package affine

import (
	"testing"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/geomtest"
)

func TestAdd(t *testing.T) {
	w := Weighted{}
	w.Add([]d3.Pt{
		{0, 0, 1},
		{1, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
	}...)

	geomtest.Equal(t, d3.Pt{0.5, 0.5, 0.25}, w.Get())
}

func TestWeight(t *testing.T) {
	w := Weighted{}
	w.Weight(d3.Pt{0, 0, 0}, 2)
	w.Add(d3.Pt{3, 0, 0})

	geomtest.Equal(t, d3.Pt{1, 0, 0}, w.Get())
}

func TestCentroid(t *testing.T) {
	c := Center{
		{3, 1, 4},
		{1, 5, 9},
	}
	geomtest.Equal(t, d3.Pt{2, 3, 6.5}, c.Centroid())
}
