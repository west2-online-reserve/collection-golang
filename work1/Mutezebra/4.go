package main

import (
	"fmt"
	"math"
)

func main() {
	var n int
	_, _ = fmt.Scan(&n)
	if isPrime(n) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}

}

func isPrime(n int) bool {
	if n <= 3 {
		return n > 1
	}
	sqrt := int(math.Sqrt(float64(n)))
	if n%2 == 0 {
		return false
	}
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}
