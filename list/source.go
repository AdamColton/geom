package list

type Source[T any] struct {
	Fn func(float64) T
	RangeStep
}

func (s Source[T]) Len() int {
	return s.RangeStep.Len()
}

func (s Source[T]) Idx(idx int) T {
	return s.Fn(s.RangeStep.Idx(idx))
}
