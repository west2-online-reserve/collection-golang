package main

import (
	"fmt"
	"math"
)

func main() {
	var n int
	fmt.Scanf("%d", &n)
	if isPrime(n) {
		print("YES")
	} else {
		print("NO")
	}
}
func isPrime(x int) bool {
	var flag = true
	for i := 2; i < int(math.Sqrt(float64(x))); i++ {
		if x%i == 0 {
			flag = false
		}
	}
	return flag
}
