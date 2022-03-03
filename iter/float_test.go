package iter_test

import (
	"testing"

	"github.com/adamcolton/geom/iter"
	"github.com/stretchr/testify/assert"
)

func TestFloat(t *testing.T) {
	expected := []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1}
	r := iter.Include(1, 0.1)
	c := 0
	for it, i := r.Iter(); !it.Done(); i = it.Next() {
		assert.InDelta(t, expected[it.Idx], i, 1e-6)
		c++
	}
	assert.Equal(t, c, len(expected))

	c = 0
	r.Each(func(f float64) {
		assert.InDelta(t, expected[c], f, 1e-6)
		c++
	})
	assert.Equal(t, c, len(expected))

	reset := iter.BufferLen
	iter.BufferLen = 3
	defer func() {
		iter.BufferLen = reset
	}()

	c = 0
	for f := range iter.FloatChan(0, 1.05, 0.1) {
		assert.InDelta(t, expected[c], f, 1e-6)
		c++
	}
	assert.Equal(t, c, len(expected))
}
