package geomerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceErrs(t *testing.T) {
	a := []int{1, 2, 0, 4}
	b := []int{1, 2, 3}
	fn := func(i int) error {
		return NewNotEqual(a[i] == b[i], a[i], b[i])
	}

	err := NewSliceErrs(len(a), len(b), fn)

	assert.Equal(t, "Lengths do not match: Expected 4 got 3\n\t2: Expected 0 got 3", err.Error())

	a = []int{1, 2, 3}
	err = NewSliceErrs(len(a), len(b), fn)
	assert.NoError(t, err)

	a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	b = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	err = NewSliceErrs(len(a), len(b), fn)
	assert.Equal(t, "Lengths do not match: Expected 11 got 12\n\t0: Expected 1 got 2\n\t1: Expected 2 got 3\n\t2: Expected 3 got 4\n\t3: Expected 4 got 5\n\t4: Expected 5 got 6\n\t5: Expected 6 got 7\n\t6: Expected 7 got 8\n\t7: Expected 8 got 9\n\t8: Expected 9 got 10\nOmitting 2 more", err.Error())

	se := SliceErrs{}
	se = se.AppendF(1, "%s is a test %d", "this", 123)
	assert.Equal(t, "\t1: this is a test 123", se.Error())

	a = []int{1, 2, 3}
	b = []int{1, 2, 3, 4, 5}
	err = NewSliceErrs(len(a), -1, fn)
	assert.NoError(t, err)
}

func TestTypeMismatch(t *testing.T) {
	err := NewTypeMismatch("test", 1.0)
	assert.Equal(t, "Types do not match: expected \"string\", got \"float64\"", err.Error())

	err = TypeMismatch("test", 1.0)
	assert.Equal(t, "Types do not match: expected \"string\", got \"float64\"", err.Error())

	err = NewTypeMismatch("test", "foo")
	assert.NoError(t, err)
}
