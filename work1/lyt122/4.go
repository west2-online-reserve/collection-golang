package main

import (
	"fmt"
	"math"
)

func isPrime(x int) bool {
	if x < 2 {
		return false
	}
	maxDivisor := int(math.Sqrt(float64(x)))
	for i := 2; i <= maxDivisor; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}
func main() {
	var number int
	fmt.Scan(&number)
	result := isPrime(number)
	if result {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
