// Package fbuf provides a set of helpers around buffers of []float64.
package fbuf

import (
	"reflect"
	"unsafe"
)

// Empty returns a zero length buffer with at least capacity c. If the provided
// buffer has capacity, it will be used otherwise a new one is created.
func Empty(c int, buf []float64) []float64 {
	if cap(buf) >= c {
		return buf[:0]
	}
	return make([]float64, 0, c)
}

// Slice returns a buffer with length c. If the provided buffer has capacity, it
// will be used otherwise a new one is created.
func Slice(c int, buf []float64) []float64 {
	if cap(buf) >= c {
		return buf[:c]
	}
	return make([]float64, c)
}

// Zeros returns a buffer with length c with all values set to 0. If the
// provided buffer has capacity, it will be used otherwise a new one is created.
func Zeros(c int, buf []float64) []float64 {
	if cap(buf) >= c {
		buf = buf[:c]
		for i := range buf {
			buf[i] = 0
		}
		return buf
	}
	return make([]float64, c)
}

// ReduceCapacity sets the capacity to a lower value. This can be useful when
// splitting a buffer to prevent use of the first part of the buffer from
// overflowing into the second part.
func ReduceCapacity(c int, buf []float64) []float64 {
	if c < cap(buf) {
		pv := reflect.ValueOf(&buf)
		sh := (*reflect.SliceHeader)(unsafe.Pointer(pv.Pointer()))
		sh.Cap = c
	}
	return buf
}

// Split a buffer returns two buffers from one. The frist buffer will have at
// least capacity c. The second buffer will have the remainder. If the provided
// buffer does not have a capacity a new buffer is created.
func Split(c int, buf []float64) ([]float64, []float64) {
	if cap(buf) < c {
		return make([]float64, 0, c), buf
	}
	buf = buf[:cap(buf)]
	return buf[:0], buf[c:c]
}
