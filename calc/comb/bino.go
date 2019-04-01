package comb

var memo = make([]int, 256)

func Binomial(n, i int) int {
	if n < 0 || i > n || i < 0 {
		return 0
	}
	if i > n/2 {
		i = n - i
	}
	if i < 2 {
		return (i * n) + (1 - i)
	}

	idx := n - 3
	idx = ((idx * idx) / 4) + i - 2

	if idx >= len(memo) {
		ln := len(memo) << 1
		for ; ln < idx; ln <<= 1 {
		}
		cp := make([]int, ln)
		copy(cp, memo)
		memo = cp
	}

	v := memo[idx]
	if v == 0 {
		v = (n * binomial(n-1, i-1)) / i
		memo[idx] = v
	}
	return v
}

func binomial(n, i int) int {
	if i < 2 {
		return (i * n) + (1 - i)
	}

	idx := n - 3
	idx = ((idx * idx) / 4) + i - 2
	v := memo[idx]
	if v == 0 {
		v = (n * binomial(n-1, i-1)) / i
		memo[idx] = v
	}
	return v
}
