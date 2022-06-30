package list

type List[T any] interface {
	Len() int
	Idx(int) T
}

type Empty[T any] struct{}

// Coefficient always returns 0
func (Empty[T]) Idx(idx int) (t T) {
	return
}

func (Empty[T]) Len() int {
	return 0
}
