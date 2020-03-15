package helper

import (
	"fmt"
	"math/rand"
)

// TODO: 分布の偏りをなくす
func RandomChoice(min int, max int, n int) []int {
	vals := make([]int, n)
	diff := max - min
	maxstep := float64(diff) / float64(n)
	if maxstep < 1.0 {
		panic(fmt.Sprintf(
			"Can not choice %d values from range min: %d, max: %d", n, min, max))
	}

	vals[0] = min
	for i := 1; i < n; i++ {
		rand.Seed(int64(i))
		vals[i] = vals[i-1] + int(maxstep*rand.Float64()+1.0)
	}

	return vals
}
