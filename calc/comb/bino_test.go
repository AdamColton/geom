package comb

import (
	"testing"
)

func Classic(n, i int) int {
	// https://math.stackexchange.com/questions/202554/how-do-i-compute-binomial-coefficients-efficiently
	if n < 0 || i > n || i < 0 {
		return 0
	} else if i == 0 {
		return 1
	} else if i > n/2 {
		return Classic(n, n-i)
	}

	return n * Classic(n-1, i-1) / i
}

type key [2]int

var mapMemo = make(map[key]int)

func MapMemo(n, i int) int {
	if n < 0 || i > n || i < 0 {
		return 0
	}
	if i > n/2 {
		i = n - i
	}
	if i < 2 {
		return (i * n) + (1 - i)
	}

	k := key{n, i}
	if v, found := mapMemo[k]; found {
		return v
	}

	v := n * MapMemo(n-1, i-1) / i
	mapMemo[k] = v
	return v
}

func TestAgainstClassic(t *testing.T) {
	for n := 0; n < 20; n++ {
		for i := 0; i <= n; i++ {
			b := Binomial(n, i)
			k := Classic(n, i)
			if b != k {
				t.Error(n, i)
			}
		}
	}
}

func BenchmarkMemo(b *testing.B) {
	n := 0
	i := -1
	for iter := 0; iter < b.N; iter++ {
		i++
		if i > n {
			n, i = n+1, 0
			if n > 100 {
				n, i = 0, 0
			}
		}
		Binomial(n, i)
	}
}

func BenchmarkMap(b *testing.B) {
	n := 0
	i := -1
	for iter := 0; iter < b.N; iter++ {
		i++
		if i > n {
			n, i = n+1, 0
			if n > 100 {
				n, i = 0, 0
			}
		}
		MapMemo(n, i)
	}
}

func BenchmarkClassic(b *testing.B) {
	n := 0
	i := -1
	for iter := 0; iter < b.N; iter++ {
		i++
		if i > n {
			n, i = n+1, 0
			if n > 100 {
				n, i = 0, 0
			}
		}
		Classic(n, i)
	}
}
