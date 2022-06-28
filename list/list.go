package list

type List[T any] interface {
	Len() int
	Idx(int) T
}

type Slice[T any] []T

func (s Slice[T]) Len() int {
	return len(s)
}

func (s Slice[T]) Idx(idx int) (t T) {
	if idx < 0 || idx >= len(s) {
		return
	}
	return s[idx]
}

type Empty[T any] struct{}

// Coefficient always returns 0
func (Empty[T]) Idx(idx int) (t T) {
	return
}

func (Empty[T]) Len() int {
	return 0
}
