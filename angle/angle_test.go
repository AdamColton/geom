package angle

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestAngle(t *testing.T) {
	for r := 0.0; r <= 1.0; r += 0.01 {
		rot := Rot(r)
		deg := Deg(rot.Deg())
		rad := Rad(rot.Rad())
		geomtest.Equal(t, rot.Rad(), deg.Rad())
		geomtest.Equal(t, rot.Rad(), deg.Rad())
		geomtest.Equal(t, rot.Rad(), rad.Rad())
		geomtest.Equal(t, rot.Deg(), deg.Deg())
		geomtest.Equal(t, rot.Deg(), rad.Deg())
		geomtest.Equal(t, rot.Rot(), deg.Rot())
		geomtest.Equal(t, rot.Rot(), rad.Rot())

		s, c := rot.Sincos()
		assert.Equal(t, s, rot.Sin())
		assert.Equal(t, c, rot.Cos())

		a := Atan(s, c)
		if a < 0 {
			a += math.Pi * 2
		}
		geomtest.Equal(t, rot.Rad(), a.Rad())
	}
}
