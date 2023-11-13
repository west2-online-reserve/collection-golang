package main

import (
	"fmt"
	"math"
)

func isPrime(x int) bool {
	n := int(math.Sqrt(float64(x)))
	for i := 2; i <= n; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}
func main() {
	var x int
	fmt.Scanf("%d", &x)
	if isPrime(x) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}

}
