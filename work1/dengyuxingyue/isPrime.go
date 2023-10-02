package main

import (
	"fmt"
	"math"
)

func isPrime(x int64) bool {

	for i := int64(2); i <= int64(math.Sqrt(float64(x))); i++ {
		if x%i == 0 {
			return false
		}
	}
	return true

}
func main() {

	var n int64
	fmt.Scan(&n)
	if isPrime(n) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
