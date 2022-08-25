package triangle

import (
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {
	for i := 1; i < 15; i++ {
		n := i + 1
		base := subdivbase(n)
		num := base / 2
		idx := n - num
		if idx%2 == 0 {
			num -= idx + 1
		} else {
			num += idx
		}

		// if idx == 0 {
		// 	x= (num / 2) + d
		// } else {
		// 	num += idx * 2
		// }
		fmt.Println(i, idx, num, base)
	}
}

func subdivbase(n int) int {
	base := 2
	for {
		if n < base {
			return base
		}
		base <<= 1
	}
}
