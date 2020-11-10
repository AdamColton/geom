package iter

import "math"

// FloatRange represets a range of float64
type FloatRange struct {
	Start, End, Step float64
}

// FloatIterator iterates over a FloatRange
type FloatIterator struct {
	*IdxIterator
	FloatRange
}

// Iter creates a new FloatIterator from a FloatRange
func (r FloatRange) Iter() (*FloatIterator, float64) {
	i := &FloatIterator{
		&IdxIterator{IdxRange(math.Ceil((r.End - r.Start) / r.Step)), 0},
		r,
	}
	return i, r.Start
}

// Reset to the start of the range
func (i *FloatIterator) Reset() float64 {
	i.Idx = 0
	return i.Start
}

// Next value from the iterator
func (i *FloatIterator) Next() float64 {
	return i.At(float64(i.IdxIterator.Next()))
}

// At returns the value at the given step.
func (i *FloatIterator) At(step float64) float64 {
	return i.Start + i.Step*float64(step)
}

// Include creates a float range that is guarenteed to include the end value.
func Include(end, step float64) FloatRange {
	d := end / step
	if d-math.Floor(d) < 1e-10 {
		end += step * 1e-10
	}
	return FloatRange{
		Start: 0,
		End:   end,
		Step:  step,
	}
}

// Each calls the function passed in for each value in the range.
func (i *FloatIterator) Each(fn func(float64)) {
	for f := i.Reset(); !i.Done(); f = i.Next() {
		fn(f)
	}
}

// Each calls the function passed in for each value in the range.
func (r FloatRange) Each(fn func(float64)) {
	i, _ := r.Iter()
	i.Each(fn)
}

// Ch returns a chan that will iterate over the floats in the range.
func (i *FloatIterator) Ch() <-chan float64 {
	ln := int(i.IdxRange)
	if ln > BufferLen {
		ln = BufferLen
	}
	ch := make(chan float64, ln)
	go func() {
		for f := i.Reset(); !i.Done(); f = i.Next() {
			ch <- f
		}
		close(ch)
	}()
	return ch
}

// Ch returns a chan that will iterate over the floats in the range.
func (r FloatRange) Ch() <-chan float64 {
	i, _ := r.Iter()
	return i.Ch()
}

// FloatChan returns a chan that will iterate over the floats in the range.
func FloatChan(start, end, step float64) <-chan float64 {
	return FloatRange{
		Start: start,
		End:   end,
		Step:  step,
	}.Ch()
}
