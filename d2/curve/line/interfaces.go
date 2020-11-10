package line

// Intersector finds the points where the interface intersects a line.
// The returned values should be relative to the line passed in. The buffer both
// provides reuse to avoid generating garbage and allows for fine tuning. If
// the buffer has a length of 0, all the intersections will be appended. If the
// buffer has a non-zero length, that output will be limited. So if a buffer
// of length 1 is passed in, only the first intersection is returned. But there
// is no guarenteed order.
type Intersector interface {
	LineIntersections(l Line, buf []float64) []float64
}
