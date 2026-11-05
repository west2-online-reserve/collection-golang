package main

import (
	"fmt"
)

func isPrime(x int) bool {
	if x < 2 {
		return false
	}
	if x == 2 {
		return true
	}
	if x%2 == 0 {
		return false
	}
	for i := 3; i*i <= x; i = i + 2 {
		if x%i == 0 {
			return false
		}
	}
	return true
}
func main() {
	var N int
	fmt.Scan(&N)
	if isPrime(N) == true {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
