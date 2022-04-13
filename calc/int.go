package calc

// GCD finds the greatest common denominator of two ints.
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM finds the least common multiple of two ints.
func LCM(a, b int) int {
	return a * b / GCD(a, b)
}
