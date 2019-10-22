package line

// Intersections with the given line, expressed as parametric values on the
// line.
type Intersections interface {
	Intersections(Line) []float64
}
