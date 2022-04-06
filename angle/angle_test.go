package angle

import (
	"math"
	"reflect"
	"testing"

	"github.com/adamcolton/geom/geomerr"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestAngle(t *testing.T) {
	for r := 0.0; r <= 1.0; r += 0.01 {
		rot := Rot(r)
		deg := Deg(rot.Deg())
		rad := Rad(rot.Rad())

		geomtest.Equal(t, rot, deg)
		geomtest.Equal(t, rot, rad)
		geomtest.Equal(t, deg, rad)

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
		geomtest.Equal(t, 1.0, s*s+c*c)

		a := Atan(s, c)
		if a < 0 {
			a += math.Pi * 2
		}
		geomtest.Equal(t, rot.Rad(), a.Rad())

		if r <= 0.5 {
			geomtest.Equal(t, rad, Acos(c))
		} else {
			geomtest.Equal(t, rad, 2*math.Pi-Acos(c), r)
		}

		geomtest.Equal(t, rad.Rad(), (rad + Tau).Normal().Rad())
		geomtest.Equal(t, rad.Rad(), (rad - Tau).Normal().Rad())
	}

	geomtest.Equal(t, Rad(0), Rot(1))
}

func TestExpectError(t *testing.T) {
	r := Rad(1)
	r2 := Rad(2)
	err := r.AssertEqual(r2, 1e-6)
	assert.Equal(t, geomerr.ErrNotEqual{r, r2}, err)

	i := 5
	err = r.AssertEqual(i, 1e-6)
	assert.Equal(t, geomerr.ErrTypeMismatch{reflect.TypeOf(r), reflect.TypeOf(i)}, err)

}
