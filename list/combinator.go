package list

import "github.com/adamcolton/geom/calc"

// Holds indexing logic for combining lists.
type Combinator interface {
	Len(a, b int) int
	Idx(n, a, b int) (int, int)
}

// CrossComb does every pairing
type CrossComb struct{}

func (CrossComb) Len(a, b int) int {
	return a * b
}

func (CrossComb) Idx(n, a, b int) (int, int) {
	return n % a, n / a
}

// ModComb does the modulus of the least common multiple.
// If one of the lists is of length one it apply that one to all of the other.
// If the lists are the same size, this is one-to-one.
// Given lists A and B where B is twice the length of A, A will be repeated
// twice.
type ModComb struct{}

func (ModComb) Len(a, b int) int {
	return calc.LCM(a, b)
}

func (ModComb) Idx(n, a, b int) (int, int) {
	return n % a, n % b
}
