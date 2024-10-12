package main

import (
	"fmt"
	"math"
)

func main() {
	var n int
	fmt.Scanf("%d", &n)
	fmt.Printf(isPrime(n))
}

func isPrime(n int) string {
	if n <= 1 {
		return "NO"
	}
	if n == 2 || n == 3 {
		return "YES"
	}
	if n%2 == 0 || n%3 == 0 {
		return "NO"
	}
	for i := 5; i <= int(math.Sqrt(float64(n))); i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return "NO"
		}
	}
	return "YES"
}
