package polygon

import (
	"testing"

	d2polygon "github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/stretchr/testify/assert"
)

func Test2DPolygon(t *testing.T) {
	tt := map[string]struct {
		Polygon
		D2  d2polygon.Polygon
		err bool
	}{
		"NilPolygon": {
			Polygon: nil,
			err:     true,
		},
		"EmptyPolygon": {
			Polygon: Polygon{},
			err:     true,
		},
		"1ptPolygon": {
			Polygon: Polygon{{0, 0, 0}},
			err:     true,
		},
		"2ptPolygon": {
			Polygon: Polygon{{0, 0, 0}, {1, 0, 0}},
			err:     true,
		},
		"Triangle": {
			Polygon: Polygon{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}},
			D2:      d2polygon.Polygon{{0, 0}, {1, 0}, {0, 1}},
		},
		"PlanarSquare": {
			Polygon: Polygon{{0, 0, 0}, {1, 0, 0}, {1, 1, 0}, {0, 1, 0}},
			D2:      d2polygon.Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
		},
		"NonPlanarSquare": {
			Polygon: Polygon{{0, 0, 0}, {1, 0, 1}, {1, 1, 0}, {0, 1, 0}},
			err:     true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			d2p, pln, err := tc.Polygon.D2()
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.D2, d2p)
				assert.Equal(t, tc.Polygon, From2D(d2p, pln))
			}
		})
	}
}
