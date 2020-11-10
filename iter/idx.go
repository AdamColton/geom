package iter

// IdxRange represents an integer range from 0 to IdxRange stepping by one.
type IdxRange int

// IdxIterator coordinates iterating over an IdxRange.
type IdxIterator struct {
	IdxRange
	Idx int
}

// To creates an IdxRange
func To(end int) IdxRange {
	return IdxRange(end)
}

// Iter creates an IdxIterator from an IdxRange
func (r IdxRange) Iter() (*IdxIterator, int) {
	i := &IdxIterator{r, 0}
	return i, i.Idx
}

// Reset the IdxIterator back to 0.
func (i *IdxIterator) Reset() int {
	i.Idx = 0
	return 0
}

// Done will return true if the current index value is greater than or equal to
// IdxRange.
func (i *IdxIterator) Done() bool {
	return i.Idx >= int(i.IdxRange)
}

// Next increments and returns the Idx
func (i *IdxIterator) Next() int {
	i.Idx++
	return i.Idx
}

// Ch returns a chan that will iterate from 0 to IdxRange.
func (i *IdxIterator) Ch() <-chan int {
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

// Ch returns a chan that will iterate from 0 to IdxRange.
func (r IdxRange) Ch() <-chan int {
	i, _ := r.Iter()
	return i.Ch()
}
