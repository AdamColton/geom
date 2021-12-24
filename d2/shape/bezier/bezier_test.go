package bezier

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/bezier"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	b := New(
		[]d2.Pt{{0, 0}, {0, 1}, {0, 1}},
		[]d2.Pt{{1, 1}, {1, 0}, {1, 0}},
	)

	assert.True(t, b.Validate())

	assert.Len(t, b[0], 4)
	assert.Len(t, b[1], 4)

	assert.Equal(t, b[0][0], b[1][3])
	assert.Equal(t, b[0][3], b[1][0])

	m, M := b.BoundingBox()
	geomtest.Equal(t, d2.Pt{0, 0}, m)
	geomtest.Equal(t, d2.Pt{1, 1}, M)

	assert.True(t, b.Contains(d2.Pt{0.5, 0.5}))
	assert.False(t, b.Contains(d2.Pt{0, 1}))

}

func TestFailValidation(t *testing.T) {
	b := BezShape{
		bezier.Bezier{{0, 0}, {0, 1}, {0, 1}},
		bezier.Bezier{{1, 1}, {1, 0}, {1, 0}},
	}
	assert.False(t, b.Validate())
}
