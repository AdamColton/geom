package line

// Intersector finds the points where the interface intersects a line.
// The returned values should be relative to the line passed in.
type Intersector interface {
	LineIntersections(Line) []float64
}
