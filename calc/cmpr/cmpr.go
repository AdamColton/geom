package cmpr

// Tolerance adjusts how close values must be to be considered equal
var Tolerance float64 = 1e-5

// EqualWithin returns true if a and b are within the specified distance of
// eachother.
func EqualWithin(a, b, within float64) bool {
	d := a - b
	return d < within && d > -within
}

// Equal returns true if a and b are within the Tolerance value of eachother.
func Equal(a, b float64) bool {
	return EqualWithin(a, b, Tolerance)
}

// Zero returns true if x is within the Tolerance value of 0.
func Zero(x float64) bool {
	return EqualWithin(x, 0, Tolerance)
}
