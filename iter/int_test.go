package iter_test

import (
	"testing"

	"github.com/adamcolton/geom/iter"
	"github.com/testify/assert"
)

func TestIdx(t *testing.T) {
	reset := iter.BufferLen
	iter.BufferLen = 3
	defer func() {
		iter.BufferLen = reset
	}()

	c := 0
	for i := range iter.To(6).Ch() {
		assert.Equal(t, c, i)
		c++
	}
	assert.Equal(t, c, 6)
}

func TestIntStep(t *testing.T) {
	expected := []int{3, 5, 7}
	r := iter.IntRange{3, 9, 2}
	c := 0
	for it, i := r.Iter(); !it.Done(); i = it.Next() {
		assert.Equal(t, expected[it.Idx], i)
		c++
	}
	assert.Equal(t, c, len(expected))
}

func TestRange(t *testing.T) {
	reset := iter.BufferLen
	iter.BufferLen = 3
	defer func() {
		iter.BufferLen = reset
	}()

	expected := []int{10, 11, 12, 13, 14}
	c := 0
	for i := range iter.Range(10, 15).Ch() {
		assert.Equal(t, expected[c], i)
		c++
	}
	assert.Equal(t, c, len(expected))
}
