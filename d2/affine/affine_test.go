package affine

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomtest"
)

func TestCenter(t *testing.T) {
	c := Center{
		{0, 0},
		{1, 0},
	}

	geomtest.Equal(t, d2.Pt{0.5, 0}, c.Centroid())

	c = append(c, d2.Pt{1, 1}, d2.Pt{0, 1})
	geomtest.Equal(t, d2.Pt{0.5, 0.5}, c.Centroid())
}

func TestWeighted(t *testing.T) {
	w := NewWeighted(3)
	w.Add(d2.Pt{1, 0}, d2.Pt{1, 1})
	w.Weight(d2.Pt{3, 1.5}, 2)

	geomtest.Equal(t, d2.Pt{2, 1}, w.Centroid())
}
