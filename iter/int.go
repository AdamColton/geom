package iter

// IntRange a range of ints.
type IntRange struct {
	Start, End, Step int
}

// IntIterator coordinates iterating over an IntRange.
type IntIterator struct {
	*IdxIterator
	IntRange
}

// Iter creates an IntIterator from an IntRange.
func (r IntRange) Iter() (*IntIterator, int) {
	i := &IntIterator{
		&IdxIterator{IdxRange((r.End - r.Start) / r.Step), 0},
		r,
	}
	return i, r.Start
}

// Next int in the range.
func (i *IntIterator) Next() int {
	return i.At(i.IdxIterator.Next())
}

// At returns the value at the provided step.
func (i *IntIterator) At(step int) int {
	return i.Start + i.Step*step
}

// Reset to the start of the range.
func (i *IntIterator) Reset() int {
	i.Idx = 0
	return i.Start
}

// Range creates a range from start to end with a step of 1.
func Range(start, end int) IntRange {
	return IntRange{
		Start: start,
		End:   end,
		Step:  1,
	}
}

// Ch returns a chan that will iterate over the IntRange.
func (i *IntIterator) Ch() <-chan int {
	ln := int(i.IdxRange)
	if ln > BufferLen {
		ln = BufferLen
	}
	ch := make(chan int, ln)
	go func() {
		for f := i.Reset(); !i.Done(); f = i.Next() {
			ch <- f
		}
		close(ch)
	}()
	return ch
}

// Ch returns a chan that will iterate over the IntRange.
func (r IntRange) Ch() <-chan int {
	i, _ := r.Iter()
	return i.Ch()
}
