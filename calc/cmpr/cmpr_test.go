package cmpr

import (
	"testing"

	"github.com/testify/assert"
)

func TestEqualWithin(t *testing.T) {
	assert.True(t, EqualWithin(1, 1.001, 0.01))
	assert.True(t, EqualWithin(1.001, 1, 0.01))
	assert.False(t, EqualWithin(1, 1.001, 0.0001))
	assert.False(t, EqualWithin(1.001, 1, 0.0001))
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal(1, 1.0+1e-6))
	assert.True(t, Equal(1.0+1e-6, 1))
	assert.False(t, Equal(1, 1.0+1e-4))
	assert.False(t, Equal(1.0+1e-4, 1))
}

func TestZero(t *testing.T) {
	assert.True(t, Zero(0))
	assert.True(t, Zero(1e-6))
	assert.True(t, Zero(-1e-6))
	assert.False(t, Zero(1e-4))
	assert.False(t, Zero(-1e-4))
}
