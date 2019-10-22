package line

// Intersections returns the t values on the line passed in.
type Intersections interface {
	Intersections(Line) []float64
}
