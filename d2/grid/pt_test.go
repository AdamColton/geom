package grid

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/geomerr"

	"github.com/stretchr/testify/assert"
)

func TestAssertEqual(t *testing.T) {
	a := Pt{1, 2}
	b := Pt{1, 2}

	err := a.AssertEqual(b, 1e-10)
	assert.NoError(t, err)

	b = Pt{2, 3}
	err = a.AssertEqual(b, 1e-10)
	assert.Equal(t, "Expected {1 2} got {2 3}", err.Error())

	err = a.AssertEqual(1.0, 1e-10)
	assert.IsType(t, geomerr.ErrTypeMismatch{}, err)

}

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

func TestMask(t *testing.T) {
	p := Pt{10, 11}

	assert.Equal(t, Pt{10, 10}, p.Mask(TwoMask))
}

func TestAspect(t *testing.T) {
	p := Widescreen.Pt(1920)
	assert.Equal(t, Pt{1920, 1080}, p)
}
