package line

// LineIntersector returns the t values on the line passed in.
type LineIntersector interface {
	LineIntersections(Line) []float64
}
