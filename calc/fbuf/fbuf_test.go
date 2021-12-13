package fbuf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	buf := Empty(10, nil)
	assert.Equal(t, 10, cap(buf))

	buf = Empty(5, buf)
	assert.True(t, cap(buf) >= 5)

	buf = Empty(12, buf)
	assert.True(t, cap(buf) >= 12)
}

func TestSplit(t *testing.T) {
	buf := Empty(15, nil)
	a, b := Split(5, buf)

	assert.True(t, cap(a) >= 5)
	assert.Equal(t, 10, cap(b))

	c, d := Split(12, b)
	assert.True(t, cap(c) >= 12)
	assert.Equal(t, 10, cap(d))
}

func TestSlice(t *testing.T) {
	buf := ([]float64{3, 1, 4})[:0]
	buf = Slice(2, buf)

	assert.Equal(t, []float64{3, 1}, buf)
	assert.Equal(t, []float64{0, 0, 0, 0}, Slice(4, buf))
}

func TestZeros(t *testing.T) {
	buf := ([]float64{3, 1, 4, 1, 5})[:3]
	buf = Zeros(5, buf)
	assert.Equal(t, []float64{0, 0, 0, 0, 0}, buf)
	buf = Zeros(6, buf)
	assert.Equal(t, []float64{0, 0, 0, 0, 0, 0}, buf)
}

func TestReduceCapacity(t *testing.T) {
	buf := []float64{1, 2, 3, 4, 5}
	sub := buf[:3]
	assert.Equal(t, 5, cap(sub))
	sub = ReduceCapacity(3, sub)
	assert.Equal(t, 3, cap(sub))
	assert.Equal(t, 5, cap(buf))
}
