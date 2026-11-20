package main

import (
	"fmt"
	"math"
)

func isPrime(N int) bool {
	ans := true
	sqrt := int(math.Sqrt(float64(N))) + 1
		for i := 2; i < sqrt; i++ {
			if N%i == 0 {
				ans = false
			}
		}
	return ans
}

func main() {
	var ans bool
	var N int
	fmt.Scan(&N)
	ans = isPrime(N)	
	if ans {
		fmt.Print("YES")
	} else {
		fmt.Print("NO")
	}
}
