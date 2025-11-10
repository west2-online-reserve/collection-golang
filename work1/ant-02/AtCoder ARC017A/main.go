package main

import (
	"fmt"
	"math"
)

func isPrime(x int) bool {
	if x == 1 {
		return false
	}
	y := int(math.Sqrt(float64(x)))
	for i := 2; i <= y; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var x int
	fmt.Scan(&x)
	if isPrime(x) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
