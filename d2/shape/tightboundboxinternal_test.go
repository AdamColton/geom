package shape

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubdiv(t *testing.T) {
	tt := []float64{
		1.0 / 2.0,
		1.0 / 4.0,
		3.0 / 4.0,
		3.0 / 8.0,
		5.0 / 8.0,
		1.0 / 8.0,
		7.0 / 8.0,
		7.0 / 16.0,
		9.0 / 16.0,
		5.0 / 16.0,
		11.0 / 16.0,
		3.0 / 16.0,
		13.0 / 16.0,
		1.0 / 16.0,
		15.0 / 16.0,
	}

	for _, s := range []string{"first_", "memo_"} {
		for n, tc := range tt {
			t.Run(s+fmt.Sprint(n), func(t *testing.T) {
				assert.Equal(t, tc, subdiv(n))
			})
		}
	}
}

func TestSubdivbase(t *testing.T) {
	tt := []int{
		2.0,
		4.0,
		4.0,
		8.0,
		8.0,
		8.0,
		8.0,
		16.0,
		16.0,
		16.0,
		16.0,
		16.0,
		16.0,
		16.0,
		16.0,
	}

	for _, s := range []string{"first_", "memo_"} {
		for i, tc := range tt {
			n := i + 1
			t.Run(s+fmt.Sprint(n), func(t *testing.T) {
				assert.Equal(t, tc, subdivbase(n))
			})
		}
	}
}
