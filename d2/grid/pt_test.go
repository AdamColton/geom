package grid

import (
	"testing"

	"github.com/adamcolton/geom/d2"

	"github.com/stretchr/testify/assert"
)

func TestBasicMethods(t *testing.T) {
	assert.Equal(t, Pt{5, 6}, Pt{5, 3}.Add(Pt{0, 3}))

	assert.Equal(t, 12, Pt{3, 4}.Area())
	assert.Equal(t, 12, Pt{-3, 4}.Area())

	assert.Equal(t, d2.D2{1, 2}, Pt{1, 2}.D2())

	assert.Equal(t, Pt{10, 20}, Pt{1, 2}.Multiply(10))
}

func TestScale(t *testing.T) {
	s := Scale{
		X:  .1,
		Y:  .2,
		DX: 5,
		DY: 7,
	}

	t0, t1 := s.T(Pt{3, 2})
	assert.Equal(t, 5.3, t0)
	assert.Equal(t, 7.4, t1)
}
