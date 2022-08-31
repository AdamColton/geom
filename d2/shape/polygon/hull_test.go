package polygon

import (
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/generate"
	"github.com/adamcolton/geom/geomtest"
)

func TestGrahamScan(t *testing.T) {
	tt := map[string][]d2.Pt{
		"triange":     {{0, 0}, {1, 0}, {1, 1}},
		"convexQuad":  {{0, 0}, {1, 0}, {1, 1}, {0, 1}},
		"concaveQuad": {{0, 0}, {1, 1}, {0, 0.5}, {-1, 1}},
		//"3OnLineQuad": {{0, 0}, {1, 1}, {0, 0.5}, {-1, -1}},
		"regression": {{264.1697, 34.9189}, {320.8487, 34.9189},
			{434.2065, 47.9243}, {448.3762, 60.9297}, {462.5460, 86.9405},
			{462.5460, 99.9459}, {448.3762, 256.0108}, {434.2065, 282.0216},
			{405.8670, 321.0379}, {391.6973, 334.0433}, {349.1881, 373.0595},
			{306.6789, 399.0703}, {278.3395, 412.0757}, {235.8303, 425.0811},
			{221.6605, 425.0811}, {65.7935, 412.0757}, {51.6238, 399.0703},
			{37.4540, 373.0595}, {37.4540, 360.0541}, {79.9632, 203.9892},
			{221.6605, 73.9351}, {250.0000, 47.9243}, {193.3211, 99.9459},
			{51.6238, 230.0000}, {51.6238, 203.9892}, {65.7935, 177.9784},
			{94.1330, 138.9621}, {108.3027, 125.9567}, {193.3211, 60.9297},
			{221.6605, 47.9243}},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			geomtest.Equal(t, AssertConvexHull(tc), GrahamScan(tc...))
		})
	}
}

func TestGrahamScanFuzz(t *testing.T) {
	pts := make([]d2.Pt, 1000)
	for i := range pts {
		p := generate.Pt()
		pts[i] = d2.Polar{M: p.X, A: angle.Rot(p.Y)}.Pt()
	}
	got := GrahamScan(pts...)
	geomtest.Equal(t, AssertConvexHull(pts), got)
}
