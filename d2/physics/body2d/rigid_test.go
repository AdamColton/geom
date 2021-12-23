package body2d

import (
	"fmt"
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomtest"
)

func TestT(t *testing.T) {
	tt := []Position{
		{
			V:     d2.V{},
			Angle: 0,
		},
		{
			V:     d2.V{1, 1},
			Angle: 0,
		},
		{
			V:     d2.V{},
			Angle: 3.1415 / 2,
		},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			geomtest.Equal(t, tc.T(), tc.T2())
		})
	}
}
