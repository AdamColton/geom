package angle

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAngle(t *testing.T) {
	for r := 0.0; r <= 1.0; r += 0.01 {
		rot := Rot(r)
		deg := Deg(rot.Deg())
		rad := Rad(rot.Rad())
		assert.InDelta(t, rot.Rad(), deg.Rad(), 1e-13)
		assert.InDelta(t, rot.Rad(), rad.Rad(), 1e-13)
		assert.InDelta(t, rot.Deg(), deg.Deg(), 1e-13)
		assert.InDelta(t, rot.Deg(), rad.Deg(), 1e-13)
		assert.InDelta(t, rot.Rot(), deg.Rot(), 1e-13)
		assert.InDelta(t, rot.Rot(), rad.Rot(), 1e-13)

		s, c := rot.Sincos()
		assert.Equal(t, s, rot.Sin())
		assert.Equal(t, c, rot.Cos())

		a := Atan(s, c)
		if a < 0 {
			a += math.Pi * 2
		}
		assert.InDelta(t, rot.Rad(), a.Rad(), 1e-10)
	}
}
