package grid

type IteratorFactory interface {
	Iter() Iterator
}

type BaseIterator interface {
	Next() (done bool)
	Pt() Pt
	Idx() int
	Done() bool
	Reset() (done bool)
	PtIdx(Pt) int
	Size() Pt
	Contains(Pt) bool
}

type Iterator interface {
	BaseIterator
	IteratorFactory
	Start() (iter Iterator, done bool)
	Slice() []Pt
	Chan() <-chan Pt
	Each(fn func(int, Pt))
	Until(fn func(int, Pt) bool) bool
}

type BaseIteratorWrapper struct {
	BaseIterator
}

// Start is helper tha can be as the first portion of a classic for loop
func (base BaseIteratorWrapper) Start() (Iterator, bool) {
	return base, base.Reset()
}

// Slice returns all the points in the iterator as a slice.
func (base BaseIteratorWrapper) Slice() []Pt {
	s := make([]Pt, base.Size().Area())
	for i, done := base.Start(); !done; done = i.Next() {
		s[base.Idx()] = i.Pt()
	}
	return s
}

// Each calls the fn for each point in the iterator.
func (base BaseIteratorWrapper) Each(fn func(int, Pt)) {
	for i, done := base.Start(); !done; done = i.Next() {
		fn(i.Idx(), i.Pt())
	}
}

// Until calls fn against each point in the iterator until a point returns true.
// The bool indicates if a value returned true. The iterator will be left at the
// point that returned true.
func (base BaseIteratorWrapper) Until(fn func(int, Pt) bool) bool {
	for i, done := base.Start(); !done; done = i.Next() {
		if fn(i.Idx(), i.Pt()) {
			return true
		}
	}
	return false
}

// Chan runs a go routine that will return the points of the iterator. When all
// the points are consumed the channel is closed. Failing to consume all the
// points will cause a Go routine leak.
func (base BaseIteratorWrapper) Chan() <-chan Pt {
	c := make(chan Pt)
	go func() {
		for i, done := base.Start(); !done; done = i.Next() {
			c <- i.Pt()
		}
		close(c)
	}()
	return c
}

func (base BaseIteratorWrapper) Iter() Iterator {
	return base
}
