package line

// Returns the t values on the line
type Intersections interface {
	Intersections(Line) []float64
}
