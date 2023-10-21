package main

import (
	"fmt"
	"math"
)

func isPrime(x int) bool {
	for i := 2; i <= int(math.Sqrt(float64(x))); i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}
func main() {
	var n int
	fmt.Scan(&n)
	if isPrime(n) {
		fmt.Print("YES\n")
	} else {
		fmt.Print("NO\n")
	}
}
